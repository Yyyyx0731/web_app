package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePostHandler(t *testing.T){
	//设置一下测试的模式
	gin.SetMode(gin.TestMode)
	//自己注册一个路由 不用原有routes包里的，因为会造成循环引包
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url,CreatePostHandler)

	body := `{
		"community_id":1,
		"title":"test",
		"content":"just a test"
	}`
	//第三个参数body是io.Reader类型
	//bytes.NewReader返回的是一个Reader对象，但注意要传字节切片类型
	req,_ := http.NewRequest(http.MethodPost,url,bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w,req)

	//判断状态码是不是200
	assert.Equal(t,200,w.Code)


	//assert.Equal(t,"pong",w.Body.String())
	//判断响应的内容是不是按预期返回了需要登录的错误
	//因为传参数时没有jwt token，所以在CreatePostHandler里
	//getCurrentUserID(c)时，从c获取当前发送请求的用户id报“需要登录”的错
	//在第三个参数 w拿到json格式的字符串

	//方法一：所以需要解析,判断w.Body.String()字符串是不是包含“需要登录”
	//assert.Contains(t,w.Body.String(),"需要登录")

	//方法二：将相应的内容反序列化到res 然后判断与预期是否一致
	res := new(ResponseData)
	//反序列化 json->object
	if err := json.Unmarshal(w.Body.Bytes(),res);err != nil {
		t.Fatalf("json.Unmarshal w.Body failed,err:%v\n",err)
	}
	assert.Equal(t, res.Code,CodeNeedLogin)
}
