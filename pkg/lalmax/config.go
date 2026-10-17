package lalmax

import (
	"context"
	"fmt"
)

type ServerConfig struct {
	FixedHeader
	Data GetServerConfigResponse `json:"data"`
}

type GetServerConfigResponse struct {
	ConfVersion                 string `json:"conf_version"`
	CheckSessionDisposeInterval uint32 `json:"check_session_dispose_interval"`
	UpdateSessionStateInterval  uint32 `json:"update_session_state_interval"`
	ManagerChanSize             uint32 `json:"manager_chan_size"`
	AdjustPts                   bool   `json:"adjust_pts"`
	MaxOpenFiles                uint64 `json:"max_open_files"`
	// 关键帧存储路径
	KeyFramePath          string                `json:"key_frame_path"`
	KeyFrameIntervalSec   int                   `json:"key_frame_interval_sec"` // 关键帧间隔秒数
	GopCacheConfig        GopCacheConfig        `json:"gop_cache_config"`
	RtmpConfig            RtmpConfig            `json:"rtmp"`
	InSessionConfig       InSessionConfig       `json:"in_session"`
	DefaultHttpConfig     DefaultHttpConfig     `json:"default_http"`
	HttpflvConfig         HttpflvConfig         `json:"httpflv"`
	HlsConfig             HlsConfig             `json:"hls"`
	HttptsConfig          HttptsConfig          `json:"httpts"`
	RtspConfig            RtspConfig            `json:"rtsp"`
	SrtConfig             SrtConfig             `json:"srt"`
	DashConfig            DashConfig            `json:"dash"`
	RtcConfig             RtcConfig             `json:"rtc"`
	Gb28181Config         SipConfig             `json:"gb28181"`
	Jt1078Config          Jt1078Config          `json:"jt1078"`
	OnvifConfig           OnvifConfig           `json:"onvif"`
	HttpFmp4Config        HttpFmp4Config        `json:"httpfmp4"`
	RoomConfig            RoomConfig            `json:"room"`
	RecordConfig          RecordConfig          `json:"record"`
	PlaybackConfig        PlaybackConfig        `json:"playback"`
	MetricsConfig         MetricsConfig         `json:"metrics"`
	RelayPushConfig       RelayPushConfig       `json:"relay_push"`
	StaticRelayPullConfig StaticRelayPullConfig `json:"static_relay_pull"`

	HttpApiConfig    HttpApiConfig    `json:"http_api"`
	ServerId         string           `json:"server_id"`
	HttpNotifyConfig HttpNotifyConfig `json:"http_notify"`
	SimpleAuthConfig SimpleAuthConfig `json:"simple_auth"`
	PprofConfig      PprofConfig      `json:"pprof"`
	LogConfig        LogConfig        `json:"log"`
	DebugConfig      DebugConfig      `json:"debug"`

	ReportStatWithFrameRecord bool `json:"report_stat_with_frame_record"`
}

type SipConfig struct {
	Enable                bool        `json:"enable"`                     // gb28181使能标志
	ListenAddr            string      `json:"listen_addr"`                // gb28181监听地址
	SipIP                 string      `json:"sip_ip"`                     // sip 服务器公网IP
	SipPort               uint16      `json:"sip_port"`                   // sip 服务器端口，默认 5060
	Serial                string      `json:"serial"`                     // sip 服务器 id, 默认 34020000002000000001
	Realm                 string      `json:"realm"`                      // sip 服务器域，默认 3402000000
	Username              string      `json:"username"`                   // sip 服务器账号
	Password              string      `json:"password"`                   // sip 服务器密码
	KeepaliveInterval     int         `json:"keepalive_interval"`         // 心跳包时长
	QuickLogin            bool        `json:"quick_login"`                // 快速登陆,有keepalive就认为在线
	SipLogClose           bool        `json:"sip_log_close"`              // 关闭sip日志
	MediaConfig           MediaConfig `json:"media_config"`               // 媒体服务器配置
	PubNotSubAutoCloseSec int         `json:"pub_not_sub_auto_close_sec"` // 发布但无订阅自动关闭秒数
}

// MediaConfig GB28181 媒体配置
type MediaConfig struct {
	MediaIP               string `json:"media_ip"`                 // 媒体服务器IP
	ListenPort            int    `json:"listen_port"`              // 监听端口
	MultiPortMaxIncrement int    `json:"multi_port_max_increment"` // 多端口最大增量
}

type GopCacheConfig struct {
	GopNum               int `json:"gop_cache_num"`
	SingleGopMaxFrameNum int `json:"single_gop_max_frame_num"`
}

