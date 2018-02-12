package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

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

func isFriend(userName string, friendName string) (bool, error) {
	res, err := redisClient.SIsMember("yoaccount_friend_"+userName, friendName).Result()

	if err != nil {
		return false, err
	}

	return res, nil
}

func poke(userName, friendName string) string {
	return userName + " poked " + friendName
}

func accountHandler(in io.Reader, out io.Writer) {
	var user struct {
		Name       string `json:"name"`
		Token      string `json:"token"`
		FriendName string `json:"friendName"`
	}

	json.NewDecoder(in).Decode(&user)

	token, err := redisClient.HGet("yoaccount", user.Name).Result()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("Cant get token %+v", err))
		return
	}

	if token != user.Token {
		io.WriteString(out, "INVALID TOKEN")
		return
	}

	friend, err := isFriend(user.Name, user.FriendName)

	if err != nil {
		io.WriteString(out, fmt.Sprintf("isFriend %v", err))
		return
	}

	if !friend {
		io.WriteString(out, "USER NOT FRIEND")
		return
	}

	io.WriteString(out, poke(user.Name, user.FriendName))
}

func main() {
	accountHandler(os.Stdin, os.Stdout)
}
