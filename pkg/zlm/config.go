package zlm

import "encoding/json"

const (
	getServerConfig = "/index/api/getServerConfig" // 获取配置
	setServerConfig = "/index/api/setServerConfig" // 设置配置
	restartServer   = "/index/api/restartServer"   // 重启 MediaServer 进程
)

type FixedHeader struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"` // 仅 code 发生错误时，此参数才有效
}

type SetServerConfigReponse struct {
	FixedHeader
	Changed int `json:"changed"`
}

type GetServerConfigResponse struct {
	FixedHeader
	Data []GetServerConfigData `json:"data"`
}

type GetServerConfigData struct {
	APIAPIDebug                        string `json:"api.apiDebug"`
	APIDefaultSnap                     string `json:"api.defaultSnap"`
	APIDownloadRoot                    string `json:"api.downloadRoot"`
	APISecret                          string `json:"api.secret"`
	APISnapRoot                        string `json:"api.snapRoot"`
	ClusterOriginURL                   string `json:"cluster.origin_url"`
	ClusterRetryCount                  string `json:"cluster.retry_count"`
	ClusterTimeoutSec                  string `json:"cluster.timeout_sec"`
	FfmpegBin                          string `json:"ffmpeg.bin"`
	FfmpegCmd                          string `json:"ffmpeg.cmd"`
	FfmpegLog                          string `json:"ffmpeg.log"`
	FfmpegRestartSec                   string `json:"ffmpeg.restart_sec"`
	FfmpegSnap                         string `json:"ffmpeg.snap"`
	GeneralBroadcastPlayerCountChanged string `json:"general.broadcast_player_count_changed"`
	GeneralCheckNvidiaDev              string `json:"general.check_nvidia_dev"`
	GeneralEnableVhost                 string `json:"general.enableVhost"`
	GeneralEnableFfmpegLog             string `json:"general.enable_ffmpeg_log"`
	GeneralFlowThreshold               string `json:"general.flowThreshold"`
	GeneralListenIP                    string `json:"general.listen_ip"`
	GeneralMaxStreamWaitMS             string `json:"general.maxStreamWaitMS"`
	GeneralMediaServerID               string `json:"general.mediaServerId"`
	GeneralMergeWriteMS                string `json:"general.mergeWriteMS"`
	GeneralResetWhenRePlay             string `json:"general.resetWhenRePlay"`
	GeneralStreamNoneReaderDelayMS     string `json:"general.streamNoneReaderDelayMS"`
	GeneralUnreadyFrameCache           string `json:"general.unready_frame_cache"`
	GeneralWaitAddTrackMs              string `json:"general.wait_add_track_ms"`
	GeneralWaitAudioTrackDataMs        string `json:"general.wait_audio_track_data_ms"`
	GeneralWaitTrackReadyMs            string `json:"general.wait_track_ready_ms"`
	HlsBroadcastRecordTs               string `json:"hls.broadcastRecordTs"`
	HlsDeleteDelaySec                  string `json:"hls.deleteDelaySec"`
	HlsFastRegister                    string `json:"hls.fastRegister"`
	HlsFileBufSize                     string `json:"hls.fileBufSize"`
	HlsSegDelay                        string `json:"hls.segDelay"`
	HlsSegDur                          string `json:"hls.segDur"`
	HlsSegKeep                         string `json:"hls.segKeep"`
	HlsSegNum                          string `json:"hls.segNum"`
	HlsSegRetain                       string `json:"hls.segRetain"`
	HookAliveInterval                  string `json:"hook.alive_interval"`
	HookEnable                         string `json:"hook.enable"`
	HookOnFlowReport                   string `json:"hook.on_flow_report"`
	HookOnHTTPAccess                   string `json:"hook.on_http_access"`
	HookOnPlay                         string `json:"hook.on_play"`
	HookOnPublish                      string `json:"hook.on_publish"`
	HookOnRecordMp4                    string `json:"hook.on_record_mp4"`
	HookOnRecordTs                     string `json:"hook.on_record_ts"`
	HookOnRtpServerTimeout             string `json:"hook.on_rtp_server_timeout"`
	HookOnRtspAuth                     string `json:"hook.on_rtsp_auth"`
	HookOnRtspRealm                    string `json:"hook.on_rtsp_realm"`
	HookOnSendRtpStopped               string `json:"hook.on_send_rtp_stopped"`
	HookOnServerExited                 string `json:"hook.on_server_exited"`
	HookOnServerKeepalive              string `json:"hook.on_server_keepalive"`
	HookOnServerStarted                string `json:"hook.on_server_started"`
	HookOnShellLogin                   string `json:"hook.on_shell_login"`
	HookOnStreamChanged                string `json:"hook.on_stream_changed"`
	HookOnStreamNoneReader             string `json:"hook.on_stream_none_reader"`
	HookOnStreamNotFound               string `json:"hook.on_stream_not_found"`
	HookRetry                          string `json:"hook.retry"`
	HookRetryDelay                     string `json:"hook.retry_delay"`
	HookStreamChangedSchemas           string `json:"hook.stream_changed_schemas"`
	HookTimeoutSec                     string `json:"hook.timeoutSec"`
	HTTPAllowCrossDomains              string `json:"http.allow_cross_domains"`
	HTTPAllowIPRange                   string `json:"http.allow_ip_range"`
	HTTPCharSet                        string `json:"http.charSet"`
	HTTPDirMenu                        string `json:"http.dirMenu"`
	HTTPForbidCacheSuffix              string `json:"http.forbidCacheSuffix"`
	HTTPForwardedIPHeader              string `json:"http.forwarded_ip_header"`
	HTTPKeepAliveSecond                string `json:"http.keepAliveSecond"`
	HTTPMaxReqSize                     string `json:"http.maxReqSize"`
	HTTPNotFound                       string `json:"http.notFound"`
	HTTPPort                           int    `json:"http.port,string"`
	HTTPRootPath                       string `json:"http.rootPath"`
	HTTPSendBufSize                    string `json:"http.sendBufSize"`
	HTTPSslport                        int    `json:"http.sslport,string"`
	HTTPVirtualPath                    string `json:"http.virtualPath"`
	MulticastAddrMax                   string `json:"multicast.addrMax"`
	MulticastAddrMin                   string `json:"multicast.addrMin"`
	MulticastUDPTTL                    string `json:"multicast.udpTTL"`
	ProtocolAddMuteAudio               string `json:"protocol.add_mute_audio"`
	ProtocolAutoClose                  string `json:"protocol.auto_close"`
	ProtocolContinuePushMs             string `json:"protocol.continue_push_ms"`
	ProtocolEnableAudio                string `json:"protocol.enable_audio"`
	ProtocolEnableFmp4                 string `json:"protocol.enable_fmp4"`
	ProtocolEnableHls                  string `json:"protocol.enable_hls"`
	ProtocolEnableHlsFmp4              string `json:"protocol.enable_hls_fmp4"`
	ProtocolEnableMp4                  string `json:"protocol.enable_mp4"`
	ProtocolEnableRtmp                 string `json:"protocol.enable_rtmp"`
	ProtocolEnableRtsp                 string `json:"protocol.enable_rtsp"`
	ProtocolEnableTs                   string `json:"protocol.enable_ts"`
	ProtocolFmp4Demand                 string `json:"protocol.fmp4_demand"`
	ProtocolHlsDemand                  string `json:"protocol.hls_demand"`
	ProtocolHlsSavePath                string `json:"protocol.hls_save_path"`
	ProtocolModifyStamp                string `json:"protocol.modify_stamp"`
	ProtocolMp4AsPlayer                string `json:"protocol.mp4_as_player"`
	ProtocolMp4MaxSecond               string `json:"protocol.mp4_max_second"`
	ProtocolMp4SavePath                string `json:"protocol.mp4_save_path"`
	ProtocolPacedSenderMs              string `json:"protocol.paced_sender_ms"`
	ProtocolRtmpDemand                 string `json:"protocol.rtmp_demand"`
	ProtocolRtspDemand                 string `json:"protocol.rtsp_demand"`
	ProtocolTsDemand                   string `json:"protocol.ts_demand"`
	RecordAppName                      string `json:"record.appName"`
	RecordEnableFmp4                   string `json:"record.enableFmp4"`
	RecordFastStart                    string `json:"record.fastStart"`
	RecordFileBufSize                  string `json:"record.fileBufSize"`
	RecordFileRepeat                   string `json:"record.fileRepeat"`
	RecordSampleMS                     string `json:"record.sampleMS"`
	RtcDatachannelEcho                 string `json:"rtc.datachannel_echo"`
	RtcExternIP                        string `json:"rtc.externIP"`
	RtcMaxRtpCacheMS                   string `json:"rtc.maxRtpCacheMS"`
	RtcMaxRtpCacheSize                 string `json:"rtc.maxRtpCacheSize"`
	RtcMaxBitrate                      string `json:"rtc.max_bitrate"`
	RtcMinBitrate                      string `json:"rtc.min_bitrate"`
	RtcNackIntervalRatio               string `json:"rtc.nackIntervalRatio"`
	RtcNackMaxCount                    string `json:"rtc.nackMaxCount"`
	RtcNackMaxMS                       string `json:"rtc.nackMaxMS"`
	RtcNackMaxSize                     string `json:"rtc.nackMaxSize"`
	RtcNackRtpSize                     string `json:"rtc.nackRtpSize"`
	RtcPort                            string `json:"rtc.port"`
	RtcPreferredCodecA                 string `json:"rtc.preferredCodecA"`
	RtcPreferredCodecV                 string `json:"rtc.preferredCodecV"`
	RtcRembBitRate                     string `json:"rtc.rembBitRate"`
	RtcStartBitrate                    string `json:"rtc.start_bitrate"`
	RtcTCPPort                         string `json:"rtc.tcpPort"`
	RtcTimeoutSec                      string `json:"rtc.timeoutSec"`
	RtmpDirectProxy                    string `json:"rtmp.directProxy"`
	RtmpEnhanced                       string `json:"rtmp.enhanced"`
	RtmpHandshakeSecond                string `json:"rtmp.handshakeSecond"`
	RtmpKeepAliveSecond                string `json:"rtmp.keepAliveSecond"`
	RtmpPort                           int    `json:"rtmp.port,string"`
	RtmpSslport                        int    `json:"rtmp.sslport,string"`
	RtpAudioMtuSize                    string `json:"rtp.audioMtuSize"`
	RtpH264StapA                       string `json:"rtp.h264_stap_a"`
	RtpLowLatency                      string `json:"rtp.lowLatency"`
	RtpRtpMaxSize                      string `json:"rtp.rtpMaxSize"`
	RtpVideoMtuSize                    string `json:"rtp.videoMtuSize"`
	RtpProxyDumpDir                    string `json:"rtp_proxy.dumpDir"`
	RtpProxyGopCache                   string `json:"rtp_proxy.gop_cache"`
	RtpProxyH264Pt                     string `json:"rtp_proxy.h264_pt"`
	RtpProxyH265Pt                     string `json:"rtp_proxy.h265_pt"`
	RtpProxyOpusPt                     string `json:"rtp_proxy.opus_pt"`
	RtpProxyPort                       int    `json:"rtp_proxy.port,string"`
	RtpProxyPortRange                  string `json:"rtp_proxy.port_range"`
	RtpProxyPsPt                       string `json:"rtp_proxy.ps_pt"`
	RtpProxyRtpG711DurMs               string `json:"rtp_proxy.rtp_g711_dur_ms"`
	RtpProxyTimeoutSec                 string `json:"rtp_proxy.timeoutSec"`
	RtpProxyUDPRecvSocketBuffer        string `json:"rtp_proxy.udp_recv_socket_buffer"`
	RtspAuthBasic                      string `json:"rtsp.authBasic"`
	RtspDirectProxy                    string `json:"rtsp.directProxy"`
	RtspHandshakeSecond                string `json:"rtsp.handshakeSecond"`
	RtspKeepAliveSecond                string `json:"rtsp.keepAliveSecond"`
	RtspLowLatency                     string `json:"rtsp.lowLatency"`
	RtspPort                           int    `json:"rtsp.port,string"`
	RtspRtpTransportType               string `json:"rtsp.rtpTransportType"`
	RtspSslport                        int    `json:"rtsp.sslport,string"`
	ShellMaxReqSize                    string `json:"shell.maxReqSize"`
	ShellPort                          string `json:"shell.port"`
	SrtLatencyMul                      string `json:"srt.latencyMul"`
	SrtPktBufSize                      string `json:"srt.pktBufSize"`
	SrtPort                            string `json:"srt.port"`
	SrtTimeoutSec                      string `json:"srt.timeoutSec"`
}

