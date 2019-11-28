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
    p1.V = 12

    p2 := &packet.PropUserProperty{}
    s, _ := packet.NewStringPair("test1", "test2")
    p2.V = s

    buf := bytes.NewBuffer(nil)
    n, err := packet.WriteProperties(buf, []packet.Property{
        p1, p2,
    })
    if err != nil {
        t.Fatal(err)
    }
    t.Log("write n ", n)

    t.Log(buf.Bytes())

    props, n, err := packet.ReadPropertyMap(buf)
    if err != nil {
        t.Fatal(err)
    }
    t.Log("read n ", n)
    p := props[packet.SessionExpiryInterval]
    if p != nil {
        t.Log(p.Get().(uint32))
    }

    p = props[packet.UserProperty]
    if p != nil {
        pair := p.Get().(packet.StringPair)
        t.Log(pair[0].String())
        t.Log(pair[1].String())
    }
}

func TestProperty2(t *testing.T) {
    p1 := &packet.PropSessionExpiryInterval{}
    p1.V = 12

    p2 := &packet.PropUserProperty{}
    s, _ := packet.NewStringPair("test1", "test2")
    p2.V = s

    buf := bytes.NewBuffer(nil)
    n, err := packet.WriteProperties(buf, []packet.Property{
        p1, p2,
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

    p3 := packet.FindPropValue(packet.SessionExpiryInterval, props).(*packet.PropSessionExpiryInterval)
    if p1.V != p3.V {
        t.Fatal("not match")
    }
}

func TestProperty3(t *testing.T) {
    p1 := &packet.PropSessionExpiryInterval{}
    p1.V = 12

    p2 := &packet.PropUserProperty{}
    s, _ := packet.NewStringPair("test1", "test2")
    p2.V = s

    buf := bytes.NewBuffer(nil)
    n, err := packet.WriteProperties(buf, []packet.Property{
        p1, p2,
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

    p3 := &packet.PropSessionExpiryInterval{}
    packet.FindAndSetPropValue(p3, props)
    if p1.V != p3.V {
        t.Fatal("not match")
    }

    p4 := &packet.PropUserProperty{}
    packet.FindAndSetPropValue(p4, props)
    if !p2.V.Equals(p4.V) {
        t.Fatal("not match")
    }
}
