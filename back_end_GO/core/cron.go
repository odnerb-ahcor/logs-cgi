package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/odnerb-ahcor/logs-cgi/back_end_GO/data"
	"github.com/odnerb-ahcor/logs-cgi/back_end_GO/format"
)

func Analytical() {
	for log := range db.Log {
		if validarMetodo(log.Metodo) {
			go func() {
				fmt.Println(log.Metodo, ": Iniciando tratamento")

				//TratarSQL
				fmt.Println(log.Metodo, ": Tratando sql")
				validarSQL(&log.Sql)
				removerDuplicidadeSQL(&log.Sql)
				identarSQL(&log.Sql)

				//TratarXML
				fmt.Println(log.Metodo, ": Tratando xml")
				formatarRequisicao(&log.Requisicao)
				err := formatarResposta(&log.Resposta)
				if err != nil {
					fmt.Printf("Erro ao formatar resposta %s: %s\n", log.NameFile, err)
					os.Rename(log.NameFile, log.Metodo+".error")
					return
				}

				fmt.Println(log.Metodo, ": Salvando Alteracoes")
				push(log)
				fmt.Println(log.Metodo, ": Processo Concluido!")

				//Exclur arquivo
				if log.NameFile != "" {
					err = os.Remove(log.NameFile)
					if err != nil {
						fmt.Printf("Erro ao excluir arquivo %s: %s\n", log.NameFile, err)
					}
				}
			}()
		} else {
			if log.NameFile != "" {
				err := os.Remove(log.NameFile)
				if err != nil {
					fmt.Printf("Erro ao excluir arquivo  %s: %s\n", log.NameFile, err)
				}
			}
		}
	}
}

func validarMetodo(metodo string) bool {
	for _, item := range abrirValidador("config/ignoreMetodo.txt") {
		if len(item) > 0 && strings.Contains(metodo, item) {
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
			if len(item) > 0 && strings.Contains((*sqls)[i].Script, item) {
				*sqls = append((*sqls)[:i], (*sqls)[i+1:]...)
				break
			}
		}
	}
}

func abrirValidador(path string) (linhas []string) {
	conteudo, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	if string(conteudo) != "" {
		lista := strings.ReplaceAll(string(conteudo), "\r", "")
		linhas = strings.Split(lista, "\n")
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

func formatarResposta(xml *data.Formated) error {
	retornoXML, err := format.XMLtoJson(xml.Script)
	if err != nil {
		return err
	}
	xml.Script = retornoXML.Data
	xml.Linhas = retornoXML.Lines
	return nil
}

func push(log data.Log) {
	db.Logs = append(db.Logs, log)
}
