//go:build android
// +build android

package device_code

import (
	"fyne.io/fyne/v2"
	"github.com/xtls/xray-core/core"
)

var xrayCore *core.Instance

//socket 入站端口30808、http入站端口30809
//socks2tun库
//https://github.com/heiher/hev-socks5-tunnel
//https://github.com/heiher/sockstun

func XrayStartService(w fyne.Window, radioText string, selectIndex int, start bool) {
}

// 如何golang调用呢
func SetDeviceProxy(proxy string) error {
}
