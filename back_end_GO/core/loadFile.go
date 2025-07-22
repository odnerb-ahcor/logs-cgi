package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/odnerb-ahcor/logs-cgi/back_end_GO/data"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const DIR = "logs/file_logs/"

func LerArquivos() {
	fileNames := make(chan string, 5)

	go carregarArquivos(fileNames)

	loadFile(fileNames)
}

func carregarArquivos(fileNames chan<- string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("Erro carregar arquivos" + err.Error())
	}
	defer watcher.Close()

	err = watcher.Add(DIR)
	if err != nil {
		panic(fmt.Sprintf("Erro ao adicionar diretório ao watcher: %v", err))
	}

	for {
		select {
		case evento, ok := <-watcher.Events:
			if !ok {
				return
			}

			if evento.Op&(fsnotify.Create) != 0 {
				if filepath.Ext(evento.Name) == ".log" {
					newName := strings.TrimSuffix(evento.Name, filepath.Ext(evento.Name)) + ".data"
					os.Rename(evento.Name, newName)
					fileNames <- newName
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			panic(fmt.Sprintf("Erro no watcher:", err))
		}
	}
}

func loadFile(fileNames <-chan string) {
	for fileNames := range fileNames {
		arq, err := os.Open(fileNames)
		if err != nil {
			fmt.Println("Erro ao abrir o arquivo", err)
		}

		defer arq.Close()

		scanner := bufio.NewScanner(arq)

		decoder := transform.NewReader(arq, charmap.Windows1250.NewDecoder())
		scanner = bufio.NewScanner(decoder)

		db := data.GetInstance()
		log := data.NewLog()
		log.NameFile = fileNames

		const maxTokenSize = 10 * 1024 * 1024
		scanner.Buffer(make([]byte, maxTokenSize), maxTokenSize)

		for scanner.Scan() {
			ler_logs(scanner.Text(), log)
		}

		db.Log <- *log
	}
}

func ler_logs(linha string, log *data.Log) {
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
