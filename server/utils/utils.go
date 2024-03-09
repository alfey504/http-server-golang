package utils

import "unicode"

func SplitFirst(s string, del string) []string {
	words := []string{}
	word := ""
	split := false
	for _, b := range s {
		if !split && string(b) == del {
			words = append(words, word)
			word = ""
			split = !split
			continue
		}
		word += string(b)
	}
	words = append(words, word)
	return words
}

func Split(request string, del string) []string {
	lines := []string{}
	line := ""
	for _, b := range request {
		if string(b) == del {
			lines = append(lines, line)
			line = ""
			continue
		}
		line += string(b)
	}
	lines = append(lines, line)
	return lines
}

func StripStart(s string) string {
	line := ""
	for pos, val := range s {
		if unicode.IsLetter(val) {
			return s[pos:]
		}
	}
	return line
}

func StripWhiteSpace(s string) string {
	frontStrip := ""
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			frontStrip = s[i:]
			break
		}
	}
	strippedString := frontStrip
	for j := len(frontStrip) - 1; j >= 0; j-- {
		if frontStrip[j] != ' ' {
			strippedString = frontStrip[:j+1]
			break
		}
	}
	return strippedString
}
