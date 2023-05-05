package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

var filename string = "users.json"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	//Deleted   string `json:"deleted"`
	//CreatedAt string `json:"createdat"`
	//UpdatedAt string `json:"updatedat"`
}

type UserRepositoryI interface {
	Create(user *User) (int, error)
	GetByEmail(email string) (*User, error)
	GetAll() (*[]User, error)
	Update(user *User) (*User, error)
	Delete(id int) error
}

type UserRepository struct{}

func (r *UserRepository) Create(user *User) (id int, err error) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}
	defer file.Close()
	_, err = file.Write(append(userJSON, '\n'))

	return user.ID, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		user := new(User)
		err := json.Unmarshal([]byte(scanner.Text()), &user)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if user.Email == email {
			return user, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return nil, nil
}

func (r *UserRepository) GetAll() (*[]User, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var users []User

	for scanner.Scan() {
		var user User
		err := json.Unmarshal([]byte(scanner.Text()), &user)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	//fmt.Println(users)
	return &users, nil
}

func main() {
	// create a file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	// initialize a repository
	r := UserRepository{}

	user := &User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "password",
		//Deleted:   "",
		//CreatedAt: "2021-05-01",
		//UpdatedAt: "2021-05-02",
	}

	user2 := &User{
		ID:        2,
		FirstName: "Will",
		LastName:  "Smith",
		Email:     "willsmith@gmail.com",
		Password:  "password",
	}

	r.Create(user)
	r.Create(user2)
	//user2, _ := r.GetByEmail("johndoe@example.com")
	//if user2 != nil {
	//	fmt.Println(user2.ID)
	//}

	//r.GetAll()
	//if err != nil {
	//	panic(err)
	//}
	users, _ := r.GetAll()
	for _, user := range *users {
		fmt.Println(user)
	}
	for i := 0; i < len(*users); i++ {
		fmt.Println((*users)[i].ID)
	}
	//fmt.Println(users)

	fmt.Println("Success")
	//fmt.Printf("%d\n", u.ID)
}
