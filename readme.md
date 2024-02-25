# GWNG-AutoLogin-Client For OpenWrt

基于Golang+lua实现的 广东外语外贸大学南国商学院 校园网 OpenWrt路由器自动登录插件

~~**该插件由于特殊原因未经过验证，请谨慎使用。**~~

**强烈建议该插件运行在x86-64环境下，其他环境由于设备缺乏暂未验证。**

目前已验证可使用平台有：

1.x86_64（基于iStore OS）

2.mipsle（基于e8820s，mt7621需下专版可执行文件）

## 目前已实现的功能：

- [x] 校园网登录
- [x] 基于macvlan的单线双拨实现

## Install / Uninstall：

请自行查看Wiki。

## 以下内容仅供开发者观看：

由于个人水平问题，该项目代码包括但不局限于以下血压高内容：

1. 三角稳定：

```python
if:
    if:
        if:
            if:
                if:
                    if:
                    else:
                else:
            else:
        else:
    else:
else:
```

2. 视而不见：

```python
try:
except:
try:
except:
try:
except:
```

3. 热辣炸串：

```python
if:
elif:
elif:
elif:
elif:
else
```

4. 真假李逵：

```python
def mix(a,b):
    ......
def mix_reverse(a,b):
    ......
```


## 实现原理：

通过f12开发者模式进行浏览器抓包，可得知广外南国商学院验证方式为curl，具体格式如下：
```bash
curl "http://10.244.0.13/quickauth.do?userid={学号}&passwd={密码}&wlanuserip={当前设备的内网IP地址}&wlanacname=gwng&wlanacIp=10.244.0.1&ssid=&vlan=502&mac={当前设备的mac地址}&version=0&portalpageid=5&timestamp={时间戳}&uuid={uuid}&portaltype=1&hostname={设备名}" \  
-H "Accept: application/json, text/javascript, */*; q=0.01" \  
-H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6" \  
-H "Connection: keep-alive" \  
-H "Cookie: macAuth={mac}|{另一台已被验证的mac地址}|{当前设备的mac地址}; ABMS=c0c701f6-ea15-448e-a48b-6230c9117317" \  
-H "Referer: http://10.244.0.13/portal.do?wlanuserip={当前设备的内网IP地址}&wlanacname=gwng&mac={当前设备的mac地址}&vlan=502&hostname={设备名}&rand={个人猜测为随机种子}&url=http%3A%2F%2Fwww.msftconnecttest.com%2Fredirec" \  
-H {请求头} \  
-H "X-Requested-With: XMLHttpRequest"
```
如果用单线多拨的方式则需要多次验证，一次只能验证一个设备。

校内个人允许两台设备在线。

## 开源协议/License：

版权所有 (c) 2024 冰蓝莓barryblueice

根据 Apache 许可证 2.0 版本（以下简称 "许可证"）获得许可；
您不得使用此文件，除非符合许可证的规定。
您可以在以下位置获得许可证的副本：
```
http://www.apache.org/licenses/LICENSE-2.0
```

除非适用法律要求或书面同意，否则软件按 "原样" 提供，不提供任何形式的担保或条件，无论是明示的还是默示的。
有关许可证下的特定语言和限制，请参阅许可证。
