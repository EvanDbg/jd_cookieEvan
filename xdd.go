package jd_cookie

import "github.com/beego/beego/v2/client/httplib"

//对接xdd
func xdd(cookie string, qq string) {
	xdd_url := jd_cookie.Get("xdd_url")
	xdd_token := jd_cookie.Get("xdd_token")
	if xdd_url != "" {
		req := httplib.Post(xdd_url)
		req.Param("ck", cookie)
		req.Param("token", xdd_token)
		req.Param("qq", qq)
		req.Response()
	}
}
