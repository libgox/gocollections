package syncx

import (
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	m := &Map[string, int]{}

	// Define some test data
	testData := []struct {
		key   string
		value int
	}{
		{"key1", 10},
		{"key2", 20},
		{"key3", 30},
	}

	var wg sync.WaitGroup
	wg.Add(len(testData))

	for i := 0; i < len(testData); i++ {
		go func(i int) {
			m.Store(testData[i].key, testData[i].value)
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < len(testData); i++ {
		if v, ok := m.Load(testData[i].key); !ok || v != testData[i].value {
			t.Fatalf("Got %v, %v, expected %v, %v", v, ok, testData[i].value, true)
		}
	}
}

func TestMap_LoadOrStoreLazy(t *testing.T) {
	m := &Map[string, int]{}

	key := "lazyKey"
	value := 42

	fn := func() int {
		return value
	}

	// Case 1: Key does not exist, so LoadOrStoreLazy should store and return the lazy value
	if v, loaded := m.LoadOrStoreLazy(key, fn); loaded || v != value {
		t.Fatalf("Expected (%v, %v), got (%v, %v)", value, false, v, loaded)
	}

	// Case 2: Key now exists, so LoadOrStoreLazy should return the existing value without calling fn
	// Set up a different function that should not be called if the key already exists
	fnUnexpected := func() int {
		return 99
	}

	if v, loaded := m.LoadOrStoreLazy(key, fnUnexpected); !loaded || v != value {
		t.Fatalf("Expected (%v, %v), got (%v, %v)", value, true, v, loaded)
	}
}
