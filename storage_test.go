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
