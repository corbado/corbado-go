package project

import (
	"context"

	"github.com/corbado/corbado-go/pkg/sdk/assert"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/servererror"
)

type Project interface {
	CreateSecret(ctx context.Context, req api.ProjectSecretCreateReq, editors ...api.RequestEditorFn) (*api.ProjectSecretCreateRsp, error)
	ConfigGet(ctx context.Context, editors ...api.RequestEditorFn) (*api.ProjectConfigGetRsp, error)
	ConfigUpdate(ctx context.Context, req api.ProjectConfigSaveReq, editors ...api.RequestEditorFn) error
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ Project = &Impl{}

// New returns new projects client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

// CreateSecret creates an API secret
func (i *Impl) CreateSecret(ctx context.Context, req api.ProjectSecretCreateReq, editors ...api.RequestEditorFn) (*api.ProjectSecretCreateRsp, error) {
	res, err := i.client.ProjectSecretCreateWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// ConfigGet retrieves project config by projectID inferred from authentication
func (i *Impl) ConfigGet(ctx context.Context, editors ...api.RequestEditorFn) (*api.ProjectConfigGetRsp, error) {
	res, err := i.client.ProjectConfigGetWithResponse(ctx, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// ConfigUpdate saves project config
func (i *Impl) ConfigUpdate(ctx context.Context, req api.ProjectConfigSaveReq, editors ...api.RequestEditorFn) error {
	res, err := i.client.ProjectConfigSaveWithResponse(ctx, req, editors...)
	if err != nil {
		return err
	}

	if res.JSONDefault != nil {
		return servererror.New(res.JSONDefault)
	}

	return nil
}
