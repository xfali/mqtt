// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "bytes"
    "mqtt/message"
    "testing"
)

func TestConnect(t *testing.T) {
    conn := message.NewConnectMessage()
    conn.SetWillEnable(true)
    conn.SetUsername("test")
    conn.SetPassword([]byte("123"))
    conn.SetWillTopic("123")
    conn.Count()

    t.Log("before")
    t.Log(conn)

    buf := bytes.NewBuffer(nil)
    n, e := conn.WriteVarHeader(buf)
    if e != nil {
        t.Fatal(e)
    }
    n2, e2 := conn.WritePayload(buf)
    if e2 != nil {
        t.Fatal(e2)
    }

    conn2 := message.NewConnectMessage()
    n3, e3 := conn2.ReadVarHeader(buf)
    if e3 != nil {
        t.Fatal(e3)
    }
    n4, e4 := conn2.ReadPayload(buf)
    if e4 != nil {
        t.Fatal(e4)
    }

    if n + n2 != n3 + n4 {
        t.Fatal("not match")
    }

    t.Log("after")
    t.Log(conn2)
}
