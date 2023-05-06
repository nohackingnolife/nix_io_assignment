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
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}

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
		var u User
		err := json.Unmarshal([]byte(scanner.Text()), &u)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if u.Email == email {
			return &u, nil
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
		var u User
		err := json.Unmarshal([]byte(scanner.Text()), &u)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		users = append(users, u)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &users, nil
}

func (r *UserRepository) Update(user *User) (*User, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var output []byte
	isFound := false

	for scanner.Scan() {
		var u User
		err := json.Unmarshal([]byte(scanner.Text()), &u)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}

		if u.ID == user.ID {
			data, err := json.Marshal(user)
			if err != nil {
				fmt.Println("Error: ", err)
				return nil, err
			}
			output = append(output, data...)
			isFound = true
		} else {
			output = append(output, []byte(scanner.Text())...)
		}

		output = append(output, '\n')
	}

	file.Close()

	file, err = os.Create(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer file.Close()

	_, err = file.Write(output)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	if isFound {
		return user, nil
	}
	return nil, nil
}

func (r *UserRepository) Delete(id int) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var output []byte

	for scanner.Scan() {
		var u User
		err := json.Unmarshal([]byte(scanner.Text()), &u)
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		if u.ID == id {
			continue
		} else {
			output = append(output, []byte(scanner.Text())...)
			output = append(output, '\n')
		}
	}

	file.Close()

	file, err = os.Create(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(output)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return err
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

	// create instances
	user := &User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "password",
	}

	user2 := &User{
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
