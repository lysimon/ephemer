package whisker

// Get basic configuration to show to the user
import (
	"net/http"

	"../../pkg/config"
)

func Configuration() {
	http.HandleFunc("/configuration", configuration_endpoint)
}

// json structure returning the current configuration
func configuration_endpoint(w http.ResponseWriter, r *http.Request) {
	js := config.GlobalConfig.To_json()
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
