package util

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
		log.Fatalln("not a valid grpc error")
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
