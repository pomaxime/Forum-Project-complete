package repository

import (
    "database/sql"
    "forum/models"
)

type PostRepo struct {
    db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
    return &PostRepo{db: db}
}

func (r *PostRepo) GetAll(category string) ([]models.Post, error) {
    query := `SELECT p.id, p.user_id, u.username, p.title, p.content, p.category, p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id`
    var rows *sql.Rows
    var err error

    if category != "" {
        query += " WHERE p.category = ?"
        query += " ORDER BY p.created_at DESC"
        rows, err = r.db.Query(query, category)
    } else {
        query += " ORDER BY p.created_at DESC"
        rows, err = r.db.Query(query)
    }
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var p models.Post
        if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.Category, &p.CreatedAt); err != nil {
            continue
        }
        posts = append(posts, p)
    }
    return posts, nil
}

func (r *PostRepo) GetByID(id int) (*models.Post, error) {
    var p models.Post
    err := r.db.QueryRow(
        `SELECT p.id, p.user_id, u.username, p.title, p.content, p.category, p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = ?`,
        id,
    ).Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.Category, &p.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *PostRepo) Create(userID int, title, content, category string) (int, error) {
    result, err := r.db.Exec(
        "INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)",
        userID, title, content, category,
    )
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    return int(id), nil
}

func (r *PostRepo) Delete(id, userID int) error {
    _, err := r.db.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", id, userID)
    return err
}
