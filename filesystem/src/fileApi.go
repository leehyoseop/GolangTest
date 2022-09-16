package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func PrintFileList() { // C:/tmp 하위의 파일 리스트 터미널에 출력
	fmt.Println("현재 파일 리스트")
	targerDir := "C:/tmp"
	files, err := ioutil.ReadDir(targerDir)
	check(err)
	var myFileList []myFile
	for i, file := range files {
		idNum := i + 1
		myfile := &myFile{Id: idNum, Name: file.Name()}
		myFileList = append(myFileList, *myfile)
		idNum++
	}
	fmt.Println("=>", myFileList)
}

func MappingFile() map[int]string {
	targerDir := "C:/tmp"
	files, err := ioutil.ReadDir(targerDir)
	check(err)
	var idMap map[int]string
	idMap = make(map[int]string)
	for i, file := range files {
		idNum := i + 1
		idMap[idNum] = file.Name()
		idNum++
	}
	return idMap
}

func MappingFileStruct() []myFile {
	targerDir := "C:/tmp"
	files, err := ioutil.ReadDir(targerDir)
	check(err)
	var myFileList []myFile
	for i, file := range files {
		idNum := i + 1
		myfile := &myFile{Id: idNum, Name: file.Name()}
		myFileList = append(myFileList, *myfile)
		//fmt.Println(myFileList)
		idNum++
	}
	return myFileList
}

func GetFileList(wr http.ResponseWriter, r *http.Request) {
	var myFileList = MappingFileStruct()
	//j, _ := json.Marshal(myFileList)
	//json.NewEncoder(wr).Encode(string(j))
	wr.Header().Add("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(myFileList)
}

func GetFile(wr http.ResponseWriter, r *http.Request) {
	var idMap = MappingFile()

	vars := mux.Vars(r)
	parseid := vars["id"]
	intparseid, _ := strconv.Atoi(parseid)
	_, exists := idMap[intparseid]

	var myfile myFile

	if !exists {
		fmt.Println("No Matching File")
	} else {
		readFileName := idMap[intparseid]
		dat, err := ioutil.ReadFile("C:/tmp/" + readFileName)
		check(err)
		myfile.Name = readFileName
		myfile.Content = string(dat)
		wr.Header().Add("Content-Type", "application/json")
		json.NewEncoder(wr).Encode(myfile)
	}
}

func PostFile(wr http.ResponseWriter, r *http.Request) {
	var myFileList = MappingFileStruct()

	reqBody, _ := ioutil.ReadAll(r.Body)
	var myfile myFile
	json.Unmarshal(reqBody, &myfile)

	f, err := os.Create("C:/tmp/" + myfile.Name + ".txt")
	check(err)
	f.Sync()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(myfile.Content)
	w.Flush()

	myFileList = append(myFileList, myfile)
	wr.Header().Add("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(myfile)
}

func PutFile(wr http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var myfile myFile
	json.Unmarshal(reqBody, &myfile)

	var idMap = MappingFile()

	vars := mux.Vars(r)
	parseid := vars["id"]
	intparseid, _ := strconv.Atoi(parseid)
	_, exists := idMap[intparseid]

	if !exists {
		fmt.Println("No Matching File")
	} else {
		putFileName := idMap[intparseid]
		// Append
		// f, err := os.OpenFile("C:/tmp/"+putFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		//    fmt.Println(err)
		// }
		// defer f.Close()
		// if _, err := f.WriteString(myfile.Content); err != nil {
		//    fmt.Println(err)
		// }
		os.Remove("C:/tmp/" + putFileName)
		fmt.Println(putFileName + " File Deleted")

		f, err := os.Create("C:/tmp/" + putFileName)
		check(err)
		f.Sync()
		w := bufio.NewWriter(f)
		_, err = w.WriteString(myfile.Content)
		w.Flush()
		myfile.Id = intparseid
		fmt.Println(putFileName + " File Created")

		wr.Header().Add("Content-Type", "application/json")
		json.NewEncoder(wr).Encode(myfile)
	}
}

func DeleteFile(wr http.ResponseWriter, r *http.Request) {
	var idMap = MappingFile()

	vars := mux.Vars(r)
	parseid := vars["id"]
	intparseid, _ := strconv.Atoi(parseid)
	_, exists := idMap[intparseid]

	if !exists {
		fmt.Println("No Matching File")
	} else {
		delFileName := idMap[intparseid]
		os.Remove("C:/tmp/" + delFileName)
		fmt.Println(delFileName + " File Deleted")
	}
}
