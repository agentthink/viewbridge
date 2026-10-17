package ipc

// ChannelDisplayName 返回通道展示名。
func ChannelDisplayName(ch *Channel) string {
	if ch == nil {
		return ""
	}
	if ch.Name != "" {
		return ch.Name
	}
	if ch.ChannelID != "" {
		return ch.ChannelID
	}
	return ch.ID
}

// ProfileDisplayName 返回 ONVIF Profile 展示名：通道名-设备名。
func ProfileDisplayName(ch *Channel, deviceName string) string {
	return ChannelDisplayName(ch) + "-" + deviceName
}

// ChannelReachable 判断通道是否可用：通道与所属设备均在线才算在线（与设备列表里对子通道的处理一致）。
func ChannelReachable(ch *Channel, deviceOnline bool) bool {
	if ch == nil {
		return false
	}
	return ch.IsOnline && deviceOnline
}

// ONVIFProfileName 返回 ONVIF/HA 上展示的 Profile 名。
// RTSP 拉流通道用 BUSY/IDLE 语义，不在名称标离线；其余类型需通道与设备同时在线，否则加「(离线)」。
func ONVIFProfileName(ch *Channel, deviceName string, deviceOnline bool) string {
	name := ProfileDisplayName(ch, deviceName)
	if ch != nil && ch.IsRTSP() {
		return name
	}
	if !ChannelReachable(ch, deviceOnline) {
		name += " (离线)"
	}
	return name
}
