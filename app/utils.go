package app

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	SecretKey 	    []byte
	VirusTotalAPI 	string
	Ctx 			= context.Background()
	RDB 			*redis.Client
)

func must(errType string, err error) {
	if err != nil {
		log.Fatal(errType, err)
	}
}

func DotEnvInit() {
	err := godotenv.Load()
	must("[ERROR]DOTENV-LOADING:", err)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
	if len(SecretKey) == 0 {
		log.Fatal("[ERROR]DOTENV-SECRET_KEY: empty value")
	}

	VirusTotalAPI = string(os.Getenv("VIRUS_TOTAL_API"))
	if (len(SecretKey) == 0) {
		log.Fatal("[ERROR]DOTENV-VIRUS_TOTAL_API: empty value")
	}
}

func RedisInit() {
	RDB = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
  })
	must("[ERROR]REDIS-PING:", RDB.Ping(Ctx).Err())
}

func AddClient(id string, client ClientDBO) error {
  data, err := json.Marshal(client)
  
	if err != nil {
		return err
	}
  if err := RDB.Set(Ctx, "name:"+client.Name, id, 0).Err(); err != nil {
		return err
	}
	if err := RDB.Set(Ctx, "client:"+id, data, 0).Err(); err != nil {
		return err
	}
	return  nil
}

func GetClient(id string, clientDBO *ClientDBO) error {
  val, err := RDB.Get(Ctx, "client:"+id).Result()
  if err != nil {
    return err
  }
  return json.Unmarshal([]byte(val), clientDBO)
}

func UpdateClient(id string, data *ClientDBO) error {
    var old ClientDBO
    err := GetClient(id, &old)
    if err != nil && err != redis.Nil {
        return err
    }

    if old.Name != "" && old.Name != data.Name {
        RDB.Del(Ctx, "name:"+old.Name)
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    if err := RDB.Set(Ctx, "client:"+id, jsonData, 0).Err(); err != nil {
        return err
    }
		
    if data.Name != "" {
        if err := RDB.Set(Ctx, "name:"+data.Name, id, 0).Err(); err != nil {
            return err
        }
    }

    return nil
}

func RemoveClient(id string) error {
    var client ClientDBO

    val, err := RDB.Get(Ctx, "client:"+id).Result()
    if err != nil {
        if err == redis.Nil {
            return nil
        }
        return err
    }

    if err := json.Unmarshal([]byte(val), &client); err != nil {
        return err
    }

    if client.Name != "" {
        if _, err := RDB.Del(Ctx, "name:"+client.Name).Result(); err != nil {
            fmt.Println("[WARN] failed to delete name index:", err)
        }
    }

    if _, err := RDB.Del(Ctx, "client:"+id).Result(); err != nil {
        return err
    }

    return nil
}



func IsNameInDB(name string) bool {
  exists, _ := RDB.Exists(Ctx, "name:"+name).Result()
  return exists == 1
}

func ClearDB()  {
	err := RDB.FlushDB(Ctx).Err()
	if err != nil {
	  log.Println("Redis flush error:", err)
	}
}
