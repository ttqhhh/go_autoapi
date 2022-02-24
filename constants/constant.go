package constants

const CookieSecretKey = "zuiyou"
const TimeFormat = "2006-01-02 15:04:05"

// 0：最右，1：皮皮，2：海外，3：中东，4：妈妈，5：商业化，6：海外-US
const (
	ZuiyYou = iota
	PiPi
	HaiWai
	ZhongDong
	Matuan
	ShangYeHua
	HaiWaiUS
)

const (
	ZuiYouOnlineDomain      = "http://api.izuiyou.com/"
	ZuiYouH5OnlineDomain    = "https://h5.izuiyou.com"
	PiPiOnlineDomain        = "http://api.ippzone.com/"
	HaiWaiOnlineDomain      = "http://api.icocofun.com/"
	ZhongDongOnlineDomain   = "http://api.mehiya.com/"
	ZhongDongH5OnlineDomain = "https://h5.mehiya.com/"
	MatuanOnlineDomain      = "http://api.isupermama.com/"
	ShangYeHuaOnlineDomain  = "https://adapi.izuiyou.com/"
	HaiWaiUSOnlineDomain    = "http://usapi.icocofun.com/"

	ZuiYouTestDomain      = "http://test.izuiyou.com/"
	ZuiYouH5TestDomain    = "https://h5test.izuiyou.com/"
	PiPiTestDomain        = "http://testapi.ippzone.com/"
	HaiWaiTestDomain      = "http://test.icocofun.com/"
	ZhongDongTestDomain   = "http://me-live-test.gifgif.cn/"
	ZhongDongH5TestDomain = "http://melive-web-test.ixiaochuan.cn/"
	MatuanTestDomain      = ""
	ShangYeHuaTestDomain  = "http://test.izuiyou.com/"
	HaiWaiUSTestDomain    = "http://magatest.icocofun.com/"
)

const ONLINE_DOMAIN = ZuiYouOnlineDomain + ZuiYouH5OnlineDomain + PiPiOnlineDomain + HaiWaiOnlineDomain + ZhongDongOnlineDomain + ZhongDongH5OnlineDomain + MatuanOnlineDomain + ShangYeHuaOnlineDomain + HaiWaiUSOnlineDomain
const TEST_DOMAIN = ZuiYouTestDomain + ZuiYouH5TestDomain + PiPiTestDomain + HaiWaiTestDomain + ZhongDongTestDomain + ZhongDongH5TestDomain + MatuanTestDomain + ShangYeHuaTestDomain + HaiWaiUSTestDomain
