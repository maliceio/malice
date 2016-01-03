package persist

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/utils/clitable"
	"github.com/rakyll/magicmime"
	// "github.com/dutchcoders/gossdeep"
)

// File is a file object
type File struct {
	Name   string
	Path   string
	Valid  bool
	Size   int64
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

// Init initializes the File object
func (file *File) Init() {
	if file.Path != "" {
		file.GetName()
		file.GetSize()
		file.getData()
		file.GetCRC32()
		file.GetMD5()
		file.GetSHA1()
		file.GetSHA256()
		file.GetSHA512()
		file.GetFileMimeType()
		file.Data = nil
	} else {
		log.Fatalf("error occured during file.Init() because file.Path was not set.")
	}
}

// GetName returns file name
func (file *File) GetName() (name string, err error) {
	fileHandle, err := os.Open(file.Path)
	if err != nil {
		return
	}
	defer fileHandle.Close()

	stat, err := fileHandle.Stat()
	if err != nil {
		return
	}

	name = stat.Name()

	file.Name = name

	return
}

// GetSize calculates file Size
func (file *File) GetSize() (bytes int64, err error) {
	fileHandle, err := os.Open(file.Path)
	if err != nil {
		return
	}
	defer fileHandle.Close()

	stat, err := fileHandle.Stat()
	if err != nil {
		return
	}
	stat.Name()
	bytes = stat.Size()

	file.Size = bytes

	return
}

// getData loads file into []bytes
// TODO: This is probably pretty dumb to keep this data in memory
func (file *File) getData() {
	dat, err := ioutil.ReadFile(file.Path)
	assert(err)
	file.Data = dat
}

// GetCRC32 calculates the Files CRC32
func (file *File) GetCRC32() (hCRC32Sum string, err error) {

	var castagnoliTable = crc32.MakeTable(crc32.Castagnoli)

	crc := crc32.New(castagnoliTable)
	_, err = crc.Write(file.Data)
	assert(err)
	hCRC32Sum = fmt.Sprintf("%x", crc.Sum32())

	file.CRC32 = hCRC32Sum

	return
}

// GetMD5 calculates the Files md5sum
func (file *File) GetMD5() (hMd5Sum string, err error) {

	hmd5 := md5.New()
	_, err = hmd5.Write(file.Data)
	assert(err)
	hMd5Sum = fmt.Sprintf("%x", hmd5.Sum(nil))

	file.MD5 = hMd5Sum

	return
}

// GetSHA1 calculates the Files sha256sum
func (file *File) GetSHA1() (h1Sum string, err error) {

	h1 := sha1.New()
	_, err = h1.Write(file.Data)
	assert(err)
	h1Sum = fmt.Sprintf("%x", h1.Sum(nil))

	file.SHA1 = h1Sum

	return
}

// GetSHA256 calculates the Files sha256sum
func (file *File) GetSHA256() (h256Sum string, err error) {

	h256 := sha256.New()
	_, err = h256.Write(file.Data)
	assert(err)
	h256Sum = fmt.Sprintf("%x", h256.Sum(nil))

	file.SHA256 = h256Sum

	return
}

// GetSHA512 calculates the Files sha256sum
func (file *File) GetSHA512() (h512Sum string, err error) {

	h512 := sha512.New()
	_, err = h512.Write(file.Data)
	assert(err)
	h512Sum = fmt.Sprintf("%x", h512.Sum(nil))

	file.SHA512 = h512Sum

	return
}

// GetFileMimeType returns the mime-type of a file path
func (file *File) GetFileMimeType() (mimetype string, err error) {

	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		log.Fatal(err)
	}
	defer magicmime.Close()

	mimetype, err = magicmime.TypeByFile(file.Path)
	if err != nil {
		log.Fatalf("error occured during type lookup: %v", err)
	}

	file.Mime = mimetype
	// log.Printf("mime-type: %v", mimetype)
	return
}

// PrintFileDetails prints file details
func (file *File) PrintFileDetails() {
	table := clitable.New([]string{"Field", "Value"})
	table.AddRow(map[string]interface{}{"Field": "Name", "Value": file.Name})
	table.AddRow(map[string]interface{}{"Field": "Path", "Value": file.Path})
	table.AddRow(map[string]interface{}{"Field": "Valid", "Value": file.Valid})
	table.AddRow(map[string]interface{}{"Field": "Size", "Value": file.Size})
	table.AddRow(map[string]interface{}{"Field": "CRC32", "Value": file.CRC32})
	table.AddRow(map[string]interface{}{"Field": "MD5", "Value": file.MD5})
	table.AddRow(map[string]interface{}{"Field": "SHA1", "Value": file.SHA1})
	table.AddRow(map[string]interface{}{"Field": "SHA256", "Value": file.SHA256})
	table.AddRow(map[string]interface{}{"Field": "SHA512", "Value": file.SHA512})
	table.AddRow(map[string]interface{}{"Field": "Ssdeep", "Value": file.Ssdeep})
	table.AddRow(map[string]interface{}{"Field": "Mime", "Value": file.Mime})
	table.Markdown = true
	table.Print()
	// fmt.Println("Name: ", file.Name)
	// fmt.Println("Path: ", file.Path)
	// fmt.Println("Valid: ", file.Valid)
	// fmt.Println("Size: ", file.Size)
	// fmt.Println("CRC32: ", file.CRC32)
	// fmt.Println("MD5: ", file.MD5)
	// fmt.Println("SHA1: ", file.SHA1)
	// fmt.Println("SHA256: ", file.SHA256)
	// fmt.Println("SHA512: ", file.SHA512)
	// fmt.Println("Ssdeep: ", file.Ssdeep)
	// fmt.Println("Mime: ", file.Mime)
}

// // GetSsdeep calculates the Files ssdeep
// func (file File) GetSsdeep(data []byte) (hssdeep string, err error) {
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

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
