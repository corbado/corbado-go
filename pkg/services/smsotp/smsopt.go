package smsotp

import (
	"context"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
)

type SmsOTP interface {
	Send(ctx context.Context, req api.SmsCodeSendReq, editors ...api.RequestEditorFn) (*api.SmsCodeSendRsp, error)
	Validate(ctx context.Context, emailCodeID api.SmsCodeID, req api.SmsCodeValidateReq, editors ...api.RequestEditorFn) (*api.SmsCodeValidateRsp, error)
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ SmsOTP = &Impl{}

// New returns new SMS OTP client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

// Send sends OTP SMS to given phone number
func (i *Impl) Send(ctx context.Context, req api.SmsCodeSendReq, editors ...api.RequestEditorFn) (*api.SmsCodeSendRsp, error) {
	res, err := i.client.SmsCodeSendWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Validate validates SMS OTP
func (i *Impl) Validate(ctx context.Context, smsCodeID api.SmsCodeID, req api.SmsCodeValidateReq, editors ...api.RequestEditorFn) (*api.SmsCodeValidateRsp, error) {
	res, err := i.client.SmsCodeValidateWithResponse(ctx, smsCodeID, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
