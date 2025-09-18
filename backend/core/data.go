package main

import (
	"fmt"
	"wechat-secretary/backend/core/service"
	"wechat-secretary/backend/wechat"
)

func main() {
	manager := wechat.NewManager()
	manager.Init()
	keyService := service.NewKeyService(manager)

	keyService.Init()
	// err := keyService.DecryptKey()
	account := keyService.GetAccounts()[0].WeChatAccountKey.Account

	coreService := service.NewCoreService(keyService)
	coreService.Decrypt(account)

	fmt.Println("err")
}
