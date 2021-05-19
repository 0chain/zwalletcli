package util

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJSON(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Failed to convert data to json format : %v", err)
	}
	fmt.Println(string(b))
}

func PrettyPrintJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert data to json format : %v", err)
	}
	fmt.Println(string(b))
}
