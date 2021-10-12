package jd_cookie

import (
	"fmt"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/develop/qinglong"
	"github.com/gin-gonic/gin"
)

var pinQQ = core.NewBucket("pinQQ")
var pinTG = core.NewBucket("pinTG")
var pinWXMP = core.NewBucket("pinWXMP")
var pin = func(class string) core.Bucket {
	return core.Bucket("pin" + strings.ToUpper(class))
}

func init() {
	// 
	core.Server.POST("/send_wx_msg", func(c *gin.Context) {
		user_id := c.Query("user_id")
		access_token := c.Query("access_token")
		message := c.PostForm("message")
		type Result struct {
			Code    int         `json:"retcode"`
			Data    interface{} `json:"data"`
			Message string      `json:"message"`
			ErrMsg  string      `json:"errmsg"`
		}
		result := Result{
			Data: nil,
			Code: 0,
		}
		notify_token_cfg := jd_cookie.Get("access_token")
		if access_token != notify_token_cfg {
			result.ErrMsg = "notify_token未设置或设置不正确。"
			result.Code = 100
			c.JSON(200, result)
			return
		}
		if push, ok := core.Pushs["wx"]; ok {
			push(user_id, message)
		}
		result.Message = "发送给WX[" + user_id + "] : " + message // "一句mmp，不知当讲不当讲。"
		c.JSON(200, result)
		return
	})
	core.Server.POST("/send_group_msg", func(c *gin.Context) {
		group_id := c.Query("group_id")
		access_token := c.Query("access_token")
		message := c.PostForm("message")
		type Result struct {
			Code    int         `json:"retcode"`
			Data    interface{} `json:"data"`
			Message string      `json:"message"`
			ErrMsg  string      `json:"errmsg"`
		}
		result := Result{
			Data: nil,
			Code: 0,
		}
		notify_token_cfg := jd_cookie.Get("access_token")
		if access_token != notify_token_cfg {
			result.ErrMsg = "notify_token未设置或设置不正确。"
			result.Code = 100
			c.JSON(200, result)
			return
		}
		if push, ok := core.GroupPushs["qq"]; ok {
			push(core.Int(group_id), int(0), message)
		}
		result.Message = "发送给QQ群[" + group_id + "] : " + message // "一句mmp，不知当讲不当讲。"
		c.JSON(200, result)
		return
	})
	core.Server.POST("/send_private_msg", func(c *gin.Context) {
		user_id := c.Query("user_id")
		access_token := c.Query("access_token")
		message := c.PostForm("message")
		type Result struct {
			Code    int         `json:"retcode"`
			Data    interface{} `json:"data"`
			Message string      `json:"message"`
			ErrMsg  string      `json:"errmsg"`
		}
		result := Result{
			Data: nil,
			Code: 0,
		}
		notify_token_cfg := jd_cookie.Get("access_token")
		if access_token != notify_token_cfg {
			result.ErrMsg = "notify_token未设置或设置不正确。"
			result.Code = 100
			c.JSON(200, result)
			return
		}
		core.Push("qq", core.Int(user_id), message)
		result.Message = "发送给QQ[" + user_id + "] : " + message // "一句mmp，不知当讲不当讲。"
		c.JSON(200, result)
		return
	})
	// 可改造成发送通知的机器人
	core.Server.POST("/notify", func(c *gin.Context) {
		msg := c.PostForm("msg")
		qq := c.PostForm("qq")
		notify_token := c.PostForm("notify_token")
		type Result struct {
			Code    int         `json:"code"`
			Data    interface{} `json:"data"`
			Message string      `json:"message"`
		}
		result := Result{
			Data: nil,
			Code: 300,
		}
		notify_token_cfg := jd_cookie.Get("notify_token")
		if notify_token != notify_token_cfg {
			result.Message = "notify_token未设置或设置不正确。"
			c.JSON(200, result)
			return
		}
		core.Push("qq", core.Int(qq), msg)
		result.Message = "发送给QQ[" + qq + "] : " + msg // "一句mmp，不知当讲不当讲。"
		c.JSON(200, result)
		return
	})
	core.Server.POST("/cookie", func(c *gin.Context) {
		cookie := c.PostForm("ck")
		qq := c.PostForm("qq")
		ck := &JdCookie{
			PtKey: core.FetchCookieValue(cookie, "pt_key"),
			PtPin: core.FetchCookieValue(cookie, "pt_pin"),
		}
		type Result struct {
			Code    int         `json:"code"`
			Data    interface{} `json:"data"`
			Message string      `json:"message"`
		}
		result := Result{
			Data: nil,
			Code: 300,
		}
		if ck.PtPin == "" || ck.PtKey == "" {
			result.Message = cookie // "一句mmp，不知当讲不当讲。"
			c.JSON(200, result)
			return
		}
		if !ck.Available() {
			result.Message = "无效的ck，请重试。"
			c.JSON(200, result)
			return
		}
		value := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
		envs, err := qinglong.GetEnvs("JD_COOKIE")
		if err != nil {
			result.Message = err.Error()
			c.JSON(200, result)
			return
		}
		find := false
		for _, env := range envs {
			if strings.Contains(env.Value, fmt.Sprintf("pt_pin=%s;", ck.PtPin)) {
				envs = []qinglong.Env{env}
				find = true
				break
			}
		}
		if !find {
			if err := qinglong.AddEnv(qinglong.Env{
				Name:  "JD_COOKIE",
				Value: value,
				Remarks: "QQ=" + qq + ";",
			}); err != nil {
				result.Message = err.Error()
				c.JSON(200, result)
				return
			}
			rt := ck.Nickname + "，添加成功。"
			core.NotifyMasters(rt)
			result.Message = rt
			result.Code = 200
			c.JSON(200, result)
			return
		} else {
			env := envs[0]
			env.Value = value
			if env.Status != 0 {
				if err := qinglong.Config.Req(qinglong.PUT, qinglong.ENVS, "/enable", []byte(`["`+env.ID+`"]`)); err != nil {
					result.Message = err.Error()
					c.JSON(200, result)
					return
				}
			}
			env.Status = 0
			if err := qinglong.UdpEnv(env); err != nil {
				result.Message = err.Error()
				c.JSON(200, result)
				return
			}
			rt := ck.Nickname + "，更新成功。"
			core.NotifyMasters(rt)
			result.Message = rt
			result.Code = 200
			c.JSON(200, result)
			return
		}
	})
	core.AddCommand("jd", []core.Function{
		{
			Rules: []string{`unbind ?`},
			Handle: func(s core.Sender) interface{} {
				s.Disappear(time.Second * 40)
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "暂时无法操作。"
				}
				for _, env := range envs {
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					pin(s.GetImType()).Foreach(func(k, v []byte) error {
						if string(k) == pt_pin && string(v) == s.Get() {
							s.Reply(fmt.Sprintf("已解绑，%s。", pt_pin))
							defer func() {
								pinQQ.Set(string(k), "")
							}()
						}
						return nil
					})
				}
				return "操作完成"
			},
		},
		{
			Rules:   []string{`raw pt_key=([^;=\s]+);\s*pt_pin=([^;=\s]+)`},
			FindAll: true,
			Handle: func(s core.Sender) interface{} {
				s.Reply(s.Delete())
				s.Disappear(time.Second * 20)
				for _, v := range s.GetAllMatch() {
					ck := &JdCookie{
						PtKey: v[0],
						PtPin: v[1],
					}
					if len(ck.PtKey) <= 20 {
						s.Reply("再捣乱我就报警啦！")
						continue
					}
					if !ck.Available() {
						s.Reply("请先到app内设置好账号昵称。") //有瞎编ck的嫌疑
						continue
					}
					if ck.Nickname == "" {
						s.Reply("再捣乱我就报警啦！")
					}
					value := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
					envs, err := qinglong.GetEnvs("JD_COOKIE")
					if err != nil {
						s.Reply(err)
						continue
					}
					find := false
					for _, env := range envs {
						if strings.Contains(env.Value, fmt.Sprintf("pt_pin=%s;", ck.PtPin)) {
							envs = []qinglong.Env{env}
							find = true
							break
						}
					}
					pin(s.GetImType()).Set(ck.PtPin, s.GetUserID())
					if !find {
						if err := qinglong.AddEnv(qinglong.Env{
							Name:  "JD_COOKIE",
							Value: value,
							Remarks: "QQ=" + fmt.Sprintf("%d", s.GetUserID()) + ";",
						}); err != nil {
							s.Reply(err)
							continue
						}
						rt := ck.Nickname + "，添加成功。"
						core.NotifyMasters(rt)
						s.Reply(rt)
						continue
					} else {
						env := envs[0]
						env.Value = value
						if env.Status != 0 {
							if err := qinglong.Config.Req(qinglong.PUT, qinglong.ENVS, "/enable", []byte(`["`+env.ID+`"]`)); err != nil {
								s.Reply(err)
								continue
							}
						}
						env.Status = 0
						if err := qinglong.UdpEnv(env); err != nil {
							s.Reply(err)
							continue
						}
						rt := ck.Nickname + "，更新成功。"
						core.NotifyMasters(rt)
						s.Reply(rt)
						continue
					}
				}
				return nil
			},
		},
		{
			Rules:   []string{`raw pin=([^;=\s]+);\s*wskey=([^;=\s]+)`},
			FindAll: true,
			Handle: func(s core.Sender) interface{} {
				s.Reply(s.Delete())
				s.Disappear(time.Second * 20)
				value := fmt.Sprintf("pin=%s;wskey=%s;", s.Get(0), s.Get(1))

				pt_key, err := getKey(value)
				if err == nil {
					if strings.Contains(pt_key, "fake") {
						return "无效的wskey，请重试。"
					}
				} else {
					s.Reply(err)
				}
				ck := &JdCookie{
					PtKey: pt_key,
					PtPin: s.Get(0),
				}
				ck.Available()
				envs, err := qinglong.GetEnvs("pin=")
				if err != nil {
					return err
				}
				pin(s.GetImType()).Set(ck.PtPin, s.GetUserID())
				var envCK *qinglong.Env
				var envWsCK *qinglong.Env
				for i := range envs {
					if strings.Contains(envs[i].Value, fmt.Sprintf("pin=%s;wskey=", ck.PtPin)) && envs[i].Name == "JD_WSCK" {
						envWsCK = &envs[i]
					} else if strings.Contains(envs[i].Value, fmt.Sprintf("pt_pin=%s;", ck.PtPin)) && envs[i].Name == "JD_COOKIE" {
						envCK = &envs[i]
					}
				}
				value2 := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
				if envCK == nil {
					qinglong.AddEnv(qinglong.Env{
						Name:  "JD_COOKIE",
						Value: value2,
						Remarks: "QQ=" + fmt.Sprintf("%d", s.GetUserID()) + ";",
					})
				} else {
					envCK.Value = value2
					if err := qinglong.UdpEnv(*envCK); err != nil {
						return err
					}
				}
				if envWsCK == nil {
					if err := qinglong.AddEnv(qinglong.Env{
						Name:  "JD_WSCK",
						Value: value,
						Remarks: "QQ=" + fmt.Sprintf("%d", s.GetUserID()) + ";",
					}); err != nil {
						return err
					}
					return ck.Nickname + ",添加成功。"
				} else {
					env := envs[0]
					env.Value = value
					if env.Status != 0 {
						if err := qinglong.Config.Req(qinglong.PUT, qinglong.ENVS, "/enable", []byte(`["`+env.ID+`"]`)); err != nil {
							return err
						}
					}
					env.Status = 0
					if err := qinglong.UdpEnv(env); err != nil {
						return err
					}
					return ck.Nickname + ",更新成功。"
				}
			},
		},
	})
}
