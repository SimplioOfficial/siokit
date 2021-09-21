package error

type Code string

const (
	CodeUnknown         = "Unknown"
	CodeInvalidArgument = "InvalidArgument"
	CodeNotFound        = "NotFound"
	CodeJwtInvalid      = "JwtInvalid"
	CodeJwtMalformed    = "JwtMalformed"
	CodeDbDuplication   = "DuplicatedValue"
	CodeDbReference     = "InvalidReference"
	CodeValidation      = "Validation"
)
