package errcode

// Tips level code
const (
	ParamsInvalidCode        = 10100
	DetectRepeatedUserIDCode = 10101
	NoSuchUserCode           = 10102
)

// Tips level msg
const (
	ParamsInvalidMsg        = "[tips]params invalid, feel free to check and try again"
	DetectRepeatedUserIDMsg = "[tips]this userID is using by other user, feel free to check and try again"
	NoSuchUserMsg           = "[tips]cannot find such user, recheck id and password"
)

// Error level code
const (
	DatabaseOperationErrCode = 50100
)

// Error level msg
const (
	DatabaseOperationErrMsg = "[error]database operation failed, err msg: %v"
)
