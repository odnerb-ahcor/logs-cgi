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

func RetornarLog(id int, horas string) string {
	var logs []data.Log
	index := -1

	if id > 0 && horas != "" {
		index = db.FindLogs(func(log data.Log) bool {
			return log.Id == id && log.Horas == horas
		})
	}

	if index >= 0 {
		logs = db.FilterLogs(func(log data.Log) bool {
			return log.Id > id
		})
	} else {
		logs = db.Logs
	}

	jsonLogs, err := json.Marshal(logs)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	return string(jsonLogs)
}

func LerArquivo(body io.Reader) string {
	arq := &arquivo{}

	err := json.NewDecoder(body).Decode(arq)
	if err != nil {
		fmt.Println("json erro: ", err)
	}

	log := data.NewLog()

	for _, linha := range arq.Linhas {
		ler_log(linha, log)
	}

	db.LogsB = append(db.LogsB, *log)
	//time.Sleep(8 * time.Second)
	return "nome"
}

func LimparLogs() {
	db.Logs = nil
	db.Id = 0
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
