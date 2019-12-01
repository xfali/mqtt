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

type ConnectVarHeader struct {
    ProtocolName    packet.String
    ProtocolVersion byte
    Flag            byte
    KeepAlive       uint16
    props           []packet.Property
}

type ConnectPayload struct {
    ClientId    packet.String
    WillProps   []packet.Property
    WillTopic   packet.String
    WillPayload packet.Bytes
    Username    packet.String
    Password    packet.Bytes
}

//客户端到服务端的网络连接建立后，客户端发送给服务端的第一个报文必须是CONNECT报文 [MQTT-3.1.0-1]。
//在一个网络连接上，客户端只能发送一次CONNECT报文。服务端必须将客户端发送的第二个CONNECT报文当作协议违规处理并断开客户端的连接
type ConnectMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   ConnectVarHeader
    payload     ConnectPayload
}

func (msg *ConnectMessage) ReadVariableHeader(r io.Reader) (int, error) {
    s, n, err := packet.ParseString(r)
    if err != nil {
        return n, err
    }
    if n != 6 || packet.MqttProtocolName != s.String() {
        return n, errcode.ProtocolNameError
    }
    msg.varHeader.ProtocolName = *s
    buf := make([]byte, 4)
    n2, err2 := r.Read(buf)
    if err2 != nil {
        return n + n2, err2
    }
    msg.varHeader.ProtocolVersion = buf[0]
    msg.varHeader.Flag = buf[1]
    msg.varHeader.KeepAlive = uint16(buf[2]<<8) | uint16(buf[3])

    props, n3, err3 := packet.ReadProperties(r)
    if err3 != nil {
        return n + n2 + n3, err3
    }
    msg.varHeader.props = props

    return n + n2 + n3, nil
}

func (msg *ConnectMessage) WriteVariableHeader(w io.Writer) (int, error) {
    n, err := packet.WriteString(w, msg.varHeader.ProtocolName)
    if err != nil {
        return n, err
    }
    n2, err2 := w.Write([]byte{
        msg.varHeader.ProtocolVersion,
        msg.varHeader.Flag,
        byte(msg.varHeader.KeepAlive >> 8),
        byte(msg.varHeader.KeepAlive & 0xFF),
    })

    if err2 != nil {
        return n + n2, err2
    }

    n3, err3 := packet.WriteProperties(w, msg.varHeader.props)
    return n + n2 + n3, err3
}

func (msg *ConnectMessage) ReadPayload(r io.Reader) (int, error) {
    s, n, err := packet.ParseString(r)
    if err != nil {
        return n, err
    }
    msg.payload.ClientId = *s
    if msg.IsWillEnable() {
        props, n2, err2 := packet.ReadProperties(r)
        n += n2
        if err2 != nil {
            return n, err2
        }
        msg.payload.WillProps = props

        s, n3, err3 := packet.ParseString(r)
        n += n3
        if err3 != nil {
            return n, err3
        }
        msg.payload.WillTopic = *s

        b, n4, err4 := packet.ParseBytes(r)
        n += n4
        if err4 != nil {
            return n, err4
        }
        msg.payload.WillPayload = *b
    }

    if msg.haveUsername() {
        s, n5, err5 := packet.ParseString(r)
        n += n5
        if err5 != nil {
            return n, err5
        }
        msg.payload.Username = *s
    }

    if msg.havePassword() {
        s, n6, err6 := packet.ParseBytes(r)
        n += n6
        if err6 != nil {
            return n, err6
        }
        msg.payload.Password = *s
    }

    return n, err
}

func (msg *ConnectMessage) WritePayload(w io.Writer) (int, error) {
    n, err := packet.WriteString(w, msg.payload.ClientId)
    if err != nil {
        return n, err
    }

    if msg.IsWillEnable() {
        n2, err2 := packet.WriteProperties(w, msg.payload.WillProps)
        n += n2
        if err2 != nil {
            return n, err2
        }

        n3, err3 := packet.WriteString(w, msg.payload.WillTopic)
        n += n3
        if err3 != nil {
            return n, err3
        }

        n4, err4 := packet.WriteBytes(w, msg.payload.WillPayload)
        n += n4
        if err4 != nil {
            return n, err4
        }
    }

    if msg.haveUsername() {
        n5, err5 := packet.WriteString(w, msg.payload.Username)
        n += n5
        if err5 != nil {
            return n, err5
        }
    }

    if msg.havePassword() {
        n6, err6 := packet.WriteBytes(w, msg.payload.Password)
        n += n6
        if err6 != nil {
            return n, err6
        }
    }

    return n, err
}

func NewConnectMessage() *ConnectMessage {
    ret := &ConnectMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypeCONNECT, packet.PktFlagCONNECT, 0),
    }
    ret.varHeader.ProtocolName = packet.MqttProtocolNameString
    ret.varHeader.ProtocolVersion = packet.MqttProtocolVersion
    return ret
}

