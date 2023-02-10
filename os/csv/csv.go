package csv

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type Csv struct {
	Path   string
	Data   [][]string
	File   *os.File
	Writer *csv.Writer
}

// New 创建对象
func New(path string) *Csv {
	return &Csv{
		Path: path,
	}
}

// Create 创建目录
func (c *Csv) Create() *Csv {
	file, err := os.Create(c.Path) //创建文件
	if err != nil {
		log.Fatalf("can not create file, err is %v", err)
	}
	defer file.Close()
	c.File = file

	if _, err := file.WriteString("\xEF\xBB\xBF"); err != nil {
		log.Fatalf("can not write UTF-8 BOM, err is %v", err)
	}
	return c
}

// Insert 插入数据
func (c *Csv) Insert(data [][]string) error {
	file, err := os.OpenFile(c.Path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(c.File)

	for i := 0; i < len(data); i++ {
		err = w.Write(data[i])
		if err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

// InsertOne 插入数据
func (c *Csv) InsertOne(data []string) error {
	file, err := os.OpenFile(c.Path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(c.File)
	err = w.Write(data)
	if err != nil {
		return err
	}

	c.Writer = w
	w.Flush()

	return nil
}

// Read 读取大文件
func (c *Csv) Read() [][]string {
	//准备读取文件
	fs, err := os.Open(c.Path)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()

	r := csv.NewReader(fs)
	//针对大文件，一行一行的读取文件
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		c.Data = append(c.Data, row)
	}
	return c.Data
}

// GetCount 获取总共多少条
func (c *Csv) GetCount() (int, error) {
	fs, err := os.Open(c.Path)
	defer fs.Close()

	if err != nil {
		return 0, err
	}
	r := csv.NewReader(fs)
	content, err := r.ReadAll()

	return len(content), err
}

// ReadSmallFile 读取小文件
func (c *Csv) ReadSmallFile() [][]string {
	fs, err1 := os.Open(c.Path)
	defer fs.Close()

	if err1 != nil {
		panic(err1)
	}
	r1 := csv.NewReader(fs)
	content, err2 := r1.ReadAll()
	if err2 != nil {
		panic(err2)
	}
	for i := 0; i < len(content); i++ {
		c.Data = append(c.Data, content[i])
	}

	return c.Data
}
