package gcp

import (
	"fmt"
	"strings"
	"testing"
)

func TestStorage_AllBuckets(t *testing.T) {
	ctx := NewContext()
	buckets, err := ctx.Storage.AllBuckets(defaultProjectID, defaultAccessToken)
	if err != nil {
		t.Fatalf("failed to get all buckets with default credentials, details: %v", err)
	}

	if len(buckets) > 0 {
		fmt.Println(strings.Join(buckets, "\n"))
	}
}

func TestStorage_AllObjects(t *testing.T) {
	ctx := NewContext()

	buckets, err := ctx.Storage.AllBuckets(defaultProjectID, defaultAccessToken)
	if err != nil {
		t.Fatalf("failed to get all buckets with default credentials, details: %v", err)
	}

	if len(buckets) == 0 {
		t.Fatalf("there is nothing to test, project ID: %q have no buckets", defaultProjectID)
	}

	objects, err := ctx.Storage.AllObjects(buckets[0], defaultAccessToken)
	if err != nil {
		t.Fatalf("failed to get all buckets with default credentials, details: %v", err)
	}

	if len(objects) > 0 {
		fmt.Println(strings.Join(objects, "\n"))
	}
}

func TestStorage_CreateAndGetObject(t *testing.T) {
	ctx := NewContext()
	data := `1,2,3,4,5,6,7,8,9,0`

	err := ctx.Storage.CreateObject("dev_dms_bucket", "20210522/sample.csv", []byte(data))
	if err != nil {
		t.Fatalf("failed to create object data with default credentials, details: %v", err)
	}

	result, err := ctx.Storage.GetObject("dev_dms_bucket", "20210522/sample.csv")
	if err != nil {
		t.Fatalf("failed to get object data with default credentials, details: %v", err)
	}

	if data != string(result) {
		t.Error("data do not match!")
	}
}
