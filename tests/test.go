package main

import (
	"fmt"
	"strings"
)

func test1 () {
	if strings.Contains("线上回归测试", "sb") {
		fmt.Print("ok")
	}
	fmt.Println("fin")

}

func main (){
	test1()
}