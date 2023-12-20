package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"web_app/models"
)

const secret = "liwenzhou.com"
//var (
//	ErrorUserExist = errors.New("用户已存在")
//	ErrorUserNotExist = errors.New("用户不存在")
//	ErrorInvalidPassword = errors.New("用户名或密码错误")
//)

// CheckUserExist 查询该用户名的用户是否存在
func CheckUserExist(username string)(err error){
	sqlStr := "select count(user_id) from user where username=?"
	var count int
	if err=db.Get(&count,sqlStr,username);err!=nil { //有error
		return err
	}
	if count >0 { //count>0说明查到了该条数据，即存在
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中加入一条新的用户记录
func InsertUser(user *models.User)(err error){
	//对密码进行加密
	password := encryptPassword(user.Password)
	sqlStr := `insert into user(user_id, username, password) values (?,?,?)`
	_,err = db.Exec(sqlStr,user.UserID,user.Username,password)
	return
}

// encryptPassword  md5密码加密
func encryptPassword(oPassword string)string{
	h := md5.New()
	h.Write([]byte(secret))
	//h.Sum([]byte(oPassword))
	//把h.Sum的返回值转换成16进制的字符串
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login  登录，判断用户名和密码是否正确
func Login(user *models.User)(err error){
	oPassword := user.Password//oPassword用户输入的密码
	sqlStr := "select user_id,username,password from user where username=?"
	err = db.Get(user,sqlStr,user.Username)	//Get(将sql语句执行结果读到此变量,sql语句,多个参数来填占位符?...)
	//用户不存在
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	//查询数据库失败
	if err != nil {
		return err
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)//把原始用户输入的密码加密，准备与数据库中加密的对比
	//现在的user结构体的内容，都是从数据库中读出来的
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

func GetUserByID(uid int64)(user *models.User,err error){
	user = new(models.User)
	sqlStr := "select user_id,username from user where user_id=?"
	err = db.Get(user,sqlStr,uid)
	return
}