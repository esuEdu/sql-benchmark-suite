package benchmark

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/esuEdu/sql-benchmark-suite/internal/db"
)

func RunAndSave(b db.BenchmarkDB, ops int, resultsDir string) (*db.Result, error) {
	wDur, err := b.WriteTest(ops)
	if err != nil {
		return nil, err
	}
	rDur, err := b.ReadTest(ops)
	if err != nil {
		return nil, err
	}

	res := &db.Result{
		DB:          b.Name(),
		Ops:         ops,
		WritesMs:    wDur.Milliseconds(),
		ReadsMs:     rDur.Milliseconds(),
		AvgWriteMs:  float64(wDur.Milliseconds()) / float64(ops),
		AvgReadMs:   float64(rDur.Milliseconds()) / float64(ops),
		GeneratedAt: time.Now().Format(time.RFC3339),
	}

	// ensure results dir
	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		return res, err
	}

	filename := fmt.Sprintf("%s_%d.json", b.Name(), ops)
	fullpath := filepath.Join(resultsDir, filename)
	f, err := os.Create(fullpath)
	if err != nil {
		return res, err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		return res, err
	}

	return res, nil
}
