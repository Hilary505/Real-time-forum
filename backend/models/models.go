package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
	Type       string    `json:"type"` // message, typing, stop_typing
	Online     bool      `json:"online"`
	Nickname   string    `json:"nickname"`
}

type Reaction struct {
	ID        int  `json:"id,omitempty"`
	UserID    int  `json:"user_id"`
	PostID    int  `json:"post_id,omitempty"`
	CommentID int  `json:"comment_id,omitempty"`
	IsLike    bool `json:"is_like"`
}

type User struct {
	ID                   int       `json:"id"`
	UUID                 string    `json:"user_uuid"`
	Nickname             string    `json:"nickname"`
	Age                  int       `json:"age"`
	Gender               string    `json:"gender"`
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	Email                string    `json:"email"`
	Password             string    `json:"Password"`
	CreatedAt            time.Time `json:"createdAt"`
	IsOnline             bool      `json:"isOnline"`
	LastPrivateMessageAt time.Time `json:"last_private_message_at"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	UserLiked     int       `json:"UserLiked"`
	Post_id       int       `json:"post_id"`
	Title         string    `json:"title"`
	Category      string    `json:"Category"`
	Content       string    `json:"content"`
	User_uuid     string    `json:"user_uuid"`
	CreatedAt     time.Time `json:"created_at"`
	Nickname      string    `json:"nickname"`
	Categories    []string  `json:"categories"`
	Likes         int       `json:"likes"`
	Dislikes      int       `json:"dislikes"`
	ImageUrl      string    `json:"image_url"`
	CommentsCount int       `json:"comments_count"`
}

type Comment struct {
	Comment_id int       `json:"comment_id"`
	User_uuid  string    `json:"user_uuid"`
	Post_id    int       `json:"post_id"`
	CreatedAt  time.Time `json:"created_at"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	Content    string    `json:"content"`
	UserLiked  int       `json:"UserLiked"`
	Author     string    `json:"author"`
}

// =====  hashes the user's password before storing it ====
func (user *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("|hashpassword method| ---> {%v}", err)
		return err
	}
	user.Password = string(hashed)
	return nil
}
