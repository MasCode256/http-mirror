package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Установлено соединение!")

	// Получаем параметр "get" из запроса
	target := r.URL.Query().Get("get")
	if target == "" {
		http.Error(w, "Missing 'get' parameter", http.StatusBadRequest)
		return
	}

	var resp *http.Response
	var err error

	// Обработка POST-запросов
	if r.Method == http.MethodPost {
		// Читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Создаем новый POST-запрос
		req, err := http.NewRequest(http.MethodPost, target, bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
			return
		}

		// Копируем заголовки из исходного запроса
		for key, value := range r.Header {
			req.Header[key] = value
		}

		// Выполняем запрос
		resp, err = http.DefaultClient.Do(req)
	} else {
		// Выполняем GET-запрос к указанному URL
		resp, err = http.Get(target)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %s", resp.Status), resp.StatusCode)
		return
	}

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Возвращаем содержимое клиенту
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
