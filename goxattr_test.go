/*
Goxattr
-------
Simple go bindings to the crossxattr library for portable access to xattrs

Copyright (c) Morgan Hill 2016 <morgan@pcwizzltd.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.


*/
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
	SetAttr(file_name, TEST_ATTR_NAME, ([]byte)(TEST_VALUE))
	data, _ := GetAttr(file_name, TEST_ATTR_NAME)
	tear_down(file_name)
	if (string)(data) != TEST_VALUE {
		T.Fatal("The value we put in the attr is not the one we got out.")
	}
}

func TestSetAttrOO1(T *testing.T) {
	file_name := setup()
	SetAttr(file_name, TEST_ATTR_NAME, ([]byte)(TEST_VALUE))
	_, numAttrs, _ := ListAttrs(file_name)
	tear_down(file_name)
	if numAttrs != 1 {
		T.Fatal("Set attr did not add an attr")
	}
}

func TestSetAttr002(T *testing.T) {
	file_name := setup()
	SetAttr(file_name, TEST_ATTR_NAME, ([]byte)(TEST_VALUE))
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
	SetAttr(file_name, TEST_ATTR_NAME, ([]byte)(TEST_VALUE))
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
