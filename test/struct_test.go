package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"testing"
)

// Define a simple struct to test with
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestToStruct(t *testing.T) {
	// Test valid conversion from map to struct
	data := map[string]any{
		"name": "John",
		"age":  30,
	}
	var p Person
	err := conv.ToStruct(data, &p)
	if err != nil {
		t.Errorf("ToStruct failed: %v", err)
	}
	if p.Name != "John" {
		t.Errorf("Expected Name to be 'John', got '%s'", p.Name)
	}
	if p.Age != 30 {
		t.Errorf("Expected Age to be 30, got %d", p.Age)
	}

	// Test invalid conversion with nil data
	err = conv.ToStruct(nil, &p)
	if err == nil {
		t.Error("ToStruct should return an error when data is nil")
	}

	// Test invalid conversion with nil model
	err = conv.ToStruct(data, nil)
	if err == nil {
		t.Error("ToStruct should return an error when model is nil")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	// Test valid JSON deserialization into struct
	jsonData := []byte(`{"name": "Alice", "age": 25}`)
	var p Person
	err := conv.UnmarshalJSON(jsonData, &p)
	if err != nil {
		t.Errorf("UnmarshalJSON failed: %v", err)
	}
	if p.Name != "Alice" {
		t.Errorf("Expected Name to be 'Alice', got '%s'", p.Name)
	}
	if p.Age != 25 {
		t.Errorf("Expected Age to be 25, got %d", p.Age)
	}

	// Test invalid JSON deserialization with nil data
	err = conv.UnmarshalJSON(nil, &p)
	if err == nil {
		t.Error("UnmarshalJSON should return an error when data is nil")
	}

	// Test invalid JSON deserialization with nil model
	err = conv.UnmarshalJSON(jsonData, nil)
	if err == nil {
		t.Error("UnmarshalJSON should return an error when model is nil")
	}
}

func TestToJSON(t *testing.T) {
	// Test valid struct to JSON serialization
	p := Person{Name: "Bob", Age: 40}
	jsonData, err := conv.ToJSON(p)
	if err != nil {
		t.Errorf("ToJSON failed: %v", err)
	}
	if jsonData == nil {
		t.Error("ToJSON should return a valid JSON byte slice, but got nil")
	}

	// Test invalid data to JSON serialization (nil data)
	_, err = conv.ToJSON(nil)
	if err == nil {
		t.Error("ToJSON should return an error when data is nil")
	}
}