// SetServerConfigRequest
// https://github.com/zlmediakit/ZLMediaKit/wiki/MediaServer%E6%94%AF%E6%8C%81%E7%9A%84HTTP-HOOK-API
type SetServerConfigRequest struct {
	APIAPIDebug                        *string `json:"api.apiDebug,omitempty"`
	APIDefaultSnap                     *string `json:"api.defaultSnap,omitempty"`
	APIDownloadRoot                    *string `json:"api.downloadRoot,omitempty"`
	APISecret                          *string `json:"api.secret,omitempty"`
	APISnapRoot                        *string `json:"api.snapRoot,omitempty"`
	ClusterOriginURL                   *string `json:"cluster.origin_url,omitempty"`
	ClusterRetryCount                  *string `json:"cluster.retry_count,omitempty"`
	ClusterTimeoutSec                  *string `json:"cluster.timeout_sec,omitempty"`
	FfmpegBin                          *string `json:"ffmpeg.bin,omitempty"`
	FfmpegCmd                          *string `json:"ffmpeg.cmd,omitempty"`
	FfmpegLog                          *string `json:"ffmpeg.log,omitempty"`
	FfmpegRestartSec                   *string `json:"ffmpeg.restart_sec,omitempty"`
	FfmpegSnap                         *string `json:"ffmpeg.snap,omitempty"`
	GeneralBroadcastPlayerCountChanged *string `json:"general.broadcast_player_count_changed,omitempty"`
	GeneralCheckNvidiaDev              *string `json:"general.check_nvidia_dev,omitempty"`
	GeneralEnableVhost                 *string `json:"general.enableVhost,omitempty"`
	GeneralEnableFfmpegLog             *string `json:"general.enable_ffmpeg_log,omitempty"`
	GeneralFlowThreshold               *string `json:"general.flowThreshold,omitempty"`
	GeneralListenIP                    *string `json:"general.listen_ip,omitempty"`
	GeneralMaxStreamWaitMS             *string `json:"general.maxStreamWaitMS,omitempty"`
	GeneralMediaServerID               *string `json:"general.mediaServerId,omitempty"`
	GeneralMergeWriteMS                *string `json:"general.mergeWriteMS,omitempty"`
	GeneralResetWhenRePlay             *string `json:"general.resetWhenRePlay,omitempty"`
	GeneralStreamNoneReaderDelayMS     *string `json:"general.streamNoneReaderDelayMS,omitempty"`
	GeneralUnreadyFrameCache           *string `json:"general.unready_frame_cache,omitempty"`
	GeneralWaitAddTrackMs              *string `json:"general.wait_add_track_ms,omitempty"`
	GeneralWaitAudioTrackDataMs        *string `json:"general.wait_audio_track_data_ms,omitempty"`
	GeneralWaitTrackReadyMs            *string `json:"general.wait_track_ready_ms,omitempty"`
	HlsBroadcastRecordTs               *string `json:"hls.broadcastRecordTs,omitempty"`
	HlsDeleteDelaySec                  *string `json:"hls.deleteDelaySec,omitempty"`
	HlsFastRegister                    *string `json:"hls.fastRegister,omitempty"`
	HlsFileBufSize                     *string `json:"hls.fileBufSize,omitempty"`
	HlsSegDelay                        *string `json:"hls.segDelay,omitempty"`
	HlsSegDur                          *string `json:"hls.segDur,omitempty"`
	HlsSegKeep                         *string `json:"hls.segKeep,omitempty"`
	HlsSegNum                          *string `json:"hls.segNum,omitempty"`
	HlsSegRetain                       *string `json:"hls.segRetain,omitempty"`
	HookAliveInterval                  *string `json:"hook.alive_interval,omitempty"`
	HookEnable                         *string `json:"hook.enable,omitempty"`
	HookOnFlowReport                   *string `json:"hook.on_flow_report,omitempty"`
	HookOnHTTPAccess                   *string `json:"hook.on_http_access,omitempty"`
	HookOnPlay                         *string `json:"hook.on_play,omitempty"`
	HookOnPublish                      *string `json:"hook.on_publish,omitempty"`
	HookOnRecordMp4                    *string `json:"hook.on_record_mp4,omitempty"`
	HookOnRecordTs                     *string `json:"hook.on_record_ts,omitempty"`
	HookOnRtpServerTimeout             *string `json:"hook.on_rtp_server_timeout,omitempty"`
	HookOnRtspAuth                     *string `json:"hook.on_rtsp_auth,omitempty"`
	HookOnRtspRealm                    *string `json:"hook.on_rtsp_realm,omitempty"`
	HookOnSendRtpStopped               *string `json:"hook.on_send_rtp_stopped,omitempty"`
	HookOnServerExited                 *string `json:"hook.on_server_exited,omitempty"`
	HookOnServerKeepalive              *string `json:"hook.on_server_keepalive,omitempty"`
	HookOnServerStarted                *string `json:"hook.on_server_started,omitempty"`
	HookOnShellLogin                   *string `json:"hook.on_shell_login,omitempty"`
	HookOnStreamChanged                *string `json:"hook.on_stream_changed,omitempty"`
	HookOnStreamNoneReader             *string `json:"hook.on_stream_none_reader,omitempty"`
	HookOnStreamNotFound               *string `json:"hook.on_stream_not_found,omitempty"`
	HookRetry                          *string `json:"hook.retry,omitempty"`
	HookRetryDelay                     *string `json:"hook.retry_delay,omitempty"`
	HookStreamChangedSchemas           *string `json:"hook.stream_changed_schemas,omitempty"`
	HookTimeoutSec                     *string `json:"hook.timeoutSec,omitempty"` // 事件触发 http post 超时时间。
	HTTPAllowCrossDomains              *string `json:"http.allow_cross_domains,omitempty"`
	HTTPAllowIPRange                   *string `json:"http.allow_ip_range,omitempty"`
	HTTPCharSet                        *string `json:"http.charSet,omitempty"`
	HTTPDirMenu                        *string `json:"http.dirMenu,omitempty"`
	HTTPForbidCacheSuffix              *string `json:"http.forbidCacheSuffix,omitempty"`
	HTTPForwardedIPHeader              *string `json:"http.forwarded_ip_header,omitempty"`
	HTTPKeepAliveSecond                *string `json:"http.keepAliveSecond,omitempty"`
	HTTPMaxReqSize                     *string `json:"http.maxReqSize,omitempty"`
	HTTPNotFound                       *string `json:"http.notFound,omitempty"`
	HTTPPort                           *string `json:"http.port,omitempty"`
	HTTPRootPath                       *string `json:"http.rootPath,omitempty"`
	HTTPSendBufSize                    *string `json:"http.sendBufSize,omitempty"`
	HTTPSslport                        *string `json:"http.sslport,omitempty"`
	HTTPVirtualPath                    *string `json:"http.virtualPath,omitempty"`
	MulticastAddrMax                   *string `json:"multicast.addrMax,omitempty"`
	MulticastAddrMin                   *string `json:"multicast.addrMin,omitempty"`
	MulticastUDPTTL                    *string `json:"multicast.udpTTL,omitempty"`
	ProtocolAddMuteAudio               *string `json:"protocol.add_mute_audio,omitempty"`
	ProtocolAutoClose                  *string `json:"protocol.auto_close,omitempty"`
	ProtocolContinuePushMs             *string `json:"protocol.continue_push_ms,omitempty"`
	ProtocolEnableAudio                *string `json:"protocol.enable_audio,omitempty"`
	ProtocolEnableFmp4                 *string `json:"protocol.enable_fmp4,omitempty"`
	ProtocolEnableHls                  *string `json:"protocol.enable_hls,omitempty"`
	ProtocolEnableHlsFmp4              *string `json:"protocol.enable_hls_fmp4,omitempty"`
	ProtocolEnableMp4                  *string `json:"protocol.enable_mp4,omitempty"`
	ProtocolEnableRtmp                 *string `json:"protocol.enable_rtmp,omitempty"`
	ProtocolEnableRtsp                 *string `json:"protocol.enable_rtsp,omitempty"`
	ProtocolEnableTs                   *string `json:"protocol.enable_ts,omitempty"`
	ProtocolFmp4Demand                 *string `json:"protocol.fmp4_demand,omitempty"`
	ProtocolHlsDemand                  *string `json:"protocol.hls_demand,omitempty"`
	ProtocolHlsSavePath                *string `json:"protocol.hls_save_path,omitempty"`
	ProtocolModifyStamp                *string `json:"protocol.modify_stamp,omitempty"`
	ProtocolMp4AsPlayer                *string `json:"protocol.mp4_as_player,omitempty"`
	ProtocolMp4MaxSecond               *string `json:"protocol.mp4_max_second,omitempty"`
	ProtocolMp4SavePath                *string `json:"protocol.mp4_save_path,omitempty"`
	ProtocolPacedSenderMs              *string `json:"protocol.paced_sender_ms,omitempty"`
	ProtocolRtmpDemand                 *string `json:"protocol.rtmp_demand,omitempty"`
	ProtocolRtspDemand                 *string `json:"protocol.rtsp_demand,omitempty"`
	ProtocolTsDemand                   *string `json:"protocol.ts_demand,omitempty"`
	RecordAppName                      *string `json:"record.appName,omitempty"`
	RecordEnableFmp4                   *string `json:"record.enableFmp4,omitempty"`
	RecordFastStart                    *string `json:"record.fastStart,omitempty"`
	RecordFileBufSize                  *string `json:"record.fileBufSize,omitempty"`
	RecordFileRepeat                   *string `json:"record.fileRepeat,omitempty"`
	RecordSampleMS                     *string `json:"record.sampleMS,omitempty"`
	RtcDatachannelEcho                 *string `json:"rtc.datachannel_echo,omitempty"`
	RtcExternIP                        *string `json:"rtc.externIP,omitempty"`
	RtcMaxRtpCacheMS                   *string `json:"rtc.maxRtpCacheMS,omitempty"`
	RtcMaxRtpCacheSize                 *string `json:"rtc.maxRtpCacheSize,omitempty"`
	RtcMaxBitrate                      *string `json:"rtc.max_bitrate,omitempty"`
	RtcMinBitrate                      *string `json:"rtc.min_bitrate,omitempty"`
	RtcNackIntervalRatio               *string `json:"rtc.nackIntervalRatio,omitempty"`
	RtcNackMaxCount                    *string `json:"rtc.nackMaxCount,omitempty"`
	RtcNackMaxMS                       *string `json:"rtc.nackMaxMS,omitempty"`
	RtcNackMaxSize                     *string `json:"rtc.nackMaxSize,omitempty"`
	RtcNackRtpSize                     *string `json:"rtc.nackRtpSize,omitempty"`
	RtcPort                            *string `json:"rtc.port,omitempty"`
	RtcPreferredCodecA                 *string `json:"rtc.preferredCodecA,omitempty"`
	RtcPreferredCodecV                 *string `json:"rtc.preferredCodecV,omitempty"`
	RtcRembBitRate                     *string `json:"rtc.rembBitRate,omitempty"`
	RtcStartBitrate                    *string `json:"rtc.start_bitrate,omitempty"`
	RtcTCPPort                         *string `json:"rtc.tcpPort,omitempty"`
	RtcTimeoutSec                      *string `json:"rtc.timeoutSec,omitempty"`
	// RTCEnableTurn                      *string `json:"rtc.enableTurn,omitempty"`
	RtmpDirectProxy             *string `json:"rtmp.directProxy,omitempty"`
	RtmpEnhanced                *string `json:"rtmp.enhanced,omitempty"`
	RtmpHandshakeSecond         *string `json:"rtmp.handshakeSecond,omitempty"`
	RtmpKeepAliveSecond         *string `json:"rtmp.keepAliveSecond,omitempty"`
	RtmpPort                    *string `json:"rtmp.port,omitempty"`
	RtmpSslport                 *string `json:"rtmp.sslport,omitempty"`
	RtpAudioMtuSize             *string `json:"rtp.audioMtuSize,omitempty"`
	RtpH264StapA                *string `json:"rtp.h264_stap_a,omitempty"`
	RtpLowLatency               *string `json:"rtp.lowLatency,omitempty"`
	RtpRtpMaxSize               *string `json:"rtp.rtpMaxSize,omitempty"`
	RtpVideoMtuSize             *string `json:"rtp.videoMtuSize,omitempty"`
	RtpProxyDumpDir             *string `json:"rtp_proxy.dumpDir,omitempty"`
	RtpProxyGopCache            *string `json:"rtp_proxy.gop_cache,omitempty"`
	RtpProxyH264Pt              *string `json:"rtp_proxy.h264_pt,omitempty"`
	RtpProxyH265Pt              *string `json:"rtp_proxy.h265_pt,omitempty"`
	RtpProxyOpusPt              *string `json:"rtp_proxy.opus_pt,omitempty"`
	RtpProxyPort                *string `json:"rtp_proxy.port,omitempty"`
	RtpProxyPortRange           *string `json:"rtp_proxy.port_range,omitempty"`
	RtpProxyPsPt                *string `json:"rtp_proxy.ps_pt,omitempty"`
	RtpProxyRtpG711DurMs        *string `json:"rtp_proxy.rtp_g711_dur_ms,omitempty"`
	RtpProxyTimeoutSec          *string `json:"rtp_proxy.timeoutSec,omitempty"`
	RtpProxyUDPRecvSocketBuffer *string `json:"rtp_proxy.udp_recv_socket_buffer,omitempty"`
	RtspAuthBasic               *string `json:"rtsp.authBasic,omitempty"`
	RtspDirectProxy             *string `json:"rtsp.directProxy,omitempty"`
	RtspHandshakeSecond         *string `json:"rtsp.handshakeSecond,omitempty"`
	RtspKeepAliveSecond         *string `json:"rtsp.keepAliveSecond,omitempty"`
	RtspLowLatency              *string `json:"rtsp.lowLatency,omitempty"`
	RtspPort                    *string `json:"rtsp.port,omitempty"`
	RtspRtpTransportType        *string `json:"rtsp.rtpTransportType,omitempty"`
	RtspSslport                 *string `json:"rtsp.sslport,omitempty"`
	ShellMaxReqSize             *string `json:"shell.maxReqSize,omitempty"`
	ShellPort                   *string `json:"shell.port,omitempty"`
	SrtLatencyMul               *string `json:"srt.latencyMul,omitempty"`
	SrtPktBufSize               *string `json:"srt.pktBufSize,omitempty"`
	SrtPort                     *string `json:"srt.port,omitempty"`
	SrtTimeoutSec               *string `json:"srt.timeoutSec,omitempty"`
}

