package data

type Formated struct {
	Linhas int
	Script string
}

type Log struct {
	Id         int
	Metodo     string
	Horas      string
	Sql        []Formated
	Outros     []string
	Requisicao Formated
	Resposta   Formated
}

type base struct {
	id     int
	Status int8
	Logs   []Log
	LogsB  []Log
}

var b *base

func GetInstance() *base {
	if b == nil {
		b = &base{id: 0, Status: 1}
	}

	return b
}

func NewLog() *Log {
	b.id++
	return &Log{Id: b.id}
}

func AddFormated(str string, linhas int) *Formated {
	return &Formated{Linhas: linhas, Script: str}
}
