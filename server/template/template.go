package templ

import (
	"fmt"
	"os"

	"tcp_http_server.com/server/utils"
)

func LoadTemplate(dir string, data map[string]string) ([]byte, error) {
	content, err := os.ReadFile(dir)
	if err != nil {
		return []byte{}, err
	}

	parsedHtml, err := parseTemplate(content, data)
	if err != nil {
		return []byte{}, err
	}
	return parsedHtml, nil
}

func parseTemplate(template []byte, data map[string]string) ([]byte, error) {
	parsedHtml := []byte{}
	for pos := 0; pos < len(template); pos++ {
		if template[pos] == '{' && template[pos+1] == '{' {
			val, err, j := parseData(template, pos, data)
			if err != nil {
				return []byte{}, err
			}
			parsedHtml = append(parsedHtml, []byte(val)...)
			pos = j + 2
		}
		parsedHtml = append(parsedHtml, template[pos])
	}
	return parsedHtml, nil
}

func parseData(template []byte, startPos int, data map[string]string) ([]byte, error, int) {
	j := startPos + 2
	param := ""

	for {
		if j >= len(template) {
			return []byte{}, fmt.Errorf("syntax error : expected }}"), -1
		}
		if template[j] == '}' && template[j+1] == '}' {
			break
		}
		param += string(template[j])
		j++
	}

	strippedParam := utils.StripWhiteSpace(param)
	val, ok := data[strippedParam]

	if !ok {
		return []byte{}, fmt.Errorf("param " + strippedParam + " not found in the data"), -1
	}
	return []byte(val), nil, j
}