type RtmpConfig struct {
	Enable                  bool   `json:"enable"`
	Addr                    string `json:"addr"`
	RtmpsEnable             bool   `json:"rtmps_enable"`
	RtmpsAddr               string `json:"rtmps_addr"`
	RtmpOverQuicEnable      bool   `json:"rtmp_over_quic_enable"`
	RtmpOverQuicAddr        string `json:"rtmp_over_quic_addr"`
	RtmpsCertFile           string `json:"rtmps_cert_file"`
	RtmpsKeyFile            string `json:"rtmps_key_file"`
	RtmpOverKcpEnable       bool   `json:"rtmp_over_kcp_enable"`
	RtmpOverKcpAddr         string `json:"rtmp_over_kcp_addr"`
	RtmpOverKcpDataShards   int    `json:"rtmp_over_kcp_data_shards"`
	RtmpOverKcpParityShards int    `json:"rtmp_over_kcp_parity_shards"`

	MergeWriteSize        int    `json:"merge_write_size"`
	PubTimeoutSec         uint32 `json:"pub_timeout_sec"`
	PullTimeoutSec        uint32 `json:"pull_timeout_sec"`
	PubNotSubAutoCloseSec int    `json:"pub_not_sub_auto_close_sec"` // 发布但无订阅自动关闭秒数
}

type InSessionConfig struct {
	AddDummyAudioEnable      bool `json:"add_dummy_audio_enable"`
	AddDummyAudioWaitAudioMs int  `json:"add_dummy_audio_wait_audio_ms"`
}

type DefaultHttpConfig struct {
	CommonHttpAddrConfig
}

type HttpflvConfig struct {
	CommonHttpServerConfig
}

type HttptsConfig struct {
	CommonHttpServerConfig
}

type HlsConfig struct {
	CommonHttpServerConfig

	UseMemoryAsDiskFlag  bool                   `json:"use_memory_as_disk_flag"`
	DiskUseMmapFlag      bool                   `json:"disk_use_mmap_flag"`
	UseM3u8MemoryFlag    bool                   `json:"use_m3u8_memory_flag"`
	FragmentDurationMs   int                    `json:"fragment_duration_ms"` // 分片时长（毫秒）
	FragmentSize         int                    `json:"fragment_size"`        // 分片大小
	FragmentNum          int                    `json:"fragment_num"`         // 分片数量
	CleanupMode          int                    `json:"cleanup_mode"`         // 清理模式
	DeleteThreshold      int                    `json:"delete_threshold"`     // 删除阈值
	BoundaryByVideo      bool                   `json:"boundary_by_video"`    // 按视频边界
	OutPath              string                 `json:"out_path"`             // 输出路径
	EnableAes128         bool                   `json:"enable_aes128"`        // 是否启用 AES128
	Aes128KeyFile        string                 `json:"aes128_key_file"`      // AES128 密钥文件
	Aes128Iv             string                 `json:"aes128_iv"`            // AES128 初始化向量
	Fmp4InitCombine      bool                   `json:"fmp4_init_combine"`    // FMP4 初始化合并
	EnableSvr            bool                   `json:"enable_svr"`           // 是否启用服务器
	EnableOnDemandHls    bool                   `json:"enable_on_demand_hls"` // 是否启用按需 HLS
	OnDemandTimeoutMs    int                    `json:"on_demand_timeout_ms"` // 按需超时（毫秒）
	SubSessionTimeoutMs  int                    `json:"sub_session_timeout_ms"`
	SubSessionHashKey    string                 `json:"sub_session_hash_key"`
	Fmp4HttpServerConfig CommonHttpServerConfig `json:"fmp4"`
}

type DashConfig struct {
	CommonHttpServerConfig
	UseMemoryAsDiskFlag bool   `json:"use_memory_as_disk_flag"`
	DiskUseMmapFlag     bool   `json:"disk_use_mmap_flag"`
	UseMpdMemoryFlag    bool   `json:"use_mpd_memory_flag"`
	OutPath             string `json:"out_path"`          // 输出路径
	FragmentNum         int    `json:"fragment_num"`      // 分片数量
	FragmentDuration    int    `json:"fragment_duration"` // 分片时长（秒）
	DeleteThreshold     int    `json:"delete_threshold"`  // 删除阈值
}

type RtcConfig struct {
	PubTimeoutSec         uint32 `json:"pub_timeout_sec"`
	PubNotSubAutoCloseSec int    `json:"pub_not_sub_auto_close_sec"` // 发布但无订阅自动关闭秒数
	CommonHttpServerConfig
	IceHostNatToIps []string `json:"iceHostNatToIps"` // ICE 主机 NAT 到 IP
	IceUdpMuxPort   int      `json:"iceUdpMuxPort"`   // ICE UDP 复用端口
	IceTcpMuxPort   int      `json:"iceTcpMuxPort"`   // ICE TCP 复用端口
}

