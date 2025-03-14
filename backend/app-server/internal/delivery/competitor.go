package delivery

import "net/http"

func CreateShot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shot created"))
}
