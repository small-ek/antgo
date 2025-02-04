package acsv

import (
	"bufio"
	"encoding/csv"
	"errors"
	"os"
	"sync"
)

// CSV 操作结构体
type CSV struct {
	FileName string
	Data     [][]string
	mu       sync.Mutex
}

// NewCSV 创建CSV实例
func New(fileName string) (*CSV, error) {
	c := &CSV{
		FileName: fileName,
	}
	if err := c.Read(); err != nil {
		return nil, err
	}
	return c, nil
}

// Create 创建目录
func (c *CSV) Create() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Create(c.FileName) //创建文件
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString("\xEF\xBB\xBF"); err != nil {
		return err
	}
	return nil
}

// Read 读取CSV文件内容
func (c *CSV) Read() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.OpenFile(c.FileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	data, err := reader.ReadAll()
	if err != nil {
		return err
	}
	c.Data = data
	return nil
}

// Write 写入CSV文件内容
func (c *CSV) Write() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.OpenFile(c.FileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(c.Data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

// AddRow 添加一行记录
func (c *CSV) AddRow(row []string) *CSV {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Data = append(c.Data, row)
	return c
}

// UpdateRow 更新一行记录
func (c *CSV) UpdateRow(index int, row []string) *CSV {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Data[index] = row
	return c
}

// DeleteRow 删除一行记录
func (c *CSV) DeleteRow(index int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if index < 0 || index >= len(c.Data) {
		return errors.New("index out of range")
	}
	c.Data = append(c.Data[:index], c.Data[index+1:]...)
	return nil
}
