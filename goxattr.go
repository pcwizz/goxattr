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

/*
#cgo LDFLAGS: -lxattr
#include <stdlib.h>
#include <libxattr.h>

void clean( void* data, size_t nbytes){
	char* bitstring = data;
	for ( int i = nbytes - 1; i >= 0; i-- ){
		bitstring[i] = '\0';
	}
	return;
}

*/
import "C"

import (
	"errors"
	"strings"
	"unsafe"
)

const LIST_BUFFER_SIZE = 1024 // A 1KB buffer is more than enough
const DATA_BUFFER_SIZE = 256

/*
Function GetAttr retrieves the xattr names attrName of the path.
The content of the xattr is returned as a slice of bytes.
The only error possible from this function is "Xattr does not exist",
as all other errors are handled by the crossxattr C library, this may
change in the future.
*/
func GetAttr(path string, attrName string) (data []byte, err error) {

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cattrName := C.CString(attrName)
	defer C.free(unsafe.Pointer(cattrName))

	var cnbytes C.size_t = (C.size_t)(DATA_BUFFER_SIZE)

	var cdata unsafe.Pointer = C.malloc(cnbytes)
	defer C.free(cdata)
	//clean our memory
	C.clean(cdata, cnbytes)
	rtrn, errno := C.getAttr(cpath, cattrName, cdata, cnbytes)

	if int(rtrn) == -1 {
		if errno != nil {
			err = errno
			return
		}
		err = errors.New("Xattr does not exist")
		return
	}

	data = C.GoBytes(cdata, DATA_BUFFER_SIZE)
	//trim the 0s
	i := len(data) - 1
	for data[i] == '\000' && i > -1 {
		i--
	}
	data = data[:i+1]
	return
}

/*
Function SetAttr creates/updates the attribute attrName of the path to data.
No errors are currently returned however the my be in future.
*/
func SetAttr(path string, attrName string, data []byte) (err error) {

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cattrName := C.CString(attrName)
	defer C.free(unsafe.Pointer(cattrName))

	cdata := C.CString((string)(data))
	defer C.free(unsafe.Pointer(cdata))

	cndata := C.size_t(len(data))

	rtrn, errno := C.setAttr(cpath, cattrName, (unsafe.Pointer)(cdata), cndata)
	if int(rtrn) == -1 {
		err = errno
	}

	return
}

/*
Function DeleteAttr deletes the attribute attrName from the path.
the function may return only one error: "Xattr does not exist"
*/
func DeleteAttr(path string, attrName string) (err error) {

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cattrName := C.CString(attrName)
	defer C.free(unsafe.Pointer(cattrName))

	rtrn := C.deleteAttr(cpath, cattrName)

	if rtrn == -1 {
		err = errors.New("Xattr does not exist")
	}

	return
}

/*
Function ListAttrs returns a list of attribute names associated with the path.
*/
func ListAttrs(path string) (list []string, numAttrs int, err error) {

	var cnbytes C.size_t = (C.size_t)(LIST_BUFFER_SIZE)

	var cdata unsafe.Pointer = C.malloc(cnbytes)
	defer C.free(cdata)
	C.clean((unsafe.Pointer)(cdata), cnbytes)

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rtrn := C.listAttrs(cpath, cdata, cnbytes)
	if rtrn == -1 {
		err = errors.New("An error occured calling C.listAttrs")
	}

	numAttrs = int(rtrn)
	data := C.GoStringN((*C.char)(cdata), C.int(cnbytes))
	list = strings.Split(data, "\000")

	return
}
