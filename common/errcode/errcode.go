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

// Warn level code
const (
	Md5CheckNotPassCode = 20100
)

// Warn level msg
const (
	Md5CheckNotPassMsg = "[warn]file md5 check didn't pass, feel free to try again"
)

// Error level code
const (
	DatabaseOperationErrCode        = 50100
	ParseHTTPRequestFormFileErrCode = 50101
	OsOperationErrCode              = 50102
)

// Error level msg
const (
	DatabaseOperationErrMsg        = "[error]database operation failed, err msg: %v"
	ParseHTTPRequestFormFileErrMsg = "[error]parse http request failed, err msg: %v"
	OsOperationErrMsg              = "[error]os operation failed, err msg: %v"
)
