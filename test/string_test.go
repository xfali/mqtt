// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "bytes"
    "io"
    "mqtt/packet"
    "testing"
)

func TestString1(t *testing.T) {
    s, err := packet.NewString("测试数据")
    if err != nil {
        t.Fatal(err)
    }
    t.Log("string:", s.String(), "length", s.Length(), "string len", len(s.String()))

    r := packet.StringReader{S: *s}
    buf := bytes.NewBuffer(nil)
    io.Copy(buf, &r)
    t.Log("buf", buf.Bytes())

    s2, _, err := packet.ParseString(buf)
    if err != nil {
        if err != io.EOF {
            t.Fatal(err)
        }
    }
    t.Log("string2", "string:", s2.String(), "length", s2.Length())
}
