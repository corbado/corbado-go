package validationerror

type Code int

const (
	CodeJWTGeneral Code = iota
	CodeJWTIssuerMismatch
	CodeJWTInvalidData
	CodeJWTInvalidSignature
	CodeJWTBefore
	CodeJWTExpired
)
