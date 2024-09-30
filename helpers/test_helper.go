package helpers

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(data ...any) {
	for _, d := range data {
		b, _ := json.MarshalIndent(d, "", "  ")
		fmt.Println(string(b))
	}
}
