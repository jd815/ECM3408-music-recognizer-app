package resources

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

func coolTown(w http.ResponseWriter, r *http.Request) {
    /*
    As mentioned in search function. One way to handle many error checks is through nested if loops. 
    Other is single nested if loops with return statements in every check as done here
    */
    data := map[string]interface{} {}
    log.Println("d")
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        w.WriteHeader(500)
        return
    }

	defer r.Body.Close()

    audio, ok := data["Audio"]
    if !ok || audio == "" {
        w.WriteHeader(400)
        return
    }

	searchBody, err := json.Marshal(map[string]interface{} {"Audio": audio,})
	if err != nil {
		w.WriteHeader(500)
		return
	}
	response, err := http.Post("http://127.0.0.1:3001/search", "application/json", bytes.NewBuffer(searchBody))
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		w.WriteHeader(response.StatusCode)
		return
	}
    
    tracksBody := map[string]interface{} {}

	if err := json.NewDecoder(response.Body).Decode(&tracksBody); err != nil {
        w.WriteHeader(500)
		return
	}

    id, ok := tracksBody["Id"]
    if !ok {
        w.WriteHeader(500)
        return
    }

    responseT, err := http.Get("http://127.0.0.1:3000/tracks/" + strings.Replace(id.(string), " ", "+", -1))
    if err != nil {
        w.WriteHeader(500)
		return
	}
    defer responseT.Body.Close()
    
	if responseT.StatusCode != 200 {
		http.Error(w, responseT.Status, responseT.StatusCode)
		return
	}

    tracksR := map[string]interface{} {}

	if err := json.NewDecoder(responseT.Body).Decode(&tracksR); err != nil {
        w.WriteHeader(500)
        return
    }

    trackAudio, ok := tracksR["Audio"]
    if !ok {
        w.WriteHeader(500)
        return
    }


    responseFinal := map[string]interface{} {"Audio" : trackAudio}
    w.WriteHeader(200)
    json.NewEncoder(w).Encode(responseFinal)

}

func Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/cooltown", coolTown).Methods(http.MethodPost)

	return r
}