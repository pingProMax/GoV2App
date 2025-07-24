package frame

import (
	"GoV2App/device_code"
	"image/color"
	"net/url"
	"strings"
	"time"

	"GoV2App/webapi"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xtls/xray-core/core"
)

func LoginHandle(app fyne.App) {
	w := app.NewWindow("登录")
	w.Canvas()

	v := webapi.DB.Get("jwt")
	if v != "" {
		jwtEx := webapi.DB.Get("jwtEx")
		if gtime.New(jwtEx).Timestamp() > time.Now().Unix() {
			// w.Hide()
			UserHandle(app)
			return
		}
	}

	userName := webapi.DB.Get("userName")
	passWd := webapi.DB.Get("passWd")

	usernameEntry := &widget.Entry{PlaceHolder: "请输入用户名", Text: userName}
	passwdEntry := &widget.Entry{PlaceHolder: "请输入密码", Password: true, Text: passWd}
	// passwdEntry.Text = "******"
	loginBut := &widget.Button{
		Text: "登录",
		Icon: theme.ConfirmIcon(),
		OnTapped: func() {
			//登录
			infoStr := webapi.Login(usernameEntry.Text, passwdEntry.Text)
			if infoStr == "" {

				webapi.DB.Set("userName", usernameEntry.Text)
				webapi.DB.Set("passWd", passwdEntry.Text)

				// w.Hide()
				UserHandle(app)
				w.Close()
				xrayC(w)
			} else {
				dialog.ShowInformation("提示", infoStr, w)

			}

		},
	}
	link, err := url.Parse(webapi.RegisterUrl)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	hyperlink := widget.NewHyperlink("注册账号", link)

	rect := canvas.NewRectangle(&color.NRGBA{128, 128, 128, 255})
	rect.SetMinSize(fyne.NewSize(300, 300))

	froms := container.NewVBox(usernameEntry, passwdEntry, container.NewBorder(nil, nil, nil, hyperlink, loginBut))

	// green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
	green := color.NRGBA{}
	canvasRectangle := canvas.NewRectangle(green)
	canvasRectangle.SetMinSize(fyne.NewSize(300, 300))

	w.Resize(fyne.NewSize(300, 300))
	w.SetContent(
		container.NewCenter(
			container.NewBorder(nil, nil, nil, nil, canvasRectangle, container.NewVBox(
				layout.NewSpacer(),
				froms,
				layout.NewSpacer(),
			)),
		),
	)
	passwdEntry.Refresh()
	w.Show()
}

func xrayC(w fyne.Window) {

	str := ``

	// if data, err := os.ReadFile("./v2ray.json"); err != nil {
	// 	panic(err.Error())
	// } else {
	// 	str = string(data)
	// }

	//设置环境变量
	// os.Setenv("xray.location.asset", getCurrentAbPathByCaller()+"/")

	// fmt.Println(getCurrentAbPathByCaller())

	// config, err := core.LoadConfig(getConfigFormat(), configFiles[0], configFiles)
	c, err := core.LoadConfig("json", strings.NewReader(str))
	if err != nil {
		dialog.ShowInformation("提示", "json加载错误："+err.Error(), w)
		return
	}

	server, err := core.New(c)
	if err != nil {
		dialog.ShowInformation("提示", "核心创建失败："+err.Error(), w)
		return
	}

	err = server.Start()
	if err != nil {
		dialog.ShowInformation("提示", "核心启动失败："+err.Error(), w)
		return
	}
	dialog.ShowInformation("提示", "启动", w)

	device_code.SetDeviceProxy("127.0.0.1:30809")
}
