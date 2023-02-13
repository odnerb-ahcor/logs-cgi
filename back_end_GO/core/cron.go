package core

import (
	"encoding/json"
	"fmt"

	"github.com/BrendoRochaDel/logs-cgi/back_end_GO/data"
	"github.com/BrendoRochaDel/logs-cgi/back_end_GO/format"
)

func Analytical() {
	for {
		if len(db.LogsB) > 0 {

			for i, log := range db.LogsB {
				fmt.Println("Existe log")
				//TratarSQL
				removerDuplicidadeSQL(&log.Sql)
				identarSQL(&log.Sql)

				//TratarXML
				formatarRequisicao(&log.Requisicao)
				formatarResposta(&log.Resposta)

				push(log, i)
				fmt.Println("")
				//imprimirStruct()
			}
		}
	}
}

func imprimirStruct() {
	b, err := json.Marshal(db.Logs)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(string(b))
}

func removerDuplicidadeSQL(sqls *[]data.Formated) {
	seen := make(map[string]bool)
	j := 0
	for i, val := range *sqls {
		if _, ok := seen[val.Script]; !ok {
			seen[val.Script] = true
			(*sqls)[j] = (*sqls)[i]
			j++
		}
	}
	*sqls = (*sqls)[:j]
}

func identarSQL(sqls *[]data.Formated) {
	for i, sql := range *sqls {
		(*sqls)[i].Script, (*sqls)[i].Linhas = format.SqlFormat(sql.Script)
	}
}

func formatarRequisicao(xml *data.Formated) {
	xml.Script, xml.Linhas = format.XMLFormat(xml.Script)
}

func formatarResposta(xml *data.Formated) {
	xml.Script, xml.Linhas = format.XMLtoJson(xml.Script)
}

func push(log data.Log, pos int) {
	db.Logs = append(db.Logs, log)
	db.LogsB = append(db.LogsB[:pos], db.LogsB[pos+1:]...)
	fmt.Println("Log Processado")
}
