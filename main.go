package main

import (
	"fmt"

	"github.com/senaphim/aggreGator-go/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("Encountered error reading config:\n%v", err))
	}

	fmt.Println(
		fmt.Sprintf("DB_URL: %v, Current_user_name: %v", conf.DbUrl, conf.CurrentUserName),
	)
}
