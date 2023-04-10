package main

import (
        "log"
        "net/http"
        "search/resources"
)

func main() {
        //main function that listens on the correct port
	log.Fatal(http.ListenAndServe(":3001", resources.Router()))
}
