package emailcode

import (
	"context"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/entity/common"
	"github.com/corbado/corbado-go/pkg/servererror"
)

type EmailCode interface {
	Send(ctx context.Context, req api.EmailCodeSendReq, editors ...api.RequestEditorFn) (*api.EmailCodeSendRsp, error)
	Validate(ctx context.Context, emailCodeID common.EmailCodeID, req api.EmailCodeValidateReq, editors ...api.RequestEditorFn) (*api.EmailCodeValidateRsp, error)
	Get(ctx context.Context, emailCodeID common.EmailCodeID, editors ...api.RequestEditorFn) (*api.EmailCodeGetRsp, error)
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ EmailCode = &Impl{}

// New returns new email code client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

// Send sends email code email to given email address
func (i *Impl) Send(ctx context.Context, req api.EmailCodeSendReq, editors ...api.RequestEditorFn) (*api.EmailCodeSendRsp, error) {
	res, err := i.client.EmailCodeSendWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Validate validates email code token
func (i *Impl) Validate(ctx context.Context, emailCodeID common.EmailCodeID, req api.EmailCodeValidateReq, editors ...api.RequestEditorFn) (*api.EmailCodeValidateRsp, error) {
	res, err := i.client.EmailCodeValidateWithResponse(ctx, emailCodeID, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Get gets email code
func (i *Impl) Get(ctx context.Context, emailCodeID common.EmailCodeID, editors ...api.RequestEditorFn) (*api.EmailCodeGetRsp, error) {
	res, err := i.client.EmailCodeGetWithResponse(ctx, emailCodeID, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
