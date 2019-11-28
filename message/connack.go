// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "fmt"
    "io"
    "mqtt/packet"
    "mqtt/util"
    "strings"
)

type ConnackVarHeader struct {
    //连接确认标志 Connect Acknowledge Flags
    AckFlag byte
    //连接原因码 Connect Reason Code
    ReasonCode byte
    //CONNACK属性 CONNACK Properties
    props []packet.Property
}

type ConnackMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   ConnackVarHeader
}

func NewConnackMessage() *ConnackMessage {
    ret := &ConnackMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypeCONNACK, packet.PktFlagCONNACK, 0),
    }
    return ret
}

func (m *ConnackMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *ConnackMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

func (msg *ConnackMessage) ReadVariableHeader(r io.Reader) (int, error) {
    buf := make([]byte, 2)
    n, err := r.Read(buf)
    if err != nil {
        return n, err
    }

    msg.varHeader.AckFlag = buf[0]
    msg.varHeader.ReasonCode = buf[1]

    props, n2, err2 := packet.ReadProperties(r)
    n += n2
    if err2 != nil {
        return n, err2
    }
    msg.varHeader.props = props

    return n, nil
}

func (msg *ConnackMessage) WriteVariableHeader(w io.Writer) (int, error) {
    n, err := w.Write([]byte{
        msg.varHeader.AckFlag,
        msg.varHeader.ReasonCode,
    })
    if err != nil {
        return n, err
    }

    n2, err2 := packet.WriteProperties(w, msg.varHeader.props)
    return n + n2, err2
}

func (msg *ConnackMessage) ReadPayload(r io.Reader) (int, error) {
    return 0, nil
}

func (msg *ConnackMessage) WritePayload(w io.Writer) (int, error) {
    return 0, nil
}

func (m *ConnackMessage) SetAckFlag(v byte) {
    m.varHeader.AckFlag = v
}

func (m *ConnackMessage) GetAckFlag() byte{
    return m.varHeader.AckFlag
}

func (m *ConnackMessage) SetReasonCode(v byte) {
    m.varHeader.ReasonCode = v
}

func (m *ConnackMessage) GetReasonCode() byte{
    return m.varHeader.ReasonCode
}

// 会话过期间隔 Session Expiry Interval,四字节整数表示的以秒为单位的会话过期间隔
func (m *ConnackMessage) SetSessionExpiryInterval(v uint32) {
    p := &packet.PropSessionExpiryInterval{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

// 会话过期间隔 Session Expiry Interval,四字节整数表示的以秒为单位的会话过期间隔
func (m *ConnackMessage) GetSessionExpiryInterval() (uint32, bool) {
    p := packet.FindPropValue(packet.SessionExpiryInterval, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropSessionExpiryInterval).V, true
}

//双字节整数表示的最大接收值
func (m *ConnackMessage) SetReceiveMaximum(v uint16) {
    p := &packet.PropReceiveMaximum{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//双字节整数表示的最大接收值
func (m *ConnackMessage) GetReceiveMaximum() (uint16, bool) {
    p := packet.FindPropValue(packet.ReceiveMaximum, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropReceiveMaximum).V, true
}

//用一个字节表示的0或1。包含多个最大服务质量（Maximum QoS）或最大服务质量既不为0也不为1将造成协议错误。
// 如果没有设置最大服务质量，客户端可使用最大QoS为2。
func (m *ConnackMessage) SetMaximumQoS(v byte) {
    p := &packet.PropMaximumQoS{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//用一个字节表示的0或1。包含多个最大服务质量（Maximum QoS）或最大服务质量既不为0也不为1将造成协议错误。
// 如果没有设置最大服务质量，客户端可使用最大QoS为2。
func (m *ConnackMessage) GetMaximumQoS() (byte, bool) {
    p := packet.FindPropValue(packet.MaximumQoS, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropMaximumQoS).V, true
}

//一个单字节字段，用来声明服务端是否支持保留消息。值为0表示不支持保留消息，为1表示支持保留消息。
//如果没有设置保留可用字段，表示支持保留消息。包含多个保留可用字段或保留可用字段值不为0也不为1将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetRetainAvailable(v byte) {
    p := &packet.PropRetainAvailable{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//一个单字节字段，用来声明服务端是否支持保留消息。值为0表示不支持保留消息，为1表示支持保留消息。
//如果没有设置保留可用字段，表示支持保留消息。包含多个保留可用字段或保留可用字段值不为0也不为1将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetRetainAvailable() (byte, bool) {
    p := packet.FindPropValue(packet.RetainAvailable, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropRetainAvailable).V, true
}

//四字节整数表示的服务端愿意接收的最大报文长度（Maximum Packet Size）。
//如果没有设置最大报文长度，则按照协议由固定报头中的剩余长度可编码最大值和协议报头对数据包的大小做限制。
func (m *ConnackMessage) SetMaximumPacketSize(v uint32) {
    p := &packet.PropMaximumPacketSize{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//四字节整数表示的服务端愿意接收的最大报文长度（Maximum Packet Size）。
//如果没有设置最大报文长度，则按照协议由固定报头中的剩余长度可编码最大值和协议报头对数据包的大小做限制。
func (m *ConnackMessage) GetMaximumPacketSize() (uint32, bool) {
    p := packet.FindPropValue(packet.MaximumPacketSize, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropMaximumPacketSize).V, true
}

//UTF-8编码的分配客户标识符（Assigned Client Identifier）字符串。
//包含多个分配客户标识符将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetAssignedClientIdentifier(v string) {
    p := &packet.PropAssignedClientIdentifier{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//UTF-8编码的分配客户标识符（Assigned Client Identifier）字符串。
//包含多个分配客户标识符将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetAssignedClientIdentifier() (string, bool) {
    p := packet.FindPropValue(packet.AssignedClientIdentifier, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropAssignedClientIdentifier).V.String(), true
}

//双字节整数表示的主题别名最大值（Topic Alias Maximum）。
// 包含多个主题别名最大值（Topic Alias Maximum）将造成协议错误（Protocol Error）。
// 没有设置主题别名最大值属性的情况下，主题别名最大值默认为零。
func (m *ConnackMessage) SetTopicAliasMaximum(v uint16) {
    p := &packet.PropTopicAliasMaximum{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//双字节整数表示的主题别名最大值（Topic Alias Maximum）。
// 包含多个主题别名最大值（Topic Alias Maximum）将造成协议错误（Protocol Error）。
// 没有设置主题别名最大值属性的情况下，主题别名最大值默认为零。
func (m *ConnackMessage) GetTopicAliasMaximum() (uint16, bool) {
    p := packet.FindPropValue(packet.TopicAliasMaximum, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropTopicAliasMaximum).V, true
}

//UTF-8编码的字符串，表示此次响应相关的原因。
// 此原因字符串（Reason String）是为诊断而设计的可读字符串，不应该被客户端所解析。
func (m *ConnackMessage) SetReasonString(v string) {
    p := &packet.PropReasonString{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//UTF-8编码的字符串，表示此次响应相关的原因。
// 此原因字符串（Reason String）是为诊断而设计的可读字符串，不应该被客户端所解析。
func (m *ConnackMessage) GetReasonString() (string, bool) {
    p := packet.FindPropValue(packet.ReasonString, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropReasonString).V.String(), true
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *ConnackMessage) SetUserProperty(props map[string]string) {
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
func (m *ConnackMessage) GetUserProperty() (map[string]string, bool) {
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

//单字节字段，用来声明服务器是否支持通配符订阅（Wildcard Subscriptions）。值为0表示不支持通配符订阅，值为1表示支持通配符订阅。
//如果没有设置此值，则表示支持通配符订阅。
//包含多个通配符订阅可用属性，或通配符订阅可用属性值不为0也不为1将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetWildcardSubscriptionAvailable(v byte) {
    p := &packet.PropWildcardSubscriptionAvailable{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//单字节字段，用来声明服务器是否支持通配符订阅（Wildcard Subscriptions）。值为0表示不支持通配符订阅，值为1表示支持通配符订阅。
//如果没有设置此值，则表示支持通配符订阅。
//包含多个通配符订阅可用属性，或通配符订阅可用属性值不为0也不为1将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetWildcardSubscriptionAvailable() (byte, bool) {
    p := packet.FindPropValue(packet.WildcardSubscriptionAvailable, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropWildcardSubscriptionAvailable).V, true
}

//单字节字段，用来声明服务端是否支持订阅标识符（Subscription Identifiers）。
// 值为0表示不支持订阅标识符，值为1表示支持订阅标识符。如果没有设置此值，则表示支持订阅标识符。
// 包含多个订阅标识符可用属性，或订阅标识符可用属性值不为0也不为1将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetSubscriptionIdentifierAvailable(v byte) {
    p := &packet.PropSubscriptionIdentifierAvailable{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//单字节字段，用来声明服务端是否支持订阅标识符（Subscription Identifiers）。
// 值为0表示不支持订阅标识符，值为1表示支持订阅标识符。如果没有设置此值，则表示支持订阅标识符。
// 包含多个订阅标识符可用属性，或订阅标识符可用属性值不为0也不为1将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetSubscriptionIdentifierAvailable() (byte, bool) {
    p := packet.FindPropValue(packet.SubscriptionIdentifierAvailable, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropSubscriptionIdentifierAvailable).V, true
}

//单字节字段，用来声明服务端是否支持共享订阅（Shared Subscription）。
//值为0表示不支持共享订阅，值为1表示支持共享订阅。如果没有设置此值，则表示支持共享订阅。
//包含多个共享订阅可用（Shared Subscription Available），或共享订阅可用属性值不为0也不为1将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetSharedSubscriptionAvailable(v byte) {
    p := &packet.PropSharedSubscriptionAvailable{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//单字节字段，用来声明服务端是否支持共享订阅（Shared Subscription）。
//值为0表示不支持共享订阅，值为1表示支持共享订阅。如果没有设置此值，则表示支持共享订阅。
//包含多个共享订阅可用（Shared Subscription Available），或共享订阅可用属性值不为0也不为1将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetSharedSubscriptionAvailable() (byte, bool) {
    p := packet.FindPropValue(packet.SharedSubscriptionAvailable, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropSharedSubscriptionAvailable).V, true
}

//保持连接（Keep Alive）时间。
func (m *ConnackMessage) SetServerKeepAlive(v uint16) {
    p := &packet.PropServerKeepAlive{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

//保持连接（Keep Alive）时间。
func (m *ConnackMessage) GetServerKeepAlive() (uint16, bool) {
    p := packet.FindPropValue(packet.ServerKeepAlive, m.varHeader.props)
    if p == nil {
        return 0, false
    }
    return p.(*packet.PropServerKeepAlive).V, true
}

//以UTF-8编码的字符串，作为创建响应主题（Response Topic）的基本信息。
// 包含多个响应信息将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetResponseInformation(v string) {
    p := &packet.PropResponseInformation{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//以UTF-8编码的字符串，作为创建响应主题（Response Topic）的基本信息。
// 包含多个响应信息将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetResponseInformation() (string, bool) {
    p := packet.FindPropValue(packet.ResponseInformation, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropResponseInformation).V.String(), true
}

//以UTF-8编码的字符串，可以被客户端用来标识其他可用的服务端。
//包含多个服务端参考（Server Reference）将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetServerReference(v string) {
    p := &packet.PropServerReference{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//以UTF-8编码的字符串，可以被客户端用来标识其他可用的服务端。
//包含多个服务端参考（Server Reference）将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetServerReference() (string, bool) {
    p := packet.FindPropValue(packet.ServerReference, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropServerReference).V.String(), true
}

//以UTF-8编码的字符串，包含了认证方法（Authentication Method）名。
//包含多个认证方法将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetAuthenticationMethod(v string) {
    p := &packet.PropAuthenticationMethod{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//以UTF-8编码的字符串，包含了认证方法（Authentication Method）名。
//包含多个认证方法将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetAuthenticationMethod() (string, bool) {
    p := packet.FindPropValue(packet.AuthenticationMethod, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropAuthenticationMethod).V.String(), true
}


//包含认证数据（Authentication Data）的二进制数据。此数据的内容由认证方法和已交换的认证数据状态定义。
//包含多个认证数据将造成协议错误（Protocol Error）。
func (m *ConnackMessage) SetAuthenticationData(v []byte) {
    p := &packet.PropAuthenticationData{}
    s, err := packet.FromString(string(v))
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//包含认证数据（Authentication Data）的二进制数据。此数据的内容由认证方法和已交换的认证数据状态定义。
//包含多个认证数据将造成协议错误（Protocol Error）。
func (m *ConnackMessage) GetAuthenticationData() ([]byte, bool) {
    p := packet.FindPropValue(packet.AuthenticationData, m.varHeader.props)
    if p == nil {
        return nil, false
    }
    return []byte(p.(*packet.PropAuthenticationData).V.String()), true
}

func (v *ConnackVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("AckFlag: %d ReasonCode: %d \nprops:\n%s",
        v.AckFlag, v.ReasonCode, builder.String())
}

func (v *ConnackMessage) String() string {
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%s\n",
        v.fixedHeader, v.varHeader.String())
}