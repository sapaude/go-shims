package shim

import (
	"mime"
)

var extMimeTypes = map[string]string{
	".md": "text/markdown",
}

func init() {
	for ext, typ := range extMimeTypes {
		mime.AddExtensionType(ext, typ)
	}
}

// GetExtensionByMimeType 基于Mime type获取对应后缀
func GetExtensionByMimeType(mimeType string) (string, error) {
	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", err
	}

	if len(exts) > 0 {
		return exts[0], nil
	}

	return "", nil
}
