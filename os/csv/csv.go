package csv

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type Csv struct {
	Path string
	Data [][]string
	File *os.File
}

//New 创建对象
func New(path string) *Csv {
	return &Csv{
		Path: path,
	}
}

//Create 创建目录
func (c *Csv) Create() *Csv {
	f, err := os.Create(c.Path) //创建文件
	if err != nil {
		panic(err)
	}

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	c.File = f
	return c
}

//Insert 插入数据
func (c *Csv) Insert(data [][]string) {
	c.Data = data
	w := csv.NewWriter(c.File) //创建一个新的写入文件流
	w.WriteAll(data)           //写入数据
	w.Flush()
	defer c.File.Close()
}

//Read 读取大文件
func (c *Csv) Read() *Csv {
	//准备读取文件
	fileName := c.Path
	fs, err := os.Open(fileName)
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
	return c
}

//ReadMini 读取小文件
func (c *Csv) ReadMini() *Csv {
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
	for _, row := range content {
		c.Data = append(c.Data, row)
	}
	return c
}

//Get 获取数据
func (c *Csv) Get() [][]string {
	return c.Data
}
