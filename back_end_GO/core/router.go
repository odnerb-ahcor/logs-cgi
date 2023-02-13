package core

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	path := strings.TrimPrefix(r.URL.Path, "/status/")
	fmt.Println("Verificacao de status iniciado...")

	if path == "" {
		status := RetornarStatus()
		fmt.Println("Status: ", status)
		fmt.Fprint(w, status)
		return
	}

	status, err := strconv.Atoi(path)
	if err != nil {
		log.Fatal(err)
	}
	AlteraStatus(int8(status))
	fmt.Fprint(w, "Success!")
}

func Arquivo(w http.ResponseWriter, r *http.Request) {
	nome := LerArquivo(r.Body)
	fmt.Fprintln(w, nome)
}

func Log(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	fmt.Println("Verificacao de logs iniciado...")
	fmt.Fprint(w, RetornarLog())
}
