package main

import (
	_ "embed"

	_ "GoV2App/webapi"

	_ "github.com/xtls/xray-core/main/distro/all"

	"GoV2App/frame"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed resource/font/siyuanyuanti.ttf
	NotoSansSC []byte
)

func main() {
	// device_code.SetDeviceProxy("") //打开就清除代理

	myApp := app.New()

	myApp.Settings().SetTheme(&MyTheme{})

	//登录
	frame.LoginHandle(myApp)

	myApp.Run()

	//用户界面

}

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

// StaticName 为 fonts 目录下的 ttf 类型的字体文件名
func (m MyTheme) Font(fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "siyuanyuanti.ttf",
		StaticContent: NotoSansSC,
	}
}

func (*MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
