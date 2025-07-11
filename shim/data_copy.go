package shim

import (
    "bytes"
    "encoding/gob"
    "encoding/json"
)

// DeepCopyByJSON 使用JSON进行深度拷贝，Warning: 未导出的便利无法进行深拷贝
func DeepCopyByJSON(src, dst interface{}) error {
    buf, err := json.Marshal(src)
    if err != nil {
        return err
    }
    return json.Unmarshal(buf, dst)
}

// DeepCopyByGob 使用Gob进行深度拷贝，Warning: 未导出的便利无法进行深拷贝
func DeepCopyByGob(src, dst interface{}) error {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    dec := gob.NewDecoder(&buf)

    if err := enc.Encode(src); err != nil {
        return err
    }
    return dec.Decode(dst)
}
