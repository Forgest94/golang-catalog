package product

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetList(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	jsonResp, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
