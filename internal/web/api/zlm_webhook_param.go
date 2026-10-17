package api

// 注销
//	{
//		"mediaServerId" : "your_server_id",
//		"app" : "live",
//		"regist" : false,
//		"schema" : "rtsp",
//		"stream" : "obs",
//		"vhost" : "__defaultVhost__"
//	}

// 注册
//
//	{
//	    "regist" : true,
//	    "aliveSecond": 0, #存活时间，单位秒
//	    "app": "live", # 应用名
//	    "bytesSpeed": 0, #数据产生速度，单位byte/s
//	    "createStamp": 1617956908,  #GMT unix系统时间戳，单位秒
//	    "mediaServerId": "your_server_id", # 服务器id
//	    "originSock": {
//	        "identifier": "000001C257D35E40",
//	        "local_ip": "172.26.20.112", # 本机ip
//	        "local_port": 50166, # 本机端口
//	        "peer_ip": "172.26.20.112", # 对端ip
//	        "peer_port": 50155 # 对端port
//	    },
//	    "originType": 8,  # 产生源类型，包括 unknown = 0,rtmp_push=1,rtsp_push=2,rtp_push=3,pull=4,ffmpeg_pull=5,mp4_vod=6,device_chn=7,rtc_push=8
//	    "originTypeStr": "rtc_push",
//	    "originUrl": "", #产生源的url
//	    "readerCount": 0, # 本协议观看人数
//	    "schema": "rtsp", # 协议
//	    "stream": "test",  # 流id
//	    "totalReaderCount": 0, # 观看总人数，包括hls/rtsp/rtmp/http-flv/ws-flv/rtc
//	    "tracks": [{
//	       "channels" : 1, # 音频通道数
//	       "codec_id" : 2, # H264 = 0, H265 = 1, AAC = 2, G711A = 3, G711U = 4
//	       "codec_id_name" : "CodecAAC", # 编码类型名称
//	       "codec_type" : 1, # Video = 0, Audio = 1
//	       "ready" : true, # 轨道是否准备就绪
//	       "sample_bit" : 16, # 音频采样位数
//	       "sample_rate" : 8000 # 音频采样率
//	    },
//	    {
//	       "codec_id" : 0, # H264 = 0, H265 = 1, AAC = 2, G711A = 3, G711U = 4
//	       "codec_id_name" : "CodecH264", # 编码类型名称
//	       "codec_type" : 0, # Video = 0, Audio = 1
//	       "fps" : 59,  # 视频fps
//	       "height" : 720, # 视频高
//	       "ready" : true,  # 轨道是否准备就绪
//	       "width" : 1280 # 视频宽
//	    }],
//	    "vhost": "__defaultVhost__"
//	}
type onStreamChangedInput struct {
	Regist           bool       `json:"regist"`
	AliveSecond      int        `json:"aliveSecond"`
	App              string     `json:"app"`
	BytesSpeed       int        `json:"bytesSpeed"`
	CreateStamp      int        `json:"createStamp"`
	MediaServerID    string     `json:"mediaServerId"`
	OriginSock       OriginSock `json:"originSock"`
	OriginType       int        `json:"originType"`
	OriginTypeStr    string     `json:"originTypeStr"`
	OriginURL        string     `json:"originUrl"`
	ReaderCount      int        `json:"readerCount"`
	Schema           string     `json:"schema"`
	Stream           string     `json:"stream"`
	TotalReaderCount int        `json:"totalReaderCount"`
	Tracks           []Tracks   `json:"tracks"`
	Vhost            string     `json:"vhost"`

	// 以下字段为 lalmax 新增
	AppName    string `json:"app_name"`    // 流应用名
	StreamName string `json:"stream_name"` // 流名称
}
type OriginSock struct {
	Identifier string `json:"identifier"`
	LocalIP    string `json:"local_ip"`
	LocalPort  int    `json:"local_port"`
	PeerIP     string `json:"peer_ip"`
	PeerPort   int    `json:"peer_port"`
}
type Tracks struct {
	Channels    int     `json:"channels,omitempty"`
	CodecID     int     `json:"codec_id"`
	CodecIDName string  `json:"codec_id_name"`
	CodecType   int     `json:"codec_type"`
	Ready       bool    `json:"ready"`
	SampleBit   int     `json:"sample_bit,omitempty"`
	SampleRate  int     `json:"sample_rate,omitempty"`
	Fps         float32 `json:"fps,omitempty"`
	Height      int     `json:"height,omitempty"`
	Width       int     `json:"width,omitempty"`
}

// 心跳
// {
// 	"data" : {
// 		"Buffer" : 12,
// 		"BufferLikeString" : 0,
// 		"BufferList" : 0,
// 		"BufferRaw" : 12,
// 		"Frame" : 0,
// 		"FrameImp" : 0,
// 		"MediaSource" : 0,
// 		"MultiMediaSourceMuxer" : 0,
// 		"RtmpPacket" : 0,
// 		"RtpPacket" : 0,
// 		"Socket" : 108,
// 		"TcpClient" : 0,
// 		"TcpServer" : 96,
// 		"TcpSession" : 0,
// 		"UdpServer" : 12,
// 		"UdpSession" : 0
// 	 },
// 	 "mediaServerId" : "192.168.255.10"
//   }

