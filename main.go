package main

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type QuiknodeProxy struct {
	app_port                           string
	downstream_fast_rpc_client_addr    string
	downstream_fast_ws_client_addr     string
	downstream_archive_rpc_client_addr string
	downstream_archive_ws_client_addr  string
}

var (
	healthy        int32
	quiknode_proxy QuiknodeProxy
)

func init() {
	atomic.StoreInt32(&healthy, 1)

	app_port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		fmt.Println("APP_PORT is not present")
		os.Exit(1)
	}
	quiknode_proxy.app_port = app_port

	// Archive fast
	downstream_fast_rpc_client_addr, ok := os.LookupEnv("CLIENT_DOWN_STREAM_FAST_RPC")
	if !ok {
		fmt.Println("CLIENT_DOWN_STREAM_FAST_RPC is not present")
		os.Exit(1)
	}
	quiknode_proxy.downstream_fast_rpc_client_addr = downstream_fast_rpc_client_addr

	downstream_fast_ws_client_addr, ok := os.LookupEnv("CLIENT_DOWN_STREAM_FAST_WS")
	if !ok {
		fmt.Println("CLIENT_DOWN_STREAM_FAST_WS is not present")
		os.Exit(1)
	}
	quiknode_proxy.downstream_fast_ws_client_addr = downstream_fast_ws_client_addr

	// Archive mode
	downstream_archive_rpc_client_addr, ok := os.LookupEnv("CLIENT_DOWN_STREAM_ARCHIVE_RPC")
	if !ok {
		fmt.Println("CLIENT_DOWN_STREAM_ARCHIVE_RPC is not present")
	}
	quiknode_proxy.downstream_archive_rpc_client_addr = downstream_archive_rpc_client_addr

	downstream_archive_ws_client_addr, ok := os.LookupEnv("CLIENT_DOWN_STREAM_ARCHIVE_WS")
	if !ok {
		fmt.Println("CLIENT_DOWN_STREAM_ARCHIVE_WS is not present")
	}
	quiknode_proxy.downstream_archive_ws_client_addr = downstream_archive_ws_client_addr

}

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", proxyHandler)
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/make_it_fail", make_it_failHandler)
	router.HandleFunc("/make_it_work", make_it_workHandler)
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	n.UseHandler(router)
	n.Run(":" + quiknode_proxy.app_port)
}
