package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type URLShortener struct {
	// shortened keys as keys and original URLs as values
	urls map[string]string
}

// Implement URL Shortening
func (urlshortener *URLShortener) HandleShorten(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := request.FormValue("url")
	// url := request.URL.Query().Get("url")  <--- Another way to get the url but from the query string
	if originalURL == "" {
		http.Error(responseWriter, "URL parameter is required", http.StatusBadRequest)
		return
	}

	// Generate a random and unique shortened key for de original URL
	shortenedKey := generateShortKey()
	// shortenedkey := fmt.Sprintf("%d", rand.Int()) <--- Another way to generate a random key
	urlshortener.urls[shortenedKey] = originalURL

	// Constuct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortenedKey)

	// Render HTML response with the shortened URL
	responseWriter.Header().Set("Content-Type", "text/html")
	responseWriter.WriteHeader(http.StatusCreated)
	responseHTML := fmt.Sprintf(`
	<html>
		<body>
			<h2>David's URL Shortener</h2>
			<p>Original URL: %s</p>
			<a href="%s">%s</a>
		</body>
	</html>
	`,originalURL, shortenedURL, shortenedURL)

	fmt.Fprintf(responseWriter, responseHTML)

	/*
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	responseWriter.Write([]byte(shortenedKey))
	*/
}

// Implement URL Redirection
func (urlshortener *URLShortener) HandleRedirect(responseWriter http.ResponseWriter, request *http.Request) {
	shortenedKey := request.URL.Path[len("/short/"):]
	if shortenedKey == "" {
		http.Error(responseWriter, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the 'urls' map using the shortened key
	originalURL, ok := urlshortener.urls[shortenedKey]
	if !ok {
		http.Error(responseWriter, "Shortened key not found", http.StatusBadRequest)
		return
	}

	// Redirect to the original URL
	http.Redirect(responseWriter, request, originalURL, http.StatusMovedPermanently)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keylength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keylength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func index(responseWriter http.ResponseWriter, request *http.Request) {
	responseHTML := fmt.Sprintf(`
	<html>
		<body>
			<h2>Welcome to David's URL Shortener</h2>
			<form method="post" action="/shorten">
        <input type="text" name="url" placeholder="Enter a URL">
        <input type="submit" value="Shorten">
      </form>
		</body>
	</html>
	`)
	fmt.Fprintf(responseWriter, responseHTML)
}

func main() {
	urlshortener := &URLShortener{
		urls: make(map[string]string),
	}
	
	http.HandleFunc("/shorten", urlshortener.HandleShorten)
	http.HandleFunc("/short/", urlshortener.HandleRedirect)
	http.HandleFunc("/", index)

	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}