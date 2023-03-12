package chat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type QueryResponse map[string]string

// SendQuery sends a query to the OpenAI ChatGPT API and returns the query
// together with the response.
func SendQuery(token, query string) (QueryResponse, error) {
	qr := QueryResponse{"query": query}

	bs, err := json.Marshal(qr)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/engines/chatbot/query", bytes.NewBuffer(bs))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	qr["response"] = string(body)

	return qr, nil
}

// EOF
