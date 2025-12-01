package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/esuEdu/sql-benchmark-suite/internal/benchmark"
	"github.com/esuEdu/sql-benchmark-suite/internal/db"
)

func main() {
	ops := 100000
	if v := os.Getenv("OPS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ops = n
		}
	}

	pgURI := getEnv("POSTGRES_URI", "postgres://postgres:postgres@localhost:5432/benchdb?sslmode=disable")
	table := getEnv("POSTGRES_TABLE", "tests")

	pg, err := db.NewPostgres(pgURI, table)
	if err != nil {
		fmt.Println("postgres connect error:", err)
		return
	}

	resultsDir := getEnv("RESULTS_DIR", "results/local")
	res, err := benchmark.RunAndSave(pg, ops, resultsDir)
	if err != nil {
		fmt.Println("benchmark error:", err)
		return
	}
	fmt.Printf("Result: %+v\n", res)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
