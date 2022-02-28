package model

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
func SaveRealName(userName, realName, idCard string) error {
	return GlobalConn.Model(new(User)).Where("name = ?", userName).
		Updates(map[string]interface{}{"real_name": realName, "id_card": idCard}).Error
}
