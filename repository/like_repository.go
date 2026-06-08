package repository

import (
	"database/sql"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Like(postID, userID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	tx.Exec("DELETE FROM dislikes WHERE post_id = ? AND user_id = ?", postID, userID)

	_, err = tx.Exec("INSERT OR IGNORE INTO likes (post_id, user_id) VALUES (?, ?)", postID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *LikeRepository) Dislike(postID, userID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	tx.Exec("DELETE FROM likes WHERE post_id = ? AND user_id = ?", postID, userID)

	_, err = tx.Exec("INSERT OR IGNORE INTO dislikes (post_id, user_id) VALUES (?, ?)", postID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *LikeRepository) CountLikes(postID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", postID).Scan(&count)
	return count, err
}
func (r *LikeRepository) LikeComment(commentID, userID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	tx.Exec("DELETE FROM comment_dislikes WHERE comment_id = ? AND user_id = ?", commentID, userID)

	_, err = tx.Exec("INSERT OR IGNORE INTO comment_likes (comment_id, user_id) VALUES (?, ?)", commentID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *LikeRepository) DislikeComment(commentID, userID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	tx.Exec("DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)

	_, err = tx.Exec("INSERT OR IGNORE INTO comment_dislikes (comment_id, user_id) VALUES (?, ?)", commentID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *LikeRepository) CountCommentLikes(commentID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?", commentID).Scan(&count)
	return count, err
}

func (r *LikeRepository) GetPostAuthorID(postID int) (int, error) {
	var userID int
	err := r.db.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&userID)
	return userID, err
}

func (r *LikeRepository) GetCommentPostID(commentID int) (int, error) {
	var postID int
	err := r.db.QueryRow("SELECT post_id FROM comments WHERE id = ?", commentID).Scan(&postID)
	return postID, err
}

func (r *LikeRepository) GetCommentAuthorID(commentID int) (int, error) {
	var userID int
	err := r.db.QueryRow("SELECT user_id FROM comments WHERE id = ?", commentID).Scan(&userID)
	return userID, err
}

func (r *LikeRepository) CountCommentDislikes(commentID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM comment_dislikes WHERE comment_id = ?", commentID).Scan(&count)
	return count, err
}
func (r *LikeRepository) CountDislikes(postID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ?", postID).Scan(&count)
	return count, err
}
