package persit

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash/crc32"
	"log"

	"github.com/rakyll/magicmime"

	// "github.com/dutchcoders/gossdeep"
)

// File is a file object
type File struct {
	Name   string
	Valid  bool
	Size   int
	CRC32  string
	MD5    string
	SHA1   string
	SHA256 string
	SHA512 string
	Ssdeep string
	Mime   string
	Arch   string
	Data   []byte
}

// GetCRC32 calculates the Files CRC32
func (file *File) GetCRC32(data []byte) (hCRC32Sum string, err error) {

	var castagnoliTable = crc32.MakeTable(crc32.Castagnoli)

	crc := crc32.New(castagnoliTable)
	_, err = crc.Write(data)
	assert(err)
	hCRC32Sum = fmt.Sprintf("%x", crc.Sum32())

	return
}

// GetMD5 calculates the Files md5sum
func (file *File) GetMD5(data []byte) (hMd5Sum string, err error) {

	hmd5 := md5.New()
	_, err = hmd5.Write(data)
	assert(err)
	hMd5Sum = fmt.Sprintf("%x", hmd5.Sum(nil))

	return
}

// GetSHA1 calculates the Files sha256sum
func (file *File) GetSHA1(data []byte) (h1Sum string, err error) {

	h1 := sha1.New()
	_, err = h1.Write(data)
	assert(err)
	h1Sum = fmt.Sprintf("%x", h1.Sum(nil))

	return
}

// GetSHA256 calculates the Files sha256sum
func (file *File) GetSHA256(data []byte) (h256Sum string, err error) {

	h256 := sha256.New()
	_, err = h256.Write(data)
	assert(err)
	h256Sum = fmt.Sprintf("%x", h256.Sum(nil))

	return
}

// GetSHA512 calculates the Files sha256sum
func (file *File) GetSHA512(data []byte) (h512Sum string, err error) {

	h256 := sha256.New()
	_, err = h256.Write(data)
	assert(err)
	h512Sum = fmt.Sprintf("%x", h256.Sum(nil))

	return
}

// GetFileMimeType returns the mime-type of a file path
func (file *File) GetFileMimeType(path string) {
	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		log.Fatal(err)
	}
	defer magicmime.Close()

	mimetype, err := magicmime.TypeByFile(path)
	if err != nil {
		log.Fatalf("error occured during type lookup: %v", err)
	}

	log.Printf("mime-type: %v", mimetype)
}

// // GetSsdeep calculates the Files ssdeep
// func (file *File) GetSsdeep(data []byte) (hssdeep string, err error) {
//
// 	hssdeep, err = ssdeep.HashFilename(file)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return
// }
//
// // CompareSsdeep returns the percent that two hashes are similar
// func CompareSsdeep(firstHash, secondHash string) (percent int, err error) {
//
// 	percent, err = ssdeep.Compare(firstHash, secondHash)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	return
// }
