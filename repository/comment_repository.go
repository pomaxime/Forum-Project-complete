package repository

import (
    "database/sql"
    "forum/models"
)

type CommentRepo struct {
    db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
    return &CommentRepo{db: db}
}

func (r *CommentRepo) GetByPostID(postID int) ([]models.Comment, error) {
    rows, err := r.db.Query(
        `SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at ASC`,
        postID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []models.Comment
    for rows.Next() {
        var comment models.Comment
        if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Username, &comment.Content, &comment.CreatedAt); err != nil {
            continue
        }
        comments = append(comments, comment)
    }
    return comments, nil
}

func (r *CommentRepo) Create(postID, userID int, content string) (int, error) {
    result, err := r.db.Exec(
        "INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)",
        postID, userID, content,
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
