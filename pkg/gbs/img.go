package gbs

import (
	"log/slog"

	"github.com/gowvp/owl/pkg/gbs/sip"
)

func (g *GB28181API) QuerySnapshot(deviceID, channelID string) error {
	slog.Debug("QuerySnapshot", "deviceID", deviceID)
	ipc, ok := g.svr.memoryStorer.Load(deviceID)
	if !ok {
		return ErrDeviceOffline
	}

	body := NewDeviceConfig(channelID).SetSnapShotConfig(&SnapShot{
		SnapNum:   1,
		Interval:  1,
		UploadURL: "http://192.168.10.31:15123/gb28181/snapshot",
		SessionID: "1234567890",
	}).Marshal()

	tx, err := g.svr.wrapRequest(ipc, sip.MethodMessage, &sip.ContentTypeXML, body)
	if err != nil {
		return err
	}
	_, err = sipResponse(tx)
	return err
}
