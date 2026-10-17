package recording

// TimeRange 时间轴数据项，表示一段录像的时间范围
type TimeRange struct {
	ID          int64   `json:"id"`           // 录像记录 ID
	StartMs     int64   `json:"start_ms"`     // 开始时间（毫秒时间戳）
	EndMs       int64   `json:"end_ms"`       // 结束时间（毫秒时间戳）
	Duration    float64 `json:"duration"`     // 时长（秒）
	ObjectCount int     `json:"object_count"` // AI检测对象数量
	DeleteFlag  bool    `json:"delete_flag"`  // 待删除标记（已被标记即将清理）
}
