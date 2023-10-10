package main

import (
	"dorm"
	"dorm/Drivers"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type User struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

type Post struct {
	Id          int
	Caption     string
	Description string
	Likes       int
}

func main() {
	driver, err := Drivers.NewPostgresDriver("postgres", "123456", "testusers", "5432", "disable")
	if err != nil {
		log.Fatalln(err)
	}
	db := dorm.NewDatabase(driver)
	table := db.Table("users", User{})

	var user []User
	err = table.FindAll(&user)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(user)
}
