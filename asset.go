package themekit

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Asset struct {
	Key        string `json:"key"`
	Value      string `json:"value,omitempty"`
	Attachment string `json:"attachment,omitempty"`
}

func (a Asset) String() string {
	result := fmt.Sprintf("key: %s", a.Key)
	if len(a.Value) > 0 {
		result = fmt.Sprintf("%s | value: %d chars", result, len(a.Value))
	}

	if len(a.Attachment) > 0 {
		result = fmt.Sprintf("%s | attachment: %d chars", result, len(a.Attachment))
	}

	return fmt.Sprintf("<%s>", result)
}

func (a Asset) IsValid() bool {
	return len(a.Key) > 0 && (len(a.Value) > 0 || len(a.Attachment) > 0)
}

// Implementing sort.Interface
type ByAsset []Asset

func (assets ByAsset) Len() int {
	return len(assets)
}

func (assets ByAsset) Swap(i, j int) {
	assets[i], assets[j] = assets[j], assets[i]
}

func (assets ByAsset) Less(i, j int) bool {
	return assets[i].Key < assets[j].Key
}

func LoadAsset(root, filename string) (asset Asset, err error) {
	asset = Asset{}
	path := toSlash(fmt.Sprintf("%s/%s", root, filename))
	file, err := os.Open(path)
	if err != nil {
		return asset, errors.New(fmt.Sprintf("LoadAsset: %s", err))
	}
	info, err := os.Stat(path)
	if err != nil {
		return asset, errors.New(fmt.Sprintf("LoadAsset: %s", err))
	}
	defer file.Close()

	if info.IsDir() {
		err = errors.New("LoadAsset: File is a directory")
		return
	}

	buffer := make([]byte, info.Size())
	_, err = file.Read(buffer)
	if err != nil {
		return asset, errors.New(fmt.Sprintf("LoadAsset: %s", err))
	}

	asset = Asset{Key: toSlash(filename)}
	if contentTypeFor(buffer) == "text" {
		asset.Value = string(buffer)
	} else {
		asset.Attachment = encode64(buffer)
	}
	return
}

func toSlash(path string) string {
	newpath := filepath.ToSlash(path)
	if strings.Index(newpath, "\\") >= 0 {
		newpath = strings.Replace(newpath, "\\", "/", -1)
	}
	return newpath
}

func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func contentTypeFor(data []byte) string {
	contentType := http.DetectContentType(data)
	if strings.Contains(contentType, "text") {
		return "text"
	} else {
		return "binary"
	}
}
