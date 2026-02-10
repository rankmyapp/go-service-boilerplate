package pdf

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/go-pdf/fpdf"

	"github.com/user/gin-microservice-boilerplate/models"
)

var (
	ErrInvalidTablePayload = errors.New("invalid table payload")
)

type tableStrategy struct{}

type tablePayload struct {
	Title   string          `json:"title,omitempty"`
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

func NewTableStrategy() *tableStrategy {
	return &tableStrategy{}
}

func (s *tableStrategy) Generate(ctx context.Context, req models.ExportRequest) (*models.ExportFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	payload, err := decodeTablePayload(req.Payload)
	if err != nil {
		return nil, err
	}

	data, err := renderTablePDF(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &models.ExportFile{
		FileName:    ensurePDFFileName(req.FileName, "table_export.pdf"),
		ContentType: "application/pdf",
		Data:        data,
	}, nil
}

func decodeTablePayload(raw interface{}) (*tablePayload, error) {
	var encoded []byte

	switch v := raw.(type) {
	case nil:
		return nil, fmt.Errorf("%w: payload is required", ErrInvalidTablePayload)
	case json.RawMessage:
		encoded = v
	case []byte:
		encoded = v
	default:
		var err error
		encoded, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("%w: malformed json", ErrInvalidTablePayload)
		}
	}

	var payload tablePayload
	if err := json.Unmarshal(encoded, &payload); err != nil {
		return nil, fmt.Errorf("%w: malformed json", ErrInvalidTablePayload)
	}

	if len(payload.Columns) == 0 {
		return nil, fmt.Errorf("%w: columns is required", ErrInvalidTablePayload)
	}

	for idx, col := range payload.Columns {
		if strings.TrimSpace(col) == "" {
			return nil, fmt.Errorf("%w: columns[%d] is required", ErrInvalidTablePayload, idx)
		}
	}

	for idx, row := range payload.Rows {
		if len(row) != len(payload.Columns) {
			return nil, fmt.Errorf(
				"%w: rows[%d] length must match columns length",
				ErrInvalidTablePayload,
				idx,
			)
		}
	}

	return &payload, nil
}

func renderTablePDF(ctx context.Context, payload *tablePayload) ([]byte, error) {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(false, 10)
	pdf.AddPage()

	title := strings.TrimSpace(payload.Title)
	if title == "" {
		title = "Exported Table"
	}

	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(31, 41, 55)
	pdf.CellFormat(0, 8, title, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	pageW, pageH := pdf.GetPageSize()
	left, _, right, _ := pdf.GetMargins()
	usableW := pageW - left - right
	bottomLimit := pageH - 12
	lineH := 5.5

	colWidths := computeColumnWidths(payload.Columns, usableW)

	drawHeader := func() {
		pdf.SetFont("Arial", "B", 10)
		pdf.SetFillColor(245, 247, 250)
		pdf.SetTextColor(55, 65, 81)
		drawPDFRow(pdf, payload.Columns, colWidths, lineH, true)
	}

	drawHeader()
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(55, 65, 81)

	for _, row := range payload.Rows {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		cells := make([]string, len(row))
		for i, cell := range row {
			cells[i] = tableCellToString(cell)
		}

		rowH := calcRowHeight(pdf, cells, colWidths, lineH)
		if pdf.GetY()+rowH > bottomLimit {
			pdf.AddPage()
			pdf.SetFont("Arial", "B", 10)
			drawHeader()
			pdf.SetFont("Arial", "", 9)
		}

		drawPDFRow(pdf, cells, colWidths, lineH, false)
	}

	var out bytes.Buffer
	if err := pdf.Output(&out); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func computeColumnWidths(columns []string, totalWidth float64) []float64 {
	widths := make([]float64, len(columns))
	weights := make([]float64, len(columns))
	totalWeight := 0.0

	for i, col := range columns {
		name := strings.ToLower(strings.TrimSpace(col))
		weight := 1.0
		switch {
		case strings.Contains(name, "explanation"), strings.Contains(name, "description"):
			weight = 3.2
		case strings.Contains(name, "expected"), strings.Contains(name, "variation"):
			weight = 1.5
		case strings.Contains(name, "date"):
			weight = 1.2
		case strings.Contains(name, "metric"):
			weight = 1.3
		}
		weights[i] = weight
		totalWeight += weight
	}

	for i := range columns {
		widths[i] = totalWidth * (weights[i] / totalWeight)
	}

	return widths
}

func calcRowHeight(pdf *fpdf.Fpdf, cells []string, colWidths []float64, lineH float64) float64 {
	maxLines := 1
	for i, text := range cells {
		lines := pdf.SplitLines([]byte(text), colWidths[i]-2)
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}
	return lineH * float64(maxLines)
}

func drawPDFRow(pdf *fpdf.Fpdf, cells []string, colWidths []float64, lineH float64, fill bool) {
	startX, startY := pdf.GetX(), pdf.GetY()
	rowH := calcRowHeight(pdf, cells, colWidths, lineH)

	x := startX
	for i, text := range cells {
		w := colWidths[i]

		pdf.SetXY(x, startY)
		if fill {
			pdf.Rect(x, startY, w, rowH, "FD")
		} else {
			pdf.Rect(x, startY, w, rowH, "D")
		}

		pdf.SetXY(x+1, startY+0.7)
		pdf.MultiCell(w-2, lineH, text, "", "L", false)
		x += w
	}

	pdf.SetXY(startX, startY+rowH)
}

func tableCellToString(v interface{}) string {
	switch cell := v.(type) {
	case nil:
		return ""
	case string:
		return cell
	case float64:
		if math.Abs(cell-math.Round(cell)) < 1e-9 {
			return strconv.FormatInt(int64(math.Round(cell)), 10)
		}
		return strconv.FormatFloat(cell, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(cell)
	default:
		return fmt.Sprintf("%v", cell)
	}
}

func ensurePDFFileName(fileName, fallback string) string {
	trimmed := strings.TrimSpace(fileName)
	if trimmed == "" {
		return fallback
	}
	if strings.HasSuffix(strings.ToLower(trimmed), ".pdf") {
		return trimmed
	}
	return trimmed + ".pdf"
}
