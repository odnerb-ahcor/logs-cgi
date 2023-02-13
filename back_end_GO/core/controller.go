package core

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/BrendoRochaDel/logs-cgi/back_end_GO/data"
)

type arquivo struct {
	Linhas []string `json:"arq"`
}

var db = data.GetInstance()

func RetornarStatus() int8 {
	return db.Status
}

func AlteraStatus(status int8) {
	db.Status = status
}

func RetornarLog() string {
	b, err := json.Marshal(db.Logs)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	return string(b)
}

func LerArquivo(body io.Reader) string {
	arq := &arquivo{}

	err := json.NewDecoder(body).Decode(arq)
	if err != nil {
		fmt.Println("json erro")
	}

	log := data.NewLog()

	for _, linha := range arq.Linhas {
		ler_log(linha, log)
	}

	db.LogsB = append(db.LogsB, *log)
	//time.Sleep(8 * time.Second)
	return "nome"
}

func ler_log(linha string, log *data.Log) {
	switch linha[0:3] {
	case "met":
		log.Metodo = linha[5:]
	case "sql":
		log.Sql = append(log.Sql, *data.AddFormated(linha[5:], 0))
	case "req":
		log.Requisicao.Script = linha[5:]
	case "res":
		log.Resposta.Script = linha[5:]
	case "hor":
		log.Horas = linha[5:]
	default:
		log.Outros = append(log.Outros, linha)
	}
}
