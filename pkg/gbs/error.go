package gbs

import "errors"

var (
	ErrXMLDecode = errors.New("xml decode error")
	ErrDatabase  = errors.New("database error")
)

var (
	ErrDeviceNotExist  = errors.New("device not exist")
	ErrChannelNotExist = errors.New("channel not exist")

	ErrDeviceOffline  = errors.New("device offline")
	ErrChannelOffline = errors.New("channel offline")
)
