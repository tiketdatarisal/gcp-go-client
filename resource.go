package gcp

import (
	"context"
	"github.com/tiketdatarisal/gcp-go-client/models"
	"golang.org/x/oauth2"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

type Resource struct {
	context *Context
}

func (r *Resource) newResourceManager(ctx context.Context, token string) (*cloudresourcemanager.Service, error) {
	var result *cloudresourcemanager.Service
	var err error
	if token != "" {
		t := &oauth2.Token{AccessToken: token, TokenType: "Bearer"}
		result, err = cloudresourcemanager.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(t)))
	} else {
		result, err = cloudresourcemanager.NewService(ctx)
	}

	if err != nil {
		status, err := parseGCPError(err)
		return nil, &HttpError{Code: status, Err: err}
	}

	return result, nil
}

func (r *Resource) getProjectPolicy(svc *cloudresourcemanager.Service, projectID string) (*cloudresourcemanager.Policy, error) {
	req := cloudresourcemanager.GetIamPolicyRequest{
		Options: &cloudresourcemanager.GetPolicyOptions{RequestedPolicyVersion: 3},
	}

	return svc.Projects.GetIamPolicy(projectID, &req).Do()
}

func (r *Resource) AllProjectIDs(token ...string) ([]string, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	s, err := r.newResourceManager(ctx, t)
	if err != nil {
		return nil, err
	}

	var result []string
	var pt = ""
	for {
		res, err := s.Projects.List().PageToken(pt).Do()
		if err != nil {
			st, err := parseGCPError(err)
			return nil, &HttpError{Code: st, Err: err}
		}

		for _, p := range res.Projects {
			if p != nil {
				result = append(result, p.ProjectId)
			}
		}

		pt = res.NextPageToken
		if pt == "" {
			break
		}
	}

	return result, nil
}

func (r *Resource) AllRoleUsers(projectID string, token ...string) (models.RoleUsers, error) {
	t := ""
	if token != nil && len(token) > 0 {
		t = token[0]
	}

	ctx := context.Background()
	s, err := r.newResourceManager(ctx, t)
	if err != nil {
		return nil, err
	}

	p, err := r.getProjectPolicy(s, projectID)
	if err != nil {
		return nil, err
	}

	if p.Bindings == nil {
		return nil, nil
	}

	var result models.RoleUsers
	for _, b := range p.Bindings {
		if b != nil {
			ru := models.RoleUser{Role: b.Role}
			for _, m := range b.Members {
				ru.Users = append(ru.Users, models.User(m))
			}
			result = append(result, ru)
		}
	}

	return result, nil
}
