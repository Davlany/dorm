package main

import (
	"dorm"
	"dorm/Drivers"
	_ "github.com/lib/pq"
	"log"
)

type User struct {
	Id       int    `db:"id" serial:"true"`
	Name     string `db:"name"`
	Password string `db:"password"`
	Posts    []Post `rel:"posts" field:"user_id"`
}

// SELECT * FROM users WHERE id = user.id
// SELECT * FROM posts WHERE user_id = user.id
//user.posts = posts

type Post struct {
	Id          int    `db:"id" serial:"true"`
	Caption     string `db:"caption"`
	Description string `db:"description"`
	Likes       int    `db:"likes"`
	UserId      int    `db:"user_id" fk:"true" field:"id" rel:"users"` // INTEGER user_id FOREIGN KEY (user_id) REFERENCES users(id)
}

func main() {
	driver, err := Drivers.NewPostgresDriver("postgres", "123456", "testusers", "5432", "disable")
	if err != nil {
		log.Fatalln(err)
	}
	db := dorm.NewDatabase(driver)
	err = db.RegisterSchemas([]interface{}{User{}, Post{}})
	if err != nil {
		log.Fatalln(err)
	}
	userTable := db.Table("users", User{})

	_ = userTable

}
