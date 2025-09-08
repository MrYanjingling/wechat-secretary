package v4

import (
	"context"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	source, err := New("D:\\wx")
	if err != nil {
		fmt.Println(err)
	}

	contacts, err := source.GetContacts(context.Background(), "", 20, 0)

	if err != nil {
		fmt.Println(err)
	}
	for _, contact := range contacts {
		fmt.Println(contact)
	}
}
