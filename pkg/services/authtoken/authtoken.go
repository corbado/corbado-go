package authtoken

import (
	"context"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
)

type AuthToken interface {
	Validate(ctx context.Context, req api.AuthTokenValidateReq, editors ...api.RequestEditorFn) (*api.AuthTokenValidateRsp, error)
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ AuthToken = &Impl{}

// New returns new auth tokens client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

func (i *Impl) Validate(ctx context.Context, req api.AuthTokenValidateReq, editors ...api.RequestEditorFn) (*api.AuthTokenValidateRsp, error) {
	res, err := i.client.AuthTokenValidateWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
