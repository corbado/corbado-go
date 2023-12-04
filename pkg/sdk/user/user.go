package user

import (
	"context"

	"github.com/corbado/corbado-go/pkg/sdk/assert"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/entity/common"
	"github.com/corbado/corbado-go/pkg/sdk/servererror"
)

type User interface {
	List(ctx context.Context, params *api.UserListParams, editors ...api.RequestEditorFn) (*api.UserListRsp, error)
	Update(ctx context.Context, userID common.UserID, req api.UserUpdateReq, editors ...api.RequestEditorFn) (*api.UserUpdateRsp, error)
	Create(ctx context.Context, req api.UserCreateReq, editors ...api.RequestEditorFn) (*api.UserCreateRsp, error)
	Get(ctx context.Context, userID common.UserID, params *api.UserGetParams, editors ...api.RequestEditorFn) (*api.UserGetRsp, error)
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

// List lists project users
func (i *Impl) List(ctx context.Context, params *api.UserListParams, editors ...api.RequestEditorFn) (*api.UserListRsp, error) {
	res, err := i.client.UserListWithResponse(ctx, params, editors...)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
