//go:build !android
// +build !android

package device_code

import (
	"GoV2App/nodep"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/tidwall/gjson"
	"github.com/xtls/xray-core/common/cmdarg"
	"github.com/xtls/xray-core/core"
)

const (
	INTERNET_PER_CONN_FLAGS               = 1
	INTERNET_PER_CONN_PROXY_SERVER        = 2
	INTERNET_PER_CONN_PROXY_BYPASS        = 3
	INTERNET_OPTION_REFRESH               = 37
	INTERNET_OPTION_SETTINGS_CHANGED      = 39
	INTERNET_OPTION_PER_CONNECTION_OPTION = 75
)

/*
	typedef struct {
	  DWORD dwOption;
	  union {
	    DWORD    dwValue;
	    LPSTR    pszValue;
	    FILETIME ftValue;
	  } Value;
	} INTERNET_PER_CONN_OPTIONA, *LPINTERNET_PER_CONN_OPTIONA;

	typedef struct _FILETIME {
	  DWORD dwLowDateTime;
	  DWORD dwHighDateTime;
	} FILETIME, *PFILETIME, *LPFILETIME;
*/
type INTERNET_PER_CONN_OPTION struct {
	dwOption uint32
	dwValue  uint64 // 注意 32位 和 64位 struct 和 union 内存对齐
}

type INTERNET_PER_CONN_OPTION_LIST struct {
	dwSize        uint32
	pszConnection *uint16
	dwOptionCount uint32
	dwOptionError uint32
	pOptions      uintptr
}

var xrayCore *core.Instance

// 开启代理
// w
// xray
// radioText代理类型 绕过大陆，全局
// selectIndex选择索引
// start是否开启
func XrayStartService(w fyne.Window, radioText string, selectIndex int, start bool) {
	path, _ := os.Getwd()

	//设置环境变量
	os.Setenv("XRAY_LOCATION_ASSET", path+"/config")

	inConfigJson, err := os.ReadFile("./config/inbounds_config.json")
	if err != nil {
		dialog.ShowInformation("提示", "inbounds_config.json 加载错误："+err.Error(), w)
		w.Close()
	}

	outboundsConfig, err := os.ReadFile("./config/outbounds_config.json")
	if err != nil {
		dialog.ShowInformation("提示", "outbounds_config.json 加载错误："+err.Error(), w)
		w.Close()
	}

	routingFileName := "routing2_config.json"
	if radioText == "绕过大陆" {
		routingFileName = "routing1_config.json"
	}
	routingConfig, err := os.ReadFile("./config/" + routingFileName)
	if err != nil {
		dialog.ShowInformation("提示", routingFileName+" 加载错误："+err.Error(), w)
		w.Close()
	}

	//生成 connet.json
	connetJsonStr := fmt.Sprintf(
		`{
			"inbounds":%s,
			"outbounds":[%s],
			"routing":%s
		}`,
		gjson.Get(string(inConfigJson), "inbounds").String(),
		gjson.Get(string(outboundsConfig), "outbounds."+strconv.Itoa(selectIndex)).String()+`,{"tag":"direct","protocol":"freedom","settings":{}},{"tag":"block","protocol":"blackhole","settings":{"response":{"type":"http"}}}`,
		gjson.Get(string(routingConfig), "routing"),
	)

	nodep.WriteText(connetJsonStr, "./config/connet.json")

	if xrayCore != nil {
		xrayCore.Close()
	}

	c, err := core.LoadConfig("json", cmdarg.Arg{"./config/connet.json"})
	if err != nil {
		dialog.ShowInformation("提示", "json加载错误："+err.Error(), w)
		return
	}

	xrayCore, err = core.New(c)
	if err != nil {
		dialog.ShowInformation("提示", "核心创建失败："+err.Error(), w)
		return
	}

	if start {
		err = xrayCore.Start()
		if err != nil {
			dialog.ShowInformation("提示", "核心启动失败："+err.Error(), w)
			return
		}
		SetDeviceProxy("127.0.0.1:30809")
	} else {
		SetDeviceProxy("")
	}

}

func SetDeviceProxy(proxy string) error {

	winInet, err := syscall.LoadLibrary("Wininet.dll")
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("LoadLibrary Wininet.dll Error: %s", err))
	}
	InternetSetOptionW, err := syscall.GetProcAddress(winInet, "InternetSetOptionW")
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("GetProcAddress InternetQueryOptionW Error: %s", err))
	}

	options := [3]INTERNET_PER_CONN_OPTION{}
	options[0].dwOption = INTERNET_PER_CONN_FLAGS
	if proxy == "" {
		options[0].dwValue = 1
	} else {
		options[0].dwValue = 2
	}

	options[1].dwOption = INTERNET_PER_CONN_PROXY_SERVER
	options[1].dwValue = uint64(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(proxy))))
	options[2].dwOption = INTERNET_PER_CONN_PROXY_BYPASS
	options[2].dwValue = uint64(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("localhost;127.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;172.32.*;192.168.*"))))

	list := INTERNET_PER_CONN_OPTION_LIST{}
	list.dwSize = uint32(unsafe.Sizeof(list))
	list.pszConnection = nil
	list.dwOptionCount = 3
	list.dwOptionError = 0
	list.pOptions = uintptr(unsafe.Pointer(&options))

	callInternetOptionW := func(dwOption uintptr, lpBuffer uintptr, dwBufferLength uintptr) error {
		r1, _, err := syscall.Syscall6(InternetSetOptionW, 4, 0, dwOption, lpBuffer, dwBufferLength, 0, 0)
		if r1 != 1 {
			return err
		}
		return nil
	}

	err = callInternetOptionW(INTERNET_OPTION_PER_CONNECTION_OPTION, uintptr(unsafe.Pointer(&list)), uintptr(unsafe.Sizeof(list)))
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_PER_CONNECTION_OPTION Error: %s", err)
	}
	err = callInternetOptionW(INTERNET_OPTION_SETTINGS_CHANGED, 0, 0)
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_SETTINGS_CHANGED Error: %s", err)
	}
	err = callInternetOptionW(INTERNET_OPTION_REFRESH, 0, 0)
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_REFRESH Error: %s", err)
	}
	return nil
}
