package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	requestURL := "https://rgt-tools.duckdns.org/drive/api/public/dl/9ziyBW4A"
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	if res.StatusCode != 200 {
		fmt.Printf("client: error: status code is not 200\n")
		os.Exit(1)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: error reading response body: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: response body: %s\n", body)
}
