package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run client.go <url>")
		os.Exit(1)
	}

	originalURL := os.Args[1]
	shortenedURL, err := getShortenedURL(originalURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Shortened URL: %s\n", shortenedURL)
}

func getShortenedURL(originalURL string) (string, error) {
	url := "http://localhost:3333/user-input" 

	payload := bytes.NewBufferString(fmt.Sprintf("url=%s", originalURL))
	resp, err := http.Post(url, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned non-OK status: %s", resp.Status)
	}

	var result string
	_, err = fmt.Fscanf(resp.Body, "Ваша укороченная ссылка: %s", &result)
	if err != nil {
		return "", err
	}

	return result, nil
}
