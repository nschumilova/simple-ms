package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nschumilova/simple-ms/healthcheck"
	"go.uber.org/zap"
)

func main() {
	port:=flag.Int("port", 8080, "application port")
	flag.Parse()

	log, _ := zap.NewProduction()
	defer log.Sync()
	sugaredLog := log.Sugar()

	httpHealthHandler := healthcheck.HeathCheckHandler{
		Log: sugaredLog,
	}
	router := mux.NewRouter()
	router.HandleFunc("/health/", httpHealthHandler.Health).Methods("GET")
	sugaredLog.Error(http.ListenAndServe(fmt.Sprintf(":%d", *port), router))
}
