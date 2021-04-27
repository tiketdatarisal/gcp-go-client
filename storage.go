package gcp

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"time"
)

type Storage struct {
	context *Context
}

// newStorageManager return a reference to storage manager client.
func (s *Storage) newStorageManager(ctx context.Context, token string) (*storage.Client, error) {
	var result *storage.Client
	var err error
	if token != "" {
		t := &oauth2.Token{AccessToken: token, TokenType: "Bearer"}
		result, err = storage.NewClient(ctx, option.WithTokenSource(oauth2.StaticTokenSource(t)))
	} else {
		result, err = storage.NewClient(ctx)
	}

	if err != nil {
		status, err := parseGCPError(err)
		return nil, &HttpError{Code: status, Err: err}
	}

	return result, nil
}

// AllBuckets returns list of buckets.
func (s *Storage) AllBuckets(projectID string, token ...string) ([]string, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	c, err := s.newStorageManager(ctx, t)
	if err != nil {
		return nil, err
	} else {
		defer c.Close()
	}

	var result []string
	it := c.Buckets(ctx, projectID)
	for {
		bs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		result = append(result, bs.Name)
	}

	return result, err
}

// CreateBucket create a new bucket.
func (s *Storage) CreateBucket(projectID, bucketName string, attrs *storage.BucketAttrs, token ...string) error  {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	c, err := s.newStorageManager(ctx, t)
	if err != nil {
		return err
	} else {
		defer c.Close()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := c.Bucket(bucketName)
	if err = bucket.Create(ctx, projectID, attrs); err != nil {
		return err
	}

	return nil
}

// AllObjects returns all objects.
func (s *Storage) AllObjects(bucketName string, token ...string) ([]string, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	c, err := s.newStorageManager(ctx, t)
	if err != nil {
		return nil, err
	} else {
		defer c.Close()
	}

	var result []string
	it := c.Bucket(bucketName).Objects(ctx, nil)
	for {
		os, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		result = append(result, os.Name)
	}

	return result, err
}

// CreateObject create a new object.
func (s *Storage) CreateObject(bucketName, objectName string, data []byte, token ...string) error {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	c, err := s.newStorageManager(ctx, t)
	if err != nil {
		return err
	} else {
		defer c.Close()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	reader := bytes.NewReader(data)
	writer := c.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}

	_ = writer.Close()
	return nil
}

// GetObject get (download) an object.
func (s *Storage) GetObject(bucketName, objectName string, token ...string) ([]byte, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	c, err := s.newStorageManager(ctx, t)
	if err != nil {
		return nil, err
	} else {
		defer c.Close()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	reader, err := c.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close()}()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
