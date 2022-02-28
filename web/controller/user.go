package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"path"

	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/micro_web/web/model"
	getCaptcha "github.com/micro_web/web/proto/getCaptcha"
	register "github.com/micro_web/web/proto/register"
	user "github.com/micro_web/web/proto/user"
	"github.com/micro_web/web/utils"
)

func GetSession(ctx *gin.Context) {
	// 初始化错误返回的 map
	resp := make(map[string]interface{})
	s := sessions.Default(ctx)
	userName := s.Get("userName")
	if userName == nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		var nameData struct {
			Name string `json:"name"`
		}
		nameData.Name = userName.(string) // 类型断言
		resp["data"] = nameData
	}

	ctx.JSON(http.StatusOK, resp)
}
func DeleteSession(ctx *gin.Context) {
	resp := make(map[string]interface{})

	s := sessions.Default(ctx)
	s.Delete("userName")
	err := s.Save()
	if err != nil {
		resp["errno"] = utils.RECODE_IOERR // 没有合适错误,使用 IO 错误!
		resp["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)

	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	}
	ctx.JSON(http.StatusOK, resp)
}

//getCaptcha微服务-获取图片验证码
func GetImageCd(ctx *gin.Context) {
	// 获取图片验证码 uuid
	uuid := ctx.Param("uuid")

	microService := utils.InitMicro()

	// 初始化客户端
	microClient := getCaptcha.NewGetCaptchaService("getCaptcha", microService.Client())

	// 调用远程函数
	resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{Uuid: uuid})
	if err != nil {
		fmt.Println("未找到远程服务...")
		return
	}

	// 将得到的数据,反序列化,得到图片数据
	var img captcha.Image
	json.Unmarshal(resp.Img, &img)

	// 将图片写出到 浏览器.
	png.Encode(ctx.Writer, img)

	fmt.Println("uuid = ", uuid)
}

//register微服务-根据用户填写的图片验证码是否正确来获取短信验证码
func GetSmscd(ctx *gin.Context) {
	phone := ctx.Param("phone")
	imgCode := ctx.Query("text")
	uuid := ctx.Query("id")

	fmt.Printf("该请求的电话为%s，图片验证码为%s，用户id为%s\n", phone, imgCode, uuid)

	microService := utils.InitMicro()

	microClient := register.NewRegisterService("register", microService.Client())

	rsp, err := microClient.SendSms(context.TODO(), &register.Request{Phone: phone, ImgCode: imgCode, Uuid: uuid})

	if err != nil {
		fmt.Println("调用短信验证码远程服务失败", err)
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

//register微服务-根据用户填写的手机号码和密码注册用户（还需要传短信验证码）
func PostRet(ctx *gin.Context) {
	var regData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	ctx.Bind(&regData)
	microService := utils.InitMicro()

	microClient := register.NewRegisterService("register", microService.Client())

	rsp, err := microClient.Register(context.TODO(), &register.RegReq{
		Mobile:   regData.Mobile,
		SmsCode:  regData.SmsCode,
		Password: regData.PassWord,
	})

	if err != nil {
		fmt.Println("注册用户时没有找到远程服务", err)
		return
	}
	s := sessions.Default(ctx)
	s.Set("userName", rsp.Name)
	s.Save()
	ctx.JSON(http.StatusOK, rsp)
}

//register微服务-根据用户填写的手机号和密码登录账号
func PostLogin(ctx *gin.Context) {
	type userData struct {
		Password string `json:"password"`
		Mobile   string `json:"mobile"`
	}
	var ud userData
	ctx.Bind(&ud)
	microService := utils.InitMicro()
	microClient := register.NewRegisterService("register", microService.Client())
	rsp, _ := microClient.Login(context.TODO(), &register.LoginReq{
		Mobile:   ud.Mobile,
		Password: ud.Password,
	})
	resp := make(map[string]interface{})
	if rsp.Errno == utils.RECODE_DBERR {
		fmt.Println("调用远程服务register的登录方法失败")
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR)
	} else if rsp.Errno == utils.RECODE_OK {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
		s := sessions.Default(ctx)
		s.Set("userName", rsp.Name)
		s.Save()
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetArea(ctx *gin.Context) {
	var areas []model.Area

	model.GlobalConn.Find(&areas)
	// fmt.Println(areas)
	conn := model.RedisPool.Get()
	areaData, _ := redis.Bytes(conn.Do("get", "areaData"))
	if len(areaData) == 0 {
		// fmt.Println("redis中目前没有数据，从mysql中拿数据中")
		model.GlobalConn.Find(&areas)
		areaBuf, _ := json.Marshal(areas)
		conn.Do("set", "areaData", areaBuf)
	} else {
		// fmt.Println("redis中有数据，从redis中拿数据中")
		json.Unmarshal(areaData, &areas)
	}

	rsp := make(map[string]interface{})

	rsp["errno"] = "0"
	rsp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	rsp["data"] = areas

	ctx.JSON(http.StatusOK, rsp)
}

//user微服务-获取用户详细信息
func GetUserInfo(ctx *gin.Context) {

	s := sessions.Default(ctx)
	userName := s.Get("userName")

	microService := utils.InitMicro()

	microClient := user.NewUserService("user", microService.Client())

	rsp, err := microClient.MicroGetUser(context.TODO(), &user.Request{
		Name: userName.(string),
	})
	if err != nil {
		fmt.Println("调用远程user服务错误", err)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	}

	ctx.JSON(http.StatusOK, rsp)
}

//user微服务-更新用户姓名
func PutUserInfo(ctx *gin.Context) {
	s := sessions.Default(ctx) // 初始化Session 对象
	userName := s.Get("userName")
	var nameData struct {
		Name string `json:"name"`
	}
	ctx.Bind(&nameData)

	microService := utils.InitMicro()

	microClient := user.NewUserService("user", microService.Client())

	resp, _ := microClient.UpdateUserName(context.TODO(), &user.UpdateReq{
		NewName: nameData.Name,
		OldName: userName.(string),
	})
	//更新session数据
	if resp.Errno == utils.RECODE_OK {
		//更新成功,session中的用户名也需要更新一下
		s.Set("userName", nameData.Name)
		s.Save()
	}

	ctx.JSON(http.StatusOK, resp)
}

//user微服务-上传用户头像
func PostAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		fmt.Println("获取avatar图片失败", err)
	}
	f, _ := file.Open()
	buf := make([]byte, file.Size)
	f.Read(buf)
	fileExt := path.Ext(file.Filename)
	microService := utils.InitMicro()

	microClient := user.NewUserService("user", microService.Client())
	userName := sessions.Default(ctx).Get("userName")
	resp, _ := microClient.UploadAvatar(context.TODO(), &user.UploadReq{
		Avatar:   buf,
		UserName: userName.(string),
		FileExt:  fileExt,
	})

	ctx.JSON(http.StatusOK, resp)
}

type AuthStu struct {
	IdCard   string `json:"id_card"`
	RealName string `json:"real_name"`
}

//user微服务-实名认证
func PutUserAuth(ctx *gin.Context) {
	fmt.Println("jinru auth")
	//获取数据
	var auth AuthStu
	err := ctx.Bind(&auth)
	//校验数据
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}
	microService := utils.InitMicro()

	microClient := user.NewUserService("user", microService.Client())
	userName := sessions.Default(ctx).Get("userName")
	fmt.Printf("username: %s", userName)
	resp, _ := microClient.AuthUpdate(context.TODO(), &user.AuthReq{
		IdCard:   auth.IdCard,
		RealName: auth.RealName,
		UserName: userName.(string),
	})
	ctx.JSON(http.StatusOK, resp)
}
