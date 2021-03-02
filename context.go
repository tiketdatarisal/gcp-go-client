package gcp

import (
	"google.golang.org/api/googleapi"
	"net/http"
)

func parseGCPError(err error) (int, error) {
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			return e.Code, e
		} else {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

type Context struct {
	Resources  *Resource
	BigQueries *BigQuery
	IAM        *IAM
	Storage    *Storage
}

func NewContext() *Context {
	ctx := &Context{}
	ctx.Resources = &Resource{context: ctx}
	ctx.BigQueries = &BigQuery{context: ctx}
	ctx.IAM = &IAM{context: ctx}
	ctx.Storage = &Storage{context: ctx}
	return ctx
}