func (e *Engine) GetServerConfig() (*GetServerConfigResponse, error) {
	var resp GetServerConfigResponse
	if err := e.post(getServerConfig, nil, &resp); err != nil {
		return nil, err
	}
	if err := e.ErrHandle(resp.Code, resp.Msg); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (e *Engine) SetServerConfig(in *SetServerConfigRequest) (*SetServerConfigReponse, error) {
	req, err := struct2map(in)
	if err != nil {
		return nil, err
	}
	var resp SetServerConfigReponse
	if err := e.post(setServerConfig, req, &resp); err != nil {
		return nil, err
	}
	if err := e.ErrHandle(resp.Code, resp.Msg); err != nil {
		return nil, err
	}

	return &resp, nil
}

// RestartServer 重启 ZLM MediaServer 进程。
//
// 为什么: rtc.port / rtc.tcpPort 等网络监听端口在 SetServerConfig 后仅更新内存配置, 不会重绑套接字,
// 必须重启进程才能让新端口生效; 典型场景是 gowvp 启动时把 rtc.port 合并到 RTP 段首位与 ZLM 默认不一致。
// ZLM 重启后会持久化的配置, 等价于启动期一次性生效, 运行期不应频繁调用。
func (e *Engine) RestartServer() error {
	var resp FixedHeader
	if err := e.post(restartServer, nil, &resp); err != nil {
		return err
	}
	return e.ErrHandle(resp.Code, resp.Msg)
}

func struct2map(in any) (map[string]any, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}
