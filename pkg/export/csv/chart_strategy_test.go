package csv

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/user/gin-microservice-boilerplate/models"
)

func TestGenerate_Success(t *testing.T) {
	strategy := NewChartStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceChart,
		Delimiter:  ",",
		Payload: []byte(`{
			"categories": ["2026-01-01","2026-01-02"],
			"series": [
				{"name":"views","values":[10,20]},
				{"name":"conversion_rate","values":[1.5,2]}
			]
		}`),
	}

	file, err := strategy.Generate(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, file)

	assert.Equal(t, "chart_export.csv", file.FileName)
	assert.Equal(t, "text/csv; charset=utf-8", file.ContentType)
	assert.Equal(
		t,
		"category,views,conversion_rate\n2026-01-01,10,1.5\n2026-01-02,20,2\n",
		string(file.Data),
	)
}

func TestGenerate_CustomDelimiterAndFileName(t *testing.T) {
	strategy := NewChartStrategy()

	req := models.ExportRequest{
		FileName:   "my-report",
		Delimiter:  ";",
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceChart,
		Payload: []byte(`{
			"categories": ["Jan","Feb"],
			"series": [
				{"name":"revenue","values":[100,200]}
			]
		}`),
	}

	file, err := strategy.Generate(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "my-report.csv", file.FileName)
	assert.Equal(t, "category;revenue\nJan;100\nFeb;200\n", string(file.Data))
}

func TestGenerate_InvalidDelimiter(t *testing.T) {
	strategy := NewChartStrategy()

	req := models.ExportRequest{
		Delimiter:  ";;",
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceChart,
		Payload: []byte(`{
			"categories": ["Jan"],
			"series": [{"name":"views","values":[1]}]
		}`),
	}

	_, err := strategy.Generate(context.Background(), req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidChartPayload))
}

func TestGenerate_InvalidPayload(t *testing.T) {
	strategy := NewChartStrategy()

	req := models.ExportRequest{
		Delimiter:  ",",
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceChart,
		Payload: []byte(`{
			"categories": ["Jan","Feb"],
			"series": [{"name":"views","values":[1]}]
		}`),
	}

	_, err := strategy.Generate(context.Background(), req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidChartPayload))
}

func TestGenerate_ContextCancelled(t *testing.T) {
	strategy := NewChartStrategy()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := models.ExportRequest{
		Delimiter:  ",",
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceChart,
		Payload: []byte(`{
			"categories": ["Jan"],
			"series": [{"name":"views","values":[1]}]
		}`),
	}

	_, err := strategy.Generate(ctx, req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
}
