package errcode

// Tips level code
const (
	ParamsInvalidCode                  = 10100
	DetectRepeatedUserIDCode           = 10101
	NoSuchUserCode                     = 10102
	NoEnoughVolCode                    = 10103
	FindNoFileInServerCode             = 10104
	PathNotExistCode                   = 10105
	DirOrFileAlreadyExistCode          = 10106
	ForbidHardDeleteFolderCode         = 10107
	NoSuchShareCode                    = 10108
	SameFileHasAlreadyBeenUploadedCode = 10109
	AuthMatchFailCode                  = 10110
	WrongSharePasswordCode             = 10111
)

// Tips level msg
const (
	ParamsInvalidMsg                  = "[tips] params invalid, feel free to check and try again"
	DetectRepeatedUserIDMsg           = "[tips] this userID is using by other user, feel free to check and try again"
	NoSuchUserMsg                     = "[tips] cannot find such user, recheck id and password"
	NoEnoughVolMsg                    = "[tips] your left volume is not enough, please contact the admin to expand it"
	FindNoFileInServerMsg             = "[tips] cannot use quick upload because this file hasn't been uploaded yet, try upload"
	PathNotExistMsg                   = "[tips] path is not existed, check the input of the path"
	DirOrFileAlreadyExistMsg          = "[tips] there is already a file or folder respond to input name and path exist, do you want to overwrite?"
	ForbidHardDeleteFolderMsg         = "[tips] force delete a folder is not allowed, delete the files one by one if you insist"
	NoSuchShareMsg                    = "[tips] cannot find such share, maybe it has expired if once existed"
	SameFileHasAlreadyBeenUploadedMsg = "[tips] there is a same file already existed in your collection, no need to upload"
	AuthMatchFailMsg                  = "[tips] auth value didn't match given user"
	WrongSharePasswordMsg             = "[tips] wrong password for this share, check params and retry"
)

// Warn level code
const (
	Md5CheckNotPassCode     = 20100
	FindCountNotMatchCode   = 20101
	UpdateCountNotMatchCode = 20102
)

// Warn level msg
const (
	Md5CheckNotPassMsg     = "[warn] file md5 check didn't pass, feel free to try again"
	FindCountNotMatchMsg   = "[warn] actual query count didn't match the expected one"
	UpdateCountNotMatchMsg = "[warn] actual update count didn't match the expected one"
)

// Error level code
const (
	DatabaseOperationErrCode        = 50100
	ParseHTTPRequestFormFileErrCode = 50101
	OsOperationErrCode              = 50102
	RPCCallErrCode                  = 50103
	JsonMarshalErrCode              = 50104
	StructpbConvertErrCode          = 50105
	JWTParseErrCode                 = 50106
)

// Error level msg
const (
	DatabaseOperationErrMsg        = "[error] database operation failed, err msg: %v"
	ParseHTTPRequestFormFileErrMsg = "[error] parse http request failed, err msg: %v"
	OsOperationErrMsg              = "[error] os operation failed, err msg: %v"
	RPCCallErrMsg                  = "[error] rpc call failed, err msg: %v"
	JsonMarshalErrMsg              = "[error] json marshal or unmarshal failed, err msg: %v"
	StructpbConvertErrMsg          = "[error] structpb conversion failed, err msg: %v"
	JWTParseErrMsg                 = "[error] jwt string signing failed, err msg: %v"
)
