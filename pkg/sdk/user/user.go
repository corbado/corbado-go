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

// Update updates user
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