type onServerKeepaliveInput struct {
	Data          Data   `json:"data"`
	HookIndex     int    `json:"hook_index"`
	MediaServerID string `json:"mediaServerId"`
}
type Data struct {
	Buffer                int `json:"Buffer"`
	BufferLikeString      int `json:"BufferLikeString"`
	BufferList            int `json:"BufferList"`
	BufferRaw             int `json:"BufferRaw"`
	Frame                 int `json:"Frame"`
	FrameImp              int `json:"FrameImp"`
	MediaSource           int `json:"MediaSource"`
	MultiMediaSourceMuxer int `json:"MultiMediaSourceMuxer"`
	RtmpPacket            int `json:"RtmpPacket"`
	RtpPacket             int `json:"RtpPacket"`
	Socket                int `json:"Socket"`
	TCPClient             int `json:"TcpClient"`
	TCPServer             int `json:"TcpServer"`
	TCPSession            int `json:"TcpSession"`
	UDPServer             int `json:"UdpServer"`
	UDPSession            int `json:"UdpSession"`
}

type onPublishInput struct {
	MediaServerID string `json:"mediaServerId"`
	App           string `json:"app"`
	ID            string `json:"id"`     // TCP 链接唯一 ID
	IP            string `json:"ip"`     // 推流器 ip
	Params        string `json:"params"` // 推流 url 参数
	Port          int    `json:"port"`   // 推流器端口号
	Schema        string `json:"schema"` // 推流的协议，可能是 rtsp、rtmp
	Stream        string `json:"stream"`
	Vhost         string `json:"vhost"` // 流虚拟主机
}

type onPublishOutput struct {
	DefaultOutput
	AddMuteAudio   *bool   `json:"add_mute_audio,omitempty"`
	ContinuePushMs *int    `json:"continue_push_ms,omitempty"`
	EnableAudio    *bool   `json:"enable_audio,omitempty"`
	EnableFmp4     *bool   `json:"enable_fmp4,omitempty"`
	EnableHls      *bool   `json:"enable_hls,omitempty"`
	EnableHlsFmp4  *bool   `json:"enable_hls_fmp4,omitempty"`
	EnableMp4      *bool   `json:"enable_mp4,omitempty"`
	EnableRtmp     *bool   `json:"enable_rtmp,omitempty"`
	EnableRtsp     *bool   `json:"enable_rtsp,omitempty"`
	EnableTs       *bool   `json:"enable_ts,omitempty"`
	HlsSavePath    *string `json:"hls_save_path,omitempty"`
	ModifyStamp    *bool   `json:"modify_stamp,omitempty"`
	Mp4AsPlayer    *bool   `json:"mp4_as_player,omitempty"`
	Mp4MaxSecond   *int    `json:"mp4_max_second,omitempty"`
	Mp4SavePath    *string `json:"mp4_save_path,omitempty"`
	AutoClose      *bool   `json:"auto_close,omitempty"`
	StreamReplace  *string `json:"stream_replace,omitempty"`
}

type DefaultOutput struct {
	Code int    `json:"code"` // 错误代码，0 代表允许推流
	Msg  string `json:"msg"`  // 不允许推流时的错误提示
}

func newDefaultOutputOK() DefaultOutput {
	return DefaultOutput{Code: 0, Msg: "success"}
}

type onStreamNoneReaderOutput struct {
	Code  int  `json:"code"`
	Close bool `json:"close"`
}

type onStreamNoneReaderInput struct {
	App           string `json:"app"`           // 流应用名
	Schema        string `json:"schema"`        // rtsp 或 rtmp
	Stream        string `json:"stream"`        // 流 ID
	Vhost         string `json:"vhost"`         // 流虚拟主机
	MediaServerID string `json:"mediaServerId"` // 服务器 id,通过配置文件设置
}

type onRTPServerTimeoutInput struct {
	LocalPort     int    `json:"local_port"`    // openRtpServer 输入的参数
	ReUsePort     bool   `json:"re_use_port"`   // openRtpServer 输入的参数
	SSRC          uint32 `json:"ssrc"`          // openRtpServer 输入的参数
	StreamID      string `json:"stream_id"`     // openRtpServer 输入的参数
	TCPMode       int    `json:"tcp_mode"`      // openRtpServer 输入的参数
	MediaServerID string `json:"mediaServerId"` // 服务器 id,通过配置文件设置
}

type onStreamNotFoundInput struct {
	MediaServerID string `json:"mediaServerId"` // 服务器 id,通过配置文件设置
	App           string `json:"app"`           // 流应用名
	ID            string `json:"id"`            // TCP链接唯一ID
	IP            string `json:"ip"`            // 播放器ip
	Params        string `json:"params"`        // 播放url参数
	Port          int    `json:"port"`          // 播放器端口号
	Schema        string `json:"schema"`        // 播放的协议，可能是rtsp、rtmp、http
	Stream        string `json:"stream"`        // 流 ID
	Vhost         string `json:"vhost"`         // 流虚拟主机

	// 以下字段为 lalmax 新增
	AppName    string `json:"app_name"`    // 流应用名
	StreamName string `json:"stream_name"` // 流名称
}

// onRecordMP4Input 录制 mp4 完成后通知事件参数
// https://docs.zlmediakit.com/zh/guide/media_server/web_hook_api.html#_8%E3%80%81on-record-mp4
type onRecordMP4Input struct {
	MediaServerID string  `json:"mediaServerId"` // 服务器 id
	App           string  `json:"app"`           // 录制的流应用名
	FileName      string  `json:"file_name"`     // 文件名
	FilePath      string  `json:"file_path"`     // 文件绝对路径
	FileSize      int64   `json:"file_size"`     // 文件大小，单位字节
	Folder        string  `json:"folder"`        // 文件所在目录路径
	StartTime     int64   `json:"start_time"`    // 开始录制时间戳（秒）
	Stream        string  `json:"stream"`        // 录制的流 ID
	TimeLen       float64 `json:"time_len"`      // 录制时长，单位秒
	URL           string  `json:"url"`           // http/rtsp/rtmp 点播相对 url 路径
	Vhost         string  `json:"vhost"`         // 流虚拟主机
}
