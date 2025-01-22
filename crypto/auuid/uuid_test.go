package auuid

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	t.Parallel()
	u := New()
	if u.Version() != 4 {
		t.Errorf("Expected version 4, got %d", u.Version())
	}
	if IsNilUUID(u) {
		t.Error("Generated nil UUID")
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		u, err := Create()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if u.Version() != 1 {
			t.Errorf("Expected version 1, got %d", u.Version())
		}
	})

	t.Run("retry logic", func(t *testing.T) {
		t.Parallel()
		// 注意：实际测试中需要模拟失败场景
	})
}

func TestCreateWithNode(t *testing.T) {
	t.Parallel()

	testNode := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}

	t.Run("valid node", func(t *testing.T) {
		t.Parallel()
		u, err := CreateWithNode(testNode)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// 验证节点ID
		if !bytes.Equal(u[10:], testNode) {
			t.Errorf("NodeID mismatch\nWant: %X\nGot:  %X", testNode, u[10:])
		}

		// 验证版本
		if u.Version() != 1 {
			t.Errorf("Expected version 1, got %d", u.Version())
		}
	})

	t.Run("invalid node length", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			name string
			node []byte
		}{
			{"short", []byte{0x01}},
			{"long", []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				_, err := CreateWithNode(tc.node)
				if err == nil {
					t.Error("Expected error but got nil")
				}
			})
		}
	})
}

func TestBatchGenerate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		count int
	}{
		{"small batch", 10},
		{"medium batch", 100},
		{"large batch", 1000},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uuids, err := BatchGenerate(tc.count)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(uuids) != tc.count {
				t.Errorf("Expected %d UUIDs, got %d", tc.count, len(uuids))
			}

			// 验证唯一性
			seen := make(map[uuid.UUID]bool)
			for _, u := range uuids {
				if seen[u] {
					t.Errorf("Duplicate UUID found: %s", u)
				}
				seen[u] = true
			}
		})
	}

	t.Run("invalid count", func(t *testing.T) {
		t.Parallel()
		_, err := BatchGenerate(-1)
		if err == nil {
			t.Error("Expected error for negative count")
		}
	})
}

func TestDCEGeneration(t *testing.T) {
	t.Parallel()

	t.Run("person domain", func(t *testing.T) {
		t.Parallel()
		u, err := NewDCEPerson()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if u.Version() != 2 {
			t.Errorf("Expected version 2, got %d", u.Version())
		}
	})

	t.Run("org domain validation", func(t *testing.T) {
		t.Parallel()
		_, err := NewDCEOrg(0)
		if err == nil {
			t.Error("Expected error for zero org ID")
		}
	})
}

func TestTimeFunctions(t *testing.T) {
	t.Parallel()

	t.Run("valid time-based UUID", func(t *testing.T) {
		t.Parallel()
		u, _ := Create()
		ts, err := GetTimestamp(u)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if time.Since(ts) > time.Minute {
			t.Errorf("Timestamp %s is too old", ts)
		}
	})

	t.Run("invalid UUID type", func(t *testing.T) {
		t.Parallel()
		u := New()
		_, err := GetTimestamp(u)
		if err == nil {
			t.Error("Expected error for non-time-based UUID")
		}
	})
}

func TestValidationFunctions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid", "f47ac10b-58cc-0372-8567-0e02b2c3d479", true},
		{"invalid", "not-a-uuid", false},
		{"nil", "00000000-0000-0000-0000-000000000000", true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if IsValidUUID(tc.input) != tc.expected {
				t.Errorf("Validation failed for %s", tc.input)
			}
		})
	}
}

func TestConversionFunctions(t *testing.T) {
	t.Parallel()

	t.Run("string roundtrip", func(t *testing.T) {
		t.Parallel()
		original := New()
		str := UUIDToString(original)
		parsed, err := StringToUUID(str)
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if original != parsed {
			t.Errorf("Roundtrip mismatch\nOriginal: %s\nParsed:   %s", original, parsed)
		}
	})

	t.Run("nil handling", func(t *testing.T) {
		t.Parallel()
		if UUIDToString(uuid.Nil) != "" {
			t.Error("Nil UUID should return empty string")
		}
	})
}
