package main

import (
	"context"
	"fmt"
	"wechat-secretary/backend/wechatdb/repo"
)

func main() {

	contactRepo, err := repo.NewContactRepo("D:\\wx2")
	if err != nil {
		fmt.Println(err)
	}

	username := contactRepo.GetContactByNickName("wxmbzshhhh")

	sessionRepo, err := repo.NewSessionRepo("D:\\wx2", contactRepo)
	if err != nil {
		fmt.Println(err)
	}

	sessions := sessionRepo.GetSessions(context.Background(), "", 30, 0)

	for _, session := range sessions {
		fmt.Println(session)
	}

	fmt.Println(username)
}
