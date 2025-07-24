package frame

import (
	"fmt"
	"image/color"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"GoV2App/device_code"
	"GoV2App/go_resource"
	"GoV2App/webapi"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/tidwall/gjson"
)

func init() {
}

func dialogInfoQuit(title, content, butText string, butFunc func(), parent fyne.Window) {
	nc := dialog.NewCustom(title, "", &widget.Label{Text: content, Alignment: fyne.TextAlignCenter}, parent)

	coList := make([]fyne.CanvasObject, 0)

	but := &widget.Button{Text: butText, Importance: widget.SuccessImportance,
		OnTapped: func() {
			if butFunc != nil {
				butFunc()
			}
			parent.Close()
		},
	}
	coList = append(coList, but)
	nc.SetButtons(coList)
	nc.Show()
}

func UserHandle(app fyne.App) {
	w := app.NewWindow("GoV2加速")
	w.Resize(fyne.NewSize(400, 600))
	w.Show()

	w.SetMaster()
	// w.SetFixedSize(true) //禁止拖放

	//获取用户订阅地址
	infoStr := webapi.GetSubscribeToken()
	if infoStr != "" {
		dialogInfoQuit("提示", infoStr, "ok", nil, w)
		return
	}

	//获取节点信息
	infoStr = webapi.GetNodeInfo()
	if infoStr != "" {
		dialogInfoQuit("提示", "账号到期或者流量用完，请充值购买:)", "ok", func() {
			webapi.DB.Del("jwtEx")
			cmd := exec.Command("cmd", "/c", "start", webapi.PlanUrl)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Start()
		}, w)
		return
	}

	//获取公告
	infoStr = webapi.GetappBul()
	if infoStr != "" {
		dialogInfoQuit("提示", infoStr, "ok", nil, w)
		return
	}

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(
			"首页",
			theme.HomeIcon(),
			makeHomeUI(w),
		),
		container.NewTabItemWithIcon(
			"用户",
			theme.AccountIcon(),
			makeUserUI(w),
		),
	)
	tabs.SetTabLocation(container.TabLocationBottom) //底部

	w.SetContent(container.NewBorder(nil, nil, nil, nil, tabs))

}

// 选择节点索引
var selectIndex int

// 模式
var radioText string

// 是否启动
var start bool

// 首页 界面
func makeHomeUI(w fyne.Window) *fyne.Container {

	dStr, err := os.ReadFile("./config/outbounds_config.json")
	if err != nil {
		dialog.ShowInformation("提示", "outbounds_config.json 加载错误："+err.Error(), w)
		w.Close()
	}

	nodeArr := make([]string, 0)

	for _, v := range gjson.Get(string(dStr), "outbounds.#.name").Array() {
		nodeArr = append(nodeArr, v.String())
	}

	//logo
	canvasImg := canvas.NewImageFromResource(go_resource.ResourceImgLogoPng)

	//选择节点
	selectEntry := widget.NewSelect(nodeArr, nil)

	selectEntry.OnChanged = func(s string) {
		selectIndex = selectEntry.SelectedIndex()
		device_code.XrayStartService(w, radioText, selectIndex, start)
	}
	selectEntry.SetSelectedIndex(0)

	//单选
	radio := widget.NewRadioGroup([]string{"绕过大陆", "全局"}, func(s string) {
		fmt.Println("selected", s)
		radioText = s
		device_code.XrayStartService(w, radioText, selectIndex, start)
	})
	radio.Horizontal = true
	radio.Selected = "绕过大陆"

	//公告文本
	coticeText := widget.NewRichTextFromMarkdown(webapi.DB.Get("app_bulletin"))
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
		makeAnimationCanvas(w),
		nil,
		nil,
		coticeText,
	)
}

func makeAnimationCanvas(w fyne.Window) fyne.CanvasObject {
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
			start = false
		} else {
			a.Start()
			toggle.SetText("停止代理")
			start = true
		}
		running = !running
		device_code.XrayStartService(w, radioText, selectIndex, start)
	})
	toggle.Resize(toggle.MinSize())
	// toggle.Move(fyne.NewPos(152, 54))

	return container.NewVBox(rect, toggle, Spacing(10))
	// return container.NewBorder(nil, container.NewVBox(toggle, Spacing(20)), nil, nil,
	// container.NewGridWithRows(1, rect))
	// return container.NewVBox(rect, toggle)
}

// 用户 界面
func makeUserUI(w fyne.Window) *fyne.Container {

	//到期时间

	userName := widget.NewLabel("欢迎: " + webapi.DB.Get("user_name"))
	subName := widget.NewLabel("套餐信息: " + webapi.DB.Get("plan_name"))
	expireTime := widget.NewLabel("到期时间: " + webapi.DB.Get("expired_at"))
	v := widget.NewLabel("软件版本: " + webapi.V)

	fprogress := widget.NewProgressBar()
	fprogress.TextFormatter = func() string {
		return fmt.Sprintf("%.2f GB / %.2f GB", fprogress.Value, fprogress.Max)
	}

	fprogress.Max = BytesToGB(gconv.Int64(webapi.DB.Get("transfer_enable")))
	fprogress.SetValue(BytesToGB(gconv.Int64(webapi.DB.Get("u")) + gconv.Int64(webapi.DB.Get("d"))))

	link, err := url.Parse(webapi.PlanUrl)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	hyperlink := widget.NewHyperlink("购买套餐", link)

	outBut := widget.NewButtonWithIcon("退出账号", theme.ContentUndoIcon(), func() {
		webapi.DB.Clear()
		w.Close()
	})
	outBut.Importance = widget.DangerImportance

	userInfo := container.NewVBox(
		userName,
		subName,
		expireTime,
		v,
		hyperlink,
		fprogress,
	)

	return container.NewBorder(userInfo, container.NewVBox(outBut, Spacing(10)), nil, nil)
}

// float64 只保留两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// bytes 转 GB
func BytesToGB(bytes int64) float64 {
	gigabytes := Decimal(float64(bytes) / 1073741824)
	return gigabytes
}
