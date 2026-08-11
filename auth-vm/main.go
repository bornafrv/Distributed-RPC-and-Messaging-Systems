package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Username string
	Password string
}

type AuthResponse struct {
	Success bool
	Message string
}

type AuthService struct {
	Users map[string]string
}

func loadUsers(filename string) (map[string]string, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var userList []User
	err = json.Unmarshal(file, &userList)
	if err != nil {
		return nil, err
	}

	users := make(map[string]string)
	for _, u := range userList {
		users[u.Username] = u.Password
	}

	return users, nil
}

func (a *AuthService) Authenticate(req AuthRequest, res *AuthResponse) error {
	log.Printf("JSON-RPC Authenticate called for username: %s", req.Username)

	password, exists := a.Users[req.Username]
	if !exists {
		res.Success = false
		res.Message = "user not found"
		return nil
	}

	if password != req.Password {
		res.Success = false
		res.Message = "invalid password"
		return nil
	}

	res.Success = true
	res.Message = "authentication successful"
	return nil
}

func main() {
	users, err := loadUsers("users.json")
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	authService := &AuthService{
		Users: users,
	}

	err = rpc.Register(authService)
	if err != nil {
		log.Fatalf("Failed to register RPC service: %v", err)
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Failed to listen on port 1234: %v", err)
	}

	fmt.Println("Auth JSON-RPC server running on port 1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}
