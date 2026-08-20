package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"forum/database"
	"forum/handlers"
	"forum/middleware"
	"forum/utils"

	"github.com/google/uuid"
)

// setupTestDB initialise une base de données SQLite en mémoire pour les tests.
// Chaque test repart d'une base vide et propre.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Impossible de créer la base de test : %v", err)
	}
	db.Exec("PRAGMA foreign_keys = ON")
	if err := database.CreateTables(db); err != nil {
		t.Fatalf("Impossible de créer les tables : %v", err)
	}
	return db
}

// createTestUser insère un utilisateur de test directement en base.
func createTestUser(t *testing.T, db *sql.DB, username, email, password string) int {
	t.Helper()
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Erreur hachage : %v", err)
	}
	res, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hash)
	if err != nil {
		t.Fatalf("Erreur insertion utilisateur : %v", err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

// createTestSession crée une session valide pour un utilisateur donné.
func createTestSession(t *testing.T, db *sql.DB, userID int) string {
	t.Helper()
	sessionID := uuid.New().String()
	_, err := db.Exec(
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, time.Now().Add(24*time.Hour),
	)
	if err != nil {
		t.Fatalf("Erreur création session : %v", err)
	}
	return sessionID
}

// =============================================
// Tests Validation
// =============================================

func TestValidateRegister_OK(t *testing.T) {
	if err := utils.ValidateRegister("alice", "alice@example.com", "password123"); err != nil {
		t.Errorf("Validation correcte échouée : %v", err)
	}
}

func TestValidateRegister_ChampsVides(t *testing.T) {
	if err := utils.ValidateRegister("", "alice@example.com", "password123"); err == nil {
		t.Error("Aurait dû échouer avec username vide")
	}
	if err := utils.ValidateRegister("alice", "", "password123"); err == nil {
		t.Error("Aurait dû échouer avec email vide")
	}
	if err := utils.ValidateRegister("alice", "alice@example.com", ""); err == nil {
		t.Error("Aurait dû échouer avec password vide")
	}
}

func TestValidateRegister_EmailInvalide(t *testing.T) {
	if err := utils.ValidateRegister("alice", "pas-un-email", "password123"); err == nil {
		t.Error("Aurait dû échouer avec email invalide")
	}
}

func TestValidateRegister_PasswordTropCourt(t *testing.T) {
	if err := utils.ValidateRegister("alice", "alice@example.com", "abc"); err == nil {
		t.Error("Aurait dû échouer avec password < 8 caractères")
	}
}

func TestValidateRegister_UsernameTropCourt(t *testing.T) {
	if err := utils.ValidateRegister("ab", "alice@example.com", "password123"); err == nil {
		t.Error("Aurait dû échouer avec username < 3 caractères")
	}
}

// =============================================
// Tests bcrypt
// =============================================

func TestHashPassword(t *testing.T) {
	hash, err := utils.HashPassword("monmotdepasse")
	if err != nil {
		t.Fatalf("Erreur hachage : %v", err)
	}
	if hash == "monmotdepasse" {
		t.Error("Le hash ne doit pas être identique au mot de passe en clair")
	}
	if len(hash) < 20 {
		t.Error("Hash trop court, attendu un hash bcrypt")
	}
}

func TestCheckPassword_OK(t *testing.T) {
	hash, _ := utils.HashPassword("secret")
	if err := utils.CheckPassword(hash, "secret"); err != nil {
		t.Errorf("Vérification correcte échouée : %v", err)
	}
}

func TestCheckPassword_Incorrect(t *testing.T) {
	hash, _ := utils.HashPassword("secret")
	if err := utils.CheckPassword(hash, "mauvais"); err == nil {
		t.Error("Aurait dû échouer avec mauvais mot de passe")
	}
}

// =============================================
// Tests Register Handler
// =============================================

func TestRegister_GET(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := handlers.NewAuthHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	rr := httptest.NewRecorder()
	h.Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /register : attendu 200, reçu %d", rr.Code)
	}
}

func TestRegister_POST_OK(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := handlers.NewAuthHandler(db)

	form := url.Values{}
	form.Set("username", "alice")
	form.Set("email", "alice@example.com")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.Register(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Register réussi : attendu 303, reçu %d", rr.Code)
	}
	if !strings.Contains(rr.Header().Get("Location"), "/login") {
		t.Error("Redirection attendue vers /login")
	}

	// Vérifie que l'utilisateur est bien en base
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", "alice@example.com").Scan(&count)
	if count != 1 {
		t.Error("Utilisateur non trouvé en base après inscription")
	}
}

func TestRegister_POST_EmailDuplique(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	createTestUser(t, db, "alice", "alice@example.com", "password123")
	h := handlers.NewAuthHandler(db)

	form := url.Values{}
	form.Set("username", "alice2")
	form.Set("email", "alice@example.com") // Email déjà pris
	form.Set("password", "password456")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.Register(rr, req)

	if rr.Code == http.StatusSeeOther {
		t.Error("Aurait dû échouer avec email dupliqué")
	}
}

// =============================================
// Tests Login Handler
// =============================================

func TestLogin_POST_OK(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	createTestUser(t, db, "alice", "alice@example.com", "password123")
	h := handlers.NewAuthHandler(db)

	form := url.Values{}
	form.Set("email", "alice@example.com")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.Login(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Login réussi : attendu 303, reçu %d", rr.Code)
	}

	// Vérifie qu'un cookie session_id a bien été créé
	cookies := rr.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session_id" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Error("Cookie session_id non créé après login")
	}
	if !sessionCookie.HttpOnly {
		t.Error("Cookie doit être HttpOnly")
	}
}

