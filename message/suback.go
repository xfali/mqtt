// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "bytes"
    "fmt"
    "io"
    "mqtt/packet"
    "mqtt/util"
    "strings"
)

type SubAckVarHeader struct {
    //报文标识符（Packet Identifier）。
    PacketIdentifier uint16

    //PUBLISH 属性 PUBLISH Properties
    props []packet.Property

    size int
}

//PUBACK报文是对QoS 1等级的PUBLISH报文的响应。
type SubAckMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   SubAckVarHeader
    payload     []byte
}

func NewSubAckMessage() *SubAckMessage {
    ret := &SubAckMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypeSUBACK, packet.PktFlagSUBACK, 0),
    }
    return ret
}

func (m *SubAckMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *SubAckMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

func (msg *SubAckMessage) ReadVariableHeader(r io.Reader) (int, error) {
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

    props, n2, err2 := packet.ReadProperties(r)
    n += n2
    if err2 != nil {
        return n, err2
    }

    msg.varHeader.props = props
    msg.varHeader.size = n

    return n, nil
}

func (msg *SubAckMessage) WriteVariableHeader(w io.Writer) (int, error) {
    n, err := w.Write([]byte{
        byte(msg.varHeader.PacketIdentifier >> 8),
        byte(msg.varHeader.PacketIdentifier & 0xFF),
    })
    if err != nil {
        return n, err
    }

    n2, err2 := packet.WriteProperties(w, msg.varHeader.props)
    return n + n2, err2
}

func (msg *SubAckMessage) ReadPayload(r io.Reader) (n int, err error) {
    size := msg.fixedHeader.RemainLength()
    size = size - int64(msg.varHeader.size)

    if size <= PayloadBufSize {
        buf := make([]byte, size)
        n, err = r.Read(buf)
        if err != nil {
            return n, err
        }
        msg.payload = buf
    } else {
        buf := bytes.NewBuffer(make([]byte, PayloadBufSize))
        x, err := util.CopyN(buf, r, size)
        if err != nil {
            return int(n), err
        }
        n = int(x)
        msg.payload = buf.Bytes()
    }

    return
}

func (msg *SubAckMessage) WritePayload(w io.Writer) (int, error) {
    return w.Write(msg.payload)
}

func (msg *SubAckMessage) Valid() bool {
    return true
}

func (msg *SubAckMessage) SetPacketIdentifier(v uint16) {
    msg.varHeader.PacketIdentifier = v
}

func (msg *SubAckMessage) GetPacketIdentifier() uint16 {
    return msg.varHeader.PacketIdentifier
}

func (m *SubAckMessage) SetPayload(v []byte) {
    m.payload = v
}

func (m *SubAckMessage) GetPayload() []byte {
    return m.payload
}

//UTF-8编码的字符串，表示此次响应相关的原因。
// 此原因字符串（Reason String）是为诊断而设计的可读字符串，不应该被客户端所解析。
func (m *SubAckMessage) SetReasonString(v string) {
    p := &packet.PropReasonString{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//UTF-8编码的字符串，表示此次响应相关的原因。
// 此原因字符串（Reason String）是为诊断而设计的可读字符串，不应该被客户端所解析。
func (m *SubAckMessage) GetReasonString() (string, bool) {
    p := packet.FindPropValue(packet.ReasonString, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropReasonString).V.String(), true
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *SubAckMessage) SetUserProperty(props map[string]string) {
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
func (m *SubAckMessage) GetUserProperty() (map[string]string, bool) {
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

func (v *SubAckVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("TopicName: %d\nprops:\n%s",
        v.PacketIdentifier, builder.String())
}

func (v *SubAckMessage) String() string {
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%s\npayload:\n%v\n",
        v.fixedHeader, v.varHeader.String(), v.payload)
}
