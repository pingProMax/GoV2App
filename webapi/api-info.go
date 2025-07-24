package webapi

type NodeInfo struct {
	Id        int
	Name      string
	Agreement string
}

var (
	ApiUrl         string = "https://127.0.0.1:8080"
	loginApi       string = "/login"                // 登录api
	getUserInfoApi string = "/user"                 // 用户信息api
	subApi         string = "/api/subscribe?token=" // 订阅api
	appBulApi      string = "/user/app_bulletin"    // 公告api

	//--------------
	// jwt   string    = ""
	// jwtEx time.Time      //jwt 到期时间

	// token string = "" //订阅用

	V string = "v 1.1"

	//注册地址
	RegisterUrl string = "https://xxxxx/register"
	//购买套餐地址
	PlanUrl string = "https://xxxxx/user/plan"
)
