package jpeg

import (
	"bytes"
	"context"
	stdjpeg "image/jpeg"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/user/gin-microservice-boilerplate/models"
)

func TestGenerate_Success(t *testing.T) {
	strategy := NewChartStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatJPEG,
		SourceType: models.ExportSourceChart,
		Payload: map[string]interface{}{
			"categories": []interface{}{"2026-01-01", "2026-01-02", "2026-01-03"},
			"series": []interface{}{
				map[string]interface{}{
					"name":   "downloads",
					"values": []interface{}{120.0, 180.0, 260.0},
				},
			},
		},
	}

	file, err := strategy.Generate(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, file)

	assert.Equal(t, "chart_export.jpeg", file.FileName)
	assert.Equal(t, "image/jpeg", file.ContentType)
	require.Greater(t, len(file.Data), 2)
	assert.Equal(t, byte(0xFF), file.Data[0])
	assert.Equal(t, byte(0xD8), file.Data[1])

	img, err := stdjpeg.Decode(bytes.NewReader(file.Data))
	require.NoError(t, err)
	assert.Equal(t, defaultWidth, img.Bounds().Dx())
	assert.Equal(t, defaultHeight, img.Bounds().Dy())
}

func TestGenerate_CustomSizeAndFileName(t *testing.T) {
	strategy := NewChartStrategy()

	req := models.ExportRequest{
		FileName:   "anomalies_chart",
		Format:     models.ExportFormatJPEG,
		SourceType: models.ExportSourceChart,
		Payload: map[string]interface{}{
			"width":      900.0,
			"height":     500.0,
			"quality":    75.0,
			"categories": []interface{}{"A", "B"},
			"series": []interface{}{
				map[string]interface{}{
					"name":   "downloads",
					"values": []interface{}{100.0, 120.0},
				},
			},
		},
	}

	file, err := strategy.Generate(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "anomalies_chart.jpeg", file.FileName)
	img, err := stdjpeg.Decode(bytes.NewReader(file.Data))
	require.NoError(t, err)
	assert.Equal(t, 900, img.Bounds().Dx())
	assert.Equal(t, 500, img.Bounds().Dy())
}

func TestGenerate_InvalidPayload(t *testing.T) {
	strategy := NewChartStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatJPEG,
		SourceType: models.ExportSourceChart,
		Payload: map[string]interface{}{
			"categories": []interface{}{"A", "B"},
			"series": []interface{}{
				map[string]interface{}{
					"name":   "downloads",
					"values": []interface{}{100.0},
				},
			},
		},
	}

	_, err := strategy.Generate(context.Background(), req)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidChartPayload)
}

func TestGenerate_ContextCancelled(t *testing.T) {
	strategy := NewChartStrategy()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := models.ExportRequest{
		Format:     models.ExportFormatJPEG,
		SourceType: models.ExportSourceChart,
		Payload: map[string]interface{}{
			"categories": []interface{}{"A"},
			"series": []interface{}{
				map[string]interface{}{
					"name":   "downloads",
					"values": []interface{}{1.0},
				},
			},
		},
	}

	_, err := strategy.Generate(ctx, req)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}
