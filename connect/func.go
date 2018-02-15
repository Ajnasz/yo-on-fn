package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	_, err := redisClient.Ping().Result()

	if err != nil {
		log.Fatal(err)
	}
}

func isExistingUser(userName string) (bool, error) {
	res, err := redisClient.HExists("yoaccount", userName).Result()

	log.Println("IS EXISTINGS", res, err)

	if err != nil {
		return false, err
	}

	return res, nil
}

func connectUser(name, friendName string) error {
	return redisClient.SAdd("yoaccount_friend_"+name, friendName).Err()
}

func accountHandler(in io.Reader, out io.Writer) {
	var user struct {
		Name       string `json:"name"`
		Password   string `json:"password"`
		FriendName string `json:"friendName"`
	}

	json.NewDecoder(in).Decode(&user)

	token, err := redisClient.HGet("yoaccount", user.Name).Result()

	if err != nil {
		io.WriteString(out, "TOKEN ERROR")
		log.Fatal(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(token), []byte(user.Password))
	if err != nil {
		io.WriteString(out, fmt.Sprintf("Wrong password %+v %s, %s", err, token, user.Password))
		return
	}

	existing, err := isExistingUser(user.FriendName)

	if err != nil {
		io.WriteString(out, "ERROR")
		log.Fatal(err)
	}

	if !existing {
		io.WriteString(out, "USER NOT FOUND")
		return
	}

	connectUser(user.Name, user.FriendName)

	if err != nil {
		io.WriteString(out, "ERROR")
		log.Fatal(err)
	} else {
		// you can write your own headers & status, if you'd like to
		io.WriteString(out, "OK")
	}
}

func main() {
	accountHandler(os.Stdin, os.Stdout)
}
