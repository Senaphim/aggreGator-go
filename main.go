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

	conf.SetUser("jonah")
	fmt.Println(fmt.Sprintf("Current user name: %v", *conf.CurrentUserName))

	conf, _ = config.Read()
	conf.SetUser("steve")
	fmt.Println(fmt.Sprintf("Db Url: %v, current username: %v", *conf.DbUrl, *conf.CurrentUserName))
}
