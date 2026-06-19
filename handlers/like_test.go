package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"database/sql"
	"forum/database"
	"forum/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatal(err)
	}
	if err := database.CreateTables(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func createTestSession(t *testing.T, db *sql.DB, userID int) string {
	res, err := db.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, datetime('now', '+1 day'))", "session-test", userID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return "session-test"
}

func TestLikeCommentRedirectsBackToPostWhenNoReferer(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	res, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", "testuser", "test@example.com", "pass")
	if err != nil {
		t.Fatal(err)
	}
	userID, _ := res.LastInsertId()

	res, err = db.Exec("INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)", userID, "Test post", "Body", "general")
	if err != nil {
		t.Fatal(err)
	}
	postID, _ := res.LastInsertId()

	res, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", "commenter", "commenter@example.com", "pass")
	if err != nil {
		t.Fatal(err)
	}
	commenterID, _ := res.LastInsertId()

	res, err = db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, commenterID, "Nice post")
	if err != nil {
		t.Fatal(err)
	}
	commentID, _ := res.LastInsertId()

	sessionID := createTestSession(t, db, int(userID))

	likeHandler := NewLikeHandler(db)
	wrapped := middleware.Auth(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		likeHandler.LikeComment(w, r)
	}))

	body := url.Values{}
	body.Set("comment_id", strconv.FormatInt(commentID, 10))
	req := httptest.NewRequest(http.MethodPost, "/comment/like", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID, Path: "/"})

	rw := httptest.NewRecorder()
	wrapped.ServeHTTP(rw, req)

	if rw.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status %d, got %d", http.StatusSeeOther, rw.Code)
	}

	expected := fmt.Sprintf("/post/%d", postID)
	if got := rw.Header().Get("Location"); got != expected {
		t.Fatalf("expected redirect to %q, got %q", expected, got)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?", commentID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 like row, got %d", count)
	}
}

func TestLikePostForbiddenForAuthor(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	res, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", "author", "author@example.com", "pass")
	if err != nil {
		t.Fatal(err)
	}
	userID, _ := res.LastInsertId()

	res, err = db.Exec("INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)", userID, "Author post", "Body", "general")
	if err != nil {
		t.Fatal(err)
	}
	postID, _ := res.LastInsertId()

	sessionID := createTestSession(t, db, int(userID))

	likeHandler := NewLikeHandler(db)
	wrapped := middleware.Auth(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		likeHandler.Like(w, r)
	}))

	body := url.Values{}
	body.Set("post_id", strconv.FormatInt(postID, 10))
	req := httptest.NewRequest(http.MethodPost, "/like", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID, Path: "/"})

	rw := httptest.NewRecorder()
	wrapped.ServeHTTP(rw, req)

	if rw.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rw.Code)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", postID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0 like rows, got %d", count)
	}
}

func TestLikeCommentForbiddenForAuthor(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	res, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", "author", "author@example.com", "pass")
	if err != nil {
		t.Fatal(err)
	}
	userID, _ := res.LastInsertId()

	res, err = db.Exec("INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)", userID, "Author post", "Body", "general")
	if err != nil {
		t.Fatal(err)
	}
	postID, _ := res.LastInsertId()

	res, err = db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, "Author comment")
	if err != nil {
		t.Fatal(err)
	}
	commentID, _ := res.LastInsertId()

	sessionID := createTestSession(t, db, int(userID))

	likeHandler := NewLikeHandler(db)
	wrapped := middleware.Auth(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		likeHandler.LikeComment(w, r)
	}))

	body := url.Values{}
	body.Set("comment_id", strconv.FormatInt(commentID, 10))
	req := httptest.NewRequest(http.MethodPost, "/comment/like", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID, Path: "/"})

	rw := httptest.NewRecorder()
	wrapped.ServeHTTP(rw, req)

	if rw.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rw.Code)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?", commentID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0 comment like rows, got %d", count)
	}
}