func TestLogin_POST_MotDePasseIncorrect(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	createTestUser(t, db, "alice", "alice@example.com", "password123")
	h := handlers.NewAuthHandler(db)

	form := url.Values{}
	form.Set("email", "alice@example.com")
	form.Set("password", "mauvais_password")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.Login(rr, req)

	if rr.Code == http.StatusSeeOther {
		t.Error("Login ne devrait pas réussir avec un mauvais mot de passe")
	}
}

func TestLogin_POST_EmailInexistant(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := handlers.NewAuthHandler(db)

	form := url.Values{}
	form.Set("email", "inexistant@example.com")
	form.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.Login(rr, req)

	if rr.Code == http.StatusSeeOther {
		t.Error("Login ne devrait pas réussir avec un email inexistant")
	}
}

// =============================================
// Tests Logout Handler
// =============================================

func TestLogout_OK(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	userID := createTestUser(t, db, "alice", "alice@example.com", "password123")
	sessionID := createTestSession(t, db, userID)
	h := handlers.NewAuthHandler(db)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rr := httptest.NewRecorder()
	h.Logout(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Logout : attendu 303, reçu %d", rr.Code)
	}

	// Vérifie que la session a été supprimée de la base
	var count int
	db.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = ?", sessionID).Scan(&count)
	if count != 0 {
		t.Error("La session aurait dû être supprimée de la base")
	}
}

// =============================================
// Tests Middleware Auth
// =============================================

func TestMiddlewareAuth_SansSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	protected := middleware.Auth(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/post/create", nil)
	rr := httptest.NewRecorder()
	protected.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Sans session : attendu redirect 303, reçu %d", rr.Code)
	}
}

func TestMiddlewareAuth_SessionValide(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	userID := createTestUser(t, db, "alice", "alice@example.com", "password123")
	sessionID := createTestSession(t, db, userID)

	protected := middleware.Auth(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/post/create", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rr := httptest.NewRecorder()
	protected.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Session valide : attendu 200, reçu %d", rr.Code)
	}
}

func TestMiddlewareAuth_SessionExpiree(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	userID := createTestUser(t, db, "alice", "alice@example.com", "password123")

	// Session déjà expirée (date dans le passé)
	sessionID := uuid.New().String()
	db.Exec(
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, time.Now().Add(-1*time.Hour),
	)

	protected := middleware.Auth(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/post/create", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rr := httptest.NewRecorder()
	protected.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Session expirée : attendu redirect 303, reçu %d", rr.Code)
	}
}

func TestMiddlewareAuth_SessionInvalide(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	protected := middleware.Auth(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/post/create", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-uuid-completement-invalide"})
	rr := httptest.NewRecorder()
	protected.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Session invalide : attendu redirect 303, reçu %d", rr.Code)
	}
}
