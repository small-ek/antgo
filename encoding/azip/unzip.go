package azip

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip 解压ZIP文件到指定目录
// Decompresses a zip archive to the specified directory
//
// 参数 src: 源ZIP文件路径
// Param src: Source zip file path
//
// 参数 dest: 目标解压目录
// Param dest: Target extraction directory
//
// 返回: 解压文件列表和错误信息
// Returns: Extracted file list and error
func Unzip(src string, dest string) ([]string, error) {
	// 打开ZIP文件
	// Open zip file
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var extractedFiles []string
	// 创建1MB缓冲区用于文件复制
	// Create 1MB buffer for file copying
	buf := make([]byte, 1<<20) // 1MB buffer

	// 并行处理控制（可根据需要调整）
	// Concurrency control (adjustable as needed)
	// sem := make(chan struct{}, runtime.NumCPU()*2)

	// 遍历ZIP文件内容
	// Iterate through zip contents
	for _, f := range r.File {
		// 安全检查和路径处理
		// Security check and path processing
		fpath, err := safeExtractPath(dest, f.Name)
		if err != nil {
			return extractedFiles, err
		}
		extractedFiles = append(extractedFiles, fpath)

		// 处理目录
		// Handle directory
		if f.FileInfo().IsDir() {
			if err := createDir(fpath); err != nil {
				return extractedFiles, err
			}
			continue
		}

		// 处理文件
		// Handle file
		if err := extractFile(f, fpath, buf); err != nil {
			return extractedFiles, err
		}

		// 可选：并行处理（需要权衡磁盘IO和CPU利用率）
		// Optional: Parallel processing (need to balance disk IO and CPU usage)
		// sem <- struct{}{}
		// go func(f *zip.File) {
		//     defer func() { <-sem }()
		//     // ...处理逻辑...
		// }(f)
	}
	return extractedFiles, nil
}

// safeExtractPath 安全路径检查和处理
// Security path validation and processing
func safeExtractPath(dest string, filename string) (string, error) {
	// 清理路径并转换分隔符
	// Clean path and convert separators
	cleanPath := filepath.ToSlash(filepath.Clean(filename))
	if strings.HasPrefix(cleanPath, "/") {
		return "", errors.New("file path error")
	}

	// 构建完整目标路径
	// Build full destination path
	fpath := filepath.Join(dest, cleanPath)

	// 检查路径穿越漏洞（ZipSlip）
	// Check for ZipSlip vulnerability
	if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
		return "", errors.New("file path error")
	}
	return fpath, nil
}

// createDir 创建目录结构
// Create directory structure
func createDir(fpath string) error {
	// 使用适当权限创建目录
	// Create directory with proper permissions
	if err := os.MkdirAll(fpath, 0755); err != nil {
		return err
	}
	// 同步目录修改（可选，确保数据持久化）
	// Sync directory changes (optional, ensure data persistence)
	if d, err := os.Open(fpath); err == nil {
		d.Sync()
		d.Close()
	}
	return nil
}

// extractFile 解压单个文件
// Extract single file
func extractFile(f *zip.File, fpath string, buf []byte) error {
	// 创建目标文件
	// Create target file
	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 打开ZIP文件条目
	// Open zip file entry
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// 使用缓冲区复制数据
	// Copy data with buffer
	if _, err := io.CopyBuffer(outFile, rc, buf); err != nil {
		return err
	}

	// 立即同步文件内容（确保数据持久化）
	// Sync file contents immediately (ensure data persistence)
	if err := outFile.Sync(); err != nil {
		return err
	}
	return nil
}
