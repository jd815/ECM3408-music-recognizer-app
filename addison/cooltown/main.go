package main

import (
	"net/http"
    "cooltown/resources"
)

func main() {
    //main function that listens on the correct port
    http.ListenAndServe(":3002", resources.Router())

}
