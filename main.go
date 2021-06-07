package main

import (
	"flag"
	"sync/atomic"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/quiknode-labs/indexer/exporter"
)

type Quiknode struct {
	app_port string
	healthy  int32
}

var (
	mode *string
)

func init() {
	mode = flag.String("mode", "stream", "mode: stream / exporter")
}

func main() {
	flag.Parse()
	if *mode == "stream" {
		q := new(Quiknode)
		q.app_port = "5000"
		atomic.StoreInt32(&q.healthy, 1)
		router := mux.NewRouter().StrictSlash(false)
		router.HandleFunc("/", q.proxyHandler)
		router.HandleFunc("/healthz", q.healthzHandler)
		n := negroni.New(
			negroni.NewRecovery(),
			negroni.NewLogger(),
		)
		n.UseHandler(router)
		n.Run(":" + q.app_port)
	} else if *mode == "exporter" {
		oct := exporter.QuiknodeExporter{
			RPC: "http://127.0.0.1:8545",
		}
		oct.Exporterlogs()
	}
}
