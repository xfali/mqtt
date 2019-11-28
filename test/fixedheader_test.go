// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "bytes"
    "mqtt/packet"
    "testing"
)

func TestFixedHeader(t *testing.T) {
    header := packet.CreateFixedHeader(packet.PktTypeCONNECT, packet.PktFlagCONNECT, 10)
    buf := bytes.NewBuffer(nil)
    packet.WriteFixedHeader(buf, header)

    rh, _, _ := packet.ReadFixedHeader(buf)
    t.Log(rh)
}

func TestFixedHeader1(t *testing.T) {
    header := packet.CreateFixedHeader(packet.PktTypeCONNECT, packet.PktFlagCONNECT, 0)
    t.Logf("create : PktTypeCONNECT: %b\n", header.TypeFlag)

    if header.TypeFlag != 0x10 {
        t.Fatal()
    }
    header.SetQos(1)
    t.Logf("SetQos 1: %b\n", header.TypeFlag)
    if header.TypeFlag != 0x12 {
        t.Fatal()
    }

    header.SetQos(2)
    t.Logf("SetQos 2: %b\n", header.TypeFlag)

    if header.TypeFlag != 0x14 {
        t.Fatal()
    }

    header.SetQos(3)
    t.Logf("SetQos 3: %b\n", header.TypeFlag)

    if header.TypeFlag != 0x16 {
        t.Fatal()
    }

    header.SetDup(true)
    t.Logf("SetDup true: %b\n", header.TypeFlag)

    if header.TypeFlag != 0x1E {
        t.Fatal()
    }

    header.SetDup(false)
    t.Logf("SetDup false: %b\n", header.TypeFlag)
    if header.TypeFlag != 0x16 {
        t.Fatal()
    }

    header.SetRetain(true)
    t.Logf("SetRetain true: %b\n", header.TypeFlag)
    if header.TypeFlag != 0x17 {
        t.Fatal()
    }

    header.SetRetain(false)
    t.Logf("SetRetain false: %b\n", header.TypeFlag)
    if header.TypeFlag != 0x16 {
        t.Fatal()
    }
}
