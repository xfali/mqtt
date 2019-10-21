// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "mqtt/packet"
    "testing"
)

func TestVarint(t *testing.T) {
    buf := make([]byte, 8)
    t.Run("1byte", func(t *testing.T) {
        n := packet.EncodeVarint(buf, 0x7F)
        t.Log("n : ", n)
        v := packet.DecodeVarint(buf)
        t.Log("v : ", v)
    })
    t.Run("2byte", func(t *testing.T) {
        n := packet.EncodeVarint(buf, 0x3FFF)
        t.Log("n : ", n)
        v := packet.DecodeVarint(buf)
        t.Log("v : ", v)
    })
    t.Run("3byte", func(t *testing.T) {
        n := packet.EncodeVarint(buf, 0x1FFFFF)
        t.Log("n : ", n)
        v := packet.DecodeVarint(buf)
        t.Log("v : ", v)
    })
    t.Run("4byte", func(t *testing.T) {
        n := packet.EncodeVarint(buf, 0xFFFFFFF)
        t.Log("n : ", n)
        v := packet.DecodeVarint(buf)
        t.Log("v : ", v)
    })
    t.Run("5byte", func(t *testing.T) {
        n := packet.EncodeVarint(buf, 0xFFFFFFF+1)
        t.Log("n : ", n)
        v := packet.DecodeVarint(buf)
        t.Log("v : ", v)
    })
}
