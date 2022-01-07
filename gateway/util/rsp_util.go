package util

import (
	"github.com/changpro/disk-service/gateway/common"
	gatewaypb "github.com/changpro/disk-service/resource/stub"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func SuccessRsp(data interface{}) *gatewaypb.CommonHttpRsp {
	msgData, ok := data.(proto.Message)
	if !ok {
		return &gatewaypb.CommonHttpRsp{
			Code: common.InvalidDataTypeCode,
			Msg:  common.InvalidDataTypeMsg,
			Data: &anypb.Any{},
		}
	}
	anyData, err := anypb.New(msgData)
	if err != nil {
		return &gatewaypb.CommonHttpRsp{
			Code: common.InvalidDataTypeCode,
			Msg:  common.InvalidDataTypeMsg,
			Data: &anypb.Any{},
		}
	}
	return &gatewaypb.CommonHttpRsp{
		Code: common.RequestSuccessCode,
		Msg:  "success",
		Data: anyData,
	}
}

func ErrorRsp(code int32, msg string, data interface{}) *gatewaypb.CommonHttpRsp {
	msgData, ok := data.(proto.Message)
	if !ok {
		return &gatewaypb.CommonHttpRsp{
			Code: common.InvalidDataTypeCode,
			Msg:  common.InvalidDataTypeMsg,
			Data: &anypb.Any{},
		}
	}
	anyData, err := anypb.New(msgData)
	if err != nil {
		return &gatewaypb.CommonHttpRsp{
			Code: common.InvalidDataTypeCode,
			Msg:  common.InvalidDataTypeMsg,
			Data: &anypb.Any{},
		}
	}
	return &gatewaypb.CommonHttpRsp{
		Code: code,
		Msg:  msg,
		Data: anyData,
	}
}
