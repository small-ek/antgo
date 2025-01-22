// array_test.go
package array

import (
	"sync"
	"testing"
)

// 测试 New 函数
func TestNew(t *testing.T) {
	t.Run("IntArray", func(t *testing.T) {
		arr := New[int](10)
		if arr.Len() != 0 {
			t.Errorf("Expected empty array, got length %d", arr.Len())
		}
	})

	t.Run("StringArray", func(t *testing.T) {
		arr := New[string](5)
		if arr.Len() != 0 {
			t.Error("Expected empty array")
		}
	})
}

// 测试 Append 和 Len
func TestAppendAndLen(t *testing.T) {
	t.Run("IntType", func(t *testing.T) {
		arr := New[int](0)
		arr.Append(42)
		arr.Append(100)

		if arr.Len() != 2 {
			t.Errorf("Expected length 2, got %d", arr.Len())
		}
	})

	t.Run("StringType", func(t *testing.T) {
		arr := New[string](2)
		arr.Append("hello")
		arr.Append("world")

		if arr.Len() != 2 {
			t.Error("Length mismatch")
		}
	})
}

// 测试 List 方法
func TestList(t *testing.T) {
	arr := New[int](3)
	arr.Append(1)
	arr.Append(2)

	list := arr.List()
	if len(list) != 2 || list[0] != 1 || list[1] != 2 {
		t.Errorf("List mismatch, got %v", list)
	}

	// 修改副本不应影响原数据
	list[0] = 99
	if val, _ := arr.Get(0); val != 1 {
		t.Error("Modifying list copy affected original data")
	}
}

// 测试 Insert 方法
func TestInsert(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		arr := New[int](3)
		arr.Append(1)
		arr.Append(3)

		err := arr.Insert(1, 2)
		if err != nil {
			t.Fatalf("Insert failed: %v", err)
		}

		if val, _ := arr.Get(1); val != 2 {
			t.Error("Insert value mismatch")
		}
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		arr := New[string](2)
		err := arr.Insert(0, "test")
		if err != nil { // 使用预定义错误变量
			t.Error(err)
		}
	})
}

// 测试 Delete 方法
func TestDelete(t *testing.T) {
	arr := New[string](3)
	arr.Append("a")
	arr.Append("b")
	arr.Append("c")

	t.Run("ValidDelete", func(t *testing.T) {
		err := arr.Delete(1)
		if err != nil {
			t.Fatal(err)
		}

		if val, _ := arr.Get(1); val != "c" {
			t.Error("Delete failed")
		}
	})

	t.Run("InvalidDelete", func(t *testing.T) {
		err := arr.Delete(99)
		if err == nil {
			t.Error("Expected index error")
		}
	})
}

// 测试 Get 和 Set 方法
func TestGetAndSet(t *testing.T) {
	arr := New[float64](2)
	arr.Append(10.5)
	arr.Append(20.3)

	t.Run("ValidGet", func(t *testing.T) {
		val, err := arr.Get(1)
		if err != nil || val != 20.3 {
			t.Error("Get failed")
		}
	})

	t.Run("InvalidGet", func(t *testing.T) {
		_, err := arr.Get(-1)
		if err == nil {
			t.Error("Expected error for invalid index")
		}
	})

	t.Run("ValidSet", func(t *testing.T) {
		err := arr.Set(0, 99.9)
		if err != nil {
			t.Fatal(err)
		}
		if val, _ := arr.Get(0); val != 99.9 {
			t.Error("Set failed")
		}
	})

	t.Run("InvalidSet", func(t *testing.T) {
		err := arr.Set(99, 0.0)
		if err == nil {
			t.Error("Expected index error")
		}
	})
}

// 测试 Search 方法
func TestSearch(t *testing.T) {
	arr := New[int](5)
	arr.Append(10)
	arr.Append(20)
	arr.Append(30)

	t.Run("Found", func(t *testing.T) {
		if idx := arr.Search(20); idx != 1 {
			t.Errorf("Expected index 1, got %d", idx)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		if idx := arr.Search(99); idx != -1 {
			t.Error("Expected -1 for missing element")
		}
	})
}

// 测试并发安全性
func TestConcurrency(t *testing.T) {
	arr := New[int](0)
	const numWorkers = 100
	var wg sync.WaitGroup

	// 并发追加
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(val int) {
			defer wg.Done()
			arr.Append(val)
		}(i)
	}
	wg.Wait()

	if arr.Len() != numWorkers {
		t.Errorf("Expected %d elements, got %d", numWorkers, arr.Len())
	}

	// 并发读取
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(idx int) {
			defer wg.Done()
			_, _ = arr.Get(idx % numWorkers)
		}(i)
	}
	wg.Wait()
}

// 测试 Clear 方法
func TestClear(t *testing.T) {
	arr := New[string](3)
	arr.Append("foo")
	arr.Append("bar")
	arr.Clear()

	if arr.Len() != 0 {
		t.Error("Clear failed")
	}
}

// 测试 WithWriteLock 和 WithReadLock
func TestLockMethods(t *testing.T) {
	arr := New[int](2)
	arr.Append(1)

	arr.WithWriteLock(func(data []int) {
		data[0] = 2
	})
	if val, _ := arr.Get(0); val != 2 {
		t.Error("WithWriteLock failed")
	}

	arr.WithReadLock(func(data []int) {
		if data[0] != 2 {
			t.Error("WithReadLock data mismatch")
		}
	})
}
