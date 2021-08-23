package redisservice

import (
	// "fmt"
	// "net/http"
	"errors"
	"thinhtd4/config"
	"thinhtd4/customlog"
	"thinhtd4/model"
	"time"

	"github.com/go-redis/redis/v7"
	"golang.org/x/crypto/bcrypt"
)

var client *redis.Client

func Init() {
	customlog.Info("Start initialize redis connect")
	client = redis.NewClient(&redis.Options{
		Addr: config.RedisConfig().Host + ":" + config.RedisConfig().Port, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SaveTokenToRedis(c model.CookieResp) error {
	status := client.Set(c.UserID, "0", time.Duration(c.TimeOut)*time.Second).Err()
	if status != nil {
		customlog.Err(errors.New("Set token error : " + status.Error()))
		return status
	}
	return nil
}

func IsTokenInRedis(userid string) (bool, error) {
	err := client.Get(userid).Err()
	if err != nil {
		customlog.Err(errors.New("Set token error : " + err.Error()))
		return false, err
	}
	return true, nil
}
