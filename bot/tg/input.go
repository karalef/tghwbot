package tg

import (
	"bytes"
	"io"
)

// InputFile represents the contents of a file to be uploaded.
type InputFile struct {
	urlID string
	name  string
	data  io.Reader
}

// Data returns the file data to send when a file does not need to be uploaded.
func (f InputFile) Data() string {
	return f.urlID
}

// UploadData returns the file name and data reader for the file.
func (f InputFile) UploadData() (string, io.Reader) {
	return f.name, f.data
}

// FileID returns InputFile that has already been uploaded to Telegram.
func FileID(fid string) *InputFile {
	return &InputFile{urlID: fid}
}

// FileURL return InputFile as URL that is used without uploading.
func FileURL(url string) *InputFile {
	return &InputFile{urlID: url}
}

// FileReader returns InputFile that needs to be uploaded.
func FileReader(name string, data io.Reader) *InputFile {
	return &InputFile{
		name: name,
		data: data,
	}
}

// FileBytes returns InputFile as bytes that needs to be uploaded.
func FileBytes(name string, data []byte) *InputFile {
	return &InputFile{
		name: name,
		data: bytes.NewReader(data),
	}
}
