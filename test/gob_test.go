package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"testing"
)

// 测试数据结构 / Test data structure
type testData2 struct {
	Name  string
	Value int
}

// TestGobEncoder_Normal 测试正常编码场景
// TestGobEncoder_Normal tests normal encoding scenarios
func TestGobEncoder_Normal(t *testing.T) {
	data := testData2{
		Name:  "test",
		Value: 123,
	}

	encoded, err := conv.GobEncoder(data)
	if err != nil {
		t.Fatalf("GobEncoder failed: %v", err)
	}

	if len(encoded) == 0 {
		t.Error("Encoded data should not be empty")
	}
}

// TestGobDecoder_Normal 测试正常解码场景
// TestGobDecoder_Normal tests normal decoding scenarios
func TestGobDecoder_Normal(t *testing.T) {
	// 准备测试数据 / Prepare test data
	src := testData2{
		Name:  "test",
		Value: 123,
	}

	encoded, err := conv.GobEncoder(src)
	if err != nil {
		t.Fatalf("GobEncoder failed: %v", err)
	}

	// 解码测试 / Decoding test
	var dest testData2
	err = conv.GobDecoder(encoded, &dest)
	if err != nil {
		t.Fatalf("GobDecoder failed: %v", err)
	}

	// 验证解码结果 / Verify decoding result
	if dest.Name != src.Name || dest.Value != src.Value {
		t.Errorf("Decoded data mismatch, got: %+v, want: %+v", dest, src)
	}
}

// TestGobEncoder_NilData 测试空数据编码
// TestGobEncoder_NilData tests encoding nil data
func TestGobEncoder_NilData(t *testing.T) {
	_, err := conv.GobEncoder(nil)
	if err == nil {
		t.Error("Expected error for nil data")
	}
	if err != conv.ErrNilData {
		t.Errorf("Expected ErrNilData, got: %v", err)
	}
}

// TestGobDecoder_NilData 测试空数据解码
// TestGobDecoder_NilData tests decoding nil data
func TestGobDecoder_NilData(t *testing.T) {
	var dest testData2
	err := conv.GobDecoder(nil, &dest)
	if err == nil {
		t.Error("Expected error for nil data")
	}
	if err != conv.ErrNilData {
		t.Errorf("Expected ErrNilData, got: %v", err)
	}
}

// TestGobDecoder_NilTarget 测试空目标对象解码
// TestGobDecoder_NilTarget tests decoding with nil target
func TestGobDecoder_NilTarget(t *testing.T) {
	src := testData2{
		Name:  "test",
		Value: 123,
	}

	encoded, err := conv.GobEncoder(src)
	if err != nil {
		t.Fatalf("GobEncoder failed: %v", err)
	}

	err = conv.GobDecoder(encoded, nil)
	if err == nil {
		t.Error("Expected error for nil target")
	}
	if err != conv.ErrNilTarget {
		t.Errorf("Expected ErrNilTarget, got: %v", err)
	}
}

// TestGobDecoder_InvalidTarget 测试无效目标类型
// TestGobDecoder_InvalidTarget tests invalid target type
func TestGobDecoder_InvalidTarget(t *testing.T) {
	src := testData2{
		Name:  "test",
		Value: 123,
	}

	encoded, err := conv.GobEncoder(src)
	if err != nil {
		t.Fatalf("GobEncoder failed: %v", err)
	}

	// 非指针类型 / Non-pointer type
	var dest testData2
	err = conv.GobDecoder(encoded, dest)
	if err == nil {
		t.Error("Expected error for non-pointer target")
	}
	if err != conv.ErrTargetNotPointer {
		t.Errorf("Expected ErrTargetNotPointer, got: %v", err)
	}
}

// TestGobEncoderDecoder_Complex 测试复杂数据结构的编解码
// TestGobEncoderDecoder_Complex tests encoding/decoding of complex data
func TestGobEncoderDecoder_Complex(t *testing.T) {
	type complexData struct {
		ID      int
		Names   []string
		Details map[string]interface{}
	}

	src := complexData{
		ID:    1,
		Names: []string{"a", "b", "c"},
		Details: map[string]interface{}{
			"key1": 123,
			"key2": "value",
		},
	}

	encoded, err := conv.GobEncoder(src)
	if err != nil {
		t.Fatalf("GobEncoder failed: %v", err)
	}

	var dest complexData
	err = conv.GobDecoder(encoded, &dest)
	if err != nil {
		t.Fatalf("GobDecoder failed: %v", err)
	}

	// 验证复杂结构 / Verify complex structure
	if dest.ID != src.ID || len(dest.Names) != len(src.Names) ||
		len(dest.Details) != len(src.Details) {
		t.Errorf("Decoded complex data mismatch, got: %+v, want: %+v", dest, src)
	}
}
