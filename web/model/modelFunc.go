package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

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
func GetUserInfo(userName string) (User, error) {
	var user User
	err := GlobalConn.Where("name = ?", userName).First(&user).Error

	return user, err
}

func UpdateUserName(newName, oldName string) error {
	// update user set name = 'itcast' where name = 旧用户名
	return GlobalConn.Model(new(User)).Where("name = ?", oldName).Update("name", newName).Error
}
func UpdateAvatar(userName, avatar string) error {
	// update user set avatar_url = avatar, where name = username
	return GlobalConn.Model(new(User)).Where("name = ?", userName).
		Update("avatar_url", avatar).Error
}
