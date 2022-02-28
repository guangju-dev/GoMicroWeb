package test

import "github.com/gin-gonic/gin"

func testCookie() {
	router := gin.Default()

	router.GET("/test", func(context *gin.Context) {
		context.SetCookie("key", "value", 60*60, "", "", true, true)
	})
	router.Run(":9999")
}
