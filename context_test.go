package gcp

import "os"

var (
	defaultProjectID string
	defaultDatasetID string
	defaultTableID string
	defaultAccessToken string
)

func init() {
	defaultProjectID = os.Getenv("DEFAULT_PROJECT_ID")
	defaultDatasetID = os.Getenv("DEFAULT_DATASET_ID")
	defaultTableID = os.Getenv("DEFAULT_TABLE_ID")
	defaultAccessToken = os.Getenv("DEFAULT_ACCESS_TOKEN")
}
