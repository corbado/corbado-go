package emaillink

import (
	"context"

	"github.com/corbado/corbado-go/pkg/sdk/assert"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/entity/common"
	"github.com/corbado/corbado-go/pkg/sdk/servererror"
)

type EmailLink interface {
	Send(ctx context.Context, req api.EmailLinkSendReq, editors ...api.RequestEditorFn) (*api.EmailLinkSendRsp, error)
	Validate(ctx context.Context, emailLinkID common.EmailLinkID, req api.EmailLinksValidateReq, editors ...api.RequestEditorFn) (*api.EmailLinkValidateRsp, error)
	Get(ctx context.Context, emailLinkID common.EmailLinkID, editors ...api.RequestEditorFn) (*api.EmailLinkGetRsp, error)
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ EmailLink = &Impl{}

// New returns new email link client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

// Send sends email link email to given email address
func (i *Impl) Send(ctx context.Context, req api.EmailLinkSendReq, editors ...api.RequestEditorFn) (*api.EmailLinkSendRsp, error) {
	res, err := i.client.EmailLinkSendWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Validate validates email link token
func (i *Impl) Validate(ctx context.Context, emailLinkID common.EmailLinkID, req api.EmailLinksValidateReq, editors ...api.RequestEditorFn) (*api.EmailLinkValidateRsp, error) {
	res, err := i.client.EmailLinkValidateWithResponse(ctx, emailLinkID, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Get gets email link
func (i *Impl) Get(ctx context.Context, emailLinkID common.EmailLinkID, editors ...api.RequestEditorFn) (*api.EmailLinkGetRsp, error) {
	res, err := i.client.EmailLinkGetWithResponse(ctx, emailLinkID, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
