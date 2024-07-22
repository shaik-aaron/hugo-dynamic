package apifolder

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// Function Interfaces
type AddNewLikeInterface interface {
	AddNewLike(title string, user string, w http.ResponseWriter)
}

type AddNewCommentInterface interface {
	AddNewComment(title string, user string, email string, comment string, w http.ResponseWriter)
}

type GetAllInterActionsInterface interface {
	GetAllInteractions(title string) (int, []interactionData, error)
}

type interactionData struct {
	User    string `json:"user"`
	Comment string `json:"comment"`
}

// Structs
type AddNewLikeFunction struct{}
type AddNewCommentFunction struct{}
type GetAllInteractionsFunction struct{}

// Main Functions
func (a AddNewLikeFunction) AddNewLike(title string, user string, w http.ResponseWriter) {
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

func (ac AddNewCommentFunction) AddNewComment(title string, user string, email string, comment string, w http.ResponseWriter) {
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

func (gal GetAllInteractionsFunction) GetAllInteractions(title string) (int, []interactionData, error) {
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

type AddNewLikeService struct {
	AddNewLikeInterface AddNewLikeInterface
}

type AddNewCommentService struct {
	AddNewCommentInterface AddNewCommentInterface
}

type GetAllInteractionsService struct {
	GetAllInterActionsInterface GetAllInterActionsInterface
}

func (anl AddNewLikeService) AddNewLikeFinal(title string, user string, w http.ResponseWriter) {
	anl.AddNewLikeInterface.AddNewLike(title, user, w)
}

func (anc AddNewCommentService) AddNewCommentFinal(title string, user string, email string, comment string, w http.ResponseWriter) {
	anc.AddNewCommentInterface.AddNewComment(title, user, email, comment, w)
}

func (gal GetAllInteractionsService) GetAllInteractionsFinal(title string) (int, []interactionData, error) {
	likeCount, comments, err := gal.GetAllInterActionsInterface.GetAllInteractions(title)
	return likeCount, comments, err
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
