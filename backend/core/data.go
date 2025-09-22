package main

import (
	"fmt"
	"time"
	"wechat-secretary/backend/core/service"
	"wechat-secretary/backend/wechat"
)

func main() {
	manager := wechat.NewManager()
	manager.Init()
	keyService := service.NewKeyService(manager)

	keyService.Init()
	// err := keyService.DecryptKey()
	account := keyService.GetWxAccount()[0].Name

	coreService := service.NewCoreService(keyService)
	coreService.Decrypt(account)

	<-time.After(30 * time.Minute)
	fmt.Println("err")
}
