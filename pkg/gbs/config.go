package gbs

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"log/slog"
	"math"

	"github.com/gowvp/owl/pkg/gbs/sip"
)

// 配置参数类型常量定义
const (
	// basicParam 基本参数配置
	basicParam = "BasicParam"
	// videoParamOpt 视频参数范围配置
	// videoParamOpt = "VideoParamOpt"
	// // SVACEncodeConfig SVAC编码配置
	// SVACEncodeConfig = "SVACEncodeConfig"
	// // SVACDecodeConfig SVAC解码配置
	// SVACDecodeConfig = "SVACDecodeConfig"
	// // videoParamAttribute 视频参数属性配置
	// videoParamAttribute = "VideoParamAttribute"
	// // videoRecordPlan 录像计划
	// videoRecordPlan = "VideoRecordPlan"
	// // videoAlarmRecord 报警录像
	// videoAlarmRecord = "VideoAlarmRecord"
	// // pictureMask 视频画面遮挡
	// pictureMask = "PictureMask"
	// // frameMirror 画面翻转
	// frameMirror = "FrameMirror"
	// // AlarmReport 报警上报开关
	// AlarmReport = "AlarmReport"
	// // OSDConfig 前端OSD设置
	// OSDConfig = "OSDConfig"
)

type ConfigDownloadRequest struct {
	XMLName        xml.Name  `xml:"Query"`
	CmdType        string    `xml:"CmdType"`    // 命令类型：设备配置查询(必选)
	SN             int32     `xml:"SN"`         // 命令序列号(必选)
	DeviceID       string    `xml:"DeviceID"`   // 目标设备编码(必选)
	ConfigType     string    `xml:"ConfigType"` // 查询配置参数类型(必选)
	SnapShotConfig *SnapShot `xml:"SnapShotConfig"`
}

type ConfigDownloadResponse struct {
	XMLName    xml.Name    `xml:"Response"`
	CmdType    string      `xml:"CmdType"`
	SN         int         `xml:"SN"`
	DeviceID   string      `xml:"DeviceID"`
	Result     string      `xml:"Result"`
	BasicParam *BasicParam `xml:"BasicParam"`
	// VideoParamOpt       *VideoParamOpt       `xml:"VideoParamOpt"`
	// SVACEncodeConfig    *SVACEncodeConfig    `xml:"SVACEncodeConfig"`
	// SVACDecodeConfig    *SVACDecodeConfig    `xml:"SVACDecodeConfig"`
	// VideoParamAttribute *VideoParamAttribute `xml:"VideoParamAttribute"`
	// VideoRecordPlan     *VideoRecordPlan     `xml:"VideoRecordPlan"`
	// VideoAlarmRecord    *VideoAlarmRecord    `xml:"VideoAlarmRecord"`
	// PictureMask         *PictureMask         `xml:"PictureMask"`
	// FrameMirror         *FrameMirror         `xml:"FrameMirror"`
	// AlarmReport         *AlarmReport         `xml:"AlarmReport"`
	// OSDConfig           *OSDConfig           `xml:"OSDConfig"`
	SnapShot *SnapShot `xml:"SnapShot"`
}

type SnapShot struct {
	SnapNum   int    `xml:"SnapNum"`   // 连拍张数(必选)，最多10张，当手动抓拍时，取值为1
	Interval  int    `xml:"Interval"`  // 单张抓拍间隔时间，单位：秒(必选)，取值范围:最短1秒
	UploadURL string `xml:"UploadURL"` // 抓拍图像上传路径(必选)
	SessionID string `xml:"SessionID"` // 会话ID，由平台生成，用于关联抓拍的图像与平台请求(必选)
}

// BasicParam 设备基本参数配置
type BasicParam struct {
	Name              string `xml:"Name"`              // 设备名称
	Expiration        int    `xml:"Expiration"`        // 注册过期时间
	HeartBeatInterval int    `xml:"HeartBeatInterval"` // 心跳间隔时间
	HeartBeatCount    int    `xml:"HeartBeatCount"`    // 心跳超时次数
}

