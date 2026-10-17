package ota

const linuxTarPath = "upgrade.tar.gz"

// var _ Upgrader = &LinuxOTA{}

// type LinuxOTA struct {
// 	err        error
// 	OnProgress func(current, total int64)
// }

// // Download implements Upgrader.
// func (l *LinuxOTA) Download(link string) Upgrader {
// 	if l.err != nil {
// 		return l
// 	}
// 	resp, err := http.Get(linuxPackage)
// 	if err != nil {
// 		l.err = err
// 		return l
// 	}
// 	defer resp.Body.Close()

// 	_ = os.RemoveAll(filepath.Join(system.Getwd(), linuxTarPath))

// 	f, err := os.OpenFile(filepath.Join(system.Getwd(), linuxTarPath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
// 	if err != nil {
// 		l.err = err
// 		return l
// 	}
// 	defer f.Close()

// 	p := NewProgressReader(resp.ContentLength, resp.Body, l.OnProgress)
// 	defer p.Close()

// 	_, err = io.Copy(f, p)
// 	if err != nil {
// 		l.err = err
// 	}
// 	return l
// }

// // Unzip implements Upgrader.
// func (l *LinuxOTA) Unzip() Upgrader {
// 	if l.err != nil {
// 		return l
// 	}

// 	// 清理旧的升级目录
// 	upgradeDir := filepath.Join(system.Getwd(), "upgrade")
// 	_ = os.RemoveAll(upgradeDir)

// 	// 打开 tar.gz 文件
// 	file, err := os.Open(filepath.Join(system.Getwd(), linuxTarPath))
// 	if err != nil {
// 		l.err = err
// 		return l
// 	}
// 	defer file.Close()

// 	// 创建 gzip reader
// 	gzr, err := gzip.NewReader(file)
// 	if err != nil {
// 		l.err = err
// 		return l
// 	}
// 	defer gzr.Close()

// 	// 创建 tar reader
// 	tr := tar.NewReader(gzr)

// 	// 找到顶层目录名称
// 	var topLevelDir string
// 	for {
// 		header, err := tr.Next()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			l.err = err
// 			return l
// 		}

// 		// 获取第一级目录名
// 		parts := strings.Split(header.Name, "/")
// 		if len(parts) > 0 && parts[0] != "" {
// 			topLevelDir = parts[0]
// 			break
// 		}
// 	}

// 	// 重新打开文件进行解压
// 	file.Close()
// 	file, err = os.Open(filepath.Join(system.Getwd(), linuxTarPath))
// 	if err != nil {
// 		l.err = err
// 		return l
// 	}
// 	defer file.Close()

// 	gzr, err = gzip.NewReader(file)
// 	if err != nil {
// 		l.err = err
// 		return l
// 	}
// 	defer gzr.Close()

// 	tr = tar.NewReader(gzr)

// 	// 解压所有文件
// 	for {
// 		header, err := tr.Next()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			l.err = err
// 			return l
// 		}

// 		if err := l.extractFile(tr, header, upgradeDir, topLevelDir); err != nil {
// 			l.err = err
// 			return l
// 		}
// 	}

// 	return l
// }

// // Backup implements Upgrader.
// func (l *LinuxOTA) Backup() Upgrader {
// 	if l.err != nil {
// 		return l
// 	}

// 	execName := os.Args[0]
// 	backupName := execName + ".bak"
// 	if err := os.RemoveAll(backupName); err != nil {
// 		l.err = err
// 		return l
// 	}
// 	if err := os.Rename(execName, backupName); err != nil {
// 		l.err = err
// 	}
// 	return l
// }

// // Replace implements Upgrader.
// func (l *LinuxOTA) Replace() Upgrader {
// 	if l.err != nil {
// 		return l
// 	}

// 	upgradeDir := filepath.Join(system.Getwd(), "upgrade")
// 	currentDir := system.Getwd()

// 	// 获取当前可执行文件名
// 	execName := filepath.Base(os.Args[0])

// 	// 替换可执行文件
// 	newExecPath := filepath.Join(upgradeDir, execName)
// 	currentExecPath := filepath.Join(currentDir, execName)

