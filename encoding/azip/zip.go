package azip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Create 压缩一个或多个文件/目录到指定的zip文件中。
// Compresses one or more files/directories into a specified zip file.
//
// 参数 filename: 输出的zip文件名。
// Param filename: Name of the output zip file.
//
// 参数 files: 需要添加到zip的文件/目录列表。
// Param files: List of files/directories to add to the zip.
//
// 返回错误信息，成功时返回nil。
// Returns error message, nil on success.
func Create(filename string, files []string) error {
	// 创建新zip文件
	// Create new zip file
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	// 创建zip.Writer用于写入压缩内容
	// Create zip.Writer for writing compressed content
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 循环处理所有输入文件/目录
	// Process all input files/directories
	for _, file := range files {
		if err := AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

// AddFileToZip 将单个文件或目录添加到zip.Writer
// Adds a single file or directory to zip.Writer
func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	// 打开目标文件/目录
	// Open target file/directory
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// 获取文件信息
	// Get file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	// 处理目录类型
	// Handle directory type
	if info.IsDir() {
		return addDirectoryToZip(zipWriter, fileToZip, filename, info)
	}

	// 处理普通文件
	// Handle regular file
	return addRegularFileToZip(zipWriter, filename, info)
}

// addDirectoryToZip 处理目录添加到zip的逻辑
// Handles directory addition to zip
func addDirectoryToZip(zipWriter *zip.Writer, dir *os.File, path string, info os.FileInfo) error {
	// 创建目录头信息
	// Create directory header
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// 规范路径格式并添加目录标识符"/"
	// Normalize path format and add directory identifier "/"
	header.Name = filepath.ToSlash(filepath.Clean(path)) + "/"

	// 移除可能的绝对路径前缀（安全性考虑）
	// Remove possible absolute path prefix (security consideration)
	if filepath.IsAbs(header.Name) && len(header.Name) > 0 {
		header.Name = header.Name[1:]
	}

	// 使用DEFLATE压缩算法（即使目录内容为空）
	// Use DEFLATE compression (even for empty directories)
	header.Method = zip.Deflate

	// 在zip中创建目录条目
	// Create directory entry in zip
	if _, err := zipWriter.CreateHeader(header); err != nil {
		return err
	}

	// 读取目录内容
	// Read directory contents
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// 递归处理子目录/文件
	// Recursively process subdirectories/files
	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		if err := AddFileToZip(zipWriter, fullPath); err != nil {
			return err
		}
	}
	return nil
}

// addRegularFileToZip 处理普通文件添加到zip的逻辑
// Handles regular file addition to zip
func addRegularFileToZip(zipWriter *zip.Writer, filename string, info os.FileInfo) error {
	// 创建文件头信息
	// Create file header
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// 规范路径格式
	// Normalize path format
	header.Name = filepath.ToSlash(filepath.Clean(filename))

	// 移除可能的绝对路径前缀（安全性考虑）
	// Remove possible absolute path prefix (security consideration)
	if filepath.IsAbs(header.Name) && len(header.Name) > 0 {
		header.Name = header.Name[1:]
	}

	// 使用DEFLATE压缩算法
	// Use DEFLATE compression
	header.Method = zip.Deflate

	// 创建zip文件条目写入器
	// Create zip entry writer
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// 重新打开文件以保证独立的文件指针
	// Reopen file to ensure independent file pointer
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用1MB缓冲区提升大文件复制性能
	// Use 1MB buffer to improve large file copy performance
	buf := make([]byte, 1<<20) // 1MB buffer
	_, err = io.CopyBuffer(writer, file, buf)
	if err != nil {
		return err
	}

	return nil
}
