package common

const RequestSuccessCode = 200

// Debug level code
const (
	InvalidDataTypeCode int32 = 1
)

// Debug level msg
const (
	InvalidDataTypeMsg = "[debug]invalid data type, proto.Message is expected"
)

// Tips level code
const (
	ParamsInvalidCode int32 = 10
)

// Tips level msg
const (
	ParamInvalidCode = "[tips]params invalid, feel free to check and try again"
)
