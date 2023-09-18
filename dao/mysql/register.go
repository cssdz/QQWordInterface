package mysql

import (
	"errors"
	"go.uber.org/zap"
	"go_web/models"
)

// CheckUser 检查用户是否注册
func CheckUser(id string) error {
	var checkId string
	sqlStr := "SELECT `id` FROM `word`.user_info WHERE `id` = ?"
	err := db.Get(&checkId, sqlStr, id)
	if err != nil {
		zap.L().Warn("CheckUser exec err", zap.Error(err), zap.String("id", id))
		return errors.New("充次错误")
	}
	return nil
}

// RegisterUser 注册用户信息
func RegisterUser(register models.Register) error {
	sqlStr := "INSERT INTO `word`.user_info (`id`, `name`, `num`, `update_time`, `register_time`) VALUES (?, ?, ?, CURRENT_TIME(), CURRENT_TIME())"
	_, err := db.Exec(sqlStr, register.ID, register.Name, register.Num)
	if err != nil {
		zap.L().Error("RegisterUser1 exec err", zap.Error(err), zap.String("id", register.ID))
		return errors.New("用户注册失败")
	}

	zap.L().Info("用户创建成功", zap.String("id", register.ID))
	return nil
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(register models.Register) error {
	if register.Num == -1 && register.Name != "" { // 只更新名字
		sqlStr := "UPDATE `word`.`user_info` SET  `name` = ? WHERE `id` = ?"
		_, err := db.Exec(sqlStr, register.Name, register.ID)
		if err != nil {
			zap.L().Error("RegisterUser1 exec err", zap.Error(err), zap.String("id", register.ID))
			return errors.New("更新名字失败")
		}
	} else if register.Num != -1 && register.Name == "" { // 只更新次数
		sqlStr := "UPDATE `word`.`user_info` SET  `num` = ? WHERE `id` = ?"
		_, err := db.Exec(sqlStr, register.Num, register.ID)
		if err != nil {
			zap.L().Error("RegisterUser2 exec err", zap.Error(err), zap.String("id", register.ID))
			return errors.New("更新次数失败")
		}
	} else { // 名字和次数都更新
		sqlStr := "UPDATE `word`.`user_info` SET  `num` = ?, `name` = ? WHERE `id` = ?"
		_, err := db.Exec(sqlStr, register.Num, register.Name, register.ID)
		if err != nil {
			zap.L().Error("RegisterUser3 exec err", zap.Error(err), zap.String("id", register.ID))
			return errors.New("更新名字和次数失败")
		}
	}
	zap.L().Info("数据更新成功", zap.String("id", register.ID))
	return nil
}
