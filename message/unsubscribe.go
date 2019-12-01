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

type UnsubscribeVarHeader struct {
    //报文标识符（Packet Identifier）。
    PacketIdentifier uint16

    //PUBLISH 属性 PUBLISH Properties
    props []packet.Property

    size int
}

type UnsubscribeMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   UnsubscribeVarHeader
    payload     []string
}

func NewUnsubscribeMessage() *UnsubscribeMessage {
    ret := &UnsubscribeMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypeUNSUBSCRIBE, packet.PktFlagUNSUBSCRIBE, 0),
    }
    return ret
}

func (m *UnsubscribeMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *UnsubscribeMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

func (msg *UnsubscribeMessage) ReadVariableHeader(r io.Reader) (int, error) {
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

func (msg *UnsubscribeMessage) WriteVariableHeader(w io.Writer) (int, error) {
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

func (msg *UnsubscribeMessage) ReadPayload(r io.Reader) (int, error) {
    size := int(msg.fixedHeader.RemainLength())
    size = size - msg.varHeader.size

    n := 0
    for n < size {
        s, rd, err := packet.ParseString(r)
        n += rd
        if err != nil {
            return n, err
        }
        msg.payload = append( msg.payload, s.String())
    }

    if n > size {
        return n, errcode.ProtocolError
    }

    return n, nil
}

func (msg *UnsubscribeMessage) WritePayload(w io.Writer) (int, error) {
    if len(msg.payload) == 0 {
        return 0, errcode.ProtocolError
    }

    n := 0
    for _, v := range msg.payload {
        s, errS := packet.FromString(v)
        if errS != nil {
            return n, errS
        }
        wt, err := packet.WriteString(w, s)
        n += wt
        if err != nil {
            return n, err
        }
    }
    return n, nil
}

func (msg *UnsubscribeMessage) Valid() bool {
    return true
}

func (msg *UnsubscribeMessage) SetPacketIdentifier(v uint16) {
    msg.varHeader.PacketIdentifier = v
}

func (msg *UnsubscribeMessage) GetPacketIdentifier() uint16 {
    return msg.varHeader.PacketIdentifier
}

//必须包含至少一个主题过滤器 [MQTT-3.10.3-2]。
//不包含有效载荷的UNSUBSCRIBE报文将造成协议错误（Protocol Error）。
func (m *UnsubscribeMessage) SetPayload(filters []string) {
    m.payload = filters
}

func (m *UnsubscribeMessage) GetPayload() []string {
    return m.payload
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *UnsubscribeMessage) SetUserProperty(props map[string]string) {
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
func (m *UnsubscribeMessage) GetUserProperty() (map[string]string, bool) {
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

func (v *UnsubscribeVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("TopicName: %d\nprops:\n%s",
        v.PacketIdentifier, builder.String())
}

func (v *UnsubscribeMessage) String() string {
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%s\npayload:\n%v\n",
        v.fixedHeader, v.varHeader.String(), v.payload)
}
