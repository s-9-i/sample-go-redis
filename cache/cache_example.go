package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"sample-go-redis/client"
)

var c Cache

const (
	keySimpleCache = "key-cache-simple"
	keyJSONCache   = "key-cache-json"
)

func main() {
	ctx := context.Background()
	cli := client.NewClient(&client.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer func() {
		cli.Del(ctx, keySimpleCache, keyJSONCache)
		_ = cli.Close()
	}()

	c = NewCache(cli)
	SampleSimpleKeyValue(ctx)
	SampleJSON(ctx)
}

func SampleSimpleKeyValue(ctx context.Context) {
	c.Set(ctx, keySimpleCache, "10", time.Hour)
	v1, _ := c.Get(ctx, keySimpleCache)
	fmt.Printf("【cache】value = %v\n", v1)

	c.Increment(ctx, keySimpleCache)
	c.Increment(ctx, keySimpleCache)
	c.Decrement(ctx, keySimpleCache)
	v2, _ := c.Get(ctx, keySimpleCache)
	fmt.Printf("【cache】value = %v\n", v2)

	c.Delete(ctx, keySimpleCache)
	v3, _ := c.Get(ctx, keySimpleCache)
	fmt.Printf("【cache】value = %v\n", v3)
}

func SampleJSON(ctx context.Context) {
	json, _ := MarshalJSON(&User{
		ID:   "userID",
		Name: "name",
		Age:  30,
	})
	c.Set(ctx, keyJSONCache, json, time.Hour)
	v, _ := c.Get(ctx, keyJSONCache)
	unmarshalled, _ := UnmarshalJSON(v)
	fmt.Printf("【cache】value = %#v\n", unmarshalled)
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func MarshalJSON(user *User) (string, error) {
	b, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func UnmarshalJSON(jsonData string) (*User, error) {
	var user User
	if err := json.Unmarshal(([]byte)(jsonData), &user); err != nil {
		return nil, err
	}
	return &user, nil
}
