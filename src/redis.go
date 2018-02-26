package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {

	InitializeRedis()
}

// Resets the redis client
func InitializeRedis() {

	client = NewClient(os.Getenv("REDIS_URL"))

	client.Set("users:nextuserid", "0", 0).Result()
	client.Del("users:usernames")
}

func NewClient(address string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	return client
}

// ClientUserCreate creates a user entry if the username is not taken
// and increments the userid counter
func ClientUserCreate(u User) error {

	val, err := client.SIsMember("users:usernames", u.Username).Result()
	if err != nil {
		fmt.Println("There was an error checking the username list")
		return err
	}

	if val {
		return errors.New("The username was already taken")
	}

	_, err = client.SAdd("users:usernames", u.Username).Result()
	if err != nil {
		fmt.Println("There was an error adding the username to the set")
		return err
	}

	intVal, err := client.Incr("users:nextuserid").Result()
	if err != nil {
		return err
	}

	newUserId := intVal - 1
	_, err = client.SAdd("users:userids", newUserId).Result()

	if err != nil {
		fmt.Println("There was an error adding the userid to the st")
		return err
	}

	userUsernameKey := fmt.Sprintf("users:%d:username", newUserId)
	_, err = client.Set(userUsernameKey, u.Username, 0).Result()
	if err != nil {
		return err
	}

	userNameKey := fmt.Sprintf("users:%d:name", newUserId)
	_, err = client.Set(userNameKey, u.Name, 0).Result()
	if err != nil {
		return err
	}

	return nil
}

func ClientUserInfo(userId int) (User, error) {

	emptyUser := User{"", "", "", ""}

	userUsernameKey := fmt.Sprintf("users:%d:username", userId)
	userUsername, err := client.Get(userUsernameKey).Result()

	if err != nil {

		if err.Error() == "redis: nil" {
			return emptyUser, errors.New(fmt.Sprintf("This user id (%d) couldn't be found!", userId))
		}
		fmt.Println("Couldn't find the user here")
		return emptyUser, err
	}

	userNameKey := fmt.Sprintf("users:%d:username", userId)
	userName, err := client.Get(userNameKey).Result()

	if err != nil {
		return emptyUser, err
	}

	return User{userUsername, userName, "", ""}, nil

}

// ClientUserGetAttribute looks up an attribute in a redis store and returns it if it exists.
func ClientUserGetAttribute(userId int, attributeRequested string) (string, error) {

	exists, err := client.SIsMember("users:userids", userId).Result()

	if exists != true {
		errorMessage := fmt.Sprintf("Couldn't find user %d", userId)
		return "", errors.New(errorMessage)
	} else if err != nil {
		return "", err
	}

	attributeKey := fmt.Sprintf("users:%d:%s", userId, attributeRequested)
	attributeResult, err := client.Get(attributeKey).Result()

	if err != nil {

		if err.Error() == "redis: nil" {
			errorMessage := fmt.Sprintf("Couldn't find attribute %s for user %d", attributeRequested, userId)
			return "", errors.New(errorMessage)
		} else {
			return "", err
		}
	}

	return attributeResult, nil
}

// ClientUserSetAttribute sets an attribute in a redis store and returns true on success
func ClientUserSetAttribute(userId int, attributeToSet string, attributeValue string) error {

	exists, err := client.SIsMember("users:userids", userId).Result()

	if exists != true {
		errorMessage := fmt.Sprintf("Couldn't find user %d", userId)
		return errors.New(errorMessage)
	} else if err != nil {
		return err
	}

	attributeKey := fmt.Sprintf("users:%d:%s", userId, attributeToSet)
	_, err = client.Set(attributeKey, attributeValue, 0).Result()

	return err
}
