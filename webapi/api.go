package webapi

import (
	"GoV2App/nodep"
	"encoding/base64"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/tidwall/gjson"
)

// 登录
func Login(uname, passwd string) string {
	if r, err := g.Client().Post(
		gctx.New(),
		ApiUrl+loginApi,
		g.Map{
			"UserName": uname,
			"Passwd":   passwd,
		},
	); err != nil {
		return err.Error()
	} else {
		defer r.Close()
		jsonStr := r.ReadAllString()

		fmt.Println(jsonStr)

		if gjson.Get(jsonStr, "code").Int() != 0 {
			return gjson.Get(jsonStr, "message").String()
		} else {

			DB.Set("jwt", gjson.Get(jsonStr, "data.token").String())
			DB.Set("jwtEx", gjson.Get(jsonStr, "data.expire").Time().String())
		}
	}

	return ""

}

func GetgClient() *gclient.Client {
	c := g.Client()
	v := DB.Get("jwt")
	c.SetCookie("jwt", string(v))
	return c

}

// 获取订阅token
func GetSubscribeToken() string {
	c := GetgClient()
	if r, err := c.Post(
		gctx.New(),
		ApiUrl+getUserInfoApi,
	); err != nil {
		return err.Error()
	} else {
		defer r.Close()
		jsonStr := r.ReadAllString()

		fmt.Println(jsonStr)

		if gjson.Get(jsonStr, "code").Int() != 0 {
			return gjson.Get(jsonStr, "message").String()
		} else {
			DB.Set("token", gjson.Get(jsonStr, "data.token").String())
			DB.Set("transfer_enable", gjson.Get(jsonStr, "data.transfer_enable").String())
			DB.Set("u", gjson.Get(jsonStr, "data.u").String())
			DB.Set("d", gjson.Get(jsonStr, "data.d").String())
			DB.Set("expired_at", gjson.Get(jsonStr, "data.expired_at").String())
			DB.Set("plan_name", gjson.Get(jsonStr, "data.plan_name").String())
			DB.Set("user_name", gjson.Get(jsonStr, "data.user_name").String())
		}
	}

	return ""
}

// 获取节点列表
func GetNodeInfo() string {
	c := GetgClient()
	token := DB.Get("token")
	if r, err := c.Get(
		gctx.New(),
		ApiUrl+subApi+token+"&flag=v2rayn&flag_info_hide=true",
		// "https://xxxx.com/api/v1/client/subscribe?token=9d21d6638329c5888ddfcbb663ff4c29&flag=v2rayn&test="+token, //测试用
	); err != nil {
		return err.Error()
	} else {
		defer r.Close()
		base64Str := r.ReadAllString()

		fmt.Println(base64Str)

		nodeBase64Str, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			return err.Error()
		}

		err = nodep.WriteBytes(nodeBase64Str, "./config/sub.txt")
		if err != nil {
			return err.Error()
		}

		return nodep.ConvertShareTextToXrayJson("./config/sub.txt", "./config/outbounds_config.json")

	}

}

// 获取公告
func GetappBul() string {
	c := GetgClient()
	if r, err := c.Get(
		gctx.New(),
		ApiUrl+appBulApi,
	); err != nil {
		return err.Error()
	} else {
		defer r.Close()
		jsonStr := r.ReadAllString()

		fmt.Println(jsonStr)

		if gjson.Get(jsonStr, "code").Int() != 0 {
			return gjson.Get(jsonStr, "message").String()
		} else {
			DB.Set("app_bulletin", gjson.Get(jsonStr, "data.data").String())
		}

	}

	return ""
}
