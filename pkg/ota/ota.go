package ota

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	linuxPackage   = `/releases/latest/download/`
	LastVersionURL = `https://api.github.com/repos/%s/releases/latest`
)

// ReleaseInfo GitHub Release 信息
type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}

// OTA 提供版本检查和下载功能的结构体
// OTA 只负责下载，不关心后续的解压、备份、替换等操作
type OTA struct {
	repoName   string
	filename   string
	err        error
	onProgress func(current, total int64)
}

// NewOTA 创建 OTA 实例
// repoName: GitHub 仓库名，如 "gowvp/owl"，也支持 "github.com/gowvp/owl" 格式
// filename: 下载的文件名
func NewOTA(repoName, filename string) *OTA {
	return &OTA{
		repoName: cleanRepoName(repoName),
		filename: filename,
	}
}

// SetProgressCallback 设置下载进度回调
func (o *OTA) SetProgressCallback(callback func(current, total int64)) *OTA {
	o.onProgress = callback
	return o
}

// GetLastVersion 从 GitHub API 获取最新版本信息
// 返回 tag_name, body(release notes), error
func (o *OTA) GetLastVersion() (string, string, error) {
	return GetLastVersion(o.repoName)
}

// Download 下载升级包到指定路径
func (o *OTA) Download() *OTA {
	if o.err != nil {
		return o
	}

	// link := o.getDownloadLink()
	// linuxOTA := &LinuxOTA{OnProgress: o.onProgress}
	// linuxOTA.Download(link)
	// o.err = linuxOTA.Error()
	return o
}

// Error 返回错误
func (o *OTA) Error() error {
	return o.err
}

// getDownloadLink 获取下载链接
func (o *OTA) getDownloadLink() string {
	repoLink := "https://github.com/" + o.repoName
	link, _ := url.JoinPath(repoLink, linuxPackage, o.filename)
	return link
}

// cleanRepoName 清理仓库名称，移除前缀
// 支持 "gowvp/owl"、"github.com/gowvp/owl" 等格式
func cleanRepoName(repoName string) string {
	repoName = strings.TrimPrefix(repoName, "https://")
	repoName = strings.TrimPrefix(repoName, "http://")
	repoName = strings.TrimPrefix(repoName, "github.com/")
	repoName = strings.TrimPrefix(repoName, "api.github.com/repos/")
	return repoName
}

// GetLastVersion 从 GitHub API 获取最新版本信息
// repoName: GitHub 仓库名，如 "gowvp/owl"
// 返回 tag_name, body(release notes), error
func GetLastVersion(repoName string) (string, string, error) {
	repoName = cleanRepoName(repoName)
	apiURL := fmt.Sprintf(LastVersionURL, repoName)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", fmt.Errorf("解析响应失败: %w", err)
	}

	return release.TagName, release.Body, nil
}
