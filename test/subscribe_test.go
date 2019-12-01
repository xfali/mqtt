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

func TestSubscribe1(t *testing.T) {
    conn := message.NewSubscribeMessage()
    conn.SetSubscriptionIdentifier(1000)
    conn.SetPayload([]message.SubscribeFilter{
        {"test", 1},
        {"test2", 2},
        {"test3", 3},
    })
    conn.GetFixedHeader()

    t.Log("before")
    t.Log(conn)

    buf := bytes.NewBuffer(nil)
    n, e := conn.WriteVariableHeader(buf)
    if e != nil {
        t.Fatal(e)
    }
    n2, e2 := conn.WritePayload(buf)
    if e2 != nil {
        t.Fatal(e2)
    }

    conn2 := message.NewSubscribeMessage()
    conn2.SetFixedHeader(conn.GetFixedHeader())
    n3, e3 := conn2.ReadVariableHeader(buf)
    if e3 != nil {
        t.Fatal(e3)
    }
    n4, e4 := conn2.ReadPayload(buf)
    if e4 != nil {
        t.Fatal(e4)
    }

    if n+n2 != n3+n4 {
        t.Fatal("not match")
    }

    t.Log("after")
    t.Log(conn2)
}

func TestSubscribe2(t *testing.T) {
    conn := message.NewSubscribeMessage()
    conn.SetPayload([]message.SubscribeFilter{
        {"test", 1},
        {"test2", 2},
        {"test3", 3},
    })
    conn.SetSubscriptionIdentifier(123)
    conn.SetUserProperty(map[string]string{
        "test1": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test2": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test3": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test4": "1234567890qwertyuiopasdfghjklzxcvbnm",
    })
    conn.GetFixedHeader()

    t.Log("before")
    t.Log(conn)

    buf := bytes.NewBuffer(nil)
    n, err := message.WriteMessage(buf, conn)
    if err != nil {
        t.Fatal(err)
    }

    conn2, n2, err2 := message.ReadMessage(buf)
    if err2 != nil {
        t.Fatal(err2)
    }

    if n != n2 {
        t.Fatal("not match")
    }

    t.Log("after")
    t.Log(conn2)

    t.Log(conn2.(*message.SubscribeMessage).GetSubscriptionIdentifier())
}
