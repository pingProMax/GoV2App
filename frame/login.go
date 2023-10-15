package frame

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func LoginHandle(app fyne.App) {
	w := app.NewWindow("登录")
	w.Canvas()

	usernameEntry := &widget.Entry{PlaceHolder: "请输入用户名", Text: "admin"}
	passwdEntry := &widget.Entry{PlaceHolder: "请输入密码", Password: true, Text: "admin"}
	loginBut := &widget.Button{
		Text: "登录",
		Icon: theme.ConfirmIcon(),
		OnTapped: func() {
			//登录
			if usernameEntry.Text == "admin" && passwdEntry.Text == "admin" {
				// w.Hide()
				UserHandle(app)
				w.Close()
			} else {
				dialog.ShowInformation("提示", "密码错误！", w)
			}

		},
	}
	link, err := url.Parse("https://fyne.io/")
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

	w.Show()
}
