package main

import (
	"fmt"
	Controller "github.com/Adrenesis/crypto-curry-print/controller"
	Model "github.com/Adrenesis/crypto-curry-print/model"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	//"strconv"
)

func main() {
	var port int64
	var errport error
	path, errwd := os.Getwd()
	if errwd != nil {
		log.Println(errwd)
	}
	fmt.Println(path)

	if len(os.Args) > 1 {
		port, errport = strconv.ParseInt(os.Args[1], 10, 64)
	} else {
		errport = strconv.ErrSyntax
	}
	//mode := Model
	//Controller.InitDB()
	Controller.InitDB()

	http.HandleFunc("/index", Controller.HandleIndex)
	http.HandleFunc("/links", Controller.HandleLinks)
	http.HandleFunc("/bscbalances", Controller.HandleBSCBalance)
	http.HandleFunc("/update/coin", Controller.HandleCoinUpdate)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	if errport != nil {
		port = 8880
	}
	var address = fmt.Sprintf("127.0.0.1:%d", port)
	fmt.Println(Model.ConvertToISO8601(time.Now()), "Binding http://"+address+"/index...")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(Model.ConvertToISO8601(time.Now()), "Successfully hosting on https://", address)
}
