package gcp

import (
	"fmt"
	"strings"
	"testing"
)

func TestBigQuery_AllProjectIDs(t *testing.T) {
	ctx := NewContext()
	projects, err := ctx.BigQueries.AllProjectIDs(defaultAccessToken)
	if err != nil {
		t.Fatalf("failed to get all project IDs with default credentials, details: %v", err)
	}

	if len(projects) > 0 {
		fmt.Println(strings.Join(projects, "\n"))
	}
}

func TestBigQuery_AllDatasetIDs(t *testing.T) {
	ctx := NewContext()
	datasets, err := ctx.BigQueries.AllDatasetIDs(defaultProjectID, defaultAccessToken)
	if err != nil {
		t.Fatalf("failed to get all dataset IDs with default credentials, details: %v", err)
	}

	if len(datasets) > 0 {
		fmt.Println(strings.Join(datasets, "\n"))
	}
}

func TestBigQuery_AllTableIDs(t *testing.T) {
	ctx := NewContext()
	tables, err := ctx.BigQueries.AllTableIDs(defaultProjectID, defaultDatasetID, defaultAccessToken)
	if err != nil {
		t.Fatalf("failed to get all table IDs with default credentials, details: %v", err)
	}

	if len(tables) > 0 {
		fmt.Println(strings.Join(tables, "\n"))
	}
}

func TestBigQuery_AllColumns(t *testing.T) {
	ctx := NewContext()
	columns, err := ctx.BigQueries.AllColumns(defaultProjectID, defaultDatasetID, defaultTableID, defaultAccessToken)
	if err != nil {
		t.Fatalf("failed to get all columns with default credentials, details: %v", err)
	}

	if len(*columns) > 0 {
		names := []string{}
		for _, c := range *columns {
			names = append(names, c.ColumnName)
		}
		fmt.Println(strings.Join(names, "\n"))
	}
}
