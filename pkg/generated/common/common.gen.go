// Package common provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package common

const (
	ProjectIDScopes = "projectID.Scopes"
)

// Defines values for AppType.
const (
	Empty  AppType = "empty"
	Native AppType = "native"
	Web    AppType = "web"
)

// Defines values for AuthMethod.
const (
	AuthMethodEmail       AuthMethod = "email"
	AuthMethodPassword    AuthMethod = "password"
	AuthMethodPhoneNumber AuthMethod = "phone_number"
	AuthMethodWebauthn    AuthMethod = "webauthn"
)

// Defines values for LoginIdentifierConfigEnforceVerification.
const (
	AtFirstLogin LoginIdentifierConfigEnforceVerification = "at_first_login"
	None         LoginIdentifierConfigEnforceVerification = "none"
	Signup       LoginIdentifierConfigEnforceVerification = "signup"
)

// Defines values for LoginIdentifierType.
const (
	LoginIdentifierTypeCustom      LoginIdentifierType = "custom"
	LoginIdentifierTypeEmail       LoginIdentifierType = "email"
	LoginIdentifierTypePhoneNumber LoginIdentifierType = "phone_number"
)

// Defines values for SessionManagement.
const (
	SessionManagementCorbado SessionManagement = "SessionManagementCorbado"
	SessionManagementOwn     SessionManagement = "SessionManagementOwn"
)

// Defines values for SocialProviderType.
const (
	Github    SocialProviderType = "github"
	Google    SocialProviderType = "google"
	Microsoft SocialProviderType = "microsoft"
)

// Defines values for Status.
const (
	StatusActive  Status = "active"
	StatusDeleted Status = "deleted"
	StatusPending Status = "pending"
)

// ID generic ID
type ID = string

// AdditionalPayload Additional payload in JSON format
type AdditionalPayload = string

// AllTypes defines model for allTypes.
type AllTypes struct {
	P1 *Paging `json:"p1,omitempty"`

	// P10 Timestamp of when the entity was deleted in yyyy-MM-dd'T'HH:mm:ss format
	P10 *Deleted `json:"p10,omitempty"`

	// P11 ID of the device
	P11 *DeviceID `json:"p11,omitempty"`

	// P12 Additional payload in JSON format
	P12 *AdditionalPayload `json:"p12,omitempty"`

	// P13 Generic status that can describe Corbado entities
	P13 *Status `json:"p13,omitempty"`

	// P14 ID of project
	P14 *ProjectID `json:"p14,omitempty"`

	// P15 Unique ID of request, you can provide your own while making the request, if not the ID will be randomly generated on server side
	P15 *RequestID   `json:"p15,omitempty"`
	P16 *ErrorRsp    `json:"p16,omitempty"`
	P17 *AuthMethods `json:"p17,omitempty"`

	// P18 User entry with emails and phone numbers
	P18 *FullUser `json:"p18,omitempty"`

	// P19 Login Identifier type
	P19 *LoginIdentifierType `json:"p19,omitempty"`
	P2  *ClientInfo          `json:"p2,omitempty"`

	// P20 ID of the email OTP
	P20 *EmailCodeID `json:"p20,omitempty"`

	// P21 Application type
	P21 *AppType `json:"p21,omitempty"`

	// P22 What session management should be used
	P22 *SessionManagement `json:"p22,omitempty"`

	// P23 High entropy values from browser
	P23 *HighEntropyValues     `json:"p23,omitempty"`
	P24 *LoginIdentifierConfig `json:"p24,omitempty"`
	P25 *SocialProviderType    `json:"p25,omitempty"`

	// P3 generic ID
	P3 *ID `json:"p3,omitempty"`

	// P4 ID of the user
	P4 *UserID `json:"p4,omitempty"`

	// P5 ID of the email
	P5 *EmailID `json:"p5,omitempty"`

	// P6 ID of the email magic link
	P6 *EmailLinkID `json:"p6,omitempty"`

	// P7 ID of the phone number
	P7 *PhoneNumberID `json:"p7,omitempty"`

	// P8 Timestamp of when the entity was created in yyyy-MM-dd'T'HH:mm:ss format
	P8 *Created `json:"p8,omitempty"`

	// P9 Timestamp of when the entity was last updated in yyyy-MM-dd'T'HH:mm:ss format
	P9 *Updated `json:"p9,omitempty"`
}

