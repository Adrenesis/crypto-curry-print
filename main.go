package main

import (
	"fmt"
	Controller "github.com/Adrenesis/crypto-curry-print/controller"
	"log"
	"net/http"
	"os"
	"strconv"
	//"strconv"
)

func main() {
	var port int64
	var errport error
	if len(os.Args) > 1 {
		port, errport = strconv.ParseInt(os.Args[1], 10, 64)
	} else {
		errport = strconv.ErrSyntax
	}

	http.HandleFunc("/index", Controller.HandleIndex)
	http.HandleFunc("/links", Controller.HandleLinks)
	if errport != nil {
		port = 8880
	}
	var address = fmt.Sprintf("127.0.0.1:%d", port)
	fmt.Println("Binding http://" + address + "/index...")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfully hosting on https://", address)
}
