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
	err := keyService.DecryptKey()

	fmt.Println(err)
}
