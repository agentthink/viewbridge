package gbs

import (
	"strings"

	"github.com/gowvp/owl/internal/core/ipc"
	"github.com/gowvp/owl/pkg/gbs/sip"
)

// QueryDeviceInfo 设备信息查询请求
// GB/T28181 81 页 A.2.4.4
func (g *GB28181API) QueryDeviceInfo(ctx *sip.Context) {
	tx, err := ctx.SendRequest(sip.MethodMessage, sip.GetDeviceInfoXML(ctx.DeviceID))
	if err != nil {
		ctx.Log.Error("sipDeviceInfo", "err", err)
		return
	}
	if _, err := sipResponse(tx); err != nil {
		ctx.Log.Error("sipResponse", "err", err)
		return
	}
}

// MessageDeviceInfoResponse 设备信息查询应答结构
type MessageDeviceInfoResponse struct {
	CmdType      string `xml:"CmdType"`
	SN           int    `xml:"SN"`
	DeviceID     string `xml:"DeviceID"`     // 目标设备的编码(必选)
	DeviceName   string `xml:"DeviceName"`   // 目标设备的名称(可选
	Manufacturer string `xml:"Manufacturer"` // 设备生产商(可选)
	Model        string `xml:"Model"`        // 设备型号(可选)
	Firmware     string `xml:"Firmware"`     // 设备固件版本(可选)
	Result       string `xml:"Result"`       // 査询结果(必选)
}

// isResultOK 判定 DeviceInfo 应答是否成功
// 为什么: 部分厂商不严格按协议返回 Result，可能为空串或大小写不一；空串按成功处理，维持与历史行为兼容。
func (m *MessageDeviceInfoResponse) isResultOK() bool {
	r := strings.TrimSpace(m.Result)
	return r == "" || strings.EqualFold(r, "OK")
}

// sipMessageDeviceInfo 设备信息查询应答
// GB/T28181 91 页 A.2.6.5
func (g *GB28181API) sipMessageDeviceInfo(ctx *sip.Context) {
	var msg MessageDeviceInfoResponse
	if err := sip.XMLDecode(ctx.Request.Body(), &msg); err != nil {
		ctx.Log.Error("sipMessageDeviceInfo", "err", err, "body", string(ctx.Request.Body()))
		ctx.String(400, ErrXMLDecode.Error())
		return
	}

	// 为什么: body 里的 DeviceID 必须与 SIP From 头解析出的 DeviceID 一致，防止越权覆盖他人设备信息。
	if msg.DeviceID != "" && msg.DeviceID != ctx.DeviceID {
		ctx.Log.Error("sipMessageDeviceInfo device id mismatch", "body_device_id", msg.DeviceID, "ctx_device_id", ctx.DeviceID)
		ctx.String(400, "device id mismatch")
		return
	}

	// 为什么: Result 非 OK 代表设备端查询失败，可选字段可能为空或旧值，不应覆盖数据库，避免清空已有厂商/型号等信息。
	if !msg.isResultOK() {
		ctx.Log.Warn("sipMessageDeviceInfo result not ok", "result", msg.Result, "sn", msg.SN)
		ctx.String(200, "OK")
		return
	}

	if err := g.core.Update(ctx.DeviceID, func(d *ipc.Device) {
		// 为什么: 可选字段为空时不覆盖，避免把上一次成功拿到的信息抹成空串。
		if msg.Firmware != "" {
			d.Ext.Firmware = msg.Firmware
		}
		if msg.Manufacturer != "" {
			d.Ext.Manufacturer = msg.Manufacturer
		}
		if msg.Model != "" {
			d.Ext.Model = msg.Model
		}
		if msg.DeviceName != "" {
			d.Ext.Name = msg.DeviceName
		}

		d.Address = ctx.Source.String()
		d.Transport = ctx.Source.Network()
	}); err != nil {
		ctx.Log.Error("Edit", "err", err)
		ctx.String(500, ErrDatabase.Error())
		return
	}

	ctx.String(200, "OK")
}
