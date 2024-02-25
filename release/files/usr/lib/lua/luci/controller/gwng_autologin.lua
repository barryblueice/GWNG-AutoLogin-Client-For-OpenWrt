module("luci.controller.gwng_autologin", package.seeall)

function index()
    entry({"admin", "services", "gwng_autologin"}, cbi("gwng_autologin"), _("GWNG AutoLogin Client"), 1)
end