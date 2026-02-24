package csv

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/user/gin-microservice-boilerplate/models"
)

func TestGenerateTable_Success(t *testing.T) {
	strategy := NewTableStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceTable,
		Delimiter:  ",",
		Payload: map[string]interface{}{
			"columns": []interface{}{"name", "email", "active"},
			"rows": []interface{}{
				[]interface{}{"Alice", "alice@example.com", true},
				[]interface{}{"Bob", "bob@example.com", false},
			},
		},
	}

	file, err := strategy.Generate(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, file)

	assert.Equal(t, "table_export.csv", file.FileName)
	assert.Equal(t, "text/csv; charset=utf-8", file.ContentType)
	assert.Equal(
		t,
		"name,email,active\nAlice,alice@example.com,true\nBob,bob@example.com,false\n",
		string(file.Data),
	)
}

func TestGenerateTable_CustomDelimiterAndFileName(t *testing.T) {
	strategy := NewTableStrategy()

	req := models.ExportRequest{
		FileName:   "users_export",
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceTable,
		Delimiter:  ";",
		Payload: map[string]interface{}{
			"columns": []interface{}{"name", "score"},
			"rows": []interface{}{
				[]interface{}{"Alice", 9.5},
				[]interface{}{"Bob", 8.0},
			},
		},
	}

	file, err := strategy.Generate(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "users_export.csv", file.FileName)
	assert.Equal(t, "name;score\nAlice;9.5\nBob;8\n", string(file.Data))
}

func TestGenerateTable_InvalidColumns(t *testing.T) {
	strategy := NewTableStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceTable,
		Payload: map[string]interface{}{
			"columns": []interface{}{"name", ""},
			"rows": []interface{}{
				[]interface{}{"Alice", "alice@example.com"},
			},
		},
	}

	_, err := strategy.Generate(context.Background(), req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidTablePayload))
}

func TestGenerateTable_InvalidRowLength(t *testing.T) {
	strategy := NewTableStrategy()

	req := models.ExportRequest{
		Format:     models.ExportFormatCSV,
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

func TestGenerateTable_ContextCancelled(t *testing.T) {
	strategy := NewTableStrategy()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := models.ExportRequest{
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceTable,
		Payload: map[string]interface{}{
			"columns": []interface{}{"name"},
			"rows": []interface{}{
				[]interface{}{"Alice"},
			},
		},
	}

	_, err := strategy.Generate(ctx, req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
}
