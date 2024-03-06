package templ

import "os"

func LoadTemplate(dir string) ([]byte, error) {
	return os.ReadFile(dir)
}
