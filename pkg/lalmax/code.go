package lalmax

type ResCode int64

const (
	CodeSuccess      ResCode = 10000
	CodeInvalidParam ResCode = 10001
	CodeServerBusy   ResCode = 10002

	CodeGroupNotFound      ResCode = 11001
	CodeSessionNotFound    ResCode = 11002
	CodeStartRelayPullFail ResCode = 11003
	CodeStartRelayPushFail ResCode = 11004
	CodeStopRelayPushFail  ResCode = 11005

	CodeGbObServerNotFound    ResCode = 12001
	CodeDeviceNotRegister     ResCode = 12002
	CodeDevicePlayError       ResCode = 12003
	CodeDeviceStopError       ResCode = 12004
	CodeDevicePtzError        ResCode = 12005
	CodeDeviceTalkError       ResCode = 12006
	CodeDeviceRecordListError ResCode = 12007

	CodeGetRoomParticipantFail  ResCode = 13001
	CodeConnectRoomFail         ResCode = 13002
	CodeListRoomParticipantFail ResCode = 13003
	CodeParticipantExist        ResCode = 13004
	CodeParticipantNotFound     ResCode = 13005
	CodeStartRoomPublishFail    ResCode = 13006
	CodeStopRoomPublishFail     ResCode = 13007

	CodeOnvifObServerNotFound       ResCode = 14001
	CodeOnvifDiscoverDeviceListFail ResCode = 14002
	CodeOnvifGetRtspPlayInfoFail    ResCode = 14003
	CodeOnvifGetSnapInfoFail        ResCode = 14004
	CodeOnvifGetPTZCapabilitiesFail ResCode = 14005
	CodeOnvifAddDeviceFail          ResCode = 14006
	CodeOnvifGetDevicesFail         ResCode = 14007
	CodeOnvifDeleteDeviceFail       ResCode = 14008
	CodeOnvifPtzDirectionFail       ResCode = 14009
	CodeOnvifPtzStopFail            ResCode = 14010
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:                     "success",
	CodeInvalidParam:                "请求参数错误",
	CodeServerBusy:                  "服务繁忙",
	CodeGroupNotFound:               "group不存在",
	CodeSessionNotFound:             "session不存在",
	CodeStartRelayPullFail:          "relay pull 失败",
	CodeGbObServerNotFound:          "gb 服务没有启动",
	CodeDeviceNotRegister:           "设备未注册",
	CodeDevicePlayError:             "gb播放错误",
	CodeDeviceStopError:             "gb停止错误",
	CodeDevicePtzError:              "gb ptz操作错误",
	CodeGetRoomParticipantFail:      "获取房间参与者失败",
	CodeConnectRoomFail:             "连接房间失败",
	CodeListRoomParticipantFail:     "获取房间参与者列表失败",
	CodeParticipantExist:            "参与者已存在",
	CodeParticipantNotFound:         "参与者不存在",
	CodeDeviceTalkError:             "对讲操作失败",
	CodeOnvifObServerNotFound:       "onvif 服务没有启动",
	CodeOnvifDiscoverDeviceListFail: "onvif 设备发现失败",
	CodeOnvifGetRtspPlayInfoFail:    "获取 onvif rtsp 播放地址失败",
	CodeOnvifGetSnapInfoFail:        "获取 onvif 快照地址失败",
	CodeOnvifGetPTZCapabilitiesFail: "获取 onvif ptz 能力失败",
	CodeOnvifAddDeviceFail:          "添加 onvif 设备失败",
	CodeOnvifGetDevicesFail:         "获取 onvif 设备列表失败",
	CodeOnvifDeleteDeviceFail:       "删除 onvif 设备失败",
	CodeOnvifPtzDirectionFail:       "onvif ptz 方向控制失败",
	CodeOnvifPtzStopFail:            "onvif ptz 停止控制失败",
}

type FixedHeader struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"` // 仅 code 发生错误时，此参数才有效
}
