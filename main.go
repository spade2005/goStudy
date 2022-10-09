package main

import (
	"fmt"
	"xstart/app/common"
	"xstart/app/common/router"
)

func main() {
	fmt.Println("hello start")
	config := common.ConfigObj
	val, _ := config.GetString("db", "driver")
	fmt.Printf("?", val)
	val1, _ := config.GetString(common.DefaultSection, "hi")
	fmt.Printf("?", val1)
	val2, _ := config.GetString("default", "you")
	fmt.Printf("?", val2)

	router.Run()
}