const CMDTypeConfigDownload = "ConfigDownload"

func NewBasicParamRequest(sn int32, deviceID string) []byte {
	c := ConfigDownloadRequest{
		CmdType:    CMDTypeConfigDownload,
		SN:         sn,
		DeviceID:   deviceID,
		ConfigType: basicParam,
	}
	xmlData, _ := sip.XMLEncode(c)
	return xmlData
}

func (g *GB28181API) QueryConfigDownloadBasic(deviceID string) error {
	slog.Debug("QueryConfigDownloadBasic", "deviceID", deviceID)
	ipc, ok := g.svr.memoryStorer.Load(deviceID)
	if !ok || !ipc.IsOnline {
		return ErrDeviceOffline
	}

	tx, err := g.svr.wrapRequest(ipc, sip.MethodMessage, &sip.ContentTypeXML, NewBasicParamRequest(1, deviceID))
	if err != nil {
		return err
	}
	_, err = sipResponse(tx)
	return err
}

func (g *GB28181API) handleDeviceConfig(ctx *sip.Context) {
	slog.Debug("handleDeviceConfig", "deviceID", ctx.DeviceID)

	b := ctx.Request.Body()
	fmt.Println(">>>", string(b))
	// var msg DeviceConfigResponse
	// if err := sip.XMLDecode(ctx.Request.Body(), &msg); err != nil {
	// 	ctx.Log.Error("handleDeviceConfig", "err", err, "body", hex.EncodeToString(ctx.Request.Body()))
	// 	ctx.String(400, ErrXMLDecode.Error())
	// 	return
	// }

	// if msg.SnapShotConfig != nil {
	// 	slog.Debug("handleDeviceConfig", "snapShotConfig", msg.SnapShotConfig)
	// }

	ctx.String(200, "OK")
}

func (g *GB28181API) sipMessageConfigDownload(ctx *sip.Context) {
	slog.Debug("sipMessageConfigDownload", "deviceID", ctx.DeviceID)

	var msg ConfigDownloadResponse
	if err := sip.XMLDecode(ctx.Request.Body(), &msg); err != nil {
		ctx.Log.Error("sipMessageConfigDownload", "err", err, "body", hex.EncodeToString(ctx.Request.Body()))
		ctx.String(400, ErrXMLDecode.Error())
		return
	}

	if msg.BasicParam != nil {
		ipc, ok := g.svr.memoryStorer.Load(ctx.DeviceID)
		if !ok {
			ctx.Log.Debug("sipMessageConfigDownload", "deviceID", ctx.DeviceID, "err", "device offline")
			return
		}

		// 确保 HeartBeatCount 在合法范围内
		if msg.BasicParam.HeartBeatCount > math.MaxUint16 {
			msg.BasicParam.HeartBeatCount = math.MaxUint16
		}
		if msg.BasicParam.HeartBeatInterval > math.MaxUint16 {
			msg.BasicParam.HeartBeatInterval = math.MaxUint16
		}
		if msg.BasicParam.HeartBeatCount <= 0 {
			msg.BasicParam.HeartBeatCount = 1
		}
		// 计算设备离线超时时间
		if msg.BasicParam.HeartBeatInterval*msg.BasicParam.HeartBeatCount > 0 {
			ipc.keepaliveInterval = uint16(msg.BasicParam.HeartBeatInterval) // nolint
			ipc.keepaliveTimeout = uint16(msg.BasicParam.HeartBeatCount)     // nolint
			ctx.Log.Debug("sipMessageConfigDownload update", "deviceID", ctx.DeviceID, "keepaliveInterval", ipc.keepaliveInterval, "keepaliveTimeout", ipc.keepaliveTimeout)
		}
	}

	ctx.String(200, "OK")
}
