package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handlePostRequest(w, r)
		return
	}

	filePath := r.URL.Path[1:]
	fullPath := filepath.Join(".", filePath)

	// Проверяем, существует ли файл
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	switch filepath.Ext(fullPath) {
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
		fmt.Println("Ошибка при определении типа файла: " + filepath.Ext(fullPath) + " (" + fullPath + ")")
	}

	http.ServeFile(w, r, fullPath)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Получен POST-запрос с данными: %s", body)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Данные успешно получены"))
}

func main() {
	http.HandleFunc("/", fileHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
