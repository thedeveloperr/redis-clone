package hashmap

import (
	"testing"
	"time"
)

func TestNonExistingKeyGet(t *testing.T) {
	hashMap := Create()
	if _, exists := hashMap.Get("NonExistingKey"); exists {
		t.Errorf("Should not fetch non existent key.")
	}
}

func TestSetGet(t *testing.T) {
	hashMap := Create()
	hashMap.Set("TestKey", "TestValue")

	value, exists := hashMap.Get("TestKey")
	if !exists {
		t.Errorf("Should not fetch non existent key.")
	}

	if value != "TestValue" {
		t.Errorf("Expected 'TestValue' but got %v", value)
	}
}

func TestSetGetOverwrite(t *testing.T) {
	hashMap := Create()
	hashMap.Set("TestKey", "TestValue")

	value, exists := hashMap.Get("TestKey")
	if !exists {
		t.Errorf("TestKey key should exist but not found.")
	}

	if value != "TestValue" {
		t.Errorf("Expected 'TestValue' but got %v", value)
	}

	hashMap.Set("TestKey", "TestValueNew")

	valueNew, ok := hashMap.Get("TestKey")
	if !ok {
		t.Errorf("TestKey key should exist but not found.")
	}

	if valueNew != "TestValueNew" {
		t.Errorf("Expected 'TestValueNew' but got %v", valueNew)
	}
}

func TestSetGetExpire(t *testing.T) {
	t.Log("Start Running slow test as involve timeout in seconds")

	hashMap := Create()
	hashMap.Set("TestKey", "TestValue")

	value, exists := hashMap.Get("TestKey")
	if !exists {
		t.Errorf("TestKey key should exist but not found.")
	}

	if value != "TestValue" {
		t.Errorf("Expected 'TestValue' but got %v", value)
	}

	hashMap.Expire("TestKey", 2)

	time.Sleep(1 * time.Second)
	value, exists = hashMap.Get("TestKey")
	if !exists {
		t.Errorf("TestKey key should exist but not found.")
	}

	if value != "TestValue" {
		t.Errorf("Expected 'TestValue' but got %v", value)
	}

	time.Sleep(1 * time.Second)
	value, exists = hashMap.Get("TestKey")
	if exists {
		t.Errorf("TestKey key should not exist but found.")
	}

	if value != "" {
		t.Errorf("Expected '' but got %v", value)
	}

	hashMap.Set("TestKey2", "TestValue2")

	value, exists = hashMap.Get("TestKey2")
	if !exists {
		t.Errorf("TestKey key should exist but not found.")
	}

	if value != "TestValue2" {
		t.Errorf("Expected 'TestValue' but got %v", value)
	}

	hashMap.Expire("TestKey", 0)

	value, exists = hashMap.Get("TestKey")
	if exists {
		t.Errorf("TestKey key should not exist but found.")
	}

	if value != "" {
		t.Errorf("Expected '' but got %v", value)
	}
	t.Log("End Running slow test as involve timeout in seconds")
}
