# Goxattr

Simple go bindings to the crossxattr library for portable access to xattrs.

## Platforms

Goxattr itself will run an any pltform that cn run go but it will only be useful on platforms supported by crossxattr which provides a standard c interface to xattrs on the following platforms:

- GNU+LINUX
- FreeBSD

## Methods

### GetAttr (path string, attrName string) (data []byte, err error)

Retrieves the attribute (attrName) about the file (at path) as a slice of bytes (data). Error is set if no attribute with that name exists.

### SetAttr(path string, attrName string, data []byte) (err error)

Sets the attr (attrName) about file (at path) value (data)

### DeleteAttr(path string, attrName) (err error)

Deletes attr (attrName) about file (at path) returns error if attr does not exist.

### ListAttr(path string) (list []string, numAttrs int, err error)

Returns a list (list) of attrs set on file (at path), the length of the list (numAttrs) and an err if crossxatt isn't catchig all the errors (if it isn't then contact me)

## Quirks

- Crossxattr only returns atters in the user namespace
- Crossxattr emmits an error to stderr and dies on most errors
- Goxattr does not check tht files exist

## Licence

This project is licenced under the General Public License version 3 or above. A copy of the GPLv3 can be found in [LICENCE](LICENCE).
