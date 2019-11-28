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

type PubAckVarHeader struct {
    //报文标识符（Packet Identifier）字段才能出现在PUBLISH报文中。
    PacketIdentifier uint16

    ReasonCode byte
    //PUBLISH 属性 PUBLISH Properties
    props []packet.Property
}

type PubAckMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   PubAckVarHeader
}

func NewPubAckMessage() *PubAckMessage {
    ret := &PubAckMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypePUBACK, packet.PktFlagPUBACK, 0),
    }
    return ret
}

func (m *PubAckMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *PubAckMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

func (msg *PubAckMessage) ReadVariableHeader(r io.Reader) (int, error) {
    buf := make([]byte, 2)
    n, err := r.Read(buf)
    if err != nil {
        return n, err
    }
    msg.varHeader.PacketIdentifier = uint16(buf[0]<<8 | buf[1])
    if err := msg.fixedHeader.CheckLen(n); err != nil {
        return n, err
    }

    if msg.fixedHeader.RemainLength() == int64(n) {
        return n, nil
    }

    n2, err2 := r.Read(buf[:1])
    n += n2
    if err2 != nil {
        return n, err2
    }

    msg.varHeader.ReasonCode = buf[0]

    props, n3, err3 := packet.ReadProperties(r)
    n += n3
    if err3 != nil {
        return n, err3
    }

    msg.varHeader.props = props

    return n, nil
}

func (msg *PubAckMessage) WriteVariableHeader(w io.Writer) (int, error) {
    n, err := w.Write([]byte{
        byte(msg.varHeader.PacketIdentifier >> 8),
        byte(msg.varHeader.PacketIdentifier & 0xFF),
    })
    if err != nil {
        return n, err
    }

    if msg.varHeader.ReasonCode == errcode.ReasonSuccess && len(msg.varHeader.props) == 0 {
        return n, nil
    }

    n2, err2 := w.Write([]byte{msg.varHeader.ReasonCode})
    n += n2
    if err2 != nil {
        return n, err2
    }

    n3, err3 := packet.WriteProperties(w, msg.varHeader.props)
    return n + n3, err3
}

func (msg *PubAckMessage) ReadPayload(r io.Reader) (n int, err error) {
    return 0, nil
}

func (msg *PubAckMessage) WritePayload(w io.Writer) (int, error) {
    return 0, nil
}

func (msg *PubAckMessage) Valid() bool {
    return true
}

func (msg *PubAckMessage) SetReasonCode(v byte) {
    msg.varHeader.ReasonCode = v
}

func (msg *PubAckMessage) GetReasonCode() byte {
    return msg.varHeader.ReasonCode
}

func (msg *PubAckMessage) SetPacketIdentifier(v uint16) {
    msg.varHeader.PacketIdentifier = v
}

func (msg *PubAckMessage) GetPacketIdentifier() uint16 {
    return msg.varHeader.PacketIdentifier
}

//UTF-8编码的字符串，表示此次响应相关的原因。
// 此原因字符串（Reason String）是为诊断而设计的可读字符串，不应该被客户端所解析。
func (m *PubAckMessage) SetReasonString(v string) {
    p := &packet.PropReasonString{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//UTF-8编码的字符串，表示此次响应相关的原因。
// 此原因字符串（Reason String）是为诊断而设计的可读字符串，不应该被客户端所解析。
func (m *PubAckMessage) GetReasonString() (string, bool) {
    p := packet.FindPropValue(packet.ReasonString, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropReasonString).V.String(), true
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *PubAckMessage) SetUserProperty(props map[string]string) {
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
func (m *PubAckMessage) GetUserProperty() (map[string]string, bool) {
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

func (v *PubAckVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("TopicName: %d ReasonCode: %d \nprops:\n%s",
        v.PacketIdentifier, v.ReasonCode, builder.String())
}

func (v *PubAckMessage) String() string {
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%s\n",
        v.fixedHeader, v.varHeader.String())
}