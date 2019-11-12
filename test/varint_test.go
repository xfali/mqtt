// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "bytes"
    "encoding/binary"
    "mqtt/packet"
    "testing"
)

func TestVarint(t *testing.T) {
    buf := make([]byte, 10)
    t.Run("1byte", func(t *testing.T) {
        n := packet.EncodeVaruint(buf, 0x7F)
        t.Log("n : ", n)
        v := packet.DecodeVaruint(buf)
        t.Log("v : ", v)
    })
    t.Run("2byte", func(t *testing.T) {
        n := packet.EncodeVaruint(buf, 0x3FFF)
        t.Log("n : ", n)
        v := packet.DecodeVaruint(buf)
        t.Log("v : ", v)
    })
    t.Run("3byte", func(t *testing.T) {
        n := packet.EncodeVaruint(buf, 0x1FFFFF)
        t.Log("n : ", n)
        v := packet.DecodeVaruint(buf)
        t.Log("v : ", v)
    })
    t.Run("4byte", func(t *testing.T) {
        n := packet.EncodeVaruint(buf, 0xFFFFFFF)
        t.Log("n : ", n)
        v := packet.DecodeVaruint(buf)
        t.Log("v : ", v)
    })
    t.Run("5byte", func(t *testing.T) {
        n := packet.EncodeVaruint(buf, 0xFFFFFFF+1)
        t.Log("n : ", n)
        v := packet.DecodeVaruint(buf)
        t.Log("v : ", v)
    })
}

func TestVarint2(t *testing.T) {
    buf := make([]byte, 10)
    n := packet.EncodeVaruint(buf, 0xFFFFFFF+1)
    i := packet.DecodeVarint(buf[:n])
    t.Log("i : ", i, "v", 0xFFFFFFF+1)
}

func TestVarint3(t *testing.T) {
    v := packet.VarInt{}
    buf := make([]byte, 10)
    n := packet.EncodeVaruint(buf, 0xFFFFFFF+1)
    for x:=1; x <= n; x++ {
        if v.Load(buf[x-1:x]) {
            break
        }
    }

    i := v.ToInt()
    t.Log("i : ", i, "v", 0xFFFFFFF+1)
}


func TestVarint4(t *testing.T) {
    buf := make([]byte, 10)
    n := packet.EncodeVaruint(buf, 0xFFFFFFF+1)
    t.Run("uint", func(t *testing.T) {
        v, err := binary.ReadUvarint(bytes.NewReader(buf[:n]))
        if err != nil {
            t.Fatal(err)
        }
        t.Log("i : ", v, "v", 0xFFFFFFF+1)
    })
    t.Run("int", func(t *testing.T) {
        v, err := binary.ReadVarint(bytes.NewReader(buf[:n]))
        if err != nil {
            t.Fatal(err)
        }
        t.Log("i : ", v, "v", 0xFFFFFFF+1)
    })
}

func TestVarint5(t *testing.T) {
    buf := make([]byte, 10)
    n := packet.EncodeVaruint(buf, 0xFFFFFFF+1)
    t.Logf("origin value %d\n", 0xFFFFFFF+1)
    t.Run("LoadFromReader", func(t *testing.T) {
        r := bytes.NewReader(buf[:n])
        v := packet.VarInt{}
        _, err := v.LoadFromReader(r)
        if err != nil {
            t.Fatal(err)
        }
        t.Logf("value %d\n", v.ToInt())
    })

    t.Run("NewFromReader", func(t *testing.T) {
        r := bytes.NewReader(buf[:n])
        v := packet.NewFromReader(r)
        if v == nil {
            t.Fatal("parse error")
        }
        t.Logf("value %d\n", v.ToInt())
    })
}
