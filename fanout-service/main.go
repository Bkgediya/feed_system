package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Fanout service is running on :8082")
	http.ListenAndServe(":8082", nil)
}
