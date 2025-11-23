package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Post service is running on :8081")
	http.ListenAndServe(":8081", nil)
}
