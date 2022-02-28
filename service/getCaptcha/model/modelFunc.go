package model

import (
	"github.com/asim/go-micro/v3/util/log"
	"github.com/gomodule/redigo/redis"
)

func SaveImgCode(code, uuid string) error {
	conn, err := redis.Dial("tcp", "172.19.0.1:6379")
	if err != nil {
		log.Error("没有连接到redis数据库")
		return err
	}
	defer conn.Close()

	_, err = conn.Do("setex", uuid, 60*5, code)
	if err != nil {
		log.Error("设置过期时间出错")
	}
	return err
}
