package csv

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/user/gin-microservice-boilerplate/models"
)

var (
	ErrInvalidChartPayload = errors.New("invalid chart payload")
)

type chartStrategy struct{}

type chartPayload struct {
	Categories []string      `json:"categories"`
	Series     []chartSeries `json:"series"`
}

type chartSeries struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

func NewChartStrategy() *chartStrategy {
	return &chartStrategy{}
}

func (s *chartStrategy) Generate(ctx context.Context, req models.ExportRequest) (*models.ExportFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	delimiter, err := delimiterRune(req.Delimiter)
	if err != nil {
		return nil, err
	}

	payload, err := decodeChartPayload(req.Payload)
	if err != nil {
		return nil, err
	}

	data, err := renderChartCSV(ctx, payload, delimiter)
	if err != nil {
		return nil, err
	}

	return &models.ExportFile{
		FileName:    ensureCSVFileName(req.FileName),
		ContentType: "text/csv; charset=utf-8",
		Data:        data,
	}, nil
}

func decodeChartPayload(raw interface{}) (*chartPayload, error) {
	var encoded []byte

	switch v := raw.(type) {
	case nil:
		return nil, fmt.Errorf("%w: payload is required", ErrInvalidChartPayload)
	case json.RawMessage:
		encoded = v
	case []byte:
		encoded = v
	default:
		var err error
		encoded, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("%w: malformed json", ErrInvalidChartPayload)
		}
	}

	var payload chartPayload
	if err := json.Unmarshal(encoded, &payload); err != nil {
		return nil, fmt.Errorf("%w: malformed json", ErrInvalidChartPayload)
	}

	if len(payload.Categories) == 0 {
		return nil, fmt.Errorf("%w: categories is required", ErrInvalidChartPayload)
	}
	if len(payload.Series) == 0 {
		return nil, fmt.Errorf("%w: series is required", ErrInvalidChartPayload)
	}

	for idx, series := range payload.Series {
		if strings.TrimSpace(series.Name) == "" {
			return nil, fmt.Errorf("%w: series[%d].name is required", ErrInvalidChartPayload, idx)
		}
		if len(series.Values) != len(payload.Categories) {
			return nil, fmt.Errorf(
				"%w: series[%d].values length must match categories length",
				ErrInvalidChartPayload,
				idx,
			)
		}
	}

	return &payload, nil
}

func renderChartCSV(ctx context.Context, payload *chartPayload, delimiter rune) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	writer.Comma = delimiter
	writer.UseCRLF = false

	header := make([]string, 0, len(payload.Series)+1)
	header = append(header, "category")
	for _, series := range payload.Series {
		header = append(header, strings.TrimSpace(series.Name))
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for row := range payload.Categories {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		record := make([]string, 0, len(payload.Series)+1)
		record = append(record, payload.Categories[row])

		for _, series := range payload.Series {
			value := strconv.FormatFloat(series.Values[row], 'f', -1, 64)
			record = append(record, value)
		}

		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func delimiterRune(delimiter string) (rune, error) {
	if delimiter == "" {
		return ',', nil
	}

	runes := []rune(delimiter)
	if len(runes) == 1 {
		return runes[0], nil
	}
	return 0, fmt.Errorf("%w: delimiter must have a single character", ErrInvalidChartPayload)
}

func ensureCSVFileName(fileName string) string {
	trimmed := strings.TrimSpace(fileName)
	if trimmed == "" {
		return "chart_export.csv"
	}
	if strings.HasSuffix(strings.ToLower(trimmed), ".csv") {
		return trimmed
	}
	return trimmed + ".csv"
}
