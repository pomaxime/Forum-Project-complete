package models

type UserRepository interface {
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(username, email, hashedPassword string) (int, error)
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
}

type PostRepository interface {
	GetAll(category string) ([]Post, error)
	GetByID(id int) (*Post, error)
	Create(userID int, title, content, category string) (int, error)
	Delete(id, userID int) error
}

type CommentRepository interface {
	GetByPostID(postID int) ([]Comment, error)
	Create(postID, userID int, content string) (int, error)
}
