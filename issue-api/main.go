package main

import (
	"issue-api/router"
	"log"
	"net/http"
)

func main() {
	r := router.NewRouter()
	log.Println("서버 실행 중... 포트: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
