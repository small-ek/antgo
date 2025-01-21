package aini

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Decode parses INI format data into a nested map structure
// Decode 将INI格式数据解析为嵌套的map结构
func Decode(data []byte) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	reader := bufio.NewReader(bytes.NewReader(data))

	var (
		currentSection string                 // Current section name 当前节名称
		sectionData    map[string]interface{} // Current section data 当前节数据
		inSection      bool                   // Flag for valid section 是否在有效节中
	)

	for {
		lineBytes, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read error: %w", err)
		}

		// Clean and check line
		// 清理并检查行内容
		line := strings.TrimSpace(string(lineBytes))
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue // Skip empty lines and comments 跳过空行和注释
		}

		// Parse section header
		// 解析节头
		if start, end := strings.Index(line, "["), strings.Index(line, "]"); start >= 0 && end > start+1 {
			currentSection = line[start+1 : end]
			sectionData = make(map[string]interface{})
			result[currentSection] = sectionData
			inSection = true
			continue
		}

		// Skip lines not in section
		// 跳过不在节中的行
		if !inSection {
			continue
		}

		// Parse key-value pairs
		// 解析键值对
		if sepIdx := strings.Index(line, "="); sepIdx > 0 {
			key := strings.TrimSpace(line[:sepIdx])
			value := strings.TrimSpace(line[sepIdx+1:])
			sectionData[key] = value
		}
	}

	if !inSection {
		return nil, errors.New("no valid section found in INI data")
	}
	return result, nil
}

// Encode converts map data to INI format bytes
// Encode 将map数据转换为INI格式字节
func Encode(data map[string]interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for sectionName, sectionContent := range data {
		// Validate section type
		// 验证节类型
		sectionData, ok := sectionContent.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid section type: %T", sectionContent)
		}

		// Write section header
		// 写入节头
		if _, err := fmt.Fprintf(buf, "[%s]\n", sectionName); err != nil {
			return nil, err
		}

		// Write key-value pairs
		// 写入键值对
		for key, value := range sectionData {
			if _, err := fmt.Fprintf(buf, "%s = %s\n", key, value); err != nil {
				return nil, err
			}
		}

		buf.WriteByte('\n') // Add section separator 添加节分隔符
	}

	return buf.Bytes(), nil
}

// ToJson converts INI data to JSON format
// ToJson 将INI数据转换为JSON格式
func ToJson(data []byte) ([]byte, error) {
	iniMap, err := Decode(data)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(iniMap, "", "  ") // Pretty-print JSON 格式化输出JSON
}
