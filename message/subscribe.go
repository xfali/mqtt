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

type SubscribeFilter struct {
    Filter string
    Opt    byte
}

type SubscribeVarHeader struct {
    //报文标识符（Packet Identifier）字段才能出现在PUBLISH报文中。
    PacketIdentifier uint16
    //PUBLISH 属性 PUBLISH Properties
    props []packet.Property

    //VarHeader size
    size int
}

type SubscribeMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   SubscribeVarHeader
    payload     []SubscribeFilter
}

func NewSubscribeMessage() *SubscribeMessage {
    ret := &SubscribeMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypeSUBSCRIBE, packet.PktFlagSUBSCRIBE, 0),
    }
    return ret
}

func (m *SubscribeMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *SubscribeMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

func (msg *SubscribeMessage) ReadVariableHeader(r io.Reader) (int, error) {
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

func (msg *SubscribeMessage) WriteVariableHeader(w io.Writer) (int, error) {
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

func (msg *SubscribeMessage) ReadPayload(r io.Reader) (int, error) {
    size := int(msg.fixedHeader.RemainLength())
    size = size - msg.varHeader.size

    n := 0
    var filters []SubscribeFilter
    buf := make([]byte, 1)
    for n < size {
        s, rd, err := packet.ParseString(r)
        n += rd
        if err != nil {
            return n, err
        }
        rd, err = r.Read(buf)
        n += rd
        if err != nil {
            return n, err
        }
        filters = append(filters, SubscribeFilter{Filter: s.String(), Opt: buf[0]})
    }

    if n > size {
        return n, errcode.ProtocolError
    }

    msg.payload = filters

    return n, nil
}

func (msg *SubscribeMessage) WritePayload(w io.Writer) (int, error) {
    if len(msg.payload) == 0 {
        return 0, errcode.ProtocolError
    }

    n := 0
    for _, v := range msg.payload {
        s, errS := packet.FromString(v.Filter)
        if errS != nil {
            return n, errS
        }
        wt, err := packet.WriteString(w, s)
        n += wt
        if err != nil {
            return n, err
        }
        wt, err = w.Write([]byte{v.Opt})
        n += wt
        if err != nil {
            return n, err
        }
    }
    return n, nil
}

func (msg *SubscribeMessage) Valid() bool {
    return true
}

//订阅标识符取值范围从1到268,435,455。
//订阅标识符的值为0或包含多个订阅标识符将造成协议错误（Protocol Error）。
func (m *SubscribeMessage) SetSubscriptionIdentifier(v uint64) {
    p := &packet.PropSubscriptionIdentifier{}
    p.V.InitFromUInt64(v)
    m.varHeader.props = append(m.varHeader.props, p)
}

//订阅标识符取值范围从1到268,435,455。
//订阅标识符的值为0或包含多个订阅标识符将造成协议错误（Protocol Error）。
func (m *SubscribeMessage) GetSubscriptionIdentifier() (uint64, bool) {
    p := packet.FindPropValue(packet.SubscriptionIdentifier, m.varHeader.props)
    return p.(*packet.PropSubscriptionIdentifier).V.ToUint(), true
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *SubscribeMessage) SetUserProperty(props map[string]string) {
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
func (m *SubscribeMessage) GetUserProperty() (map[string]string, bool) {
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

func (msg *SubscribeMessage) SetPacketIdentifier(v uint16) {
    msg.varHeader.PacketIdentifier = v
}

func (msg *SubscribeMessage) GetPacketIdentifier() uint16 {
    return msg.varHeader.PacketIdentifier
}

func (msg *SubscribeMessage) SetPayload(payload []SubscribeFilter) {
    msg.payload = payload
}

func (msg *SubscribeMessage) GetPayload() []SubscribeFilter {
    return msg.payload
}

func (v *SubscribeVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("PacketIdentifier: %d \nprops:\n%s",
        v.PacketIdentifier, builder.String())
}

func (v *SubscribeMessage) String() string {
    builder := strings.Builder{}
    for _, f := range v.payload {
        builder.WriteString(fmt.Sprintf("filter: %s, opt %d\n", f.Filter, f.Opt))
    }
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%spayload:\n%s \n",
        v.fixedHeader, v.varHeader.String(), builder.String())
}
