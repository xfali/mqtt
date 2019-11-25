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

func TestProperty1(t *testing.T) {
    p1 := &packet.PropSessionExpiryInterval{}
    p1.V  = 12

    p2 := &packet.PropUserProperty{}
    p2.V, _ = packet.NewStringPair("test1", "test2")

    buf := bytes.NewBuffer(nil)
    n, err := packet.WriteProperties(buf, []packet.Property{
        p1,p2,
    })
    if err != nil {
        t.Fatal(err)
    }
    t.Log("write n ", n)

    t.Log(buf.Bytes())

    props, n, err := packet.ReadProperties(buf)
    if err != nil {
        t.Fatal(err)
    }
    t.Log("read n ", n)
    t.Log(props)
}
