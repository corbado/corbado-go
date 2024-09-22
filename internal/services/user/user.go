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
	Create(ctx context.Context, req api.UserCreateReq, editors ...api.RequestEditorFn) (*api.User, error)
	CreateActiveByName(ctx context.Context, fullName string, editors ...api.RequestEditorFn) (*api.User, error)
	Get(ctx context.Context, userID common.UserID, editors ...api.RequestEditorFn) (*api.User, error)
	Delete(ctx context.Context, userID common.UserID, editors ...api.RequestEditorFn) (*common.GenericRsp, error)
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

// Create creates a new user
func (i *Impl) Create(ctx context.Context, req api.UserCreateReq, editors ...api.RequestEditorFn) (*api.User, error) {
	res, err := i.client.UserCreateWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Create creates a new user
func (i *Impl) CreateActiveByName(ctx context.Context, fullName string, editors ...api.RequestEditorFn) (*api.User, error) {
	req := api.UserCreateReq{
		FullName: &fullName,
		Status:   "active",
	}

	res, err := i.client.UserCreateWithResponse(ctx, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Get gets a user by ID
func (i *Impl) Get(ctx context.Context, userID common.UserID, editors ...api.RequestEditorFn) (*api.User, error) {
	res, err := i.client.UserGetWithResponse(ctx, userID, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Delete deletes a user by ID
func (i *Impl) Delete(ctx context.Context, userID common.UserID, editors ...api.RequestEditorFn) (*common.GenericRsp, error) {
	res, err := i.client.UserDeleteWithResponse(ctx, userID, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
