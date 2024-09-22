package identifier

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/internal/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/generated/common"
	"github.com/corbado/corbado-go/pkg/servererror"
)

type Identifier interface {
	Create(ctx context.Context, userID string, req api.IdentifierCreateReq, editors ...api.RequestEditorFn) (*api.Identifier, error)
	Delete(ctx context.Context, userID string, identifierID string, editors ...api.RequestEditorFn) (*common.GenericRsp, error)
	List(ctx context.Context, filter []string, sort string, page int, pageSize int, editors ...api.RequestEditorFn) (*api.IdentifierList, error)
	ListByValueAndType(ctx context.Context, identifierValue string, identifierType api.IdentifierType, sort string, page int, pageSize int, editors ...api.RequestEditorFn) (*api.IdentifierList, error)
	ListByUserID(ctx context.Context, userID string, sort string, page int, pageSize int, editors ...api.RequestEditorFn) (*api.IdentifierList, error)
	ListByUserIDAndType(ctx context.Context, userID string, identifierType api.IdentifierType, sort string, page int, pageSize int, editors ...api.RequestEditorFn) (*api.IdentifierList, error)
	UpdateIdentifier(ctx context.Context, userID string, identifierID string, req api.IdentifierUpdateReq, editors ...api.RequestEditorFn) (*api.Identifier, error)
	UpdateStatus(ctx context.Context, userID string, identifierID string, status api.IdentifierStatus, editors ...api.RequestEditorFn) (*api.Identifier, error)
}

type Impl struct {
	client *api.ClientWithResponses
}

var _ Identifier = &Impl{}

// New returns a new Identifier client
func New(client *api.ClientWithResponses) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
	}, nil
}

// Create creates a new identifier
func (i *Impl) Create(
	ctx context.Context,
	userID string,
	req api.IdentifierCreateReq,
	editors ...api.RequestEditorFn,
) (*api.Identifier, error) {
	res, err := i.client.IdentifierCreateWithResponse(ctx, userID, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// Delete deletes an identifier
func (i *Impl) Delete(
	ctx context.Context,
	userID string,
	identifierID string,
	editors ...api.RequestEditorFn,
) (*common.GenericRsp, error) {
	res, err := i.client.IdentifierDeleteWithResponse(ctx, userID, identifierID, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// List lists identifiers based on optional filters, sorting, pagination
func (i *Impl) List(
	ctx context.Context,
	filter []string,
	sort string,
	page int,
	pageSize int,
	editors ...api.RequestEditorFn,
) (*api.IdentifierList, error) {
	var req api.IdentifierListParams

	// Only set filter if it's not nil or empty
	if len(filter) > 0 {
		req.Filter = &filter
	}

	// Only set sort if it's not empty
	if sort != "" {
		req.Sort = &sort
	}

	// Only set page if it's greater than 0
	if page > 0 {
		req.Page = &page
	}

	// Only set pageSize if it's greater than 0
	if pageSize > 0 {
		req.PageSize = &pageSize
	}

	// Call the API with the prepared request
	res, err := i.client.IdentifierListWithResponse(ctx, &req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// ListByValueAndType lists identifiers by value and type
func (i *Impl) ListByValueAndType(
	ctx context.Context,
	value string,
	identifierType api.IdentifierType,
	sort string,
	page int,
	pageSize int,
	editors ...api.RequestEditorFn,
) (*api.IdentifierList, error) {
	filter := []string{`identifierValue:eq:` + value, `identifierType:eq:` + string(identifierType)}

	return i.List(ctx, filter, sort, page, pageSize, editors...)
}

// ListByUserID lists identifiers by user ID
func (i *Impl) ListByUserID(
	ctx context.Context,
	userID string,
	sort string,
	page int,
	pageSize int,
	editors ...api.RequestEditorFn,
) (*api.IdentifierList, error) {
	userID = strings.TrimPrefix(userID, "usr-")

	// Construct the filter
	filter := []string{`userID:eq:` + userID}

	return i.List(ctx, filter, sort, page, pageSize, editors...)
}

// ListByUserIDAndType lists identifiers by user ID and type
func (i *Impl) ListByUserIDAndType(
	ctx context.Context,
	userID string,
	identifierType api.IdentifierType,
	sort string,
	page int,
	pageSize int,
	editors ...api.RequestEditorFn,
) (*api.IdentifierList, error) {
	userID = strings.TrimPrefix(userID, "usr-")

	// Construct the filter
	filter := []string{`userID:eq:` + userID, `identifierType:eq:` + string(identifierType)}
	return i.List(ctx, filter, sort, page, pageSize, editors...)
}

// UpdateIdentifier updates an identifier
func (i *Impl) UpdateIdentifier(
	ctx context.Context,
	userID string,
	identifierID string,
	req api.IdentifierUpdateReq,
	editors ...api.RequestEditorFn,
) (*api.Identifier, error) {
	res, err := i.client.IdentifierUpdateWithResponse(ctx, userID, identifierID, req, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// UpdateStatus updates the status of an identifier
func (i *Impl) UpdateStatus(
	ctx context.Context,
	userID string,
	identifierID string,
	status api.IdentifierStatus,
	editors ...api.RequestEditorFn,
) (*api.Identifier, error) {
	req := api.IdentifierUpdateReq{
		Status: status,
	}

	return i.UpdateIdentifier(ctx, userID, identifierID, req, editors...)
}
