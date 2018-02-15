package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	redis "github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
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
func GenerateToken(text string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
}

func accountHandler(in io.Reader, out io.Writer) {
	var user struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Endpoint string `json:"endpoint"`
		Key      string `json:"key"`
	}

	err := json.NewDecoder(in).Decode(&user)

	if err != nil {
		io.WriteString(out, "ERR: invalid json")
	}

	io.WriteString(out, fmt.Sprintf("json %+v", err))

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

	err = redisClient.HSet("yoaccount_endpoint", user.Name, user.Endpoint).Err()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("ERR: account creation endpoint %v", err))
		return
	}

	err = redisClient.HSet("yoaccount_key", user.Name, user.Key).Err()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("ERR: account creation key %v", err))
		return
	}

	token, err := GenerateToken(user.Password)

	if err != nil {
		io.WriteString(out, fmt.Sprintf("ERR: token generation %v", err))
		return
	}

	err = redisClient.HSet("yoaccount", user.Name, string(token)).Err()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("ERR: account creation %v", err))
		return
	}

	io.WriteString(out, "OK")
}

func main() {
	accountHandler(os.Stdin, os.Stdout)
}
