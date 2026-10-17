package api

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gowvp/owl/internal/core/recording"
)

// 缓存已处理的纯视频文件，避免重复 remux
var (
	videoOnlyCache     = make(map[string]string) // 原始路径 -> 纯视频路径
	videoOnlyCacheLock sync.RWMutex
)

// serveVideoOnly 提供纯视频文件（移除音频轨道）
// 路径: /static/recordings-video/{path}
// HLS.js 无法处理 G.711 音频，需要提供纯视频版本
func (a RecordingAPI) serveVideoOnly(c *gin.Context) {
	// 获取请求路径
	requestPath := c.Param("path")
	if requestPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "path is required"})
		return
	}

	// 构建原始文件路径
	storageDir := a.conf.Server.Recording.StorageDir
	originalPath := filepath.Join(storageDir, requestPath)

	// 检查原始文件是否存在
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "file not found"})
		return
	}

	// 检查缓存
	videoOnlyCacheLock.RLock()
	cachedPath, exists := videoOnlyCache[originalPath]
	videoOnlyCacheLock.RUnlock()

	if exists {
		// 检查缓存文件是否存在
		if _, err := os.Stat(cachedPath); err == nil {
			c.File(cachedPath)
			return
		}
	}

	// 生成纯视频文件
	videoOnlyPath, err := a.createVideoOnlyFile(originalPath)
	if err != nil {
		slog.Error("创建纯视频文件失败", "path", originalPath, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	// 缓存
	videoOnlyCacheLock.Lock()
	videoOnlyCache[originalPath] = videoOnlyPath
	videoOnlyCacheLock.Unlock()

	c.File(videoOnlyPath)
}

// createVideoOnlyFile 使用 ffmpeg 创建纯视频文件（移除音频）
func (a RecordingAPI) createVideoOnlyFile(originalPath string) (string, error) {
	// 生成输出路径
	// 使用 MD5 作为文件名，避免路径问题
	hash := md5.Sum([]byte(originalPath))
	hashStr := fmt.Sprintf("%x", hash)

	// 在临时目录创建纯视频文件
	cacheDir := filepath.Join(a.conf.Server.Recording.StorageDir, ".video-cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("创建缓存目录失败: %w", err)
	}

	outputPath := filepath.Join(cacheDir, hashStr+".mp4")

	// 如果文件已存在，直接返回
	if _, err := os.Stat(outputPath); err == nil {
		return outputPath, nil
	}

	// 使用 ffmpeg 移除音频（仅复制视频流，不转码）
	// -an: 移除音频
	// -c:v copy: 视频流直接复制，不转码
	// -movflags +faststart: 将 moov 放在文件开头，支持边下载边播放
	cmd := exec.Command("ffmpeg",
		"-y",               // 覆盖输出文件
		"-i", originalPath, // 输入文件
		"-an",          // 移除音频
		"-c:v", "copy", // 视频直接复制
		"-movflags", "+faststart", // moov 放开头
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg 执行失败: %w, output: %s", err, string(output))
	}

	slog.Info("创建纯视频文件成功", "original", originalPath, "output", outputPath)
	return outputPath, nil
}

// RegisterRecordingVideoOnly 注册纯视频静态文件服务
// 需要在 RegisterRecording 之后调用
func RegisterRecordingVideoOnly(g gin.IRouter, api RecordingAPI, handler ...gin.HandlerFunc) {
	// 纯视频文件服务
	// 路径: /static/recordings-video/*path
	g.GET("/static/recordings-video/*path", append(handler, api.serveVideoOnly)...)
}

// generateM3U8VideoOnly 生成指向纯视频文件的 M3U8
// HLS.js 使用这个版本播放视频，音频由前端单独处理
func (a RecordingAPI) generateM3U8VideoOnly(recordings []*recording.Recording, baseURL, token string) string {
	// 这里直接复用原有逻辑，但路径改为 /static/recordings-video/
	// ... 实现类似 generateM3U8WithToken
	return ""
}
