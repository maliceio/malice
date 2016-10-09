// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build linux darwin

// Package magicmime detects mimetypes using libmagic.
// This package requires libmagic, install it by the following
// commands below.
//	 - Debian or Ubuntu: apt-get install libmagic-dev
//	 - RHEL, CentOS or Fedora: yum install file-devel
//	 - Mac OS X: brew install libmagic
package magicmime

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -lmagic -L/usr/local/lib
// #include <stdlib.h>
// #include <magic.h>
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

// Decoder allows for doing mimetype detection using libmagic.
type Decoder struct {
	db C.magic_t
}

// decMu and dec are used for the package level functions.
var (
	decMu sync.Mutex
	dec   *Decoder
)

type Flag int

const (
	// No special handling
	MAGIC_NONE Flag = C.MAGIC_NONE

	// Prints debugging messages to stderr
	MAGIC_DEBUG Flag = C.MAGIC_DEBUG

	// If the file queried is a symlink, follow it.
	MAGIC_SYMLINK Flag = C.MAGIC_SYMLINK

	// If the file is compressed, unpack it and look at the contents.
	MAGIC_COMPRESS Flag = C.MAGIC_COMPRESS

	// If the file is a block or character special device, then open the device
	// and try to look in its contents.
	MAGIC_DEVICES Flag = C.MAGIC_DEVICES

	// Return a MIME type string, instead of a textual description.
	MAGIC_MIME_TYPE Flag = C.MAGIC_MIME_TYPE

	// Return a MIME encoding, instead of a textual description.
	MAGIC_MIME_ENCODING Flag = C.MAGIC_MIME_ENCODING

	// A shorthand for MAGIC_MIME_TYPE | MAGIC_MIME_ENCODING.
	MAGIC_MIME Flag = C.MAGIC_MIME

	// Return all matches, not just the first.
	MAGIC_CONTINUE Flag = C.MAGIC_CONTINUE

	// Check the magic database for consistency and print warnings to stderr.
	MAGIC_CHECK Flag = C.MAGIC_CHECK

	// On systems that support utime(2) or utimes(2), attempt to preserve the
	// access time of files analyzed.
	MAGIC_PRESERVE_ATIME Flag = C.MAGIC_PRESERVE_ATIME

	// Don't translate unprintable characters to a \ooo octal representation.
	MAGIC_RAW Flag = C.MAGIC_RAW

	// Treat operating system errors while trying to open files and follow
	// symlinks as real errors, instead of printing them in the magic buffer
	MAGIC_ERROR Flag = C.MAGIC_ERROR

	// Return the Apple creator and type.
	MAGIC_APPLE Flag = C.MAGIC_APPLE

	// Don't check for EMX application type (only on EMX).
	MAGIC_NO_CHECK_APPTYPE Flag = C.MAGIC_NO_CHECK_APPTYPE

	// Don't get extra information on MS Composite Document Files.
	MAGIC_NO_CHECK_CDF Flag = C.MAGIC_NO_CHECK_CDF

	// Don't look inside compressed files.
	MAGIC_NO_CHECK_COMPRESS Flag = C.MAGIC_NO_CHECK_COMPRESS

	// Don't print ELF details.
	MAGIC_NO_CHECK_ELF Flag = C.MAGIC_NO_CHECK_ELF

	// Don't check text encodings.
	MAGIC_NO_CHECK_ENCODING Flag = C.MAGIC_NO_CHECK_ENCODING

	// Don't consult magic files.
	MAGIC_NO_CHECK_SOFT Flag = C.MAGIC_NO_CHECK_SOFT

	// Don't examine tar files.
	MAGIC_NO_CHECK_TAR Flag = C.MAGIC_NO_CHECK_TAR

	// Don't check for various types of text files.
	MAGIC_NO_CHECK_TEXT Flag = C.MAGIC_NO_CHECK_TEXT

	// Don't look for known tokens inside ascii files.
	MAGIC_NO_CHECK_TOKENS Flag = C.MAGIC_NO_CHECK_TOKENS
)

// NewDecoder creates a detector that uses libmagic. It initializes
// the opens the magicmime database with the specified flags. Upon
// success users are expected to call Close on the returned Decoder
// when it is no longer needed.
func NewDecoder(flags Flag) (*Decoder, error) {
	db := C.magic_open(C.int(0))
	if db == nil {
		return nil, errors.New("error opening magic")
	}
	d := &Decoder{db: db}
	if code := C.magic_setflags(db, C.int(flags)); code != 0 {
		d.Close()
		return nil, errors.New(C.GoString(C.magic_error(d.db)))
	}

	if code := C.magic_load(db, nil); code != 0 {
		d.Close()
		return nil, errors.New(C.GoString(C.magic_error(d.db)))
	}
	return d, nil
}

// TypeByFile looks up for a file's mimetype by its content.
// It uses a magic number database which is described in magic(5).
func (d *Decoder) TypeByFile(filename string) (string, error) {
	path := C.CString(filename)
	defer C.free(unsafe.Pointer(path))
	out := C.magic_file(d.db, path)
	if out == nil {
		return "", errors.New(C.GoString(C.magic_error(d.db)))
	}
	return C.GoString(out), nil
}

// TypeByBuffer looks up for a blob's mimetype by its contents.
// It uses a magic number database which is described in magic(5).
func (d *Decoder) TypeByBuffer(blob []byte) (string, error) {
	bytes := unsafe.Pointer(&blob[0])
	out := C.magic_buffer(d.db, bytes, C.size_t(len(blob)))
	if out == nil {
		return "", errors.New(C.GoString(C.magic_error(d.db)))
	}

	return C.GoString(out), nil
}

// Close frees up resources associated with d.
func (d *Decoder) Close() {
	C.magic_close(d.db)
	d.db = nil
}

// Open initializes a global Decoder and opens the magicmime database
// with the specified flags. Once successfully opened, users must
// call Close when they are finished using the package. This must
// be called before any of the package level functions.
func Open(flags Flag) error {
	decMu.Lock()
	defer decMu.Unlock()
	if dec != nil {
		return errors.New("cannot open; magic mime db is already open")
	}
	var err error
	dec, err = NewDecoder(flags)
	return err
}

// TypeByFile calls TypeByFile on the global Decoder. This is safe for
// concurrent use with the other package level TypeBy* functions.
func TypeByFile(filename string) (string, error) {
	decMu.Lock()
	defer decMu.Unlock()
	return dec.TypeByFile(filename)
}

// TypeByBuffer calls TypeByBuffer on the global Decoder. This is safe for
// concurrent use with the other package level TypeBy* functions.
func TypeByBuffer(blob []byte) (string, error) {
	decMu.Lock()
	defer decMu.Unlock()
	return dec.TypeByBuffer(blob)
}

// Close cleans up the resources associated with the global detector.
func Close() {
	dec.Close()
	dec = nil
}
