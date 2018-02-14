package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	redis "github.com/go-redis/redis"
)

type PushData struct {
	Name       string `json:"name"`
	FriendName string `json:"friendName"`
	Endpoint   string `json:"endpoint"`
	Key        string `json:"key"`
}

type PushUser struct {
	Name       string `json:"name"`
	FriendName string `json:"friendName"`
	Token      string `json:"token"`
}

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

func poke(pushData *PushData) (string, error) {
	return pushData.Name + " poked " + pushData.FriendName, nil
}

func getPushData(user PushUser) (*PushData, error) {
	v := make(map[string]string)
	var keys [3]string
	keys[0] = "yoaccount_endpoint"
	keys[1] = "yoaccount_key"

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		data, err := redisClient.HGet(key, user.FriendName).Result()

		if err != nil {
			return nil, err
		}

		v[key] = data
	}

	return &PushData{
		Name:       user.Name,
		FriendName: user.FriendName,
		Endpoint:   v["yoaccount_endpoint"],
		Key:        v["yoaccount_key"],
	}, nil
}

func accountHandler(in io.Reader, out io.Writer) {
	var user PushUser

	json.NewDecoder(in).Decode(&user)

	token, err := redisClient.HGet("yoaccount", user.Name).Result()

	if err != nil {
		io.WriteString(out, fmt.Sprintf("Can't get token %+v", err))
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

	pushData, err := getPushData(user)
	if err != nil {
		io.WriteString(out, fmt.Sprintf("Can't get pushData %+v", err))
		return
	}
	ret, err := poke(pushData)
	if err != nil {
		io.WriteString(out, fmt.Sprintf("Can't poke %+v", err))
		return
	}

	io.WriteString(out, ret)
}

func main() {
	accountHandler(os.Stdin, os.Stdout)
}
