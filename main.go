package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type QuiknodeProxy struct {
	app_port   string
	client_rpc string
	db         *sql.DB
}

const (
	host     = "161.35.98.164"
	port     = 5432
	user     = "postgres"
	password = "Br@sa154"
	dbname   = "ethereum_mainnet"
)

var (
	healthy        int32
	quiknode_proxy QuiknodeProxy
)

func init() {
	quiknode_proxy.client_rpc = "http://apps.zinnion.com:8547"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	quiknode_proxy.db = db

	err = quiknode_proxy.db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected postgres!")

	atomic.StoreInt32(&healthy, 1)

	app_port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		fmt.Println("APP_PORT is not present")
		os.Exit(1)
	}
	quiknode_proxy.app_port = app_port
}

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", proxyHandler)
	router.HandleFunc("/healthz", healthzHandler)
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	n.UseHandler(router)
	n.Run(":" + quiknode_proxy.app_port)
}
