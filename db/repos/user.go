package repos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hw_nix_io/db/model"
	"os"
)

type UserRepositoryI interface {
	Create(user *model.User) (int, error)
	GetByEmail(email string) (*model.User, error)
	GetAll() (*[]model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id int) error
}

type UserRepository struct {
	Filename string
}

func (r *UserRepository) Create(user *model.User) (id int, err error) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}

	file, err := os.OpenFile(r.Filename, os.O_APPEND|os.O_WRONLY, 0644)
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

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	file, err := os.Open(r.Filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var u model.User
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

func (r *UserRepository) GetAll() (*[]model.User, error) {
	file, err := os.Open(r.Filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var users []model.User

	for scanner.Scan() {
		var u model.User
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

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	file, err := os.Open(r.Filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var output []byte
	isFound := false

	for scanner.Scan() {
		var u model.User
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

	file, err = os.Create(r.Filename)
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
	file, err := os.Open(r.Filename)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var output []byte

	for scanner.Scan() {
		var u model.User
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

	file, err = os.Create(r.Filename)
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