// 	if _, err := os.Stat(newExecPath); err == nil {
// 		if err := l.copyFile(newExecPath, currentExecPath); err != nil {
// 			l.err = fmt.Errorf("替换可执行文件失败: %w", err)
// 			return l
// 		}
// 		// 设置可执行权限
// 		if err := os.Chmod(currentExecPath, 0o755); err != nil {
// 			l.err = fmt.Errorf("设置可执行权限失败: %w", err)
// 			return l
// 		}
// 	}

// 	// 替换 www 目录
// 	newWwwPath := filepath.Join(upgradeDir, "www")
// 	currentWwwPath := filepath.Join(currentDir, "www")

// 	if _, err := os.Stat(newWwwPath); err == nil {
// 		// 备份现有 www 目录
// 		backupWwwPath := filepath.Join(currentDir, "www.bak")
// 		_ = os.RemoveAll(backupWwwPath)
// 		if _, err := os.Stat(currentWwwPath); err == nil {
// 			if err := os.Rename(currentWwwPath, backupWwwPath); err != nil {
// 				l.err = fmt.Errorf("备份 www 目录失败: %w", err)
// 				return l
// 			}
// 		}

// 		// 复制新的 www 目录
// 		if err := l.copyDir(newWwwPath, currentWwwPath); err != nil {
// 			// 恢复备份
// 			_ = os.Rename(backupWwwPath, currentWwwPath)
// 			l.err = fmt.Errorf("替换 www 目录失败: %w", err)
// 			return l
// 		}

// 		// 删除备份
// 		_ = os.RemoveAll(backupWwwPath)
// 	}

// 	// 保留升级目录，下次升级的时候删除
// 	// _ = os.RemoveAll(upgradeDir)
// 	// 清理升级临时文件
// 	_ = os.RemoveAll(filepath.Join(currentDir, linuxTarPath))

// 	return l
// }

// // Error implements Upgrader.
// func (l *LinuxOTA) Error() error {
// 	return l.err
// }

// // extractFile 解压单个文件，跳过顶层目录
// func (l *LinuxOTA) extractFile(tr *tar.Reader, header *tar.Header, destDir, topLevelDir string) error {
// 	// 跳过顶层目录
// 	relativePath := header.Name
// 	if topLevelDir != "" && strings.HasPrefix(relativePath, topLevelDir+"/") {
// 		relativePath = strings.TrimPrefix(relativePath, topLevelDir+"/")
// 	}

// 	// 如果是顶层目录本身，跳过
// 	if relativePath == "" || relativePath == topLevelDir {
// 		return nil
// 	}

// 	target := filepath.Join(destDir, relativePath)

// 	switch header.Typeflag {
// 	case tar.TypeDir:
// 		if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
// 			return err
// 		}
// 	case tar.TypeReg:
// 		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
// 			return err
// 		}

// 		f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
// 		if err != nil {
// 			return err
// 		}
// 		defer f.Close()

// 		_, err = io.Copy(f, tr)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // copyFile 复制文件
// func (l *LinuxOTA) copyFile(src, dst string) error {
// 	sourceFile, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer sourceFile.Close()

// 	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
// 	if err != nil {
// 		return err
// 	}
// 	defer destFile.Close()

// 	_, err = io.Copy(destFile, sourceFile)
// 	return err
// }

// // copyDir 递归复制目录
// func (l *LinuxOTA) copyDir(src, dst string) error {
// 	srcInfo, err := os.Stat(src)
// 	if err != nil {
// 		return err
// 	}

// 	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
// 		return err
// 	}

// 	entries, err := os.ReadDir(src)
// 	if err != nil {
// 		return err
// 	}

// 	for _, entry := range entries {
// 		srcPath := filepath.Join(src, entry.Name())
// 		dstPath := filepath.Join(dst, entry.Name())

// 		if entry.IsDir() {
// 			if err := l.copyDir(srcPath, dstPath); err != nil {
// 				return err
// 			}
// 		} else {
// 			if err := l.copyFile(srcPath, dstPath); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }
