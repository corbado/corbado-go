package passkey

import (
	"context"

	"github.com/corbado/corbado-go/pkg/sdk/assert"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/servererror"
)

type Passkey interface {
	RegisterStart(ctx context.Context, req api.WebAuthnRegisterStartReq, editors ...api.RequestEditorFn) (*api.WebAuthnRegisterStartRsp, error)
	RegisterFinish(ctx context.Context, req api.WebAuthnFinishReq, editors ...api.RequestEditorFn) (*api.WebAuthnRegisterFinishRsp, error)
	AuthenticateStart(ctx context.Context, req api.WebAuthnAuthenticateStartReq, editors ...api.RequestEditorFn) (*api.WebAuthnAuthenticateStartRsp, error)
	AuthenticateFinish(ctx context.Context, req api.WebAuthnFinishReq, editors ...api.RequestEditorFn) (*api.WebAuthnAuthenticateFinishRsp, error)
	MediationStart(ctx context.Context, req api.WebAuthnMediationStartReq, editors ...api.RequestEditorFn) (*api.WebAuthnMediationStartRsp, error)
	AssociateStart(ctx context.Context, req api.WebAuthnAssociateStartReq, editors ...api.RequestEditorFn) (*api.WebAuthnAssociateStartRsp, error)
	CredentialList(ctx context.Context, params *api.WebAuthnCredentialListParams, editors ...api.RequestEditorFn) (*api.WebAuthnCredentialListRsp, error)
	CredentialUpdate(ctx context.Context, credentialID string, req api.WebAuthnCredentialReq, editors ...api.RequestEditorFn) (*api.WebAuthnCredentialRsp, error)
	CredentialExists(ctx context.Context, req api.WebAuthnCredentialExistsReq, editors ...api.RequestEditorFn) (*api.WebAuthnCredentialExistsRsp, error)
	CredentialDelete(ctx context.Context, userID string, credentialID string, req api.EmptyReq, editors ...api.RequestEditorFn) error
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ Passkey = &Impl{}

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

// RegisterStart starts registration of a user for Passkeys (Biometrics)
func (i *Impl) RegisterStart(ctx context.Context, req api.WebAuthnRegisterStartReq, editors ...api.RequestEditorFn) (*api.WebAuthnRegisterStartRsp, error) {
	res, err := i.client.WebAuthnRegisterStartWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// RegisterFinish completes registration of a user for Passkeys (Biometrics)
func (i *Impl) RegisterFinish(ctx context.Context, req api.WebAuthnFinishReq, editors ...api.RequestEditorFn) (*api.WebAuthnRegisterFinishRsp, error) {
	res, err := i.client.WebAuthnRegisterFinishWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// AuthenticateStart starts authentication of a user for Passkeys (Biometrics)
func (i *Impl) AuthenticateStart(ctx context.Context, req api.WebAuthnAuthenticateStartReq, editors ...api.RequestEditorFn) (*api.WebAuthnAuthenticateStartRsp, error) {
	res, err := i.client.WebAuthnAuthenticateStartWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// AuthenticateFinish completes authentication of a user for Passkeys (Biometrics)
func (i *Impl) AuthenticateFinish(ctx context.Context, req api.WebAuthnFinishReq, editors ...api.RequestEditorFn) (*api.WebAuthnAuthenticateFinishRsp, error) {
	res, err := i.client.WebAuthnAuthenticateFinishWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// MediationStart starts mediation for Passkeys (Biometrics)
func (i *Impl) MediationStart(ctx context.Context, req api.WebAuthnMediationStartReq, editors ...api.RequestEditorFn) (*api.WebAuthnMediationStartRsp, error) {
	res, err := i.client.WebAuthnMediationStartWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// AssociateStart starts association token flow for Passkeys (Biometrics)
func (i *Impl) AssociateStart(ctx context.Context, req api.WebAuthnAssociateStartReq, editors ...api.RequestEditorFn) (*api.WebAuthnAssociateStartRsp, error) {
	res, err := i.client.WebAuthnAssociateStartWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// CredentialList lists webauthn credentials users
func (i *Impl) CredentialList(ctx context.Context, params *api.WebAuthnCredentialListParams, editors ...api.RequestEditorFn) (*api.WebAuthnCredentialListRsp, error) {
	res, err := i.client.WebAuthnCredentialListWithResponse(ctx, params, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// CredentialUpdate updates webauthn credential
func (i *Impl) CredentialUpdate(ctx context.Context, credentialID string, req api.WebAuthnCredentialReq, editors ...api.RequestEditorFn) (*api.WebAuthnCredentialRsp, error) {
	res, err := i.client.WebAuthnCredentialUpdateWithResponse(ctx, credentialID, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// CredentialExists checks if active webauthn credential exists for provided user and device
func (i *Impl) CredentialExists(ctx context.Context, req api.WebAuthnCredentialExistsReq, editors ...api.RequestEditorFn) (*api.WebAuthnCredentialExistsRsp, error) {
	res, err := i.client.WebAuthnCredentialExistsWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// CredentialDelete deletes webauthn credential
func (i *Impl) CredentialDelete(ctx context.Context, userID string, credentialID string, req api.EmptyReq, editors ...api.RequestEditorFn) error {
	res, err := i.client.WebAuthnCredentialDeleteWithResponse(ctx, userID, credentialID, req, editors...)
	if err != nil {
		return err
	}

	if res.JSONDefault != nil {
		return servererror.New(res.JSONDefault)
	}

	return nil
}
