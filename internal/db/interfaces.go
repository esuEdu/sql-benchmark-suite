package db

import "time"

type Result struct {
	DB          string  `json:"db"`
	Ops         int     `json:"ops"`
	WritesMs    int64   `json:"writes_duration_ms"`
	ReadsMs     int64   `json:"reads_duration_ms"`
	AvgWriteMs  float64 `json:"average_write_ms"`
	AvgReadMs   float64 `json:"average_read_ms"`
	GeneratedAt string  `json:"generated_at"`
	Note        string  `json:"note,omitempty"`
}

type BenchmarkDB interface {
	WriteTest(n int) (time.Duration, error)
	ReadTest(n int) (time.Duration, error)
	Name() string
}
