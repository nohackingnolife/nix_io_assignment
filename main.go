package main

import (
	"fmt"
	"hw_nix_io/db/model"
	"hw_nix_io/db/repos"
	"os"
)

var filename string = "users.json"

func main() {
	// create a file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	// initialize a repository
	r := repos.UserRepository{Filename: "users.json"}

	// create instances
	user := &model.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "password",
	}

	user2 := &model.User{
		ID:        2,
		FirstName: "Will",
		LastName:  "Smith",
		Email:     "willsmith@gmail.com",
		Password:  "password",
	}

	// create users
	r.Create(user)
	r.Create(user2)

	email := "johndoe@example.com"
	userFoo, err := r.GetByEmail(email)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	if userFoo == nil {
		fmt.Printf("User with email %s not found\n", email)
	} else {
		fmt.Printf("User with email %s: %+v\n", email, userFoo)
	}

	fmt.Println("Fetching all the users:")
	users, _ := r.GetAll()
	for _, user := range *users {
		fmt.Println(user)
	}

	user2.Email = "smithwill@gmail.com"
	user2.Password = "newpassword"
	userFoo, err = r.Update(user2)
	if userFoo == nil {
		fmt.Printf("User with ID %d not found\n", user2.ID)
	} else {
		fmt.Printf("User with ID %d updated: %+v\n", user2.ID, userFoo)
	}

	err = r.Delete(1)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("User with ID %d was deleted\n", 1)

	fmt.Println("Fetching all the users:")
	users, _ = r.GetAll()
	for _, user := range *users {
		fmt.Println(user)
	}
}
