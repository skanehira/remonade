package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	path := filepath.Join("testdata", "all.json")
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	m := map[string]interface{}{}

	if err := json.NewDecoder(f).Decode(&m); err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		headers := map[string]string{
			"Content-Type":           "application/json",
			"X-Rate-Limit-Limit":     "30",
			"X-Rate-Limit-Remaining": "30",
			"X-Rate-Limit-Reset":     fmt.Sprintf("%d", time.Now().Unix()),
		}
		for k, v := range headers {
			w.Header().Add(k, v)
		}

		_ = r.ParseForm()

		var respBody interface{}
		switch r.URL.Path {
		case "/users/me":
			respBody = m["Users"]
		case "/appliances":
			respBody = m["Appliances"]
		case "/devices":
			respBody = m["Devices"]
		}

		if err := json.NewEncoder(w).Encode(respBody); err != nil {
			return
		}
	})
	log.Println("start http server :9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
