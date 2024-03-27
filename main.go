package main

import (
	"encoding/json"
	"eth/parser"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func getCurrentBlockHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		block := p.GetCurrentBlock()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Last processed block: %d", block)))
	}
}

func subscribeHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Address is required"))
			return
		}
		if p.Subscribe(address) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Address %s has been subscribed", address)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Address %s has already been subscribed", address)))
		}
	}
}

func getTransactiosHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Address is required"))
			return
		}
		transactions := p.GetTransactions(address)

		if len(transactions) > 0 {
			jsonData, err := json.Marshal(transactions)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error marshaling JSON: %v", err)))
				return
			}

			totalTransactiosn := `{"total": ` + fmt.Sprintf("%d", len(transactions)) + `}`
			jsonData = []byte(strings.Replace(string(jsonData), "]", ","+totalTransactiosn+"]", 1))

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}
}

func main() {

	parser := runParser()
	fmt.Println("Server running on port 8080")
	http.HandleFunc("/block", getCurrentBlockHandler(parser))
	http.HandleFunc("/subscribe", subscribeHandler(parser))
	http.HandleFunc("/transactions", getTransactiosHandler(parser))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func runParser() *parser.EthereumParser {
	parser := parser.NewParser("https://cloudflare-eth.com")
	return parser
}
