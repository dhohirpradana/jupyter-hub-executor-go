package helper

import (
	"fmt"
	"testing"
)

func TestPocketbase(t *testing.T) {
	token, err := GetToken()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Token:", token)
}
