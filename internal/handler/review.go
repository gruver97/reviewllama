package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ollama/ollama/api"
	"log"
	"net/http"
	"strings"
)

// ReviewPost
/*
curl --location 'http://localhost:8080/review' \
--header 'Content-Type: application/json' \
--header 'Cookie: x-supplier-id=fb25c9e9-cae8-52db-b68e-736c1466a3f5' \
--data '{
    "text": "cGFja2FnZSBtYXRodXRpbHMKCi8vIEFkZCByZXR1cm5zIHRoZSBzdW0gb2YgdHdvIGludGVnZXJzLgpmdW5jIEFkZChhLCBiIGludCkgaW50IHsKICAgIHJldHVybiBhICsgYgp9Cg=="
}'
*/
func ReviewPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var request ReviewRequest

	err := decoder.Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sr := strings.NewReader(request.Text)
	b := base64.NewDecoder(base64.StdEncoding, sr)
	var text []byte
	_, err = b.Read(text)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	req := &api.GenerateRequest{
		Model:  "qwen2.5-coder:32b",
		Prompt: fmt.Sprintf("generate unit test without explanation\n```go\n%s```", text),

		// set streaming to false
		Stream: new(bool),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		_, err = w.Write([]byte(resp.Response))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ReviewRequest struct {
	Text string `json:"text"`
}
