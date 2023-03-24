// Ask - Send queries to ChatGPT.
//
// Copyright (C) 2023 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

// chat sends a query to the OpenAI ChatGPT API and returns the query
// together with the response.
func chat(apiKey, model, query string) (map[string]any, error) {
	url := "https://api.openai.com/v1/chat/completions"
	qr := map[string]any{
		"model": model,
		"messages": []map[string]any{
			{"role": "user", "content": query},
		},
	}

	bs, err := json.Marshal(qr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response.
	response := map[string]any{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	qr["response"] = response

	return qr, nil
}

// main is the entry point for the program.
func main() {
	// Parse the command line flags.
	query := flag.String("query", "", "The query to ask.")
	model := flag.String("model", "gpt-3.5-turbo", "The model to use.")

	flag.Parse()

	// Read the OpenAI API key from the environment.
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		return
	}

	fmt.Println("ask v0.1.0")
	fmt.Println("Please wait...")

	// Send the query to the OpenAI API.
	qr, err := chat(apiKey, *model, *query)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Query...: %v\n", qr["query"])
	fmt.Printf("Model...: %v\n", qr["model"])
	fmt.Printf("Response: %v\n", qr["response"])
}

// EOF
