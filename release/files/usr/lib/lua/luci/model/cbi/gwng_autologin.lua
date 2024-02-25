require("luci.model.network").init()
local sys = require "luci.sys"  

local client_running = false  
local client_path = "/usr/sbin/gwng" 
local cmd = "pgrep -x -f '%s'" % client_path  
local pid = sys.exec(cmd)  
if pid and pid ~= "" then  
    client_running = true  
end  

m = Map("gwng_autologin", translate("GWNG Autologin Client"), translate("广外南国路由器自动登录<br>Powered By barryblueice<br>联系方式：1293945058（QQ），13268648195（电话，微信同号）<br>项目地址：<br>https://github.com/barryblueice/GWNG-AutoLogin-Client-For-OpenWrt（国外）<br>https://gitee.com/barryblueice/gwng-auto-login-client-for-open-wrt（国内）"))

s = m:section(TypedSection, "gwng", "")
s.addremove = false
s.anonymous = true

local status_msg = client_running and translate("<span style='color:blue'><b>自动登录Client已在运行！</b></span>") or translate("<span style='color:red'><b>自动登录Client未在运行！</b></span>")
-- local status = s:option(DummyValue, "_status", status_msg)

o = s:option(Button, "", translate(" "))
o.title = translate("状态")
o.inputtitle = client_running and translate("自动登录Client已在运行！") or translate("自动登录Client未在运行！")
o.inputstyle = "readonly"
o.inputstyle = "pointer-events: none;"
o.disabled = true

o = s:option(Button, "s", translate(" "))
o.inputtitle = translate("开始运行")
o.inputstyle = "apply"
o.write = function(self, section)
    luci.util.exec("/etc/init.d/gwng_autologin start > /dev/null &")
    luci.http.redirect(luci.dispatcher.build_url('admin', 'services', 'gwng_autologin'))
end

o = s:option(Button, "", translate(" "))
o.inputtitle = translate("停止运行")
o.inputstyle = "reset"
o.write = function(self, section)
    luci.util.exec("/etc/init.d/gwng_autologin stop > /dev/null &")
    luci.http.redirect(luci.dispatcher.build_url('admin', 'services', 'gwng_autologin'))
end

o.rawhtml = true

username = s:option(Value, "username", translate("学号"))
pass = s:option(Value, "password", translate("密码"))
pass.password = true

num_of_cards = s:option(ListValue, "num", translate("联网方式"))
num_of_cards:value("1", "正常联网")
num_of_cards:value("2", "单线双播联网")
num_of_cards.description = translate("正常联网占用一个联网设备，单线双播联网占用两个联网设备。<br>单线双播联网可叠加网速，达到近似<span style='color:blue'><strong>百兆网络</strong></span>的网速。（以校方单设备最高限速为准）<br><span style='color:red'><b>单线双播由于学校限制问题，会导致其他联网设备被校方强行T下线！</b></span>")

iface_multiselect1 = s:option(DynamicList, "interfaces1", translate("选择正常联网网卡"))
iface_multiselect1.template = "cbi/network_ifacelist"

iface_multiselect2 = s:option(DynamicList, "interfaces2", translate("选择双播联网网卡"))
iface_multiselect2.template = "cbi/network_ifacelist"
iface_multiselect2:depends("num", "2")
iface_multiselect2.description = translate("<span style='color:red'><b>如联网方式为单线双播联网，则正常联网网卡和双播联网网卡需要同时选择。<br>注意，正常联网网卡和双播联网网卡需要分别选择不同的联网网卡！</b></span>")

local network = luci.model.network
if network then
    local ifaces = network:get_interfaces()
    if ifaces then
        for _, iface in ipairs(ifaces) do
            iface_multiselect1:value(iface[".name"], iface[".name"])
            iface_multiselect2:value(iface[".name"], iface[".name"])
        end
    end
end

local apply = luci.http.formvalue("cbi.apply")
local enable = luci.http.formvalue("cbi.enable")
if apply then
    io.popen("/etc/init.d/gwng_autologin restart > /dev/null &")
end

m.scripts = [[
    <script type="text/javascript">
        $(document).ready(function() {
            $("#cbid.gwng_autologin.server\\.num").change(function() {
                var numValue = $(this).val();
                if (numValue === "2") {
                    $("#cbid.gwng_autologin.server\\.interfaces2").show();
                } else {
                    $("#cbid.gwng_autologin.server\\.interfaces2").hide();
                }
            });
            $("#cbid.gwng_autologin.server\\.num").change();
        });
    </script>
]]

return m
