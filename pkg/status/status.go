package status

import (
	"encoding/json"
	"net/http"
)

type health struct {
	Status      string
	Status_code string
}

func Status(w http.ResponseWriter, r *http.Request) {
	profile := health{"healthy", "200"}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
