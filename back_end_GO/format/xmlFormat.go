package format

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"regexp"
	"strings"
)

type RetornoXML struct {
	Data  string
	Lines int
}

func XMLFormat(str string) (xmlIndented string, linhas int) {

	type node struct {
		Attr     []xml.Attr
		XMLName  xml.Name
		Children []node `xml:",any"`
		Text     string `xml:",chardata"`
	}
	x := node{}
	_ = xml.Unmarshal([]byte(str), &x)
	buf, _ := xml.MarshalIndent(x, "", "   ")

	header := regexp.MustCompile(`<\?xml.*\?>`)
	xmlIndented = header.FindString(str) + "\n"
	xmlIndented += string(buf)
	linhas = strings.Count(xmlIndented, "\n") + 1

	return
}

func XMLtoJson(str string) (retornoXML RetornoXML, errXML error) {
	str = reaplaceAllRegex(str, `<\?xml.*\?>`, "")
	str = reaplaceAllRegex(str, `"`, "'")
	str = reaplaceAllRegex(str, `\\'`, "")

	js := ""
	lerValor := false //false = ler nome, true = ler valor
	tipo := ""

	// Criar um decoder XML
	decoder := xml.NewDecoder(strings.NewReader(str))

	// Percorrer o arquivo XML
	for {
		// Obter o próximo token
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			errXML = err
			return
		}

		// Verificar o tipo do token
		switch token := token.(type) {
		case xml.StartElement:
			// Elemento de início
			name := token.Name.Local

			//regras
			switch name {
			case "name":
				lerValor = false
			case "struct":
				js += "{"
			case "array":
				js += "["
			case "value":
				lerValor = true
			case "int", "string", "boolean", "double":
				tipo = name
			}

		case xml.EndElement:
			// Elemento de fim
			name := token.Name.Local

			//regras
			switch name {
			case "value":
				if lerValor && tipo != "" {
					js += setarValor("", tipo)
					lerValor = false
					tipo = ""
				}
			case "struct":
				js += "},"
			case "array":
				js += "],"
			}

		case xml.CharData:
			// Conteúdo do elemento
			//regras
			if lerValor {
				js += setarValor(string(token), tipo)
				lerValor = false
				tipo = ""
			} else {
				js += `"` + string(token) + `": `
			}
		}
	}

	if js == "" {
		return
	}

	js = reaplaceAllRegex(js, `\,$`, "")
	js = reaplaceAllRegex(js, `,}`, "}")
	js = reaplaceAllRegex(js, `,]`, "]")
	js = reaplaceAllRegex(js, `{},`, "")
	js = reaplaceAllRegex(js, `^\[],`, "")
	js = reaplaceAllRegex(js, `\t`, "")

	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(js), "", "   ")

	if err != nil {
		errXML = err
		return
	}

	retornoXML.Data = prettyJSON.String()
	retornoXML.Lines = strings.Count(retornoXML.Data, "\n") + 1
	return
}

func setarValor(str, tipo string) (valor string) {
	switch tipo {
	case "string":
		valor = `"` + str + `",`
	case "int", "double":
		valor = str + ","
	case "boolean":
		if str == "1" {
			valor = "true,"
		} else {
			valor = "false,"
		}
	}

	return
}
