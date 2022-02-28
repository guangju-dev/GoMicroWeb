package model

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

//直接用创建一个实例，这样所有的函数就都可以用
var RedisPool redis.Pool

//初始化Redis连接池
func InitRedis() {
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "172.19.0.1:6379")
		},
	}
}

//检验图片验证码
func CheckImgCode(uuid, imgCode string) bool {
	// 链接 redis
	conn := RedisPool.Get()
	defer conn.Close()
	fmt.Printf("当前检验的用户为%s,图片验证码为%s\n", uuid, imgCode)
	// 查询 redis 数据
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误 err:", err)
		return false
	}

	// 返回校验结果
	return code == imgCode
}

//将用户的手机号和短信验证码放进redis数据库

func SaveSmsCode(phone, code string) error {
	conn := RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("setex", phone+"_code", 60*3, code)

	return err
}

//检验短信验证码
func CheckSmsCode(phone, code string) error {
	conn := RedisPool.Get()

	smsCode, err := redis.String(conn.Do("get", phone+"_code"))

	if err != nil {
		fmt.Println("没有找到该手机号对应的验证码", err)
	}

	if smsCode != code {
		return errors.New("验证码匹配失败！")
	}
	return nil
}

//注册用户信息，写入MySql 数据库

func RegisterUser(mobile, pwd string) (string, error) {
	var user User
	user.Name = mobile
	user.Mobile = mobile
	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum(nil))

	user.Password_hash = pwd_hash
	return user.Name, GlobalConn.Create(&user).Error
}

func Login(mobile, pwd string) (string, error) {
	var user User
	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum(nil))
	err := GlobalConn.Where("mobile = ?", mobile).
		Select("name").Where("password_hash = ?", pwd_hash).Find(&user).Error
	fmt.Printf("当前登陆的用户为%s", user.Name)
	return user.Name, err
}
