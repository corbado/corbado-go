package user

import (
	"context"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/internal/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/generated/common"
	"github.com/corbado/corbado-go/pkg/servererror"
)

type User interface {
	List(ctx context.Context, params *api.UserListParams, editors ...api.RequestEditorFn) (*api.UserListRsp, error)
	Update(ctx context.Context, userID common.UserID, req api.UserUpdateReq, editors ...api.RequestEditorFn) (*api.UserUpdateRsp, error)
	Create(ctx context.Context, req api.UserCreateReq, editors ...api.RequestEditorFn) (*api.UserCreateRsp, error)
	Get(ctx context.Context, userID common.UserID, params *api.UserGetParams, editors ...api.RequestEditorFn) (*api.UserGetRsp, error)
	Delete(ctx context.Context, userID common.UserID, req api.UserDeleteReq, editors ...api.RequestEditorFn) (*common.GenericRsp, error)
	Exists(ctx context.Context, req api.UserExistsReq, editors ...api.RequestEditorFn) (*api.UserExistsRsp, error)
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ User = &Impl{}

// New returns new user client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

// List lists users
func (i *Impl) List(ctx context.Context, params *api.UserListParams, editors ...api.RequestEditorFn) (*api.UserListRsp, error) {
	res, err := i.client.UserListWithResponse(ctx, params, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Create creates a new user
func (i *Impl) Create(ctx context.Context, req api.UserCreateReq, editors ...api.RequestEditorFn) (*api.UserCreateRsp, error) {
	res, err := i.client.UserCreateWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Update updates a user
func (i *Impl) Update(ctx context.Context, userID common.UserID, req api.UserUpdateReq, editors ...api.RequestEditorFn) (*api.UserUpdateRsp, error) {
	res, err := i.client.UserUpdateWithResponse(ctx, userID, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Get gets a user by ID
func (i *Impl) Get(ctx context.Context, userID common.UserID, params *api.UserGetParams, editors ...api.RequestEditorFn) (*api.UserGetRsp, error) {
	res, err := i.client.UserGetWithResponse(ctx, userID, params, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Delete deletes a user by ID
func (i *Impl) Delete(ctx context.Context, userID common.UserID, req api.UserDeleteReq, editors ...api.RequestEditorFn) (*common.GenericRsp, error) {
	res, err := i.client.UserDeleteWithResponse(ctx, userID, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Exists checks if a confirmed user exists for provided login identifier
func (i *Impl) Exists(ctx context.Context, req api.UserExistsReq, editors ...api.RequestEditorFn) (*api.UserExistsRsp, error) {
	res, err := i.client.UserExistsWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
