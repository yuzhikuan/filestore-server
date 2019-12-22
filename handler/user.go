package handler

import (
	"fmt"
	dblayer "github.com/yuzhikuan/filestore-server/db"
	"github.com/yuzhikuan/filestore-server/util"
	"io/ioutil"
	"net/http"
	"time"
)

const pwdSalt = "*#890"
const tokenSalt = "_tokensalt"

// SignupHandler 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// 1.http GET 请求，直接返回注册页面内容
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	// 2.校验参数的有效性
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}
	// 3.加密用户密码
	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	// 4.存入数据库tbl_user表及返回结果
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// SigninHandler 登录接口
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	// 1.校验用户名及密码
	encPasswd := util.Sha1([]byte(password + pwdSalt))
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}
	// 2.生成访问凭证（token）
	token := GenToken(username)
	// 3.存储token到数据库tbl_user_token表
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}
	// 4.登录成功后返回username, token, 重定向url等信息
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

// UserInfoHandler 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	// 2.验证token是否有效
	//isValidToken := IsTokenValid(token)
	//if !isValidToken {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	// 3.查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 4.组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

// GenToken 产生一个40位的token
func GenToken(username string) string {
	// 40位字符：md5(username + timestamp + tokenSalt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + tokenSalt))
	return tokenPrefix + ts[:8]
}

// IsTokenValid 验证token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// 1.判断token的时效性，是否过期
	// 2.从数据库表tbl_user_token查询username对应的token信息
	// 3.对比两个token是否一致
	return true
}
