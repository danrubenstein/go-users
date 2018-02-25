package main 

import 

(
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var client *redis.Client;	


func init(){

	InitializeRedis()
}

// Resets the redis client
func InitializeRedis(){
	
	client = NewClient(os.Getenv("REDIS_URL"));
	
	client.Set("users:nextuserid", "0", 0).Result(); 
	client.Del("users:usernames");
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

	val, err := client.SIsMember("users:usernames", u.Username).Result(); 
	if err != nil { 
		fmt.Println("There was an error checking the username list")
		return err 
	}

	if val { 
		return errors.New("The username was already taken")
	}

	_, err = client.SAdd("users:usernames", u.Username).Result(); 
	if err != nil { 
		fmt.Println("There was an error adding the username to the set")
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

		if err.Error() == "redis: nil" { 
			return emptyUser, errors.New(fmt.Sprintf("This user id (%d) couldn't be found!", userId)); 
		}
		fmt.Println("Couldn't find the user here"); 
		return emptyUser, err
	}

	userNameKey := fmt.Sprintf("users:%d:username", userId); 
	userName, err := client.Get(userNameKey).Result(); 
	
	if err != nil { 
		return emptyUser, err 
	}


	return User{userUsername, userName, "", ""}, nil 

}



