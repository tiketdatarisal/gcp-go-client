package gcp

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/tiketdatarisal/gcp-go-client/models"
	"golang.org/x/oauth2"
	bq "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type BigQuery struct {
	context *Context
}

func (b *BigQuery) newBigQueryManager(ctx context.Context, projectID, token string) (*bigquery.Client, error) {
	var result *bigquery.Client
	var err error
	if token != "" {
		t := &oauth2.Token{AccessToken: token, TokenType: "Bearer"}
		result, err = bigquery.NewClient(ctx, projectID, option.WithTokenSource(oauth2.StaticTokenSource(t)))
	} else {
		result, err = bigquery.NewClient(ctx, projectID)
	}

	if err != nil {
		status, err := parseGCPError(err)
		return nil, &HttpError{Code: status, Err: err}
	}

	return result, nil
}

func (b *BigQuery) AllProjectIDs(token ...string) ([]string, error) {
	ctx := context.Background()
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	var svc *bq.Service
	var err error
	if t != "" {
		t := &oauth2.Token{AccessToken: t, TokenType: "Bearer"}
		svc, err = bq.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(t)))
	} else {
		svc, err = bq.NewService(ctx)
	}

	if err != nil {
		status, err := parseGCPError(err)
		return nil, &HttpError{Code: status, Err: err}
	}

	var result []string
	var pt = ""
	for {
		res, err := svc.Projects.List().PageToken(pt).Do()
		if err != nil {
			st, err := parseGCPError(err)
			return nil, &HttpError{Code: st, Err: err}
		}

		for _, p := range res.Projects {
			result = append(result, p.Id)
		}

		pt = res.NextPageToken
		if pt == "" {
			break
		}
	}

	return result, nil
}

func (b *BigQuery) AllDatasetIDs(projectID string, token ...string) ([]string, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	c, err := b.newBigQueryManager(ctx, projectID, t)
	if err != nil {
		return nil, err
	} else {
		defer func() { _ = c.Close() }()
	}

	var result []string
	i := c.Datasets(ctx)
	for {
		d, err := i.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			st, err := parseGCPError(err)
			return nil, &HttpError{Code: st, Err: err}
		}

		result = append(result, d.DatasetID)
	}

	return result, nil
}

func (b *BigQuery) AllTableIDs(projectID, datasetID string, token ...string) ([]string, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	c, err := b.newBigQueryManager(ctx, projectID, t)
	if err != nil {
		return nil, err
	} else {
		defer func() { _ = c.Close() }()
	}

	var result []string
	i := c.Dataset(datasetID).Tables(ctx)
	for {
		t, err := i.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			st, err := parseGCPError(err)
			return nil, &HttpError{Code: st, Err: err}
		}

		result = append(result, t.TableID)
	}

	return result, nil
}

func (b *BigQuery) AllColumns(projectID, datasetID, tableID string, token ...string) (*models.Columns, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	c, err := b.newBigQueryManager(ctx, projectID, t)
	if err != nil {
		return nil, err
	} else {
		defer func() { _ = c.Close() }()
	}

	m, err := c.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		st, err := parseGCPError(err)
		return nil, &HttpError{Code: st, Err: err}
	}

	result := models.Columns{}
	for _, i := range m.Schema {
		result = append(result, models.Column{ColumnName: i.Name, DataType: fmt.Sprintf("%s", i.Type)})
	}

	return &result, nil
}
