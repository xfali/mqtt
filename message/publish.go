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

const (
    PayloadBufSize = 32 * 1024
)

type PublishVarHeader struct {
    //主题名（Topic Name）用于识别有效载荷数据应该被发布到哪一个信息通道。
    TopicName packet.String
    //只有当QoS等级是1或2时，报文标识符（Packet Identifier）字段才能出现在PUBLISH报文中。
    PacketIdentifier uint16
    //PUBLISH 属性 PUBLISH Properties
    props []packet.Property
}

type PublishMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   PublishVarHeader
    payload     []byte
}

func NewPublishMessage() *PublishMessage {
    ret := &PublishMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypePUBLISH, packet.PktFlagPUBLISH, 0),
    }
    return ret
}

func (m *PublishMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *PublishMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

func (msg *PublishMessage) ReadVariableHeader(r io.Reader) (int, error) {
    s, n, err := packet.ParseString(r)
    if err != nil {
        return n, err
    }
    msg.varHeader.TopicName = *s

    if msg.HavePacketIdentifier() {
        buf := make([]byte, 2)
        n2, err2 := r.Read(buf)
        n += n2
        if err2 != nil {
            return n, err2
        }

        msg.varHeader.PacketIdentifier = uint16(buf[0]<<8 | buf[1])
    }

    props, n3, err3 := packet.ReadProperties(r)
    n += n3
    if err3 != nil {
        return n, err3
    }

    msg.varHeader.props = props

    return n, nil
}

func (msg *PublishMessage) WriteVariableHeader(w io.Writer) (int, error) {
    n, err := packet.WriteString(w, msg.varHeader.TopicName)
    if err != nil {
        return n, err
    }

    if msg.HavePacketIdentifier() {
        n2, err2 := w.Write([]byte{
            byte(msg.varHeader.PacketIdentifier >> 8),
            byte(msg.varHeader.PacketIdentifier & 0xFF),
        })
        n += n2
        if err2 != nil {
            return n, err2
        }
    }

    n3, err3 := packet.WriteProperties(w, msg.varHeader.props)
    return n + n3, err3
}

func (msg *PublishMessage) ReadPayload(r io.Reader) (n int, err error) {
    size := msg.fixedHeader.RemainLength()
    w := new(util.CountWriter)
    msg.WriteVariableHeader(w)
    size = size - w.Count()

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

func (msg *PublishMessage) WritePayload(w io.Writer) (int, error) {
    return w.Write(msg.payload)
}

//如果DUP标志被设置为0，表示这是客户端或服务端第一次请求发送这个PUBLISH报文。
//如果DUP标志被设置为1，表示这可能是一个早前报文请求的重发。
func (msg *PublishMessage) SetDup(v bool) {
    msg.fixedHeader.SetDup(v)
}

//QoS值	Bit2	Bit1	说明
//0	    0	    0	最多分发一次
//1	    0	    1	至少分发一次
//2	    1	    0	只分发一次
//-     1	    1	保留位
func (msg *PublishMessage) SetQos(v byte) {
    msg.fixedHeader.SetQos(v)
}

//只有当QoS等级是1或2时，报文标识符（Packet Identifier）字段才能出现在PUBLISH报文中
func (msg *PublishMessage) HavePacketIdentifier() bool {
    _, qos, _ := msg.fixedHeader.PubFlag()
    return qos > 0
}

//如果客户端发给服务端的PUBLISH报文的保留（Retain）标志被设置为1
//如果客户端发给服务端的PUBLISH报文的保留标志位为0，服务器不能把此消息存储为保留消息，也不能丢弃或替换任何已存在的保留消息
func (msg *PublishMessage) SetRetain(v bool) {
    msg.fixedHeader.SetRetain(v)
}

//主题名（Topic Name）用于识别有效载荷数据应该被发布到哪一个信息通道。
func (msg *PublishMessage) SetTopicName(v string) {
    s, err := packet.FromString(v)
    if err == nil {
        msg.varHeader.TopicName = s
    }
}

func (msg *PublishMessage) GetTopicName() string {
    return msg.varHeader.TopicName.String()
}

//只有当QoS等级是1或2时，报文标识符（Packet Identifier）字段才能出现在PUBLISH报文中。
func (msg *PublishMessage) SetPacketIdentifier(v uint16) {
    msg.varHeader.PacketIdentifier = v
}

func (msg *PublishMessage) GetPacketIdentifier() uint16 {
    return msg.varHeader.PacketIdentifier
}

//单字节的载荷格式指示值，可以是：
//0 (0x00)，说明载荷是未指定格式的字节，相当于没有发送载荷格式指示。
//1 (0x01)，说明载荷是UTF-8编码的字符数据。载荷中的UTF-8数据必须是按照Unicode [Unicode]的规范和RFC 3629 [RFC3629]的重申进行编码。
func (m *PublishMessage) SetPayloadFormatIndicator(v byte) {
    p := &packet.PropPayloadFormatIndicator{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//单字节的载荷格式指示值，可以是：
//0 (0x00)，说明载荷是未指定格式的字节，相当于没有发送载荷格式指示。
//1 (0x01)，说明载荷是UTF-8编码的字符数据。载荷中的UTF-8数据必须是按照Unicode [Unicode]的规范和RFC 3629 [RFC3629]的重申进行编码。
func (m *PublishMessage) GetPayloadFormatIndicator() (byte, bool) {
    p := packet.FindPropValue(packet.PayloadFormatIndicator, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropPayloadFormatIndicator).V, true
}

//如果消息过期间隔存在，四字节整数表示以秒为单位的应用消息（Application Message）生命周期。
//如果消息过期间隔（Message Expiry Interval）已过期，服务端还没开始向匹配的订阅者交付该消息，则服务端必须删除该订阅者的消息副本
func (m *PublishMessage) SetMessageExpiryInterval(v uint32) {
    p := &packet.PropMessageExpiryInterval{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//如果消息过期间隔存在，四字节整数表示以秒为单位的应用消息（Application Message）生命周期。
//如果消息过期间隔（Message Expiry Interval）已过期，服务端还没开始向匹配的订阅者交付该消息，则服务端必须删除该订阅者的消息副本
func (m *PublishMessage) GetMessageExpiryInterval() (uint32, bool) {
    p := packet.FindPropValue(packet.MessageExpiryInterval, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropMessageExpiryInterval).V, true
}

//包含多个主题别名值将造成协议错误（Protocol Error）。
//主题别名是一个整数，用来代替主题名对主题进行识别。
func (m *PublishMessage) SetTopicAlias(v uint16) {
    p := &packet.PropTopicAlias{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//包含多个主题别名值将造成协议错误（Protocol Error）。
//主题别名是一个整数，用来代替主题名对主题进行识别。
func (m *PublishMessage) GetTopicAlias() (uint16, bool) {
    p := packet.FindPropValue(packet.TopicAlias, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropTopicAlias).V, true
}

//用作响应消息的主题名。
//包含多个响应主题将造成协议错误（Protocol Error）。
func (m *PublishMessage) SetResponseTopic(v string) {
    p := &packet.PropResponseTopic{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//用作响应消息的主题名。
//包含多个响应主题将造成协议错误（Protocol Error）。
func (m *PublishMessage) GetResponseTopic() (string, bool) {
    p := packet.FindPropValue(packet.ResponseTopic, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropResponseTopic).V.String(), true
}

//对比数据被请求消息发送端在收到响应消息时用来标识相应的请求。包含多个对比数据将造成协议错误（Protocol Error）。
//如果没有设置对比数据，则请求方（Requester）不需要任何对比数据。
func (m *PublishMessage) SetCorrelationData(v []byte) {
    p := &packet.PropCorrelationData{}
    s, err := packet.FromString(string(v))
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//对比数据被请求消息发送端在收到响应消息时用来标识相应的请求。包含多个对比数据将造成协议错误（Protocol Error）。
//如果没有设置对比数据，则请求方（Requester）不需要任何对比数据。
func (m *PublishMessage) GetCorrelationData() ([]byte, bool) {
    p := packet.FindPropValue(packet.CorrelationData, m.varHeader.props)
    if p == nil {
        return nil, false
    }
    return []byte(p.(*packet.PropCorrelationData).V.String()), true
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *PublishMessage) SetUserProperty(props map[string]string) {
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
func (m *PublishMessage) GetUserProperty() (map[string]string, bool) {
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

//订阅标识符取值范围从1到268,435,455。订阅标识符的值为0将造成协议错误。
//如果某条发布消息匹配了多个订阅，则将包含多个订阅标识符。这种情况下他们的顺序并不重要。
func (m *PublishMessage) SetSubscriptionIdentifier(v uint64) {
    p := &packet.PropSubscriptionIdentifier{}
    p.V.InitFromUInt64(v)
    m.varHeader.props = append(m.varHeader.props, p)
}

//订阅标识符取值范围从1到268,435,455。订阅标识符的值为0将造成协议错误。
//如果某条发布消息匹配了多个订阅，则将包含多个订阅标识符。这种情况下他们的顺序并不重要。
func (m *PublishMessage) GetSubscriptionIdentifier() (int64, bool) {
    p := packet.FindPropValue(packet.SubscriptionIdentifier, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropSubscriptionIdentifier).V.ToInt(), true
}

//用来描述应用消息的内容。
//包含多个内容类型将造成协议错误（Protocol Error）。
//内容类型的值由发送应用程序和接收应用程序确定。
func (m *PublishMessage) SetContentType(v string) {
    p := &packet.PropContentType{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//用来描述应用消息的内容。
//包含多个内容类型将造成协议错误（Protocol Error）。
//内容类型的值由发送应用程序和接收应用程序确定。
func (m *PublishMessage) GetContentType() (string, bool) {
    p := packet.FindPropValue(packet.ContentType, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropContentType).V.String(), true
}

func (m *PublishMessage) SetPayload(v []byte) {
    m.payload = v
}

func (m *PublishMessage) GetPayload() []byte {
    return m.payload
}

func (m *PublishMessage) Valid() bool {
    return true
}

func (v *PublishVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("TopicName: %s ReasonCode: %d \nprops:\n%s",
        v.TopicName.String(), v.PacketIdentifier, builder.String())
}

func (v *PublishMessage) String() string {
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%s payload:\n %s\n",
        v.fixedHeader, v.varHeader.String(), string(v.payload))
}
