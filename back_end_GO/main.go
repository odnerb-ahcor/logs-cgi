package main

import (
	"net/http"

	"github.com/BrendoRochaDel/logs-cgi/back_end_GO/core"
)

type Pessoa struct {
	name  string
	idade int16
	sexo  string
}

func main() {
	go core.Analytical()
	HandleResquest()
}

func FilterSlice[T any](slice []T, filterFunc func(T) bool) []T {
	var FilterSlice []T
	for _, element := range slice {
		if filterFunc(element) {
			FilterSlice = append(FilterSlice, element)
		}
	}

	return FilterSlice
}

func HandleResquest() {
	http.HandleFunc("/log/", core.Log)
	http.HandleFunc("/status/", core.Status)
	http.HandleFunc("/arquivo", core.Arquivo)
	http.HandleFunc("/limpar", core.Limpar)

	http.ListenAndServe(":5000", nil)
}
