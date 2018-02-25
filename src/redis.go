package main 

import 

(
	"errors"
	"fmt"

	"github.com/go-redis/redis"

)

var client *redis.Client;	

func init(){
	InitializeRedis("localhost:6379")
}

// Resets the redis client
func InitializeRedis(redisAddress string){
	client = NewClient(redisAddress)
	client.Set("users:nextuserid", "0", 0).Result(); 
	client.Del("users:usernames");
}

func NewClient(address string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

// ClientUserCreate creates a user entry if the username is not taken
// and increments the userid counter
func ClientUserCreate(u User) error { 

	val, err := client.SIsMember("users:usernames", u.Username).Result(); 
	if err != nil { 
		return err 
	}

	if val { 
		return errors.New("The username was already taken")
	}

	_, err = client.SAdd("users:usernames", u.Username).Result(); 
	if err != nil { 
		return err 
	}
 
	intVal, err := client.Incr("users:nextuserid").Result(); 
	if err != nil { 
		return  err 
	}

	userUsernameKey := fmt.Sprintf("users:%d:username", intVal - 1)
	_, err = client.Set(userUsernameKey, u.Username, 0).Result()
	if err != nil { 
		return err
	}

	userNameKey := fmt.Sprintf("users:%d:name", intVal - 1)
	_, err = client.Set(userNameKey, u.Name, 0).Result()
	if err != nil { 
		return err
	}

	return nil;
}

func ClientUserInfo(userId int) (User, error)  { 

	emptyUser := User{"", "", "", ""}

	userUsernameKey := fmt.Sprintf("users:%d:username", userId); 
	userUsername, err := client.Get(userUsernameKey).Result(); 
	if err != nil { 
		return emptyUser, err
	}

	userNameKey := fmt.Sprintf("users:%d:username", userId); 
	userName, err := client.Get(userNameKey).Result(); 
	if err != nil { 
		return emptyUser, err 
	}


	return User{userUsername, userName, "", ""}, nil 

}