type RtspConfig struct {
	Enable                            bool   `json:"enable"`
	Addr                              string `json:"addr"`
	RtspsEnable                       bool   `json:"rtsps_enable"`
	RtspsAddr                         string `json:"rtsps_addr"`
	RtspsCertFile                     string `json:"rtsps_cert_file"`
	RtspsKeyFile                      string `json:"rtsps_key_file"`
	OutWaitKeyFrameFlag               bool   `json:"out_wait_key_frame_flag"`
	WsRtspEnable                      bool   `json:"ws_rtsp_enable"`
	WsRtspAddr                        string `json:"ws_rtsp_addr"`
	PubTimeoutSec                     uint32 `json:"pub_timeout_sec"`
	PullTimeoutSec                    uint32 `json:"pull_timeout_sec"`
	PubNotSubAutoCloseSec             int    `json:"pub_not_sub_auto_close_sec"` // 发布但无订阅自动关闭秒数
	RtspRemuxerAddSpsPps2KeyFrameFlag bool   `json:"add_sps_pps_to_key_frame_flag"`
	AuthEnable                        bool   `json:"auth_enable"` // 是否启用认证
	AuthMethod                        int    `json:"auth_method"` // 认证方法
	Username                          string `json:"username"`    // 用户名
	Password                          string `json:"password"`    // 密码
}

// SrtConfig SRT 配置
type SrtConfig struct {
	Enable                bool   `json:"enable"`                     // 是否启用
	Addr                  string `json:"addr"`                       // 监听地址
	GopNum                int    `json:"gop_num"`                    // GOP 数量
	SingleGopMaxFrameNum  int    `json:"single_gop_max_frame_num"`   // 单个 GOP 最大帧数
	EncryptEnable         bool   `json:"encrypt_enable"`             // 是否启用加密
	EncryptMode           int    `json:"encrypt_mode"`               // 加密模式
	Passphrase            string `json:"passphrase"`                 // 密码短语
	PubTimeoutSec         uint32 `json:"pub_timeout_sec"`            // 发布超时秒数
	PubNotSubAutoCloseSec int    `json:"pub_not_sub_auto_close_sec"` // 发布但无订阅自动关闭秒数
}

// OnvifConfig ONVIF 配置
type OnvifConfig struct {
	Enable bool `json:"enable"` // 是否启用
}

type Jt1078Config struct {
	Enable                bool   `json:"enable"`
	ListenIp              string `json:"listen_ip"`
	ListenPort            int    `json:"listen_port"`
	PortNum               uint16 `json:"port_num"` // 范围 ListenPort至ListenPort+PortNum
	PubTimeoutSec         uint32 `json:"pub_timeout_sec"`
	PubNotSubAutoCloseSec int    `json:"pub_not_sub_auto_close_sec"` // 发布但无订阅自动关闭秒数
	Intercom              struct {
		Enable     bool
		IP         string `json:"ip"`           // 固定外网ip
		Port       int    `json:"port"`         // 固定外网udp端口
		AudioPorts [2]int `json:"audio_ports"`  // 范围 AudioPorts[0]至AudioPorts[1]
		OnJoinURL  string `json:"on_join_url"`  // 设备对讲连接上了的url回调
		OnLeaveURL string `json:"on_leave_url"` // 设备对讲断开了的url回调
	} `json:"intercom"`
}

type HttpFmp4Config struct {
	CommonHttpServerConfig
}

type RecordConfig struct {
	EnableFlv            bool   `json:"enable_flv"`
	FlvOutPath           string `json:"flv_out_path"`
	EnableMpegts         bool   `json:"enable_mpegts"`
	MpegtsOutPath        string `json:"mpegts_out_path"`
	EnableFmp4           bool   `json:"enable_fmp4"`
	Fmp4OutPath          string `json:"fmp4_out_path"`
	RecordInterval       int    `json:"record_interval"`        // 固定间隔录制一个文件，单位秒
	EnableRecordInterval bool   `json:"enable_record_interval"` // 是否开启固定间隔录制
	EnableMpegps         bool   `json:"enable_mpegps"`          // 是否启用 MPEGPS
	MpegpsOutPath        string `json:"mpegps_out_path"`        // MPEGPS 输出路径
}

type PlaybackConfig struct {
	CommonHttpServerConfig
}

type MetricsConfig struct {
	Enable         bool   `json:"enable"`
	PushgatewayURL string `json:"pushgateway_url"`
	JobName        string `json:"job_name"`
	InstanceName   string `json:"instance_name"`
	PushInterval   int    `json:"push_interval"`
}

