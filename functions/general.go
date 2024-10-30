package functions

import (
	"encoding/json"
	"fmt"
	"log"
)

func Includes(strs []string, args ...string) bool {
	for _, str := range strs {
		for _, arg := range args {
			if str == arg {
				return true
			}
		}
	}
	return false
}

func ToStringPretty(data interface{}) {
	formatted, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println("Erro ao formatar o struct:", err)
		return
	}
	fmt.Println(string(formatted))
}

func StructuredClone[T any](original T) (T, error) {
	data, err := json.Marshal(original)
	if err != nil {
		return *new(T), err
	}
	var clone T
	err = json.Unmarshal(data, &clone)
	if err != nil {
		return *new(T), err
	}
	return clone, nil
}
