package main

import (
	"net/http"

	"github.com/BrendoRochaDel/logs-cgi/back_end_GO/core"
)

func main() {
	go core.Analytical()
	HandleResquest()
}

func HandleResquest() {
	http.HandleFunc("/log/", core.Log)
	http.HandleFunc("/status/", core.Status)
	http.HandleFunc("/arquivo", core.Arquivo)

	http.ListenAndServe(":5000", nil)
}
