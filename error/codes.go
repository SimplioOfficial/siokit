package error

type Code string

const (
	CodeUnknown         = "Unknown"
	CodeInvalidArgument = "InvalidArgument"
	CodeNotFound        = "NotFound"
	CodeJwtInvalid      = "JwtInvalid"
	CodeJwtMalformed    = "JwtMalformed"
	CodeForbidden       = "Forbidden"
	CodeDbDuplication   = "DuplicatedValue"
	CodeDbReference     = "InvalidReference"
	CodeValidation      = "Validation"
)