type RelayPushConfig struct {
	Enable   bool     `json:"enable"`
	AddrList []string `json:"addr_list"`
}

type StaticRelayPullConfig struct {
	Enable bool   `json:"enable"`
	Addr   string `json:"addr"`
}

type HttpApiConfig struct {
	CommonHttpServerConfig
}

type HttpNotifyConfig struct {
	Enable                  bool   `json:"enable"`
	KeepaliveIntervalSec    int    `json:"keepalive_interval_sec"`
	OnKeepalive             string `json:"on_keepalive"`
	UpdateIntervalSec       int    `json:"update_interval_sec"`
	OnServerStart           string `json:"on_server_start"`
	OnUpdate                string `json:"on_update"`
	OnPubStart              string `json:"on_pub_start"`
	OnPubStop               string `json:"on_pub_stop"`
	OnSubStart              string `json:"on_sub_start"`
	OnSubStop               string `json:"on_sub_stop"`
	OnSubStartWithoutStream string `json:"on_sub_start_without_stream"` // 订阅启动但无流
	OnStreamChanged         string `json:"on_stream_changed"`           // 流改变
	OnPushStart             string `json:"on_push_start"`
	OnPushStop              string `json:"on_push_stop"`
	OnRelayPullStart        string `json:"on_relay_pull_start"`
	OnRelayPullStop         string `json:"on_relay_pull_stop"`
	OnRtmpConnect           string `json:"on_rtmp_connect"`
	OnHlsMakeTs             string `json:"on_hls_make_ts"`
	OnHlsMakeFmp4           string `json:"on_hls_make_fmp4"`
	OnReportStat            string `json:"on_report_stat"`
	OnReportFrameInfo       string `json:"on_report_frame_info"`
	MaxTaskLen              int    `json:"max_task_len"`     // 最大任务数
	ClientSize              int    `json:"client_size"`      // 并发客户端
	DiscardInterval         uint32 `json:"discard_interval"` // 丢弃间隔，当队列满的时候，丢弃数量达到此值，下一个一定保留
}

type SimpleAuthConfig struct {
	Key                string `json:"key"`
	DangerousLalSecret string `json:"dangerous_lal_secret"`
	PubRtmpEnable      bool   `json:"pub_rtmp_enable"`
	SubRtmpEnable      bool   `json:"sub_rtmp_enable"`
	SubHttpflvEnable   bool   `json:"sub_httpflv_enable"`
	SubHttptsEnable    bool   `json:"sub_httpts_enable"`
	PubRtspEnable      bool   `json:"pub_rtsp_enable"`
	SubRtspEnable      bool   `json:"sub_rtsp_enable"`
	HlsM3u8Enable      bool   `json:"hls_m3u8_enable"`
	PushRtmpEnable     bool   `json:"push_rtmp_enable"`
	PushJt1078Enable   bool   `json:"push_jt1078_enable"`
	PushPsEnable       bool   `json:"push_ps_enable"`
}

type PprofConfig struct {
	CommonHttpServerConfig
}

type DebugConfig struct {
	LogGroupIntervalSec       int `json:"log_group_interval_sec"`
	LogGroupMaxGroupNum       int `json:"log_group_max_group_num"`
	LogGroupMaxSubNumPerGroup int `json:"log_group_max_sub_num_per_group"`
}

type CommonHttpServerConfig struct {
	CommonHttpAddrConfig

	Enable      bool   `json:"enable"`
	EnableHttps bool   `json:"enable_https"`
	EnableHttp3 bool   `json:"enable_http3"`
	UrlPattern  string `json:"url_pattern"`
}

type CommonHttpAddrConfig struct {
	HttpListenAddr  string `json:"http_listen_addr"`
	HttpsListenAddr string `json:"https_listen_addr"`
	Http3ListenAddr string `json:"http3_listen_addr"`
	HttpsCertFile   string `json:"https_cert_file"`
	HttpsKeyFile    string `json:"https_key_file"`
}

