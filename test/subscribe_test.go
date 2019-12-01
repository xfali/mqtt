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
    msg := message.NewSubscribeMessage()
    msg.SetSubscriptionIdentifier(1000)
    msg.SetPayload([]message.SubscribeFilter{
        {"test", 1},
        {"test2", 2},
        {"test3", 3},
    })
    msg.GetFixedHeader()

    t.Log("before")
    t.Log(msg)

    buf := bytes.NewBuffer(nil)
    n, e := msg.WriteVariableHeader(buf)
    if e != nil {
        t.Fatal(e)
    }
    n2, e2 := msg.WritePayload(buf)
    if e2 != nil {
        t.Fatal(e2)
    }

    msg2 := message.NewSubscribeMessage()
    msg2.SetFixedHeader(msg.GetFixedHeader())
    n3, e3 := msg2.ReadVariableHeader(buf)
    if e3 != nil {
        t.Fatal(e3)
    }
    n4, e4 := msg2.ReadPayload(buf)
    if e4 != nil {
        t.Fatal(e4)
    }

    if n+n2 != n3+n4 {
        t.Fatal("not match")
    }

    t.Log("after")
    t.Log(msg2)
}

func TestSubscribe2(t *testing.T) {
    msg := message.NewSubscribeMessage()
    msg.SetPayload([]message.SubscribeFilter{
        {"test", 1},
        {"test2", 2},
        {"test3", 3},
    })
    msg.SetSubscriptionIdentifier(123)
    msg.SetUserProperty(map[string]string{
        "test1": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test2": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test3": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test4": "1234567890qwertyuiopasdfghjklzxcvbnm",
    })
    msg.GetFixedHeader()

    t.Log("before")
    t.Log(msg)

    buf := bytes.NewBuffer(nil)
    n, err := message.WriteMessage(buf, msg)
    if err != nil {
        t.Fatal(err)
    }

    msg2, n2, err2 := message.ReadMessage(buf)
    if err2 != nil {
        t.Fatal(err2)
    }

    if n != n2 {
        t.Fatal("not match")
    }

    t.Log("after")
    t.Log(msg2)

    t.Log(msg2.(*message.SubscribeMessage).GetSubscriptionIdentifier())
}
