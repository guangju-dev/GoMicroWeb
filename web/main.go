package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/micro_web/web/controller"
	"github.com/micro_web/web/model"
	"github.com/micro_web/web/utils"
)

//在用户操作自己的数据前验证session
func LoginFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.Default(ctx)
		userName := s.Get("userName")
		resp := make(map[string]interface{})
		if userName == nil {
			resp["errno"] = utils.RECODE_SESSIONERR
			resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
			ctx.JSON(http.StatusOK, resp)
			ctx.Abort() // 从这里返回, 不必继续执行了
		} else {
			ctx.Next() // 继续向下
		}
	}
}
func main() {
	//初始化数据库
	model.InitDb()
	//初始化Redis
	model.InitRedis()

	//初始化路有
	router := gin.Default()
	//初始化容器，存储session数据
	store, _ := redis.NewStore(10, "tcp", "172.19.0.1:6379", "", []byte("secret"))
	//使用容器

	//做路由匹配
	router.Static("/home", "view")
	r1 := router.Group("/api/v1.0")
	{

		r1.GET("/areas", controller.GetArea)
		r1.GET("/imagecode/:uuid", controller.GetImageCd)
		r1.GET("/smscode/:phone", controller.GetSmscd)

		r1.Use(sessions.Sessions("mysession", store))
		r1.POST("/users", controller.PostRet)
		//登录业务
		r1.GET("/session", controller.GetSession)
		r1.POST("/sessions", controller.PostLogin)
		//路有过滤器，只有在登录的情况下才能执行以下路有请求
		r1.Use(LoginFilter())
		r1.DELETE("/session", controller.DeleteSession)
		r1.GET("/user", controller.GetUserInfo)
		r1.PUT("/user/name", controller.PutUserInfo)

		r1.POST("/user/avatar", controller.PostAvatar)
		r1.POST("/user/auth", controller.PutUserAuth)
		r1.GET("/user/auth", controller.GetUserInfo)

		// //获取已发布房源信息
		// r1.GET("/user/houses", controller.GetUserHouses)
		// //发布房源
		// r1.POST("/houses", controller.PostHouses)
		// //添加房源图片
		// r1.POST("/houses/:id/images", controller.PostHousesImage)
		// //展示房屋详情
		// r1.GET("/houses/:id", controller.GetHouseInfo)
		// //展示首页轮播图
		// r1.GET("/house/index", controller.GetIndex)
		// //搜索房屋
		// r1.GET("/houses", controller.GetHouses)
		// //下订单
		// r1.POST("/orders", controller.PostOrders)
		// //获取订单
		// r1.GET("/user/orders", controller.GetUserOrder)
		// //同意/拒绝订单
		// r1.PUT("/orders/:id/status", controller.PutOrders)
	}

	//3.启动运行
	router.Run(":14000")

}
