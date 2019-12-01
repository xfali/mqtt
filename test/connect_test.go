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

func TestConnect1(t *testing.T) {
    msg := message.NewConnectMessage()
    msg.SetWillEnable(true)
    msg.SetUsername("test")
    msg.SetPassword([]byte("123"))
    msg.SetWillTopic("123")
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

    msg2 := message.NewConnectMessage()
    n3, e3 := msg2.ReadVariableHeader(buf)
    if e3 != nil {
        t.Fatal(e3)
    }
    n4, e4 := msg2.ReadPayload(buf)
    if e4 != nil {
        t.Fatal(e4)
    }

    if n + n2 != n3 + n4 {
        t.Fatal("not match")
    }

    t.Log("after")
    t.Log(msg2)
}

func TestConnect2(t *testing.T) {
    msg := message.NewConnectMessage()
    msg.SetWillEnable(true)
    msg.SetUsername("test")
    msg.SetPassword([]byte("123"))
    msg.SetWillTopic("/d12t13/t43uyh/45eu/65eiu/45u34y34syhg/eg435wuyherg345syh")
    msg.SetSessionExpiryInterval(100)
    msg.SetCorrelationData([]byte("fasfasfasc3tgergsgsdgsdgds"))
    msg.SetWillPayload([]byte("fasfascxasfs"))
    msg.SetContentType("fjlasjflkjaskflasjfkl")
    msg.SetAuthenticationData([]byte("fdg23y3h54uh564u3yhjhfxju54u"))
    msg.SetUserProperty(map[string]string{
        "test1": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test2": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test3": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test4": "1234567890qwertyuiopasdfghjklzxcvbnm",
    })
    msg.SetPayloadUserProperty(map[string]string{
        "test1": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test2": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test3": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test4": "1234567890qwertyuiopasdfghjklzxcvbnm",
    })
    msg.SetWillDelayInterval(100)
    msg.SetContentType("json")
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

    t.Log(msg2.(*message.ConnectMessage).GetCorrelationData())
}
