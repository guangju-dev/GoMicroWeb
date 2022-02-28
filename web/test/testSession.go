package test

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func testSession() {
	router := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "172.19.0.1:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.GET("/test", func(c *gin.Context) {
		s := sessions.Default(c)
		s.Set("key1", "value1")
		s.Save()
		c.Writer.WriteString("测试Session")
	})
	router.Run(":9998")
}
