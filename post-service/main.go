package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Post service is running on :8084")
	http.ListenAndServe(":8084", nil)
}
