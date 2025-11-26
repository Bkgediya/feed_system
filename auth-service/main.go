package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Auth service is running on :8081")
	http.ListenAndServe(":8081", nil)
}
