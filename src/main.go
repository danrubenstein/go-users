package main

import (
	"fmt"
	"log" 
	"net/http"
)


func main() { 

   	fmt.Println("Initializing Program")
    router := NewRouter()

    log.Fatal(http.ListenAndServe(":8080", router))
    
    fmt.Println("Exiting Program")
}