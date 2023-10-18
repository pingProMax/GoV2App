gov2pane 客户端  
安卓/windows 跨平台    
还在开发...  

暂时只当作win客户端

# Windows  
直接给系统代理设置http代理即可  device_code/win.go  


# 安卓Android
小弟不才，安卓端VpnService看了好几天还不知道如何办 device_code/android.go  

如果没办法解决就全当作小白win端把~  
到时候鄙人再学习学习kotlin弄个安卓端  



# 站点设置  
系统管理 /admin/setting  添加一个code，名为 `app_bulletin` ,描述: 客户端公告  

目前windows可以用  

# 如何编译 
参考：https://developer.fyne.io/started/  
1、安装golang   
2、安装Gcc  
3、安装fyne `go install fyne.io/fyne/v2/cmd/fyne@latest`  
4、拉取此项目 git clone xxxx  
5、go get https://xxx  
6、webapi/api-info.go 更改api接口地址  
7、编译打包exe `fyne package`  

# 页面展示  
页面略丑，流下了没有技术的泪水😟    
<img src="https://github.com/pingProMax/GoV2App/blob/main/1.png?raw=true" width="30%">
<img src="https://github.com/pingProMax/GoV2App/blob/main/2.png?raw=true" width="30%">