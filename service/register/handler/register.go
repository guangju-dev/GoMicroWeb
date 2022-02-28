package handler

import (
	"context"
	"math/rand"
	"strconv"

	"register/model"
	register "register/proto"
	"register/utils"

	"github.com/asim/go-micro/v3/logger"
)

type Register struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Register) SendSms(ctx context.Context, req *register.Request, rsp *register.Response) error {
	//判断用户输入的验证码是否正确
	result := model.CheckImgCode(req.Uuid, req.ImgCode)
	if result {
		code := strconv.Itoa(rand.Intn(100))
		logger.Infof("手机号%s——短信验证码%s", req.Phone, code)
		err := model.SaveSmsCode(req.Phone, code)
		if err != nil {
			logger.Error("存储短信验证码到redis失败")
			rsp.Errno = utils.RECODE_DATAERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		} else {
			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
		}
	} else {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	}
	return nil
}

func (e *Register) Register(ctx context.Context, req *register.RegReq, rsp *register.RegResponse) error {
	err := model.CheckSmsCode(req.Mobile, req.SmsCode)
	if err == nil {
		name, err := model.RegisterUser(req.Mobile, req.Password)
		logger.Infof("当前注册人的手机号为%s，设置的密码是%s", req.Password, req.Password)
		if err != nil {
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		} else {
			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
			rsp.Name = name
		}
	} else {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	}

	return nil
}

func (e *Register) Login(ctx context.Context, req *register.LoginReq, rsp *register.LoginResponse) error {
	name, err := model.Login(req.Mobile, req.Password)
	if err != nil {
		logger.Error("登录方法-数据库操作失败")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	} else {
		rsp.Name = name
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	}
	return nil
}
