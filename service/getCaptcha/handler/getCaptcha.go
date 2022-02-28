package handler

import (
	"context"
	"encoding/json"
	"image/color"

	"getCaptcha/model"
	getCaptcha "getCaptcha/proto"

	"github.com/afocus/captcha"
	"github.com/asim/go-micro/v3/util/log"
)

type GetCaptcha struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetCaptcha) Call(ctx context.Context, req *getCaptcha.Request, rsp *getCaptcha.Response) error {

	cap := captcha.New()

	if err := cap.SetFont("/micro_web/service/getCaptcha/comic.ttf"); err != nil {
		log.Info("没有字体文件")
		panic(err.Error())
	}
	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	//SetFrontColor(colors ...color.Color)  这两个颜色设置的函数属于不定参函数
	cap.SetFrontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{G: 153, A: 255})

	img, str := cap.Create(4, captcha.NUM)
	err := model.SaveImgCode(str, req.Uuid)
	if err != nil {
		return err
	}
	imgBuf, _ := json.Marshal(img)

	rsp.Img = imgBuf

	return nil
}
