package main

import (
	"context"
	"log"
	"time"

	"github.com/valkey-io/valkey-go"
)

var valkeyClient valkey.Client

func InitValkey() {
	var err error
	valkeyClient, err = valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		log.Fatalf("Ошибка при подключении к Valkey: %v", err)
	}
}

func ValkeySet(key string, value string, ttlSeconds int) error {
	ctx := context.Background()
	return valkeyClient.Do(ctx, valkeyClient.B().Set().
		Key(key).
		Value(value).
		Ex(time.Duration(ttlSeconds)*time.Second).
		Build()).Error()
}
func ValkeyGet(key string) (string, error) {
	ctx := context.Background()
	res := valkeyClient.Do(ctx, valkeyClient.B().Get().Key(key).Build())
	if res.Error() != nil {
		return "", res.Error()
	}
	return res.ToString()
}

func ValkeyExpire(key string, ttlSeconds int) error {
	ctx := context.Background()
	return valkeyClient.Do(ctx,
		valkeyClient.B().Expire().Key(key).Seconds(int64(ttlSeconds)).Build(),
	).Error()
}

//	func InitValkey() {
//		var err error
//		valkeyClient, err = valkey.NewClient(valkey.ClientOption{
//			Address: "127.0.0.1:6379",
//		})
//		if err != nil {
//			log.Fatalf("Ошибка подключения к Valkey: %v", err)
//		}
//		log.Println("Valkey успешно инициализирован")
//	}
//
//	func ValkeySet(key string, value string) {
//		ctx := context.Background()
//		err := valkeyClient.Do(ctx, valkey.StringCommand("SET", key, value)).Error()
//		if err != nil {
//			log.Printf("Ошибка при SET %s: %v", key, err)
//		}
//	}
// func ValkeyGet(key string) (string, error) {
// 	ctx := context.Background()
// 	res := valkeyClient.Do(ctx, valkey.StringCommand("GET", key))
// 	if res.Error() != nil {
// 		return "", res.Error()
// 	}
// 	val, err := res.ToString()
// 	if err != nil {
// 		return "", err
// 	}
// 	return val, nil
// }
