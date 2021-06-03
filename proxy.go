package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

func httpPost(url string, data []byte) (string, error) {
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logrus.Error("Problem to create request", err.Error())
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("Problem making request: ", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Problem to read response", err.Error())
		return "", err
	}
	return string(respBody), nil
}

func proxyHandler(rw http.ResponseWriter, r *http.Request) {
	// Request payload
	var request_payload Request

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message.
	err := json.NewDecoder(r.Body).Decode(&request_payload)

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
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(rw, string(cresb))
		return
	}

	if request_payload.Method == "eth_getLogs" {
		get_logs(rw, r, request_payload)
		return
	}

	cbody := json.RawMessage(`{"code":-32000,"message":"method not supported by this service"}`)
	var cresp = Response{
		Jsonrpc: "2.0",
		Error:   cbody,
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	json.NewEncoder(rw).Encode(cresp)

}

func healthzHandler(rw http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(rw, `{"alive": true}`)
		return
	}
	rw.WriteHeader(http.StatusServiceUnavailable)
}
