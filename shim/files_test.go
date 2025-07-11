package shim

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExtensionByMimeType(t *testing.T) {
	type args struct {
		mimeType string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"t1", args{"text/markdown"}, ".md", assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetExtensionByMimeType(tt.args.mimeType)
			if !tt.wantErr(t, err, fmt.Sprintf("GetExtensionByMimeType(%v)", tt.args.mimeType)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetExtensionByMimeType(%v)", tt.args.mimeType)
		})
	}
}
