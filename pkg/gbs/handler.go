package gbs

import (
	"github.com/gowvp/owl/pkg/gbs/sip"
)

func (g GB28181API) handlerMessage(ctx *sip.Context) {
	// req := ctx.Request
	// tx := ctx.Tx

	// 	case "Catalog":
	// 		// 设备列表
	// 		sipMessageCatalog(u, body)
	// 		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	// 		return
	// 	case "Keepalive":
	// 		// heardbeat
	// 		if err := sipMessageKeepalive(u, body); err == nil {
	// 			tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	// 			// 心跳后同步注册设备列表信息
	// 			// sipCatalog(u)
	// 			return
	// 		}
	// 	case "RecordInfo":
	// 		// 设备音视频文件列表
	// 		sipMessageRecordInfo(u, body)
	// 		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	// 	case "DeviceInfo":
	// 		// 主设备信息
	// 		sipMessageDeviceInfo(u, body)
	// 		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	// 		return
	// 	}
	// 	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
	// }

	// func (g GB28181API) handlerMessage2(ctx *sip.Context) {
	// 	req := ctx.Request
	// 	tx := ctx.Tx

	// 	u, ok := parserDevicesFromReqeust(req)
	// 	if !ok {
	// 		// 未解析出来源用户返回错误
	// 		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
	// 		return
	// 	}
	// 	// 判断是否存在body数据
	// 	if len, have := req.ContentLength(); !have || len.Equals(0) {
	// 		// 不存在就直接返回的成功
	// 		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	// 		return
	// 	}
	// 	body := req.Body()
	// 	message := &MessageReceive{}

	//	if err := sip.XMLDecode(body, message); err != nil {
	//		// logrus.Warnln("Message Unmarshal xml err:", err, "body:", string(body))
	//		// 有些body xml发送过来的不带encoding ，而且格式不是utf8的，导致xml解析失败，此处使用gbk转utf8后再次尝试xml解析
	//		body, err = sip.GbkToUtf8(body)
	//		if err != nil {
	//			// logrus.Errorln("message gbk to utf8 err", err)
	//		}
	//		if err := sip.XMLDecode(body, message); err != nil {
	//			// logrus.Errorln("Message Unmarshal xml after gbktoutf8 err:", err, "body:", string(body))
	//			tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
	//			return
	//		}
	//	}
	//
	// switch message.CmdType {
	// case "Catalog":
	//
	//	// 设备列表
	//	sipMessageCatalog(u, body)
	//	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	//	return
	//
	// case "Keepalive":
	//
	//	// heardbeat
	//	if err := sipMessageKeepalive(u, body); err == nil {
	//		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	//		// 心跳后同步注册设备列表信息
	//		sipCatalog(u)
	//		return
	//	}
	//
	// case "RecordInfo":
	//
	//	// 设备音视频文件列表
	//	sipMessageRecordInfo(u, body)
	//	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	//
	// case "DeviceInfo":
	//
	//		// 主设备信息
	//		sipMessageDeviceInfo(u, body)
	//		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	//		return
	//	}
	//
	// tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
}