// AppType Application type
type AppType string

// AuthMethod Authentication methods
type AuthMethod string

// AuthMethods defines model for authMethods.
type AuthMethods = []AuthMethod

// ClientInfo defines model for clientInfo.
type ClientInfo struct {
	// RemoteAddress client's IP address
	RemoteAddress string `json:"remoteAddress"`

	// UserAgent client's User Agent
	UserAgent string `json:"userAgent"`
}

// Created Timestamp of when the entity was created in yyyy-MM-dd'T'HH:mm:ss format
type Created = string

// Deleted Timestamp of when the entity was deleted in yyyy-MM-dd'T'HH:mm:ss format
type Deleted = string

// DeviceID ID of the device
type DeviceID = string

// EmailCodeID ID of the email OTP
type EmailCodeID = string

// EmailID ID of the email
type EmailID = string

// EmailLinkID ID of the email magic link
type EmailLinkID = string

// ErrorRsp defines model for errorRsp.
type ErrorRsp struct {
	Data  *map[string]interface{} `json:"data,omitempty"`
	Error struct {
		// Details Details of error
		Details *string `json:"details,omitempty"`

		// Links Additional links to help understand the error
		Links []string `json:"links"`

		// Type Type of error
		Type string `json:"type"`

		// Validation Validation errors per field
		Validation *[]struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"validation,omitempty"`
	} `json:"error"`

	// HttpStatusCode HTTP status code of operation
	HttpStatusCode int32  `json:"httpStatusCode"`
	Message        string `json:"message"`

	// RequestData Data about the request itself, can be used for debugging
	RequestData RequestData `json:"requestData"`

	// Runtime Runtime in seconds for this request
	Runtime float32 `json:"runtime"`
}

// FullUser User entry with emails and phone numbers
type FullUser struct {
	// ID ID of the user
	ID UserID `json:"ID"`

	// Created Timestamp of when the entity was created in yyyy-MM-dd'T'HH:mm:ss format
	Created        Created             `json:"created"`
	Emails         []UserEmail         `json:"emails"`
	FullName       string              `json:"fullName"`
	Name           string              `json:"name"`
	PhoneNumbers   []UserPhoneNumber   `json:"phoneNumbers"`
	SocialAccounts []UserSocialAccount `json:"socialAccounts"`

	// Status Generic status that can describe Corbado entities
	Status Status `json:"status"`

	// Updated Timestamp of when the entity was last updated in yyyy-MM-dd'T'HH:mm:ss format
	Updated   Updated        `json:"updated"`
	Usernames []UserUsername `json:"usernames"`
}

// GenericRsp defines model for genericRsp.
type GenericRsp struct {
	// HttpStatusCode HTTP status code of operation
	HttpStatusCode int32  `json:"httpStatusCode"`
	Message        string `json:"message"`

	// RequestData Data about the request itself, can be used for debugging
	RequestData RequestData `json:"requestData"`

	// Runtime Runtime in seconds for this request
	Runtime float32 `json:"runtime"`
}

// HighEntropyValues High entropy values from browser
type HighEntropyValues struct {
	// Mobile Mobile
	Mobile bool `json:"mobile"`

	// Platform Platform
	Platform string `json:"platform"`

	// PlatformVersion Platform version
	PlatformVersion string `json:"platformVersion"`
}

// LoginIdentifierConfig defines model for loginIdentifierConfig.
type LoginIdentifierConfig struct {
	EnforceVerification LoginIdentifierConfigEnforceVerification `json:"enforceVerification"`
	Metadata            *map[string]interface{}                  `json:"metadata,omitempty"`

	// Type Login Identifier type
	Type                 LoginIdentifierType `json:"type"`
	UseAsLoginIdentifier bool                `json:"useAsLoginIdentifier"`
}

// LoginIdentifierConfigEnforceVerification defines model for LoginIdentifierConfig.EnforceVerification.
type LoginIdentifierConfigEnforceVerification string

// LoginIdentifierType Login Identifier type
type LoginIdentifierType string

// Paging defines model for paging.
type Paging struct {
	// Page current page returned in response
	Page int `json:"page"`

	// TotalItems total number of items available
	TotalItems int `json:"totalItems"`

	// TotalPages total number of pages available
	TotalPages int `json:"totalPages"`
}

// PhoneNumberID ID of the phone number
type PhoneNumberID = string

// ProjectID ID of project
type ProjectID = string

// RequestData Data about the request itself, can be used for debugging
type RequestData struct {
	// Link Link to dashboard with details about request
	Link string `json:"link"`

	// RequestID Unique ID of request, you can provide your own while making the request, if not the ID will be randomly generated on server side
	RequestID RequestID `json:"requestID"`
}

// RequestID Unique ID of request, you can provide your own while making the request, if not the ID will be randomly generated on server side
type RequestID = string

// SessionManagement What session management should be used
type SessionManagement string

// SocialProviderType defines model for socialProviderType.
type SocialProviderType string

// Status Generic status that can describe Corbado entities
type Status string

// Updated Timestamp of when the entity was last updated in yyyy-MM-dd'T'HH:mm:ss format
type Updated = string

// UserEmail User's email
type UserEmail struct {
	// ID generic ID
	ID ID `json:"ID"`

	// Created Timestamp of when the entity was created in yyyy-MM-dd'T'HH:mm:ss format
	Created Created `json:"created"`
	Email   string  `json:"email"`

	// Status Generic status that can describe Corbado entities
	Status Status `json:"status"`

	// Updated Timestamp of when the entity was last updated in yyyy-MM-dd'T'HH:mm:ss format
	Updated Updated `json:"updated"`
}

// UserID ID of the user
type UserID = string

// UserPhoneNumber User's phone number
type UserPhoneNumber struct {
	// ID generic ID
	ID ID `json:"ID"`

	// Created Timestamp of when the entity was created in yyyy-MM-dd'T'HH:mm:ss format
	Created     Created `json:"created"`
	PhoneNumber string  `json:"phoneNumber"`

	// Status Generic status that can describe Corbado entities
	Status Status `json:"status"`

	// Updated Timestamp of when the entity was last updated in yyyy-MM-dd'T'HH:mm:ss format
	Updated Updated `json:"updated"`
}

// UserSocialAccount User's social account
type UserSocialAccount struct {
	AvatarUrl       string             `json:"avatarUrl"`
	FullName        string             `json:"fullName"`
	IdentifierValue string             `json:"identifierValue"`
	ProviderType    SocialProviderType `json:"providerType"`
}

// UserUsername User's username
type UserUsername struct {
	// ID generic ID
	ID ID `json:"ID"`

	// Created Timestamp of when the entity was created in yyyy-MM-dd'T'HH:mm:ss format
	Created Created `json:"created"`

	// Status Generic status that can describe Corbado entities
	Status Status `json:"status"`

	// Updated Timestamp of when the entity was last updated in yyyy-MM-dd'T'HH:mm:ss format
	Updated  Updated `json:"updated"`
	Username string  `json:"username"`
}

// Filter defines model for filter.
type Filter = []string

// Page defines model for page.
type Page = int

// PageSize defines model for pageSize.
type PageSize = int

// RemoteAddress defines model for remoteAddress.
type RemoteAddress = string

// SessionID defines model for sessionID.
type SessionID = string

// Sort defines model for sort.
type Sort = string

// UserAgent defines model for userAgent.
type UserAgent = string

// UnusedParams defines parameters for Unused.
type UnusedParams struct {
	// RemoteAddress Client's remote address
	RemoteAddress *RemoteAddress `form:"remoteAddress,omitempty" json:"remoteAddress,omitempty"`

	// UserAgent Client's user agent
	UserAgent *UserAgent `form:"userAgent,omitempty" json:"userAgent,omitempty"`

	// Sort Field sorting
	Sort *Sort `form:"sort,omitempty" json:"sort,omitempty"`

	// Filter Field filtering
	Filter *Filter `form:"filter[],omitempty" json:"filter[],omitempty"`

	// Page Page number
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`
}
