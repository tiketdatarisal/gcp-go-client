package gcp

import "os"

var (
	defaultProjectID string
	defaultAccessToken string
)

func init() {
	defaultProjectID = os.Getenv("DEFAULT_PROJECT_ID")
	defaultAccessToken = os.Getenv("DEFAULT_ACCESS_TOKEN")
}
