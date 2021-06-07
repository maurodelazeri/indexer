package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

func (q *Quiknode) proxyHandler(rw http.ResponseWriter, r *http.Request) {
	call := &rpcCall{}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message.
	err := json.NewDecoder(r.Body).Decode(&call)
	if err != nil {
		logrus.Error("wrong method, only POST available.")
		cbody := json.RawMessage(`{"code":-32000,"message":"not able to decode body"}`)
		var cresp = Response{
			Jsonrpc: "2.0",
			Error:   cbody,
		}
		cresb, _ := json.Marshal(cresp)

		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		io.WriteString(rw, string(cresb))
		return

	}
	if r.Method != "POST" {
		logrus.Error("wrong method, only POST available.")
		cbody := json.RawMessage(`{"code":-32000,"message":"wrong method, only POST available"}`)
		var cresp = Response{
			Jsonrpc: "2.0",
			Error:   cbody,
		}
		cresb, _ := json.Marshal(cresp)
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		io.WriteString(rw, string(cresb))
		return
	}

	if call.Method == "eth_getLogs" {
		q.get_logs(r.Context(), rw, r, call)
		return
	}

	cbody := json.RawMessage(`{"code":-32000,"message":"method not supported by this service"}`)
	var cresp = Response{
		Jsonrpc: "2.0",
		Error:   cbody,
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	json.NewEncoder(rw).Encode(cresp)

}

func (q *Quiknode) healthzHandler(rw http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&q.healthy) == 1 {
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(rw, `{"alive": true}`)
		return
	}
	rw.WriteHeader(http.StatusServiceUnavailable)
}
