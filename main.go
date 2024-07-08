package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// URLMapping stores the mapping between short and long URLs
var URLMapping = make(map[string]string)
var mutex = &sync.Mutex{}

// GenerateShortURL is a simple function to generate a short URL
// This is a placeholder and should be replaced with a more robust solution
func GenerateShortURL(url string) string {
	h := fmt.Sprintf("%x", md5.Sum([]byte(url)))            // 假设使用MD5哈希（需要导入"crypto/md5"）
	return base64.URLEncoding.EncodeToString([]byte(h[:8])) // 返回哈希的前8个字符作为短URL
}

// ShortenHandler handles requests to shorten URLs
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	shortURL := GenerateShortURL(longURL)
	mutex.Lock()
	URLMapping[shortURL] = longURL
	mutex.Unlock()

	fmt.Println("Short URL: {}, Long URL: {}", shortURL, longURL)

	fmt.Fprintf(w, "Your short URL is: http://yourdomain.com/%s", shortURL)
}

// RedirectHandler handles requests to redirect from short to long URLs
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:] // 去掉路径前的'/'
	mutex.Lock()
	longURL, exists := URLMapping[shortURL]
	mutex.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/shorten", ShortenHandler).Methods("POST")
	router.HandleFunc("/{shortURL}", RedirectHandler).Methods("GET")

	fmt.Println("begin to listen on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
