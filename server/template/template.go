package templ

import (
	"fmt"
	"os"

	"tcp_http_server.com/server/utils"
)

func LoadTemplate(dir string, data map[string]interface{}) ([]byte, error) {
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

func parseTemplate(template []byte, data map[string]interface{}) ([]byte, error) {
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

func parseData(template []byte, startPos int, data map[string]interface{}) ([]byte, error, int) {
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

	stringVal, ok := val.(string)
	if !ok {
		fmt.Errorf("type of of ", strippedParam, " is not string")
	}
	return []byte(stringVal), nil, j
}

func LoadTemplateV2(templateDir string, data map[string]interface{}) ([]byte, error) {
	content, err := os.ReadFile(templateDir)
	if err != nil {
		println(err.Error())
	}
	componentMap := MapComponents(content)
	component := RenderComponentMap(componentMap, "main")

	parsedHtml, err := parseTemplate(component, data)
	if err != nil {
		return []byte{}, err
	}

	mappedHtml, err := renderMap(parsedHtml, data)
	if err != nil {
		println(err.Error())
	}

	return mappedHtml, nil
}

func MapComponents(content []byte) map[string][]byte {
	lines := utils.Split(string(content), "\n")

	component := []byte{}
	currentComponent := ""

	componentMap := make(map[string][]byte)

	for _, l := range lines {
		isInstruction, instruction, err := ParseComponent(l)
		if err != nil {
			println(err.Error())
		}
		if isInstruction {
			operation, operand := ParseInstruction(instruction)
			switch operation {
			case "Component":
				component = []byte{}
				currentComponent = operand
			case "ComponentEnd":
				componentMap[currentComponent] = component
				component = []byte{}
				currentComponent = ""
			default:
				println("Unknown Instruction -> ", instruction)
			}
			continue
		}
		line := l + " \n"
		component = append(component, []byte(line)...)
	}
	return componentMap
}

func ParseComponent(line string) (bool, string, error) {
	for pos := 0; pos < len(line); pos++ {
		if line[pos] == '@' && line[pos+1] == '<' {
			instruction := ""
			i := pos + 2
			for line[i] != '>' {
				instruction += string(line[i])
				i++
			}
			return true, instruction, nil
		}
	}
	return false, "", nil
}

func ParseInstruction(instruction string) (string, string) {
	operation := ""
	operand := ""

	isOperand := false
	for _, ch := range instruction {
		if ch == '(' {
			isOperand = true
			continue
		}

		if ch == ')' {
			isOperand = false
			continue
		}

		if isOperand {
			operand += string(ch)
			continue
		}

		operation += string(ch)
	}
	return operation, operand
}

func RenderComponentMap(componentMap map[string][]byte, componentName string) []byte {
	rootComponent := componentMap[componentName]
	rootComponentLines := utils.Split(string(rootComponent), "\n")
	component := []byte{}
	for _, lines := range rootComponentLines {
		isInstruction, instruction, err := ParseComponentCall(lines)
		if err != nil {
			println(err.Error())
		}

		if isInstruction {
			operation, operand := ParseInstruction(instruction)
			switch operation {
			case "Component":
				println("ComponentCall -> ", operand)
				renderedComponent := RenderComponentMap(componentMap, operand)
				renderedComponent = append([]byte{'\n'}, renderedComponent...)
				component = append(component, renderedComponent...)
			default:
				println("unknown operation -> ", operation)
			}
			continue
		}
		newLine := lines + "\n"
		component = append(component, []byte(newLine)...)
	}
	return component
}

func ParseComponentCall(line string) (bool, string, error) {
	for pos := 0; pos < len(line); pos++ {
		if line[pos] == '!' && line[pos+1] == '<' {
			instruction := ""
			i := pos + 2
			for line[i] != '>' {
				instruction += string(line[i])
				i++
			}
			return true, instruction, nil
		}
	}
	return false, "", nil
}

func renderMap(template []byte, data map[string]interface{}) ([]byte, error) {
	parsedHtml := []byte{}
	for pos := 0; pos < len(template); pos++ {
		if template[pos] == '{' && template[pos+1] == '%' {
			val, j, err := parseMap(template, pos, data)
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

func parseMap(template []byte, startPos int, data map[string]interface{}) ([]byte, int, error) {
	j := startPos + 2
	param := ""

	for {
		if j >= len(template) {
			return []byte{}, -1, fmt.Errorf("syntax error : expected -> ")
		}
		if template[j] == '-' && template[j+1] == '>' {
			break
		}
		param += string(template[j])
		j++
	}
	j += 2

	param = utils.StripWhiteSpace(param)
	paramValue, ok := data[param]

	if !ok {
		return []byte{}, -1, fmt.Errorf("param " + param + " not found in the data")
	}

	paramValueSlice, ok := paramValue.([]string)
	if !ok {
		return []byte{}, -1, fmt.Errorf("param " + param + " is not a slice of type string")
	}

	variableName := ""
	for {
		if j >= len(template) {
			return []byte{}, -1, fmt.Errorf("syntax error : expected -> ")
		}
		if template[j] == ':' {
			break
		}
		variableName += string(template[j])
		j++
	}
	j += 1
	variableName = utils.StripWhiteSpace(variableName)

	htmlToMap := ""
	for {
		if j >= len(template) {
			return []byte{}, -1, fmt.Errorf("syntax error : expected -> ")
		}
		if template[j] == '%' && template[j+1] == '}' {
			break
		}
		htmlToMap += string(template[j])
		j++
	}
	j++

	println("Html -> ", htmlToMap, "\n Param -> ", param, "\n Variable ->", variableName)
	parsedHtml := []byte{}
	for _, val := range paramValueSlice {
		html, err := renderSliceValueToHtml(variableName, htmlToMap, val)
		if err != nil {
			return []byte{}, -1, err
		}
		parsedHtml = append(parsedHtml, []byte(html)...)
	}
	return []byte(parsedHtml), j, nil
}

func renderSliceValueToHtml(varName string, html string, val string) (string, error) {
	println("renderSliceToHtml() -> ")
	newHtml := ""
	i := 0

	for {

		if i >= len(html) {
			break
		}

		variableName := ""
		if html[i] == '$' && html[i+1] == '{' {
			j := i + 2
			for html[j] != '}' {
				variableName += string(html[j])
				j++
			}
			if variableName != varName {
				return "", fmt.Errorf("undefined variable " + variableName)
			}
			newHtml += val
			i = j + 1
			continue
		}
		newHtml += string(html[i])
		i++
	}
	return newHtml, nil
}
