package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// User represents a user from the API
type User struct {
	Name string `json:"name"`
}

// HTTPClient interface defines the methods we need from http.Client
// This is the key to testability - we define a minimal interface
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// GetUser fetches a user by ID from the API
func GetUser(client HTTPClient, id int) (User, error) {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", id)

	resp, err := client.Get(url)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, fmt.Errorf("failed to read response: %w", err)
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return User{}, fmt.Errorf("failed to parse user: %w", err)
	}

	return user, nil
}

func main() {
	user, err := GetUser(&http.Client{}, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User: %s\n", user.Name)
}
