package apifolder

import (
	"log"
	"net/http"
)

func StartServer() {
	InitDB()

	log.Println("Starting API server on :8080")

	http.HandleFunc("/api/interactions", GetInteractions)
	http.HandleFunc("/api/interactions/add-like", AddLike)
	http.HandleFunc("/api/interactions/add-comment", AddComment)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("API server failed to start: %v", err)
	}
}
