package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"encoding/base64"
	"golang.org/x/crypto/bcrypt"

	redis "github.com/go-redis/redis"
)

var redisClient *redis.Client

func init() {
	redisAddr := os.Getenv("REDIS_URL")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       9,
	})
}

// GenerateToken returns a unique token based on the provided user string
func GenerateToken(user string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// log.Println("Hash to store:", string(hash))

	return base64.StdEncoding.EncodeToString(hash), nil
}

func accountHandler(in io.Reader, out io.Writer) {
	var user struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
		Key      string `json:"key"`
	}

	json.NewDecoder(in).Decode(&user)

	if user.Name == "" {
		io.WriteString(out, "ERR: invalid username")
		return
	}

	if user.Endpoint == "" {
		io.WriteString(out, "ERR: invalid endpoint")
		return
	}

	if user.Key == "" {
		io.WriteString(out, "ERR: invalid key")
		return
	}

	err := redisClient.HSet("yoaccount_endpoint", user.Name, user.Endpoint).Err()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("ERR: account creation endpoint %v", err))
		return
	}

	err = redisClient.HSet("yoaccount_key", user.Name, user.Key).Err()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("ERR: account creation key %v", err))
		return
	}

	token, err := GenerateToken(user.Name)

	if err != nil {
		io.WriteString(out, fmt.Sprintf("ERR: token generation %v", err))
		return
	}

	err = redisClient.HSet("yoaccount", user.Name, token).Err()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("ERR: account creation %v", err))
		return
	}

	io.WriteString(out, token)
}

func main() {
	accountHandler(os.Stdin, os.Stdout)
}
