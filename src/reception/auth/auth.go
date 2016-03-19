package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Setup loads the auth json and validates it
func Setup(filepath string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	var aj authJSON
	err = json.Unmarshal(file, &aj)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(aj, "", "  ")
	fmt.Println(string(b))

}
