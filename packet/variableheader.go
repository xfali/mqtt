// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

import (
    "bytes"
    "encoding/binary"
    "io"
)

const (
    PayloadFormatIndicator          = 1  //0x01	载荷格式说明	字节	PUBLISH, Will Properties
    MessageExpiryInterval           = 2  //0x02	消息过期时间	四字节整数	PUBLISH, Will Properties
    ContentType                     = 3  //0x03	内容类型	UTF-8编码字符串	PUBLISH, Will Properties
    ResponseTopic                   = 8  //0x08	响应主题	UTF-8编码字符串	PUBLISH, Will Properties
    CorrelationData                 = 9  //0x09	相关数据	二进制数据	PUBLISH, Will Properties
    SubscriptionIdentifier          = 11 //0x0B	定义标识符	变长字节整数	PUBLISH, SUBSCRIBE
    SessionExpiryInterval           = 17 //0x11	会话过期间隔	四字节整数	CONNECT, CONNACK, DISCONNECT
    AssignedClientIdentifier        = 18 //0x12	分配客户标识符	UTF-8编码字符串	CONNACK
    ServerKeepAlive                 = 19 //0x13	服务端保活时间	双字节整数	CONNACK
    AuthenticationMethod            = 21 //0x15	认证方法	UTF-8编码字符串	CONNECT, CONNACK, AUTH
    AuthenticationData              = 22 //0x16	认证数据	二进制数据	CONNECT, CONNACK, AUTH
    RequestProblemInformation       = 23 //0x17	请求问题信息	字节	CONNECT
    WillDelayInterval               = 24 //0x18	遗嘱延时间隔	四字节整数	Will Properties
    RequestResponseInformation      = 25 //0x19	请求响应信息	字节	CONNECT
    ResponseInformation             = 26 //0x1A	请求信息	UTF-8编码字符串	CONNACK
    ServerReference                 = 28 //0x1C	服务端参考	UTF-8编码字符串	CONNACK, DISCONNECT
    ReasonString                    = 31 //0x1F	原因字符串	UTF-8编码字符串	CONNACK, PUBACK, PUBREC, PUBREL, PUBCOMP, SUBACK, UNSUBACK, DISCONNECT, AUTH
    ReceiveMaximum                  = 33 //0x21	接收最大数量	双字节整数	CONNECT, CONNACK
    TopicAliasMaximum               = 34 //0x22	主题别名最大长度	双字节整数	CONNECT, CONNACK
    TopicAlias                      = 35 //0x23	主题别名	双字节整数	PUBLISH
    MaximumQoS                      = 36 //0x24	最大QoS	字节	CONNACK
    RetainAvailable                 = 37 //0x25	保留属性可用性	字节	CONNACK
    UserProperty                    = 38 //0x26	用户属性	UTF-8字符串对	CONNECT, CONNACK, PUBLISH, Will Properties, PUBACK, PUBREC, PUBREL, PUBCOMP, SUBSCRIBE, SUBACK, UNSUBSCRIBE, UNSUBACK, DISCONNECT, AUTH
    MaximumPacketSize               = 39 //0x27	最大报文长度	四字节整数	CONNECT, CONNACK
    WildcardSubscriptionAvailable   = 40 //0x28	通配符订阅可用性	字节	CONNACK
    SubscriptionIdentifierAvailable = 41 //0x29	订阅标识符可用性	字节	CONNACK
    SharedSubscriptionAvailable     = 42 //0x2A	共享订阅可用性	字节	CONNACK
)

type Identifier uint16

type Property struct {
    DecHex []byte
    Data   []byte
}

type PropOpt func(io.Writer) error

func BuildProperty(w io.Writer, opts ... PropOpt) (err error) {
    buf := bytes.NewBuffer(make([]byte, 1024))
    for _, opt := range opts {
        err = opt(buf)
        if err != nil {
            return err
        }
    }
    length := buf.Len()
    data := buf.Bytes()
    varBuf := make([]byte, MaxVarUintBufSize)
    n := EncodeVaruint(varBuf, uint64(length))
    _, err = w.Write(varBuf[:n])
    if err != nil {
        return err
    }
    _, err = w.Write(data)
    return
}

func SetSessionExpiryInterval(v uint32) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 5)
        b[0] = SessionExpiryInterval
        binary.BigEndian.PutUint32(b[1:], v)
        _, err := w.Write(b)
        return err
    }
}

func SetReceiveMaximum(v uint16) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 3)
        b[0] = ReceiveMaximum
        binary.BigEndian.PutUint16(b[1:], v)
        _, err := w.Write(b)
        return err
    }
}

func SetMaximumPacketSize(v uint32) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 5)
        b[0] = MaximumPacketSize
        binary.BigEndian.PutUint32(b[1:], v)
        _, err := w.Write(b)
        return err
    }
}

func SetTopicAliasMaximum(v uint16) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 3)
        b[0] = TopicAliasMaximum
        binary.BigEndian.PutUint16(b[1:], v)
        _, err := w.Write(b)
        return err
    }
}

func SetRequestResponseInformation(v byte) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 2)
        b[0] = RequestResponseInformation
        b[1] = v
        _, err := w.Write(b)
        return err
    }
}

func SetRequestProblemInformation(v byte) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 2)
        b[0] = RequestProblemInformation
        b[1] = v
        _, err := w.Write(b)
        return err
    }
}

func SetUserProperty(v map[string]string) PropOpt {
    return func(w io.Writer) error {
        _, err := w.Write([]byte{UserProperty})
        if err != nil {
            return err
        }
        for k, v := range v {
            _, err1 := EncodeString(w, k)
            if err1 != nil {
                return err1
            }
            _, err2 := EncodeString(w, v)
            if err2 != nil {
                return err2
            }
        }
        return nil
    }
}

func SetAuthenticationMethod(v string) PropOpt {
    return func(w io.Writer) error {
        _, err := w.Write([]byte{AuthenticationMethod})
        if err != nil {
            return err
        }
        _, err1 := EncodeString(w, v)
        return err1
    }
}

func SetAuthenticationData(v []byte) PropOpt {
    return func(w io.Writer) error {
        _, err := w.Write([]byte{AuthenticationData})
        if err != nil {
            return err
        }
        _, err1 := w.Write(v)
        return err1
    }
}