type RoomConfig struct {
	Enable    bool   `json:"enable"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level               int    `json:"level"`                  // 日志级别
	Filename            string `json:"filename"`               // 日志文件名
	IsToStdout          bool   `json:"is_to_stdout"`           // 是否输出到标准输出
	IsRotateDaily       bool   `json:"is_rotate_daily"`        // 是否按天轮转
	IsRotateHourly      bool   `json:"is_rotate_hourly"`       // 是否按小时轮转
	BuffSize            int    `json:"buff_size"`              // 缓冲大小
	AssertBehavior      int    `json:"assert_behavior"`        // 断言行为
	MaxBackupDates      int    `json:"max_backup_dates"`       // 最大备份天数
	MaxSize             int    `json:"max_size"`               // 最大大小
	MaxBackups          int    `json:"max_backups"`            // 最大备份数
	MaxAge              int    `json:"max_age"`                // 最大保留时间
	Compress            bool   `json:"compress"`               // 是否压缩
	IsRotateBySize      bool   `json:"is_rotate_by_size"`      // 是否按大小轮转
	MaxFileSizeMb       int    `json:"max_file_size_mb"`       // 最大文件大小（MB）
	MaxFileNum          int    `json:"max_file_num"`           // 最大文件数
	ShortFileFlag       bool   `json:"short_file_flag"`        // 短文件名标志
	TimestampFlag       bool   `json:"timestamp_flag"`         // 时间戳标志
	TimestampWithMsFlag bool   `json:"timestamp_with_ms_flag"` // 带毫秒的时间戳标志
	LevelFlag           bool   `json:"level_flag"`             // 级别标志
}

// GetServerConfig 获取服务器配置信息
// 用法示例：
//
//	engine := lalmax.NewEngine().SetConfig(lalmax.Config{URL: "http://localhost:8080"})
//	config, err := engine.GetServerConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("服务器ID: %s, 配置版本: %s\n", config.ServerId, config.ConfVersion)
func (e *Engine) GetServerConfig(ctx context.Context) (*GetServerConfigResponse, error) {
	const api = `/api/config/svr_config`
	var resp ServerConfig
	if err := e.get(ctx, api, &resp); err != nil {
		return nil, err
	}
	if err := ErrHandle(resp.Code, resp.Msg); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func ErrHandle(code int, msg string) error {
	switch code {
	case 0, 10000:
		return nil
	default:
		return fmt.Errorf("lalmax: %s", msg)
	}
}

// SetServerConfigResponse 设置服务器配置响应
type SetServerConfigResponse struct {
	FixedHeader
}

// SetHttpNotifyConfig 设置 HTTP 通知配置（合并模式，只更新传入的字段）
// 用法示例：
//
//	engine := lalmax.NewEngine().SetConfig(lalmax.Config{URL: "http://localhost:8080"})
//	err := engine.SetHttpNotifyConfig(ctx, lalmax.HttpNotifyConfig{
//	    Enable:        true,
//	    OnPubStart:    "http://your-server/callback/on_pub_start",
//	    OnPubStop:     "http://your-server/callback/on_pub_stop",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
func (e *Engine) SetHttpNotifyConfig(ctx context.Context, config HttpNotifyConfig, gb28181 MediaConfig) error {
	// 使用合并模式，只更新 http_notify 配置
	data := map[string]any{
		"http_notify": config,
		"rtsp": map[string]any{
			"pub_not_sub_auto_close_sec": 30,
		},
		"gb28181": map[string]any{
			"enable":                     false,
			"pub_not_sub_auto_close_sec": 30,
			"media_config":               gb28181,
		},
	}
	return e.setServerConfig(ctx, data, true)
}

// SetServerConfig 设置服务器配置（完整替换模式）
// 用法示例：
//
//	engine := lalmax.NewEngine().SetConfig(lalmax.Config{URL: "http://localhost:8080"})
//	err := engine.SetServerConfig(ctx, map[string]any{
//	    "http_notify": lalmax.HttpNotifyConfig{Enable: true},
//	})
func (e *Engine) SetServerConfig(ctx context.Context, config map[string]any) error {
	return e.setServerConfig(ctx, config, false)
}

// SetServerConfigMerge 设置服务器配置（合并模式，只更新传入的字段）
// 用法示例：
//
//	engine := lalmax.NewEngine().SetConfig(lalmax.Config{URL: "http://localhost:8080"})
//	err := engine.SetServerConfigMerge(ctx, map[string]any{
//	    "rtmp": map[string]any{"enable": false},
//	})
func (e *Engine) SetServerConfigMerge(ctx context.Context, config map[string]any) error {
	return e.setServerConfig(ctx, config, true)
}

// setServerConfig 内部方法：设置服务器配置
// merge: true=合并模式（只更新传入的字段），false=替换模式（完整配置替换）
func (e *Engine) setServerConfig(ctx context.Context, config map[string]any, merge bool) error {
	api := `/api/config/set_server_config`
	if merge {
		api += "?merge=true"
	}
	var resp SetServerConfigResponse
	if err := e.post(ctx, api, config, &resp); err != nil {
		return err
	}
	return ErrHandle(resp.Code, resp.Msg)
}
