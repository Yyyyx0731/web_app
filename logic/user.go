package logic

// 业务逻辑层

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp)(err error){
	//1.判断用户是否存在
	//var exist bool
	if err = mysql.CheckUserExist(p.Username);err!=nil {
		return err
	}
	if err!=nil {//查询出错
		return err
	}

	//2.生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//3.保存进数据库
	return  mysql.InsertUser(user)
}


func Login(p *models.ParamLogin) (user *models.User,err error){
	//user是指针类型
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//return mysql.Login(user)//返回的结果就是登录结果
	if err = mysql.Login(user);err != nil {
		return nil,err
	}
	//mysql校验正确后
	// GenToken 生成JWT
	token,err :=  jwt.GenToken(user.UserID,user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
