package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

		version, ok := os.LookupEnv("APP_VERSION")
		if ok {
			rw.Header().Set("App-Version", version)
		}
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
		version, ok := os.LookupEnv("APP_VERSION")
		if ok {
			rw.Header().Set("App-Version", version)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(rw, string(cresb))
		return
	}

	// ####### bunch of things we make here ################
	data, _ := json.Marshal(request_payload)

	// #### Passing to the client the req ######

	downstream := quiknode_proxy.downstream_fast_rpc_client_addr
	network, ok := os.LookupEnv("NETWORK")
	if ok {
		if network == "mainnet" {
			if request_payload.Method == "eth_getBalance" {
				downstream = quiknode_proxy.downstream_archive_rpc_client_addr
			}
		}
	}

	respBody, err := httpPost(downstream, data)
	if err != nil {
		logrus.Error("Problem to connect to the local client: ", err.Error())
		cbody := json.RawMessage(`{"code":-32000,"message":"problem to connect downstream"}`)
		var cresp = Response{
			Jsonrpc: "2.0",
			Error:   cbody,
		}
		cresb, _ := json.Marshal(cresp)
		version, ok := os.LookupEnv("APP_VERSION")
		if ok {
			rw.Header().Set("App-Version", version)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		logrus.Infoln("what is here", string(cresb))
		io.WriteString(rw, string(cresb))
		return
	}
	version, ok := os.LookupEnv("APP_VERSION")
	if ok {
		rw.Header().Set("App-Version", version)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, string(respBody))
}

func healthzHandler(rw http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		version, ok := os.LookupEnv("APP_VERSION")
		if ok {
			rw.Header().Set("App-Version", version)
		}
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(rw, `{"alive": true}`)
		return
	}
	rw.WriteHeader(http.StatusServiceUnavailable)
}

func make_it_failHandler(rw http.ResponseWriter, r *http.Request) {
	atomic.StoreInt32(&healthy, 0)
	version, ok := os.LookupEnv("APP_VERSION")
	if ok {
		rw.Header().Set("App-Version", version)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, `{"done": "requests will return 503}`)
}

func make_it_workHandler(rw http.ResponseWriter, r *http.Request) {
	atomic.StoreInt32(&healthy, 1)
	version, ok := os.LookupEnv("APP_VERSION")
	if ok {
		rw.Header().Set("App-Version", version)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, `{"done": "requests will return 200}`)
}
