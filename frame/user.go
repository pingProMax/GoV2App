package frame

import (
	"fmt"
	"image/color"
	"net/url"
	"time"

	"GoV2App/go_resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func init() {
}

func UserHandle(app fyne.App) {

	w := app.NewWindow("GoV2加速")
	w.SetMaster()
	w.SetFixedSize(true) //禁止拖放

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(
			"首页",
			theme.HomeIcon(),
			makeHomeUI(),
		),
		container.NewTabItemWithIcon(
			"用户",
			theme.AccountIcon(),
			makeUserUI(),
		),
	)
	tabs.SetTabLocation(container.TabLocationBottom) //底部

	w.SetContent(container.NewBorder(nil, nil, nil, nil, tabs))

	w.Resize(fyne.NewSize(300, 500))

	w.Show()
}

// 首页 界面
func makeHomeUI() *fyne.Container {

	//logo
	canvasImg := canvas.NewImageFromResource(go_resource.ResourceImgLogoPng)

	//选择节点
	selectEntry := widget.NewSelect([]string{"节点 1", "节点 2", "节点 3"}, func(s string) {
		fmt.Println("selected", s)

	})
	selectEntry.Selected = "节点 1"

	//单选
	radio := widget.NewRadioGroup([]string{"绕过大陆", "全局"}, func(s string) { fmt.Println("selected", s) })
	radio.Horizontal = true
	radio.Selected = "绕过大陆"

	//公告文本
	coticeText := widget.NewRichTextFromMarkdown("公告：我是公告~")
	coticeText.Wrapping = fyne.TextWrapBreak
	coticeText.Scroll = container.ScrollBoth
	coticeText.Refresh()

	return container.NewBorder(
		container.NewVBox(
			container.NewHBox(
				layout.NewSpacer(),
				container.NewGridWrap(
					fyne.NewSize(100, 100),
					canvasImg,
				),
				layout.NewSpacer(),
			),

			//间距
			Spacing(10),

			widget.NewLabel("选择节点"), // 上部分内容
			selectEntry,
			radio,
		),
		makeAnimationCanvas(),
		nil,
		nil,
		coticeText,
	)
}

func makeAnimationCanvas() fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Black)
	// rect.Resize()
	rect.SetMinSize(fyne.NewSize(300, 30))

	a := canvas.NewColorRGBAAnimation(theme.PrimaryColorNamed(theme.ColorBlue), theme.PrimaryColorNamed(theme.ColorGreen),
		time.Second*3, func(c color.Color) {
			rect.FillColor = c
			canvas.Refresh(rect)
		},
	)
	a.RepeatCount = fyne.AnimationRepeatForever
	a.AutoReverse = true

	running := false
	var toggle *widget.Button
	toggle = widget.NewButton("启动代理", func() {
		if running {
			a.Stop()
			toggle.SetText("启动代理")
		} else {
			a.Start()
			toggle.SetText("停止代理")
		}
		running = !running
	})
	toggle.Resize(toggle.MinSize())
	// toggle.Move(fyne.NewPos(152, 54))

	return container.NewVBox(rect, toggle, Spacing(10))
	// return container.NewBorder(nil, container.NewVBox(toggle, Spacing(20)), nil, nil,
	// container.NewGridWithRows(1, rect))
	// return container.NewVBox(rect, toggle)
}

// 用户 界面
func makeUserUI() *fyne.Container {

	//到期时间
	userName := widget.NewLabel("欢迎: admin")
	subName := widget.NewLabel("套餐信息: VIP")
	expireTime := widget.NewLabel("到期时间: 2023-11-12 12:11:23")

	fprogress := widget.NewProgressBar()
	fprogress.TextFormatter = func() string {
		return fmt.Sprintf("%.2f GB / %.2f GB", fprogress.Value, fprogress.Max)
	}
	fprogress.Max = 100
	fprogress.SetValue(10)

	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	hyperlink := widget.NewHyperlink("购买套餐", link)

	outBut := widget.NewButtonWithIcon("退出账号", theme.ContentUndoIcon(), func() {

	})
	outBut.Importance = widget.DangerImportance

	userInfo := container.NewVBox(
		userName,
		subName,
		expireTime,
		hyperlink,
		fprogress,
	)

	return container.NewBorder(userInfo, container.NewVBox(outBut, Spacing(10)), nil, nil)
}
