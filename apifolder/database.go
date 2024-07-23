package apifolder

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// Structs
type interactionData struct {
	User    string `json:"user"`
	Comment string `json:"comment"`
}
type DBSvc struct {
}

func NewDbSvc() *DBSvc {
	return &DBSvc{}
}

// Main Functions
func (dbs *DBSvc) AddNewLike(title string, user string, w http.ResponseWriter) {
	stmt, err := db.Prepare("INSERT INTO likes (title, user) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (dbs *DBSvc) AddNewComment(title string, user string, email string, comment string, w http.ResponseWriter) {
	stmt, err := db.Prepare("INSERT INTO comments (title, user, email, comment) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, user, email, comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (dbs *DBSvc) GetAllInteractions(title string) (int, []interactionData, error) {
	// Query for like count.
	likeCountQuery := `
		SELECT COUNT(*) 
		FROM likes 
		WHERE title = ?
	`
	var likeCount int
	err := db.QueryRow(likeCountQuery, title).Scan(&likeCount)
	if err != nil {
		return 0, nil, fmt.Errorf("error querying like count: %w", err)
	}

	// Query for comments.
	commentsQuery := `
		SELECT user, comment 
		FROM comments 
		WHERE title = ?
	`
	rows, err := db.Query(commentsQuery, title)
	if err != nil {
		return 0, nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	var comments []interactionData
	for rows.Next() {
		var comment interactionData
		if err := rows.Scan(&comment.User, &comment.Comment); err != nil {
			return 0, nil, fmt.Errorf("error scanning comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, fmt.Errorf("error iterating over comments: %w", err)
	}

	return likeCount, comments, nil
}

func InitDB() {
	var err error
	db, err = sql.Open("sqlite", "./comments.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create tables if they don't exist
	createTables := `
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		user TEXT
	);
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		user TEXT,
		email TEXT,
		comment TEXT
	);
	`

	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	log.Println("Database initialized successfully")
}
