package apifolder

import (
	"encoding/json"
	"net/http"
)

func AddLike(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var like struct {
		Title string `json:"title"`
		User  string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbSvc := NewDbSvc()
	dbSvc.AddNewLike(like.Title, like.User, w)

	w.WriteHeader(http.StatusCreated)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var comment struct {
		Title   string `json:"title"`
		User    string `json:"user"`
		Email   string `json:"email"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbSvc := NewDbSvc()
	dbSvc.AddNewComment(comment.Title, comment.User, comment.Email, comment.Comment, w)

	w.WriteHeader(http.StatusCreated)
}

func GetInteractions(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	title := r.URL.Query().Get("title")

	dbSvc := NewDbSvc()
	likeCount, comments, err := dbSvc.GetAllInteractions(title)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare and send the response.
	data := struct {
		LikeCount int               `json:"like_count"`
		Comments  []interactionData `json:"comments"`
	}{
		LikeCount: likeCount,
		Comments:  comments,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
