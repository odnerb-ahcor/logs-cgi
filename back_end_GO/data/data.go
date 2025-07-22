package data

type Formated struct {
	Linhas int
	Script string
}

type Log struct {
	Id         int
	NameFile   string
	Metodo     string
	Horas      string
	Sql        []Formated
	Outros     []string
	Requisicao Formated
	Resposta   Formated
}

type base struct {
	Id     int
	Status int8
	Logs   []Log
	Log    chan Log
}

var b *base

func GetInstance() *base {
	if b == nil {
		b = &base{
			Id:     0,
			Status: 1,
			Log:    make(chan Log, 5),
		}
	}

	return b
}

func NewLog() *Log {
	b.Id++
	return &Log{Id: b.Id}
}

func AddFormated(str string, linhas int) *Formated {
	return &Formated{Linhas: linhas, Script: str}
}

func (b *base) FilterLogs(filterFunc func(Log) bool) []Log {
	var FilterSlice []Log
	for _, element := range b.Logs {
		if filterFunc(element) {
			FilterSlice = append(FilterSlice, element)
		}
	}

	return FilterSlice
}

func (b *base) FindLogs(filterFunc func(Log) bool) int {
	index := 0
	for i, element := range b.Logs {
		if filterFunc(element) {
			index = i
		}
	}

	return index
}

func FilterSlice[T any](slice []T, filterFunc func(T) bool) []T {
	var FilterSlice []T
	for _, element := range slice {
		if filterFunc(element) {
			FilterSlice = append(FilterSlice, element)
		}
	}

	return FilterSlice
}

func IndexFunc[E any](s []E, f func(E) bool) int {
	for i, v := range s {
		if f(v) {
			return i
		}
	}
	return -1
}
