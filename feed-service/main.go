package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Feed service is running on :8083")
	http.ListenAndServe(":8083", nil)
}
