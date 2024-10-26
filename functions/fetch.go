package functions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func FetchFile[T any](path string) T {
	var result T
	data, err := os.ReadFile(fmt.Sprintf("%s.json", path))
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo %s.json: %s", path, err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result
	}
	return result
}
