package core

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BrendoRochaDel/logs-cgi/back_end_GO/data"
	"github.com/BrendoRochaDel/logs-cgi/back_end_GO/format"
)

func Analytical() {
	for {
		if len(db.LogsB) > 0 {
			for i := len(db.LogsB) - 1; i >= 0; i-- {
				if validarMetodo(db.LogsB[i].Metodo) {
					go func(log data.Log) {
						fmt.Println(log.Metodo, ": Iniciando tratamento")

						//TratarSQL
						fmt.Println(log.Metodo, ": Tratando sql")
						validarSQL(&log.Sql)
						removerDuplicidadeSQL(&log.Sql)
						identarSQL(&log.Sql)

						//TratarXML
						fmt.Println(log.Metodo, ": Tratando xml")
						formatarRequisicao(&log.Requisicao)
						formatarResposta(&log.Resposta)

						fmt.Println(log.Metodo, ": Salvando Alteracoes")
						push(log)
						fmt.Println(log.Metodo, ": Processo Concluido!")
					}(db.LogsB[i])
				}

				remove(i)
			}
		}
	}
}

func validarMetodo(metodo string) bool {
	for _, item := range abrirValidador("config/ignoreMetodo.txt") {
		if strings.Contains(metodo, item) {
			fmt.Println(metodo, ": Removido")
			return false
		}
	}

	return true
}

func validarSQL(sqls *[]data.Formated) {
	config := abrirValidador("config/ignoreSQL.txt")

	for i := len(*sqls) - 1; i >= 0 && len(config) > 0; i-- {
		for _, item := range config {
			if strings.Contains((*sqls)[i].Script, item) {
				*sqls = append((*sqls)[:i], (*sqls)[i+1:]...)
				break
			}
		}
	}
}

func abrirValidador(path string) (linhas []string) {
	conteudo, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	if string(conteudo) != "" {
		linhas = strings.Split(string(conteudo), "\n")
	}

	return
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

func push(log data.Log) {
	db.Logs = append(db.Logs, log)
}

func remove(pos int) {
	db.LogsB = append(db.LogsB[:pos], db.LogsB[pos+1:]...)
}
