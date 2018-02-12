package main

import (
	"encoding/json"
	"io"
	"log"
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

	pong, err := redisClient.Ping().Result()
	log.Println("PONG", pong)

	if err != nil {
		log.Fatal(err)
	}
}

// GenerateToken returns a unique token based on the provided user string
func GenerateToken(user string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(user), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Hash to store:", string(hash))

	return base64.StdEncoding.EncodeToString(hash)
}

func accountHandler(in io.Reader, out io.Writer) {
	var user struct {
		Name string `json:"name"`
	}

	json.NewDecoder(in).Decode(&user)

	token := GenerateToken(user.Name)
	err := redisClient.HSet("yoaccount", user.Name, token).Err()

	if err != nil {
		io.WriteString(out, "ERROR")
		log.Fatal(err)
	}

	io.WriteString(out, token)
}

func main() {
	accountHandler(os.Stdin, os.Stdout)
}
