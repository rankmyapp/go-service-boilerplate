package csv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/user/gin-microservice-boilerplate/models"
)

var (
	ErrInvalidTablePayload = errors.New("invalid table payload")
)

type tableStrategy struct{}

type tablePayload struct {
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

	delimiter, err := delimiterRune(req.Delimiter)
	if err != nil {
		return nil, err
	}

	payload, err := decodeTablePayload(req.Payload)
	if err != nil {
		return nil, err
	}

	data, err := renderTableCSV(ctx, payload, delimiter)
	if err != nil {
		return nil, err
	}

	return &models.ExportFile{
		FileName:    ensureCSVFileName(req.FileName, "table_export.csv"),
		ContentType: "text/csv; charset=utf-8",
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

func renderTableCSV(ctx context.Context, payload *tablePayload, delimiter rune) ([]byte, error) {
	rows := make([][]string, 0, len(payload.Rows)+1)

	header := make([]string, len(payload.Columns))
	copy(header, payload.Columns)
	rows = append(rows, header)

	for _, row := range payload.Rows {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		record := make([]string, 0, len(row))
		for _, cell := range row {
			record = append(record, tableCellToString(cell))
		}
		rows = append(rows, record)
	}

	return writeCSV(rows, delimiter)
}

func tableCellToString(v interface{}) string {
	switch cell := v.(type) {
	case nil:
		return ""
	case string:
		return cell
	case float64:
		return strconv.FormatFloat(cell, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(cell)
	default:
		return fmt.Sprintf("%v", cell)
	}
}
