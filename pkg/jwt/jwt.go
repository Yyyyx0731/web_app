package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)


// TokenExpireDuration 定义JWT的过期时间，这里以365天为例
const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("夏天夏天悄悄过去")


// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID int64 `json:"userID"`
	Username string `json:"username"`
	jwt.StandardClaims
}


// GenToken 生成JWT
func GenToken(userID int64,username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID,
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire"))*time.Hour).Unix(), // 过期时间
			Issuer:    "web_app",                               // 签发人
		},
		//jwt.StandardClaims{
		//	ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
		//	Issuer:    "web_app",                               // 签发人
		//},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c) //用SigningMethodHS256算法加密c
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//// 解析token
	////把 tokenString 解析到变量 MyClaims，第三个参数是个函数 告诉他用什么去解
	//token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
	//	return MySecret, nil
	//})
	//if err != nil { //解析出错 返回错误
	//	return nil, err
	//}
	////token.Claims转成 *MyClaims 类型
	//if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
	//	return claims, nil
	//}
	//return nil, errors.New("invalid token")

	//无论返回值有没有声明名为mc的此变量，都要new一个，否则程序执行有问题
	//因为在返回值处声明的变量不会自动申请内存，所以需要手动
	var mc = new(MyClaims)
	//fmt.Println(tokenString)
	token,err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil { //解析出错 返回错误
		//fmt.Println("1")
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, nil  //校验结果为有效
	}
	//无效
	//fmt.Println("2")
	return nil, errors.New("invalid token") //错误：无效的token
}