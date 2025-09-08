package main

import (
	"fmt"
	"wechat-secretary/backend/wechat/repo"
)

func main() {

	contactRepo, err := repo.NewContactRepo("D:\\wx2")
	if err != nil {
		fmt.Println(err)
	}

	username := contactRepo.GetContactByNickName("wxmbzshhhh")

	fmt.Println(username)
}
