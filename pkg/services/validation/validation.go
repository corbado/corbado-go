package validation

import (
	"context"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
)

type Validation interface {
	ValidateEmail(ctx context.Context, req api.ValidateEmailReq, editors ...api.RequestEditorFn) (*api.ValidateEmailRsp, error)
	ValidatePhoneNumber(ctx context.Context, req api.ValidatePhoneNumberReq, editors ...api.RequestEditorFn) (*api.ValidatePhoneNumberRsp, error)
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ Validation = &Impl{}

// New returns new email link client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

// ValidateEmail validates provided email address with optional additional checks
func (i *Impl) ValidateEmail(ctx context.Context, req api.ValidateEmailReq, editors ...api.RequestEditorFn) (*api.ValidateEmailRsp, error) {
	res, err := i.client.ValidateEmailWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// ValidatePhoneNumber validates provided phone number
func (i *Impl) ValidatePhoneNumber(ctx context.Context, req api.ValidatePhoneNumberReq, editors ...api.RequestEditorFn) (*api.ValidatePhoneNumberRsp, error) {
	res, err := i.client.ValidatePhoneNumberWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
