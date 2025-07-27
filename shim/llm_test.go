package shim

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestExtractPotentialJSON(t *testing.T) {
    type args struct {
        rawResponse string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {"t1", args{rawResponse: "xxx```json\n{\"a\":\"100\"}"}, `{"a":"100"}`},
        {"t2", args{rawResponse: "xxx```json\n{\"a\":\"100\"}ddd"}, `{"a":"100"}`},
        {"t3", args{rawResponse: "[{\"a\":\"100\"}]"}, `[{"a":"100"}]`},
        {"t4", args{rawResponse: "{\"a\":\"100\"}"}, `{"a":"100"}`},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equalf(t, tt.want, ExtractPotentialJSON(tt.args.rawResponse), "ExtractPotentialJSON(%v)", tt.args.rawResponse)
        })
    }
}
