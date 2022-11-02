package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

var keyMap = make(map[string]string)

func setKey(w http.ResponseWriter, r *http.Request) {
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	key := vars["key"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	keyMap[key] = string(reqBody)
	fmt.Println("SET " + key + ": " + keyMap[key])
	save()
}

func getKey(w http.ResponseWriter, r *http.Request) {
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	key := vars["key"]
	value, ok := keyMap[key]
	if !ok {
		fmt.Println("GET " + key + ": NULL")
		fmt.Fprintf(w, "")
	} else {
		fmt.Println("GET " + key + ": " + value)
		fmt.Fprintf(w, value)
	}
}

func loadKeyMap() {
	contentBytes, err := ioutil.ReadFile("data.txt")
	if err == nil {
		b := new(bytes.Buffer)
		b.Write(contentBytes)
		d := gob.NewDecoder(b)
		err = d.Decode(&keyMap)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Data Loaded: ", len(keyMap), " Keys")
			i := 1
			for k, v := range keyMap {
				fmt.Printf("\n%d. [%s] = [%s]", i, k, v)
				i++
			}
		}
	} else {
		fmt.Println("Data file not found => Reset All Keys")
	}
}

func save() {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(keyMap)
	if err != nil {
		panic(err)
	} else {
		ioutil.WriteFile("data.txt", b.Bytes(), 0644)
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func handleRequests() {
	loadKeyMap()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/string/{key}", setKey).Methods("PUT", "POST")
	router.HandleFunc("/string/{key}", getKey).Methods("GET")
	ip := "localhost" //GetOutboundIP().String()
	log.Fatal(http.ListenAndServe(ip+":24500", router))
}

func main() {
	handleRequests()
}
