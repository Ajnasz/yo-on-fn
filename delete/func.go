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

func accountHandler(in io.Reader, out io.Writer) {
	var user struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	}

	json.NewDecoder(in).Decode(&user)

	token, err := redisClient.HGet("yoaccount", user.Name).Result()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("Cant get token %+v", err))
		return
	}

	if token == user.Token {
		err := redisClient.HDel("yoaccount", user.Name).Err()

		if err != nil {
			io.WriteString(out, fmt.Sprintf("delete error %+v", err))
			return
		}

		redisSetName := "yoaccount_friend_" + user.Name

		for {
			mem, err := redisClient.SRandMember(redisSetName).Result()

			if err == redis.Nil {
				break
			} else if err != nil {
				io.WriteString(out, fmt.Sprintf("get rand member: %s, %+v", user.Name, err))
				return
			}

			err = redisClient.SRem(redisSetName, mem).Err()

			if err != nil {
				io.WriteString(out, fmt.Sprintf("SRem: %+v", err))
				return
			}
		}

	} else {
		io.WriteString(out, "INVALID TOKEN")
		return
	}

	if err != nil {
		io.WriteString(out, fmt.Sprintf("%+v", err))
		return
	} else {
		// you can write your own headers & status, if you'd like to
		io.WriteString(out, "OK\n")
	}
}

func main() {
	accountHandler(os.Stdin, os.Stdout)
}
