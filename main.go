package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type JsonResponse struct {
	Data interface{} `json:"data"`
}

type RequestBody struct {
	CustInfo struct {
		IDNo string `json:"idNo"`
	} `json:"custInfo"`
}

var fileNames = map[string]string{
	"1": "scenario/scenario1.json",
	"2": "scenario/scenario2.json",
	// Add more file names here...
}

func main() {
	http.HandleFunc("/api/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var body RequestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fileName, ok := fileNames[body.CustInfo.IDNo]
		if !ok {
			http.Error(w, "Invalid idNo", http.StatusBadRequest)
			return
		}

		data, err := os.ReadFile(fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var jsonResponse JsonResponse
		err = json.Unmarshal(data, &jsonResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseData, err := json.Marshal(jsonResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseData)
	})

	log.Println("Listening on :9999...")
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal(err)
	}
}
