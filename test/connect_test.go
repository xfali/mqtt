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
    conn := message.NewConnectMessage()
    conn.SetWillEnable(true)
    conn.SetUsername("test")
    conn.SetPassword([]byte("123"))
    conn.SetWillTopic("123")
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

    conn2 := message.NewConnectMessage()
    n3, e3 := conn2.ReadVariableHeader(buf)
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

func TestConnect2(t *testing.T) {
    conn := message.NewConnectMessage()
    conn.SetWillEnable(true)
    conn.SetUsername("test")
    conn.SetPassword([]byte("123"))
    conn.SetWillTopic("/d12t13/t43uyh/45eu/65eiu/45u34y34syhg/eg435wuyherg345syh")
    conn.SetSessionExpiryInterval(100)
    conn.SetCorrelationData([]byte("fasfasfasc3tgergsgsdgsdgds"))
    conn.SetWillPayload([]byte("fasfascxasfs"))
    conn.SetContentType("fjlasjflkjaskflasjfkl")
    conn.SetAuthenticationData([]byte("fdg23y3h54uh564u3yhjhfxju54u"))
    conn.SetUserProperty(map[string]string{
        "test1": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test2": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test3": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test4": "1234567890qwertyuiopasdfghjklzxcvbnm",
    })
    conn.SetPayloadUserProperty(map[string]string{
        "test1": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test2": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test3": "1234567890qwertyuiopasdfghjklzxcvbnm",
        "test4": "1234567890qwertyuiopasdfghjklzxcvbnm",
    })
    conn.SetWillDelayInterval(100)
    conn.SetContentType("json")
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

    t.Log(conn2.(*message.ConnectMessage).GetCorrelationData())
}
