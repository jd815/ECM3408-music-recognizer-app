package main

import (
        "log"
        "net/http"
        "tracks/resources"
        "tracks/repository"
)

func main() {
        //functon calls to initialise, create, and clear the database. Listening on the correct port
        repository.Init()
        repository.Create()
        repository.Clear()
        log.Fatal(http.ListenAndServe(":3000", resources.Router()))
}
