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

	stmt, err := db.Prepare("INSERT INTO likes (title, user) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(like.Title, like.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	stmt, err := db.Prepare("INSERT INTO comments (title, user, email, comment) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.Title, comment.User, comment.Email, comment.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetInteractions(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	title := r.URL.Query().Get("title")

	// Get like count
	likeCountQuery := `
		SELECT COUNT(*) 
		FROM likes 
		WHERE title = ?
	`
	var likeCount int
	err := db.QueryRow(likeCountQuery, title).Scan(&likeCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get comments
	commentsQuery := `
		SELECT user, comment 
		FROM comments 
		WHERE title = ?
	`
	rows, err := db.Query(commentsQuery, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type interactionData struct {
		User    string `json:"user"`
		Comment string `json:"comment"`
	}

	data := struct {
		LikeCount int               `json:"like_count"`
		Comments  []interactionData `json:"comments"`
	}{
		LikeCount: likeCount,
	}

	for rows.Next() {
		var comment interactionData
		err := rows.Scan(&comment.User, &comment.Comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.Comments = append(data.Comments, comment)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}
