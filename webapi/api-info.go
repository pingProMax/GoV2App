package webapi

type NodeInfo struct {
	Id        int
	Name      string
	Agreement string
}

var (
	ApiUrl         string = "http://127.0.0.1:8080"
	loginApi       string = "/login"                // 登录api
	getUserInfoApi string = "/user"                 // 用户信息api
	subApi         string = "/api/subscribe?token=" // 订阅api
	appBulApi      string = "/user/app_bulletin"    // 公告api

	//--------------
	// jwt   string    = ""
	// jwtEx time.Time      //jwt 到期时间

	// token string = "" //订阅用

	V string = "v 1.0"
)
