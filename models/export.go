package models

import "time"

type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJPEG ExportFormat = "jpeg"
	ExportFormatPDF  ExportFormat = "pdf"
)

type ExportSourceType string

const (
	ExportSourceChart ExportSourceType = "chart"
	ExportSourceTable ExportSourceType = "table"
)

type ExportMode string

const (
	ExportModeSync  ExportMode = "sync"
	ExportModeAsync ExportMode = "async"
)

type ExportJobStatus string

const (
	ExportJobStatusQueued     ExportJobStatus = "queued"
	ExportJobStatusProcessing ExportJobStatus = "processing"
	ExportJobStatusCompleted  ExportJobStatus = "completed"
	ExportJobStatusFailed     ExportJobStatus = "failed"
)

type ExportRequest struct {
	Format     ExportFormat     `json:"format" binding:"required,oneof=csv jpeg pdf"`
	SourceType ExportSourceType `json:"source_type" binding:"required,oneof=chart table"`
	Mode       ExportMode       `json:"mode,omitempty" binding:"omitempty,oneof=sync async"`
	FileName   string           `json:"file_name,omitempty"`
	Locale     string           `json:"locale,omitempty"`
	Timezone   string           `json:"timezone,omitempty"`
	Delimiter  string           `json:"delimiter,omitempty"`
	Payload    interface{}      `json:"payload" binding:"required" swaggertype:"object"`
}

type ExportFile struct {
	FileName    string
	ContentType string
	Data        []byte
}

type ExportResult struct {
	Mode   ExportMode      `json:"mode"`
	Status ExportJobStatus `json:"status"`
	JobID  string          `json:"job_id,omitempty"`
	File   *ExportFile     `json:"-"`
}

type ExportJob struct {
	ID          string           `json:"id" bson:"_id,omitempty"`
	Format      ExportFormat     `json:"format" bson:"format"`
	SourceType  ExportSourceType `json:"source_type" bson:"source_type"`
	Status      ExportJobStatus  `json:"status" bson:"status"`
	RequestedAt time.Time        `json:"requested_at" bson:"requested_at"`
	CompletedAt *time.Time       `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Location    string           `json:"location,omitempty" bson:"location,omitempty"`
	Error       string           `json:"error,omitempty" bson:"error,omitempty"`
}
