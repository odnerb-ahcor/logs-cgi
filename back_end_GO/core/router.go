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

func Limpar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	LimparLogs()
	fmt.Println("Base de dados apagado")
	fmt.Fprint(w, "Success!")
}

func Log(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")

	fmt.Println("Verificacao de logs iniciado...")
	path := strings.TrimPrefix(r.URL.Path, "/log/")
	params := strings.Split(path, "/")

	id, err := strconv.Atoi(params[0])
	if len(params) > 1 && err == nil {
		fmt.Fprint(w, RetornarLog(id, params[1]))
	} else {
		fmt.Fprint(w, RetornarLog(0, ""))
	}
}
