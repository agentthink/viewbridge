package gbs

import "encoding/xml"

const snapShotConfig = "SnapShotConfig" // 图像抓拍配置

// 设备配置 A.2.3.2.1
type DeviceConfigRequest struct {
	XMLName        xml.Name  `xml:"Control"`
	CmdType        string    `xml:"CmdType"`  // 命令类型：设备配置查询(必选)
	SN             int32     `xml:"SN"`       // 命令序列号(必选)
	DeviceID       string    `xml:"DeviceID"` // 目标设备编码(必选)
	SnapShotConfig *SnapShot `xml:"SnapShotConfig"`
}

func NewDeviceConfig(deviceID string) *DeviceConfigRequest {
	return &DeviceConfigRequest{
		CmdType:  "DeviceConfig",
		SN:       1,
		DeviceID: deviceID,
	}
}

func (d *DeviceConfigRequest) SetSN(sn int32) *DeviceConfigRequest {
	d.SN = sn
	return d
}

func (d *DeviceConfigRequest) SetSnapShotConfig(snapShot *SnapShot) *DeviceConfigRequest {
	d.SnapShotConfig = snapShot
	return d
}

func (d *DeviceConfigRequest) Marshal() []byte {
	b, _ := xml.Marshal(d)
	return b
}
