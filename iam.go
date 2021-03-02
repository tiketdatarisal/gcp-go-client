package gcp

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

type IAM struct {
	context *Context
}

func (i *IAM) newIamManager(ctx context.Context, token string) (*iam.Service, error) {
	var result *iam.Service
	var err error
	if token != "" {
		t := &oauth2.Token{AccessToken: token, TokenType: "Bearer"}
		result, err = iam.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(t)))
	} else {
		result, err = iam.NewService(ctx)
	}

	if err != nil {
		status, err := parseGCPError(err)
		return nil, &HttpError{Code: status, Err: err}
	}

	return result, nil
}

func (i *IAM) getGrantableRoles(svc *iam.Service, r string) ([]string, error) {
	var tk = ""
	var result []string
	for {
		req := iam.QueryGrantableRolesRequest{
			FullResourceName: r,
			PageToken:        tk,
			View:             "BASIC",
		}

		res, err := svc.Roles.QueryGrantableRoles(&req).Do()
		if err != nil {
			status, err := parseGCPError(err)
			return nil, &HttpError{Code: status, Err: err}
		}

		for _, ro := range res.Roles {
			if ro != nil {
				result = append(result, ro.Name)
			}
		}

		tk := res.NextPageToken
		if tk == "" {
			break
		}
	}

	return result, nil
}

func (i *IAM) AllProjectRoles(projectID string, token ...string) ([]string, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	s, err := i.newIamManager(ctx, t)
	if err != nil {
		return nil, err
	}

	return i.getGrantableRoles(s, fmt.Sprintf("//cloudresourcemanager.googleapis.com/projects/%s", projectID))
}

func (i *IAM) AllDatasetRoles(projectID, datasetID string, token ...string) ([]string, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	s, err := i.newIamManager(ctx, t)
	if err != nil {
		return nil, err
	}

	return i.getGrantableRoles(s, fmt.Sprintf("//bigquery.googleapis.com/projects/%s/datasets/%s", projectID, datasetID))
}

func (i *IAM) AllBucketRoles(bucketID string, token ...string) ([]string, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	s, err := i.newIamManager(ctx, t)
	if err != nil {
		return nil, err
	}

	return i.getGrantableRoles(s, fmt.Sprintf("//storage.googleapis.com/projects/_/buckets/%s", bucketID))
}