func (m *ConnectMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *ConnectMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

//如果遗嘱标志（Will Flag）设置为0，遗嘱服务质量（Will QoS）必须也设置为0（0x00）
//如果遗嘱标志设置为1，遗嘱服务质量可以被设置为0（0x00），1（0x01）或2（0x02）
//设置为3（0x03）的报文是无效报文。
func (m *ConnectMessage) SetWillQos(v byte) {
    m.varHeader.Flag |= (v & 0x3) << 3
}

//表示协议修订级别
func (m *ConnectMessage) SetVersion(v byte) {
    m.varHeader.ProtocolVersion = v
}

//二进制位表明此次连接是一个新的会话还是一个已存在的会话的延续。
func (m *ConnectMessage) SetCleanStart(v bool) {
    if v {
        m.varHeader.Flag |= 1 << 1
    } else {
        m.varHeader.Flag &= 0xFF & ^(1 << 1)
    }
}

func (m *ConnectMessage) IsCleanStart() bool{
    return (m.varHeader.Flag) >> 1 & 0xFE != 0
}

//如果遗嘱标志（Will Flag）被设置为1，表示遗嘱消息必须已存储在服务端与此客户标识符相关的会话中
func (m *ConnectMessage) SetWillEnable(v bool) {
    if v {
        m.varHeader.Flag |= 1 << 2
    } else {
        m.varHeader.Flag &= 0xFF & ^(1 << 2)
    }
}

func (m *ConnectMessage) IsWillEnable() bool {
    return int(m.varHeader.Flag) & ^(1 << 2) != 0
}

func (m *ConnectMessage) SetClientId(v string) {
    m.payload.ClientId.Reset(v)
}

func (m *ConnectMessage) GetClientId() string {
    return m.payload.ClientId.String()
}

//保持连接（Keep Alive）使用双字节整数来表示以秒为单位的时间间隔。
//它是指在客户端传输完成一个MQTT控制报文的时刻到发送下一个报文的时刻，两者之间允许空闲的最大时间间隔。
func (m *ConnectMessage) SetKeepAlive(v uint16) {
    m.varHeader.KeepAlive = v
}

//此位指定遗嘱消息（Will Message）在发布时是否会被保留。
func (m *ConnectMessage) SetWillRetain(v bool) {
    m.varHeader.Flag |= 1 << 5
}

func (m *ConnectMessage) SetWillTopic(v string) {
    m.payload.WillTopic.Reset(v)
}

func (m *ConnectMessage) SetWillPayload(v []byte) {
    m.payload.WillPayload.Reset(v)
}

//如果用户名标志（User Name Flag）被设置为0，有效载荷中不能包含用户名字段。
//如果用户名标志被设置为1，有效载荷中必须包含用户名字段。
func (m *ConnectMessage) SetUsername(v string) {
    m.varHeader.Flag |= 1 << 7
    m.payload.Username.Reset(v)
}

func (m *ConnectMessage) GetUsername() string {
    return m.payload.Username.String()
}

func (m *ConnectMessage) haveUsername() bool {
    return int(m.varHeader.Flag) & ^(1 << 7) != 0
}

//如果密码标志（Password Flag）被设置为0，有效载荷中不能包含密码字段 [MQTT-3.1.2-18]。
//如果密码标志被设置为1，有效载荷中必须包含密码字段 [MQTT-3.1.2-19]。
func (m *ConnectMessage) SetPassword(v []byte) {
    m.varHeader.Flag |= 1 << 6
    m.payload.Password.Reset(v)
}

func (m *ConnectMessage) GetPassword() []byte {
    return m.payload.Password.Get()
}

func (m *ConnectMessage) havePassword() bool {
    return int(m.varHeader.Flag) & ^(1 << 6) != 0
}
//以秒为单位的会话过期间隔（Session Expiry Interval）。包含多个会话过期间隔（Session Expiry Interval）将造成协议错误（Protocol Error）。
//如果会话过期间隔（Session Expiry Interval）值未指定，则使用0。如果设置为0或者未指定，会话将在网络连接（Network Connection）关闭时结束。
//如果会话过期间隔（Session Expiry Interval）为0xFFFFFFFF (UINT_MAX)，则会话永不过期。
//如果网络连接关闭时会话过期间隔（Session Expiry Interval）大于0，则客户端与服务端必须存储会话状态
func (m *ConnectMessage) SetSessionExpiryInterval(v uint32) {
    p := &packet.PropSessionExpiryInterval{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//以秒为单位的会话过期间隔
func (m *ConnectMessage) GetSessionExpiryInterval() (uint32, bool) {
    p := packet.FindPropValue(packet.SessionExpiryInterval, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropSessionExpiryInterval).V, true
}

//客户端使用此值限制客户端愿意同时处理的QoS为1和QoS为2的发布消息最大数量。没有机制可以限制服务端试图发送的QoS为0的发布消息。
//接收最大值只将被应用在当前网络连接。如果没有设置最大接收值，将使用默认值65535。
func (m *ConnectMessage) SetReceiveMaximum(v uint16) {
    p := &packet.PropReceiveMaximum{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

func (m *ConnectMessage) GetReceiveMaximum() (uint16, bool) {
    p := packet.FindPropValue(packet.ReceiveMaximum, m.varHeader.props)
    if p == nil {
        return 65535, false
    }
    return p.(*packet.PropReceiveMaximum).V, true
}

//服务端愿意接收的最大报文长度（Maximum Packet Size）。
//如果没有设置最大报文长度，则按照协议由固定报头中的剩余长度可编码最大值和协议报头对数据包的大小做限制。
func (m *ConnectMessage) SetMaximumPacketSize(v uint32) {
    p := &packet.PropMaximumPacketSize{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//服务端愿意接收的最大报文长度（Maximum Packet Size）。
//如果没有设置最大报文长度，则按照协议由固定报头中的剩余长度可编码最大值和协议报头对数据包的大小做限制。
func (m *ConnectMessage) GetMaximumPacketSize() (uint32, bool) {
    p := packet.FindPropValue(packet.MaximumPacketSize, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropMaximumPacketSize).V, true
}

//双字节整数表示的主题别名最大值（Topic Alias Maximum）。
// 包含多个主题别名最大值（Topic Alias Maximum）将造成协议错误（Protocol Error）。
// 没有设置主题别名最大值属性的情况下，主题别名最大值默认为零。
func (m *ConnectMessage) SetTopicAliasMaximum(v uint16) {
    p := &packet.PropTopicAliasMaximum{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//双字节整数表示的主题别名最大值（Topic Alias Maximum）。
// 包含多个主题别名最大值（Topic Alias Maximum）将造成协议错误（Protocol Error）。
// 没有设置主题别名最大值属性的情况下，主题别名最大值默认为零。
func (m *ConnectMessage) GetTopicAliasMaximum() (uint16, bool) {
    p := packet.FindPropValue(packet.TopicAliasMaximum, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropTopicAliasMaximum).V, true
}

func (m *ConnectMessage) SetRequestResponseInformation(v byte) {
    p := &packet.PropRequestResponseInformation{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

func (m *ConnectMessage) SetRequestProblemInformation(v byte) {
    p := &packet.PropRequestProblemInformation{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *ConnectMessage) SetUserProperty(props map[string]string) {
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
func (m *ConnectMessage) GetUserProperty() (map[string]string, bool) {
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

//以UTF-8编码的字符串，包含了认证方法（Authentication Method）名。
//包含多个认证方法将造成协议错误（Protocol Error）。
func (m *ConnectMessage) SetAuthenticationMethod(v string) {
    p := &packet.PropAuthenticationMethod{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//以UTF-8编码的字符串，包含了认证方法（Authentication Method）名。
//包含多个认证方法将造成协议错误（Protocol Error）。
func (m *ConnectMessage) GetAuthenticationMethod() (string, bool) {
    p := packet.FindPropValue(packet.AuthenticationMethod, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropAuthenticationMethod).V.String(), true
}

//包含认证数据（Authentication Data）的二进制数据。此数据的内容由认证方法和已交换的认证数据状态定义。
//包含多个认证数据将造成协议错误（Protocol Error）。
func (m *ConnectMessage) SetAuthenticationData(v []byte) {
    p := &packet.PropAuthenticationData{}
    s, err := packet.FromString(string(v))
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//包含认证数据（Authentication Data）的二进制数据。此数据的内容由认证方法和已交换的认证数据状态定义。
//包含多个认证数据将造成协议错误（Protocol Error）。
func (m *ConnectMessage) GetAuthenticationData() ([]byte, bool) {
    p := packet.FindPropValue(packet.AuthenticationData, m.varHeader.props)
    if p == nil {
        return nil, false
    }
    return []byte(p.(*packet.PropAuthenticationData).V.String()), true
}

func (m *ConnectMessage) GetRequestResponseInformation() (byte, bool) {
    p := packet.FindPropValue(packet.RequestResponseInformation, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropRequestResponseInformation).V, true
}

func (m *ConnectMessage) GetRequestProblemInformation() (byte, bool) {
    p := packet.FindPropValue(packet.RequestProblemInformation, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropRequestProblemInformation).V, true
}

func (m *ConnectMessage) SetWillDelayInterval(v uint32) {
    p := &packet.PropWillDelayInterval{}
    p.V = v
    m.payload.WillProps = append(m.payload.WillProps, p)
}

func (m *ConnectMessage) GetWillDelayInterval() (uint32, bool) {
    p := packet.FindPropValue(packet.WillDelayInterval, m.payload.WillProps)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropWillDelayInterval).V, true
}

func (m *ConnectMessage) SetPayloadFormatIndicator(v byte) {
    p := &packet.PropPayloadFormatIndicator{}
    p.V = v
    m.payload.WillProps = append(m.payload.WillProps, p)
}

func (m *ConnectMessage) GetPayloadFormatIndicator() (byte, bool) {
    p := packet.FindPropValue(packet.PayloadFormatIndicator, m.payload.WillProps)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropPayloadFormatIndicator).V, true
}

func (m *ConnectMessage) SetMessageExpiryInterval(v uint32) {
    p := &packet.PropMessageExpiryInterval{}
    p.V = v
    m.payload.WillProps = append(m.payload.WillProps, p)
}

func (m *ConnectMessage) GetMessageExpiryInterval() (uint32, bool) {
    p := packet.FindPropValue(packet.MessageExpiryInterval, m.payload.WillProps)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropMessageExpiryInterval).V, true
}

//用来描述应用消息的内容。
//包含多个内容类型将造成协议错误（Protocol Error）。
//内容类型的值由发送应用程序和接收应用程序确定。
func (m *ConnectMessage) SetContentType(v string) {
    p := &packet.PropContentType{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.payload.WillProps = append(m.payload.WillProps, p)
    }
}

//用来描述应用消息的内容。
//包含多个内容类型将造成协议错误（Protocol Error）。
//内容类型的值由发送应用程序和接收应用程序确定。
func (m *ConnectMessage) GetContentType() (string, bool) {
    p := packet.FindPropValue(packet.ContentType, m.payload.WillProps)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropContentType).V.String(), true
}

func (m *ConnectMessage) SetResponseTopic(v string) {
    p := &packet.PropResponseTopic{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.payload.WillProps = append(m.payload.WillProps, p)
    }
}

func (m *ConnectMessage) GetResponseTopic() (string, bool) {
    p := packet.FindPropValue(packet.ResponseTopic, m.payload.WillProps)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropResponseTopic).V.String(), true
}

//对比数据被请求消息发送端在收到响应消息时用来标识相应的请求。包含多个对比数据将造成协议错误（Protocol Error）。
//如果没有设置对比数据，则请求方（Requester）不需要任何对比数据。
func (m *ConnectMessage) SetCorrelationData(v []byte) {
    p := &packet.PropCorrelationData{}
    s, err := packet.FromString(string(v))
    if err == nil {
        p.V = s
        m.payload.WillProps = append(m.payload.WillProps, p)
    }
}

//对比数据被请求消息发送端在收到响应消息时用来标识相应的请求。包含多个对比数据将造成协议错误（Protocol Error）。
//如果没有设置对比数据，则请求方（Requester）不需要任何对比数据。
func (m *ConnectMessage) GetCorrelationData() ([]byte, bool) {
    p := packet.FindPropValue(packet.CorrelationData, m.payload.WillProps)
    if p == nil {
        return nil, false
    }
    return []byte(p.(*packet.PropCorrelationData).V.String()), true
}

func (m *ConnectMessage) SetPayloadUserProperty(props map[string]string) {
    for k, v := range props {
        p := &packet.PropUserProperty{}
        pair, err := packet.NewStringPair(k, v)
        if err == nil {
            p.V = pair
            m.payload.WillProps = append(m.payload.WillProps, p)
        }
    }
}

func (m *ConnectMessage) GetPayloadUserProperty() (map[string]string, bool) {
    ret := map[string]string{}
    packet.FindPropValues(packet.UserProperty, m.payload.WillProps, func(property packet.Property) bool {
        if property != nil {
            p := property.(*packet.PropUserProperty)
            ret[p.V[0].String()] = p.V[1].String()
        }
        return false
    })
    return ret, len(ret) > 0
}

func (m *ConnectMessage) Valid() bool {
    return true
}

func (v *ConnectVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("protocal name: %s version: %d flag: %b keepAlive: %d \nprops:\n%s",
        v.ProtocolName.String(), v.ProtocolVersion, v.Flag, v.KeepAlive, builder.String())
}

func (v *ConnectPayload) String() string {
    builder := strings.Builder{}
    for i := range v.WillProps {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.WillProps[i]))
    }
    return fmt.Sprintf("ClientId: %s WillTopic: %s WillPayload: %v Username: %s Password: %v \nprops:\n%s",
        v.ClientId.String(), v.WillTopic.String(), v.WillPayload, v.Username.String(), v.Password, builder.String())
}

func (v *ConnectMessage) String() string {
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%spayload:\n%s \n",
        v.fixedHeader, v.varHeader.String(), v.payload.String())
}
