package test

import (
	"fmt"
	"github.com/small-ek/antgo/utils/conv"
	"testing"
)

// TestMap tests the Map function.
func TestMap(t *testing.T) {
	// Test 1: JSON string input
	jsonStr := `{"key1": "value1", "key2": 2}`
	expected := map[string]interface{}{
		"key1": "value1",
		"key2": float64(2),
	}

	result := conv.Map(jsonStr)
	if !equalMaps(result, expected) {
		t.Errorf("Map(jsonStr) = %v; want %v", result, expected)
	}

	// Test 2: Map input
	testMap := map[string]interface{}{"key1": "value1", "key2": 2}
	result = conv.Map(testMap)
	if !equalMaps(result, testMap) {
		t.Errorf("Map(testMap) = %v; want %v", result, testMap)
	}
}

// TestMaps tests the Maps function.
func TestMaps(t *testing.T) {
	// Test 1: JSON string input
	jsonStr := `[{"key1": "value1", "key2": 2}, {"key3": "value3"}]`
	expected := []map[string]interface{}{
		{"key1": "value1", "key2": float64(2)},
		{"key3": "value3"},
	}

	result := conv.Maps(jsonStr)
	if !equalMapsSlices(result, expected) {
		t.Errorf("Maps(jsonStr) = %v; want %v", result, expected)
	}

	// Test 2: Slice input
	testSlice := []map[string]interface{}{
		{"key1": "value1", "key2": 2},
		{"key3": "value3"},
	}
	result = conv.Maps(testSlice)
	if !equalMapsSlices(result, testSlice) {
		t.Errorf("Maps(testSlice) = %v; want %v", result, testSlice)
	}
}

// TestMapString tests the MapString function.
func TestMapString(t *testing.T) {
	// Test 1: JSON string input
	jsonStr := `{"key1": "value1", "key2": 2}`
	expected := map[string]string{
		"key1": "value1",
		"key2": "2",
	}

	result := conv.MapString(jsonStr)
	if !equalMapStrings(result, expected) {
		t.Errorf("MapString(jsonStr) = %v; want %v", result, expected)
	}

	// Test 2: Map input
	testMap := map[string]interface{}{"key1": "value1", "key2": 2}
	result = conv.MapString(testMap)
	if !equalMapStrings(result, map[string]string{"key1": "value1", "key2": "2"}) {
		t.Errorf("MapString(testMap) = %v; want %v", result, map[string]string{"key1": "value1", "key2": "2"})
	}
}

// TestMapInt tests the MapInt function.
func TestMapInt(t *testing.T) {
	// Test 1: JSON string input
	jsonStr := `{"123": "value1", "456": 2}`
	expected := map[int]interface{}{
		123: "value1",
		456: float64(2),
	}

	result := conv.MapInt(jsonStr)
	if !equalMapInts(result, expected) {
		t.Errorf("MapInt(jsonStr) = %v; want %v", result, expected)
	}

	// Test 2: Map input with invalid key
	invalidJSON := `{"key1": "value1"}`
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MapInt(invalidJSON) should panic")
		}
	}()
	conv.MapInt(invalidJSON) // This should panic due to non-integer keys
}

// Helper functions for comparison (since Go doesn't have deep equality check)
func equalMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || fmt.Sprintf("%v", v) != fmt.Sprintf("%v", bv) {
			return false
		}
	}
	return true
}

func equalMapsSlices(a, b []map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, mapA := range a {
		if !equalMaps(mapA, b[i]) {
			return false
		}
	}
	return true
}

func equalMapStrings(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || v != bv {
			return false
		}
	}
	return true
}

func equalMapInts(a, b map[int]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || fmt.Sprintf("%v", v) != fmt.Sprintf("%v", bv) {
			return false
		}
	}
	return true
}
