package persist

import (
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/docker/go-units"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/maliceio/malice/malice/malutils"
	"github.com/rakyll/magicmime"
	// "github.com/dutchcoders/gossdeep"
)

// File is a file object
type File struct {
	Name string `json:"name" gorethink:"name"`
	Path string `json:"path" gorethink:"path"`
	// Valid bool   `json:"valid"`
	Size string `json:"size" gorethink:"size"`
	// Size   int64
	// CRC32  string
	MD5    string `json:"md5" gorethink:"md5"`
	SHA1   string `json:"sha1" gorethink:"sha1"`
	SHA256 string `json:"sha256" gorethink:"sha256"`
	SHA512 string `json:"sha512" gorethink:"sha512"`
	// Ssdeep string `json:"ssdeep"`
	Mime string `json:"mime" gorethink:"mime"`
	// Arch string `json:"arch"`
	Data []byte `json:"data" gorethink:"data,omitempty"`
}

// Init initializes the File object
func (file *File) Init() {
	if file.Path != "" {
		file.GetName()
		file.GetSize()
		file.getData()
		// file.GetCRC32()
		file.GetMD5()
		file.GetSHA1()
		file.GetSHA256()
		file.GetSHA512()
		file.GetFileMimeType()
		file.CopyToSamples()
		file.gzipSample()
		file.Data = nil
	} else {
		log.Fatalf("error occured during file.Init() because file.Path was not set.")
	}
}

// CopyToSamples copys input file to samples folder
func (file *File) CopyToSamples() error {

	// Make .malice directory if it doesn't exist
	if _, err := os.Stat(maldirs.GetSampledsDir()); os.IsNotExist(err) {
		os.MkdirAll(maldirs.GetSampledsDir(), 0777)
	}

	if _, err := os.Stat(path.Join(maldirs.GetSampledsDir(), file.SHA256)); os.IsNotExist(err) {
		err := malutils.CopyFile(file.Path, path.Join(maldirs.GetSampledsDir(), file.SHA256))
		log.WithFields(log.Fields{
			"sample": file.SHA256,
		}).Debug("Copied sample to sample dir: ", maldirs.GetSampledsDir())
		er.CheckError(err)
		return err
	}

	return nil
}

func (file *File) gzipSample() error {
	reader, err := os.Open(file.Path)
	er.CheckError(err)

	destination := filepath.Join(maldirs.GetSampledsDir(), file.SHA256+".tar.gz")
	writer, err := os.Create(destination)
	er.CheckError(err)
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = file.Name
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	er.CheckError(err)
	return err
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

	bytes = stat.Size()

	file.Size = units.HumanSize(float64(bytes))

	return
}

// getData loads file into []bytes
// TODO: This is probably pretty dumb to keep this data in memory
func (file *File) getData() {
	dat, err := ioutil.ReadFile(file.Path)
	er.CheckError(err)
	file.Data = dat
}

// // GetCRC32 calculates the Files CRC32
// func (file *File) GetCRC32() (hCRC32Sum string, err error) {
//
// 	var castagnoliTable = crc32.MakeTable(crc32.Castagnoli)
//
// 	crc := crc32.New(castagnoliTable)
// 	_, err = crc.Write(file.Data)
// 	er.CheckError(err)
// 	hCRC32Sum = fmt.Sprintf("%x", crc.Sum32())
//
// 	file.CRC32 = hCRC32Sum
//
// 	return
// }

// GetMD5 calculates the Files md5sum
func (file *File) GetMD5() (hMd5Sum string, err error) {

	hmd5 := md5.New()
	_, err = hmd5.Write(file.Data)
	er.CheckError(err)
	hMd5Sum = fmt.Sprintf("%x", hmd5.Sum(nil))

	file.MD5 = hMd5Sum

	return
}

// GetSHA1 calculates the Files sha256sum
func (file *File) GetSHA1() (h1Sum string, err error) {

	h1 := sha1.New()
	_, err = h1.Write(file.Data)
	er.CheckError(err)
	h1Sum = fmt.Sprintf("%x", h1.Sum(nil))

	file.SHA1 = h1Sum

	return
}

// GetSHA256 calculates the Files sha256sum
func (file *File) GetSHA256() (h256Sum string, err error) {

	h256 := sha256.New()
	_, err = h256.Write(file.Data)
	er.CheckError(err)
	h256Sum = fmt.Sprintf("%x", h256.Sum(nil))

	file.SHA256 = h256Sum

	return
}

// GetSHA512 calculates the Files sha256sum
func (file *File) GetSHA512() (h512Sum string, err error) {

	h512 := sha512.New()
	_, err = h512.Write(file.Data)
	er.CheckError(err)
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

	// filetype := http.DetectContentType(file.Data)
	// file.Mime = filetype

	file.Mime = mimetype
	// log.Printf("mime-type: %v", mimetype)
	return
}

// ToJSON converts File object to []byte JSON
func (file *File) ToJSON() []byte {
	fileJSON, err := json.Marshal(file)
	er.CheckError(err)
	return fileJSON
}

// ToMarkdownTable converts File object Markdown table
func (file *File) ToMarkdownTable() {
	fmt.Println("#### File")
	table := clitable.New([]string{"Field", "Value"})
	table.AddRow(map[string]interface{}{"Field": "Name", "Value": file.Name})
	table.AddRow(map[string]interface{}{"Field": "Path", "Value": file.Path})
	table.AddRow(map[string]interface{}{"Field": "Size", "Value": file.Size})
	table.AddRow(map[string]interface{}{"Field": "MD5", "Value": file.MD5})
	table.AddRow(map[string]interface{}{"Field": "SHA1", "Value": file.SHA1})
	table.AddRow(map[string]interface{}{"Field": "SHA256", "Value": file.SHA256})
	// table.AddRow(map[string]interface{}{"Field": "SHA512", "Value": file.SHA512})
	table.AddRow(map[string]interface{}{"Field": "Mime", "Value": file.Mime})
	table.Markdown = true
	table.Print()
}

// PrintFileDetails prints file details
func (file *File) PrintFileDetails() {
	table := clitable.New([]string{"Field", "Value"})
	table.AddRow(map[string]interface{}{"Field": "Name", "Value": file.Name})
	table.AddRow(map[string]interface{}{"Field": "Path", "Value": file.Path})
	// table.AddRow(map[string]interface{}{"Field": "Valid", "Value": file.Valid})
	table.AddRow(map[string]interface{}{"Field": "Size", "Value": file.Size})
	// table.AddRow(map[string]interface{}{"Field": "CRC32", "Value": file.CRC32})
	table.AddRow(map[string]interface{}{"Field": "MD5", "Value": file.MD5})
	table.AddRow(map[string]interface{}{"Field": "SHA1", "Value": file.SHA1})
	table.AddRow(map[string]interface{}{"Field": "SHA256", "Value": file.SHA256})
	table.AddRow(map[string]interface{}{"Field": "SHA512", "Value": file.SHA512})
	// table.AddRow(map[string]interface{}{"Field": "Ssdeep", "Value": file.Ssdeep})
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
