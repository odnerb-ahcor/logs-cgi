package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/odnerb-ahcor/logs-cgi/back_end_GO/data"
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

	err := os.WriteFile("logs/config.conf", []byte(strconv.Itoa(int(db.Status))), 0644)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo", err)
		return
	}
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

func LimparLogs() {
	db.Logs = nil
	db.Id = 0
}
