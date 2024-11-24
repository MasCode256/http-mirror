package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Установлено соединение!")

	target := r.URL.Query().Get("get")
	if target == "" {
		http.Error(w, "Missing 'get' parameter", http.StatusBadRequest)
		return
	}

	var resp *http.Response
	var err error

	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		req, err := http.NewRequest(http.MethodPost, target, bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
			return
		}

		for key, value := range r.Header {
			req.Header[key] = value
		}

		resp, err = http.DefaultClient.Do(req)
	} else {
		resp, err = http.Get(target)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %s", resp.Status), resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response body: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Write(body)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting mirror on :6060")
	if err := http.ListenAndServe(":6060", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
