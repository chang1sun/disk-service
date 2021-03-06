package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type errorBody struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func CustomErrorHandler(c context.Context, sm *runtime.ServeMux, m runtime.Marshaler,
	rw http.ResponseWriter, r *http.Request, e error) {
	status, ok := status.FromError(e)
	if !ok {
		log.Println("not a valid grpc error")
	}
	rw.Header().Set("ContentType", "application/json")
	errbody := &errorBody{
		Code: uint32(status.Code()),
		Msg:  status.Message(),
	}
	err := json.NewEncoder(rw).Encode(errbody)
	if err != nil {
		rw.Write([]byte("json marshal failed"))
	}
}

// used in middleware and custom route handler to build error response body
func ErrorResp(code uint32, msg string, err error, w *http.ResponseWriter) {
	log.Printf("err occured, code: %v, msg: %v", code, fmt.Sprintf(msg, err))
	if err == nil {
		rsp, _ := json.Marshal(map[string]string{"code": strconv.Itoa(int(code)), "msg": msg})
		(*w).Write(rsp)
	} else {
		rsp, _ := json.Marshal(map[string]string{"code": strconv.Itoa(int(code)), "msg": fmt.Sprintf(msg, err)})
		(*w).Write(rsp)
	}
}

func LogErr(err error, intName string) {
	s, _ := status.FromError(err)
	if uint32(s.Code()) < uint32(20000) {
		return
	}
	log.Printf("[%v] request failed, err msg: %v", intName, err)
}
