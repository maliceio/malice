package persist

import (
	"bytes"
	"compress/gzip"
	"context"
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
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-units"
	"github.com/dustin/go-jsonpointer"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/maliceio/malice/malice/malutils"
	// "github.com/dutchcoders/gossdeep"
)

// File is a file object
type File struct {
	Name string `json:"name,omitempty" structs:"name"`
	Path string `json:"path,omitempty" structs:"path"`
	// Valid bool   `json:"valid"`
	Size string `json:"size,omitempty" structs:"size"`
	// CRC32  string
	MD5    string `json:"md5,omitempty" structs:"md5"`
	SHA1   string `json:"sha1,omitempty" structs:"sha1"`
	SHA256 string `json:"sha256,omitempty" structs:"sha256"`
	SHA512 string `json:"sha512,omitempty" structs:"sha512"`
	// Ssdeep string `json:"ssdeep"`
	// Arch string `json:"arch"`
}

// Init initializes the File object
func (file *File) Init() {

	if file.Path == "" {
		log.Fatalf("error occured during file.Init() because file.Path was not set.")
	}

	file.GetName()
	file.GetSize()

	// Read in file data
	dat, err := ioutil.ReadFile(file.Path)
	utils.Assert(err)

	file.GetMD5(dat)
	file.GetSHA1(dat)
	file.GetSHA256(dat)
	file.GetSHA512(dat)
}

// GetMimeType returns file's mime type
func GetMimeType(docker *client.Docker, arg string) (string, error) {

	// Create Container
	createContConf := &container.Config{
		Image: "malice/fileinfo",
		Cmd:   []string{"-m", arg},
	}
	resources := container.Resources{
		Memory:   config.Conf.Docker.Memory, // Memory    int64 // Memory limit (in bytes)
		NanoCPUs: config.Conf.Docker.CPU,    // NanoCPUs  int64 `json:"NanoCpus"` // CPU quota in units of 10<sup>-9</sup> CPUs.
	}
	log.WithFields(log.Fields{
		"func": "persist.GetMimeType",
		"mem":  config.Conf.Docker.Memory,
		"cpu":  config.Conf.Docker.CPU,
	}).Debug("setting container resources")
	hostConfig := &container.HostConfig{
		Privileged:  false,
		Binds:       []string{config.Conf.Docker.Binds},
		NetworkMode: "none",
		Resources:   resources,
		AutoRemove:  true,
	}
	networkingConfig := &network.NetworkingConfig{}

	contResponse, err := docker.Client.ContainerCreate(context.Background(), createContConf, hostConfig, networkingConfig, "getmimetype")
	if err != nil {
		return "", err
	}

	// Start Container
	err = docker.Client.ContainerStart(context.Background(), contResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	log.WithFields(log.Fields{
		"id":   contResponse.ID,
		"env":  config.Conf.Environment.Run,
		"func": "persist.GetMimeType",
	}).Debug("malice/fileinfo Container Started")

	defer func() {
		// remove container when done
		contRmOpts := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   false,
			Force:         true,
		}
		er.CheckError(docker.Client.ContainerRemove(context.Background(), "getmimetype", contRmOpts))
		log.WithFields(log.Fields{
			"id":   contResponse.ID,
			"env":  config.Conf.Environment.Run,
			"func": "persist.GetMimeType",
		}).Debug("malice/fileinfo Container Removed")
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}
	// Catch Container's Output
	reader, err := docker.Client.ContainerLogs(ctx, contResponse.ID, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var buf1 bytes.Buffer
	w := io.Writer(&buf1)

	_, err = stdcopy.StdCopy(w, os.Stderr, reader)
	if err != nil && err != io.EOF {
		return "", err
	}

	mimetype := strings.TrimSpace(buf1.String())
	log.Debug("File has mimetype: ", mimetype)

	return strings.TrimSpace(buf1.String()), nil
}

// GetFileInfo start malice/fileinfo container and extract certain fields with a search string
func GetFileInfo(docker *client.Docker, arg string, search string) (string, error) {

	// Create Container
	createContConf := &container.Config{
		Image: "malice/fileinfo",
		Cmd:   []string{arg},
	}
	hostConfig := &container.HostConfig{
		Privileged:  false,
		Binds:       []string{config.Conf.Docker.Binds},
		NetworkMode: "none",
	}
	networkingConfig := &network.NetworkingConfig{}

	contResponse, err := docker.Client.ContainerCreate(context.Background(), createContConf, hostConfig, networkingConfig, "")
	if err != nil {
		return "", err
	}

	// Start Container
	err = docker.Client.ContainerStart(context.Background(), contResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}
	// Catch Container's Output
	reader, err := docker.Client.ContainerLogs(ctx, contResponse.ID, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var buf1 bytes.Buffer
	w := io.Writer(&buf1)

	_, err = stdcopy.StdCopy(w, os.Stderr, reader)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Find fields in JSON output with search string
	found, err := jsonpointer.Find(buf1.Bytes(), search)
	if err != nil {
		return "", err
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		contRmOpts := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   true,
			Force:         true,
		}
		er.CheckError(docker.Client.ContainerRemove(ctx, contResponse.ID, contRmOpts))
		log.WithFields(log.Fields{
			"id":   contResponse.ID,
			"env":  config.Conf.Environment.Run,
			"func": "persist.GetFileInfo",
		}).Debug("malice/fileinfo Container Removed")
	}()

	return string(found), nil
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

	destination := filepath.Join(maldirs.GetSampledsDir(), file.SHA256+".gz")
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
func (file *File) GetMD5(data []byte) (hMd5Sum string, err error) {

	hmd5 := md5.New()
	_, err = hmd5.Write(data)
	utils.Assert(err)
	hMd5Sum = fmt.Sprintf("%x", hmd5.Sum(nil))

	file.MD5 = hMd5Sum

	return
}

// GetSHA1 calculates the Files sha256sum
func (file *File) GetSHA1(data []byte) (h1Sum string, err error) {

	h1 := sha1.New()
	_, err = h1.Write(data)
	utils.Assert(err)
	h1Sum = fmt.Sprintf("%x", h1.Sum(nil))

	file.SHA1 = h1Sum

	return
}

// GetSHA256 calculates the Files sha256sum
func (file *File) GetSHA256(data []byte) (h256Sum string, err error) {

	h256 := sha256.New()
	_, err = h256.Write(data)
	utils.Assert(err)
	h256Sum = fmt.Sprintf("%x", h256.Sum(nil))

	file.SHA256 = h256Sum

	return
}

// GetSHA512 calculates the Files sha256sum
func (file *File) GetSHA512(data []byte) (h512Sum string, err error) {

	h512 := sha512.New()
	_, err = h512.Write(data)
	utils.Assert(err)
	h512Sum = fmt.Sprintf("%x", h512.Sum(nil))

	file.SHA512 = h512Sum

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
	// table.AddRow(map[string]interface{}{"Field": "Mime", "Value": file.Mime})
	// table.AddRow(map[string]interface{}{"Field": "Magic", "Value": file.Magic})
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
	// table.AddRow(map[string]interface{}{"Field": "Mime", "Value": file.Mime})
	// table.AddRow(map[string]interface{}{"Field": "Magic", "Value": file.Magic})
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
