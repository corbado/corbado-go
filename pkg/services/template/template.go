package template

import (
	"context"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
)

type Template interface {
	CreateEmailTemplate(ctx context.Context, req api.EmailTemplateCreateReq, editors ...api.RequestEditorFn) (*api.EmailTemplateCreateRsp, error)
	CreateSMSTemplate(ctx context.Context, req api.SmsTemplateCreateReq, editors ...api.RequestEditorFn) (*api.SmsTemplateCreateRsp, error)
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ Template = &Impl{}

// New returns new templates client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

// CreateEmailTemplate creates a new email template
func (i *Impl) CreateEmailTemplate(ctx context.Context, req api.EmailTemplateCreateReq, editors ...api.RequestEditorFn) (*api.EmailTemplateCreateRsp, error) {
	res, err := i.client.EmailTemplateCreateWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// CreateSMSTemplate creates a new SMS template
func (i *Impl) CreateSMSTemplate(ctx context.Context, req api.SmsTemplateCreateReq, editors ...api.RequestEditorFn) (*api.SmsTemplateCreateRsp, error) {
	res, err := i.client.SmsTemplateCreateWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
