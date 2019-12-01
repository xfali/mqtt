// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "fmt"
    "io"
    "mqtt/errcode"
    "mqtt/packet"
    "mqtt/util"
    "strings"
)

type DisconnectVarHeader struct {
    //断开原因码
    ReasonCode byte
    props      []packet.Property
}

//DISCONNECT报文是客户端发给服务端的最后一个MQTT控制报文。
//表示客户端为什么断开网络连接的原因。客户端和服务端在关闭网络连接之前可以发送一个DISCONNECT报文。
//如果在客户端没有首先发送包含原因码为0x00（正常断开）DISCONNECT报文并且连接包含遗嘱消息的情况下，遗嘱消息会被发布。
type DisconnectMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   DisconnectVarHeader
}

func (m *DisconnectMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *DisconnectMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

func (msg *DisconnectMessage) ReadVariableHeader(r io.Reader) (int, error) {
    if msg.fixedHeader.RemainLength() < 1 {
        return 0, nil
    }
    buf := make([]byte, 1)
    n, err := r.Read(buf)
    if err != nil {
        return n, err
    }

    msg.varHeader.ReasonCode = buf[0]

    if msg.fixedHeader.RemainLength() == int64(n) {
        return n, nil
    }

    props, n2, err2 := packet.ReadProperties(r)
    n += n2
    if err2 != nil {
        return n, err2
    }

    msg.varHeader.props = props

    return n, nil
}

func (msg *DisconnectMessage) WriteVariableHeader(w io.Writer) (int, error) {
    if msg.varHeader.ReasonCode == 0 {
        if  len(msg.varHeader.props) == 0 {
            return 0, nil
        } else {
            return 0, errcode.ProtocolError
        }
    }

    n, err := w.Write([]byte{msg.varHeader.ReasonCode})
    if err != nil {
        return n, err
    }

    n2, err2 := packet.WriteProperties(w, msg.varHeader.props)
    return n + n2, err2
}

func (msg *DisconnectMessage) ReadPayload(r io.Reader) (int, error) {
    return 0, nil
}

func (msg *DisconnectMessage) WritePayload(w io.Writer) (int, error) {
    return 0, nil
}

func (m *DisconnectMessage) Valid() bool {
    return true
}

func NewDisconnectMessage() *DisconnectMessage {
    ret := &DisconnectMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypeDISCONNECT, packet.PktFlagDISCONNECT, 0),
    }
    return ret
}

func (m *DisconnectMessage) SetReasonCode(v byte) {
    m.varHeader.ReasonCode = v
}

func (m *DisconnectMessage) GetReasonCode() byte {
    return m.varHeader.ReasonCode
}

// 会话过期间隔 Session Expiry Interval,四字节整数表示的以秒为单位的会话过期间隔
func (m *DisconnectMessage) SetSessionExpiryInterval(v uint32) {
    p := &packet.PropSessionExpiryInterval{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

// 会话过期间隔 Session Expiry Interval,四字节整数表示的以秒为单位的会话过期间隔
func (m *DisconnectMessage) GetSessionExpiryInterval() (uint32, bool) {
    p := packet.FindPropValue(packet.SessionExpiryInterval, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropSessionExpiryInterval).V, true
}

//UTF-8编码的字符串，表示此次响应相关的原因。
// 此原因字符串（Reason String）是为诊断而设计的可读字符串，不应该被客户端所解析。
func (m *DisconnectMessage) SetReasonString(v string) {
    p := &packet.PropReasonString{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//UTF-8编码的字符串，表示此次响应相关的原因。
// 此原因字符串（Reason String）是为诊断而设计的可读字符串，不应该被客户端所解析。
func (m *DisconnectMessage) GetReasonString() (string, bool) {
    p := packet.FindPropValue(packet.ReasonString, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropReasonString).V.String(), true
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *DisconnectMessage) SetUserProperty(props map[string]string) {
    for k, v := range props {
        p := &packet.PropUserProperty{}
        pair, err := packet.NewStringPair(k, v)
        if err == nil {
            p.V = pair
            m.varHeader.props = append(m.varHeader.props, p)
        }
    }
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *DisconnectMessage) GetUserProperty() (map[string]string, bool) {
    ret := map[string]string{}
    packet.FindPropValues(packet.UserProperty, m.varHeader.props, func(property packet.Property) bool {
        if property != nil {
            p := property.(*packet.PropUserProperty)
            ret[p.V[0].String()] = p.V[1].String()
        }
        return false
    })
    return ret, len(ret) > 0
}

//服务端参考（Server Reference）标识符
//客户端可以使用它来识别其他要使用的服务端。
// 包含多个服务端参考将造成协议错误（Protocol Error）。
func (m *DisconnectMessage) SetServerReference(v string) {
    p := &packet.PropServerReference{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//服务端参考（Server Reference）标识符
//客户端可以使用它来识别其他要使用的服务端。
// 包含多个服务端参考将造成协议错误（Protocol Error）。
func (m *DisconnectMessage) GetServerReference() (string, bool) {
    p := packet.FindPropValue(packet.ServerReference, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropServerReference).V.String(), true
}

func (v *DisconnectVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("ReasonCode: %d \nprops:\n%s",
        v.ReasonCode, builder.String())
}

func (v *DisconnectMessage) String() string {
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%s\n",
        v.fixedHeader, v.varHeader.String())
}
