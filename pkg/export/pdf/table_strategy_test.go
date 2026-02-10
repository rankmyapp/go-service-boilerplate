package pdf

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/user/gin-microservice-boilerplate/models"
)

func TestGenerate_Success(t *testing.T) {
	strategy := NewTableStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatPDF,
		SourceType: models.ExportSourceTable,
		FileName:   "anomalies_export",
		Payload: map[string]interface{}{
			"title":   "Detailed historical anomalies",
			"columns": []interface{}{"date", "metric", "value", "intensity", "variation_percent", "expected_value", "explanation"},
			"rows": []interface{}{
				[]interface{}{
					"21/01/2026",
					"Downloads",
					"1,560",
					"High",
					"+506.6%",
					"185 - 381",
					"Observed value above historical pattern.",
				},
			},
		},
	}

	file, err := strategy.Generate(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, file)

	assert.Equal(t, "anomalies_export.pdf", file.FileName)
	assert.Equal(t, "application/pdf", file.ContentType)
	require.Greater(t, len(file.Data), 4)
	assert.True(t, bytes.HasPrefix(file.Data, []byte("%PDF")))
}

func TestGenerate_DefaultFileName(t *testing.T) {
	strategy := NewTableStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatPDF,
		SourceType: models.ExportSourceTable,
		Payload: map[string]interface{}{
			"columns": []interface{}{"name"},
			"rows":    []interface{}{[]interface{}{"Alice"}},
		},
	}

	file, err := strategy.Generate(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "table_export.pdf", file.FileName)
}

func TestGenerate_InvalidPayload(t *testing.T) {
	strategy := NewTableStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatPDF,
		SourceType: models.ExportSourceTable,
		Payload: map[string]interface{}{
			"columns": []interface{}{"name", "email"},
			"rows": []interface{}{
				[]interface{}{"Alice"},
			},
		},
	}

	_, err := strategy.Generate(context.Background(), req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidTablePayload))
}

func TestGenerate_ContextCancelled(t *testing.T) {
	strategy := NewTableStrategy()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := models.ExportRequest{
		Format:     models.ExportFormatPDF,
		SourceType: models.ExportSourceTable,
		Payload: map[string]interface{}{
			"columns": []interface{}{"name"},
			"rows":    []interface{}{[]interface{}{"Alice"}},
		},
	}

	_, err := strategy.Generate(ctx, req)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}
