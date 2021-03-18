package controllers

import (
	"fmt"
	"go_autoapi/models"
)

func (c *UserController) profile() {
	fmt.Println("xxxxxx", "profile")
	amc := models.AdMockCase{}
	data, _ := amc.QueryByUUid("5eb62275-9818-4101-a477-6fef0bb9c7bd")
	fmt.Println(&data)
	c.SuccessJson(&data)
}
