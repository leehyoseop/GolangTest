package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type myFile struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func handleRequests() {
	PrintFileList()

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/files", func(wr http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: // 목록 전체 조회
			GetFileList(wr, r)
		case http.MethodPost: // 등록
			PostFile(wr, r)
		}
		PrintFileList()
	})
	myRouter.HandleFunc("/files/{id}", func(wr http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: // 개별 조회
			GetFile(wr, r)
		case http.MethodPut: // 수정
			PutFile(wr, r)
		case http.MethodDelete: // 삭제
			DeleteFile(wr, r)
		}
		PrintFileList()
	})
	http.ListenAndServe(":8080", myRouter)
}

func main() {
	handleRequests()
}
