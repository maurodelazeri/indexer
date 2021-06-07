package main

import (
	"sync/atomic"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type QuiknodeIndexer struct {
	app_port string
	healthy  int32
}

func main() {
	q := new(QuiknodeIndexer)
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
}
