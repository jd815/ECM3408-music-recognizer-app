package resources

import (
	"encoding/json"
	"net/http"
    "bytes"
	"github.com/gorilla/mux"
	"io"
)

const(
	KEY = "33cd85f28fc866972f96bc13e89cf03e"
)

type Result struct{
	Fir string `json:"status"`
	ResultS Id `json:"result"`
}
type Id struct{
	Sec string `json:"artist"`
	IdS string `json:"title"`


}
func searchAPI(w http.ResponseWriter, r *http.Request) {
	body := map[string]interface{} {}
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        w.WriteHeader(500)
	}
	audio, ok := body["Audio"]
	if !ok{
		w.WriteHeader(400)
	}
	//There are two ways to handle many error checks, one is nested loops as done here
	apiRequest := map[string]interface{} {"api_token": KEY, "audio" : audio}
	if jsonBody, err := json.Marshal(apiRequest); err == nil{
		send := bytes.NewBuffer(jsonBody)
		if apiResult, err := http.Post("https://api.audd.io/recognize", "application/json", send); err == nil{
			if apiResult.StatusCode == 200{
				defer apiResult.Body.Close()
				if apiResultBody, err := io.ReadAll(apiResult.Body); err == nil{
					result := Result{}
					err = json.Unmarshal(apiResultBody, &result)
					if err == nil{
						if result.ResultS.IdS == ""{
							w.WriteHeader(404)
						}else{
							response := map[string]interface{}{"Id": result.ResultS.IdS}
							w.WriteHeader(200)
							json.NewEncoder(w).Encode(response)
						}
					}else{
						w.WriteHeader(500)
					}
				}else{
					w.WriteHeader(500)
				}
			}else{
				w.WriteHeader(500)
			}
		}else{
			w.WriteHeader(500)
		}
	}else{
		w.WriteHeader(500)
	}
	
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* controller */
	//search
	r.HandleFunc("/search", searchAPI).Methods("POST")
	return r
}

