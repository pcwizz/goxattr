package goxattr

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

const TEST_VALUE = "test value"
const TEST_ATTR_NAME = "test"

func TestGetAttr001(T *testing.T) {
	file_name := setup()
	SetAttr(file_name, TEST_ATTR_NAME, TEST_VALUE)
	data, _ := GetAttr(file_name, TEST_ATTR_NAME)
	tear_down(file_name)
	if (string)(data) != TEST_VALUE {
		T.Fatal("The value we put in the attr is not the one we got out.")
	}
}

func TestSetAttrOO1(T *testing.T) {
	file_name := setup()
	SetAttr(file_name, TEST_ATTR_NAME, TEST_VALUE)
	_, numAttrs, _ := ListAttrs(file_name)
	tear_down(file_name)
	if numAttrs != 1 {
		T.Fatal("Set attr did not add an attr")
	}
}

func TestSetAttr002(T *testing.T) {
	file_name := setup()
	SetAttr(file_name, TEST_ATTR_NAME, TEST_VALUE)
	attrs, _, _ := ListAttrs(file_name)
	tear_down(file_name)
	if string(attrs[0][:]) != TEST_ATTR_NAME {
		T.Fatal("SetAttr didn't call the attr what we told it to")
	}
}

func TestListAttrs001(T *testing.T) {
	file_name := setup()
	_, numAttrs, _ := ListAttrs(file_name)
	tear_down(file_name)
	if numAttrs != 0 {
		T.Fatal("ListAttrs is reporting attrs we haven't created")
	}
}

func TestDeleteAttr(T *testing.T) {
	file_name := setup()
	SetAttr(file_name, TEST_ATTR_NAME, TEST_VALUE)
	DeleteAttr(file_name, TEST_ATTR_NAME)
	_, numAttrs, _ := ListAttrs(file_name)
	tear_down(file_name)
	if numAttrs == 1 {
		T.Fatal("DeleteAtrr did not delete the attribute we created")
	}
}

func setup() string {
	rand.Seed(time.Now().UnixNano())
	file_name_bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		file_name_bytes[i] = (byte)(rand.Int())
	}
	_, fn, _, _ := runtime.Caller(1)
	file_name := path.Join(path.Dir(fn), fmt.Sprintf("%x", sha1.Sum(file_name_bytes)))
	file, _ := os.Create(file_name)
	file.Close()
	return file_name
}

func tear_down(file_name string) {
	os.Remove(file_name)
}
