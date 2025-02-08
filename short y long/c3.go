package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var lastVersion int = 0

func checkUpdates() {
	for {
		resp, err := http.Get("http://localhost:4000/wait_changes/" + strconv.Itoa(lastVersion))
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			var result map[string]int
			json.Unmarshal(body, &result)

			if newVersion, exists := result["new_version"]; exists {
				fmt.Println("Servidor A: Detectado cambio en Servidor B, nueva versión:", newVersion)
				lastVersion = newVersion
			}
			resp.Body.Close()
		}

		time.Sleep(2 * time.Second) // Short polling cada 2 segundos
	}
}

func main() {
	go checkUpdates()

	http.ListenAndServe(":3000", nil)
}