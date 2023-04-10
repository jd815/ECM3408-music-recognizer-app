package resources

import (
	"encoding/json"
	"tracks/repository"
	"github.com/gorilla/mux"
	"net/http"
)

func addTrack (w http.ResponseWriter, r *http.Request){
	//takes in an http request and adds track. If any error arises proper status code is returned
	data := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil{
		w.WriteHeader(500)/* Internal Server Error */
		return
	}

	id, ok := data["Id"]
	if !ok {
		w.WriteHeader(400)/* Bad request */
        return
	}
	audio, ok := data["Audio"]
	if !ok {
		w.WriteHeader(400)
        return
	}

	if audio == ""{
		w.WriteHeader(400)
        return
	}

	track := repository.Track{Id: id.(string), Audio: audio.(string)}
	n := repository.AddTrack(track) 
	if n > 0{
		w.WriteHeader(201)/* Created */
	}else if n == 0{
		w.WriteHeader(204)/* No content */
	}else{
		w.WriteHeader(500)
	}
}

func getAllTracks(w http.ResponseWriter, r *http.Request){
	//function that returns all the tracks to the user
	tracks, numTracks := repository.GetAllTracks()
	ids := make([]string, 0)

	if numTracks == -1{
		w.WriteHeader(500)
	}else if numTracks == 0{
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ids)
	}else{
		w.WriteHeader(200)
		for _, i := range tracks {
			ids = append(ids, i.Id)
			json.NewEncoder(w).Encode(ids)
		}
	}
}
func getTrack(w http.ResponseWriter, r *http.Request){
	//function that takes in a name of a song and returns the name and id of song
	data := mux.Vars(r)
	id, ok := data["id"]
    if !ok {
        w.WriteHeader(404)/* Not Found */
    }
	if c, n := repository.GetTrack(id); n > 0 {
		w.WriteHeader(200) /* OK */
		text := map[string]interface{}{"Id": c.Id, "Audio": c.Audio}
		json.NewEncoder(w).Encode(text)
	} else if n == 0 {
		w.WriteHeader(404) 
	} else {
		w.WriteHeader(500) 
	}
}
func deleteTrack(w http.ResponseWriter, r *http.Request){
	//function that deletes a track from the database
	data := mux.Vars(r)

    id, ok := data["id"]
    if !ok {
        w.WriteHeader(404)
    }

    response := repository.DeleteTrack(id)
    
    if response > 0 {
        w.WriteHeader(204)
    } else if response == 0 {
        w.WriteHeader(404)
    } else {
        w.WriteHeader(500)
    }
}

func Router() http.Handler {
	//router function that handles the request and calls the appropriate function
	r := mux.NewRouter()
	/* controller */

	//add track to DB
	r.HandleFunc("/tracks/{id}", addTrack).Methods("PUT")

	//show all tracks
	r.HandleFunc("/tracks", getAllTracks).Methods("GET")

	//get Track by ID
	r.HandleFunc("/tracks/{id}", getTrack).Methods("GET")

	//delete track
	r.HandleFunc("/tracks/{id}", deleteTrack).Methods("DELETE")
	return r
}
