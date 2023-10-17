package main

import (
	"dorm"
	"dorm/Drivers"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type User struct {
	Id       int    `db:"id" serial:"true"`
	Name     string `db:"name"`
	Password string `db:"password"`
	Posts    []Post `rel:"posts" rel_type:"one_to_many" field:"user_id"`
}

type Post struct {
	Id          int    `db:"id" serial:"true"`
	Caption     string `db:"caption"`
	Description string `db:"description"`
	Likes       int    `db:"likes"`
	UserId      int    `db:"user_id" fk:"true"`
}

func main() {
	driver, err := Drivers.NewPostgresDriver("postgres", "123456", "testusers", "5432", "disable")
	if err != nil {
		log.Fatalln(err)
	}
	db := dorm.NewDatabase(driver)
	userTable := db.Table("users", User{})
	var user User
	err = userTable.FindOne(1111, user)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(user.Posts)
}
