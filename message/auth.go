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

type AuthVarHeader struct {
    //断开原因码
    ReasonCode byte
    props      []packet.Property
}

//AUTH报文被从客户端发送给服务端，或从服务端发送给客户端，作为扩展认证交换的一部分，比如质询/响应认证。
//如果CONNECT报文不包含相同的认证方法，则客户端或服务端发送AUTH报文将造成协议错误（Protocol Error）。
type AuthMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   DisconnectVarHeader
}

func (m *AuthMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *AuthMessage) GetFixedHeader() packet.FixedHeader {
    w := new(util.CountWriter)
    m.WriteVariableHeader(w)
    m.WritePayload(w)
    m.fixedHeader.Len = w.Count()
    return m.fixedHeader
}

func (msg *AuthMessage) ReadVariableHeader(r io.Reader) (int, error) {
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

func (msg *AuthMessage) WriteVariableHeader(w io.Writer) (int, error) {
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

func (msg *AuthMessage) ReadPayload(r io.Reader) (int, error) {
    return 0, nil
}

func (msg *AuthMessage) WritePayload(w io.Writer) (int, error) {
    return 0, nil
}

func (m *AuthMessage) Valid() bool {
    return true
}

func NewAuthMessage() *AuthMessage {
    ret := &AuthMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypeAUTH, packet.PktFlagAUTH, 0),
    }
    return ret
}

func (m *AuthMessage) SetReasonCode(v byte) {
    m.varHeader.ReasonCode = v
}

func (m *AuthMessage) GetReasonCode() byte {
    return m.varHeader.ReasonCode
}

//跟随其后的是UTF-8字符串对。此属性可用于向客户端提供包括诊断信息在内的附加信息。
//如果加上用户属性之后的CONNACK报文长度超出了客户端指定的最大报文长度，则服务端不能发送此属性
func (m *AuthMessage) SetUserProperty(props map[string]string) {
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
func (m *AuthMessage) GetUserProperty() (map[string]string, bool) {
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
func (m *AuthMessage) SetAuthenticationMethod(v string) {
    p := &packet.PropAuthenticationMethod{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//以UTF-8编码的字符串，包含了认证方法（Authentication Method）名。
//包含多个认证方法将造成协议错误（Protocol Error）。
func (m *AuthMessage) GetAuthenticationMethod() (string, bool) {
    p := packet.FindPropValue(packet.AuthenticationMethod, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropAuthenticationMethod).V.String(), true
}


//包含认证数据（Authentication Data）的二进制数据。此数据的内容由认证方法和已交换的认证数据状态定义。
//包含多个认证数据将造成协议错误（Protocol Error）。
func (m *AuthMessage) SetAuthenticationData(v []byte) {
    p := &packet.PropAuthenticationData{}
    s, err := packet.FromString(string(v))
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//包含认证数据（Authentication Data）的二进制数据。此数据的内容由认证方法和已交换的认证数据状态定义。
//包含多个认证数据将造成协议错误（Protocol Error）。
func (m *AuthMessage) GetAuthenticationData() ([]byte, bool) {
    p := packet.FindPropValue(packet.AuthenticationData, m.varHeader.props)
    if p == nil {
        return nil, false
    }
    return []byte(p.(*packet.PropAuthenticationData).V.String()), true
}


//表示断开原因。此原因字符串是为诊断而设计的可读字符串，不应该被接收端所解析。
func (m *AuthMessage) SetReasonString(v string) {
    p := &packet.PropReasonString{}
    s, err := packet.FromString(v)
    if err == nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

//表示断开原因。此原因字符串是为诊断而设计的可读字符串，不应该被接收端所解析。
func (m *AuthMessage) GetReasonString() (string, bool) {
    p := packet.FindPropValue(packet.ReasonString, m.varHeader.props)
    if p == nil {
        return "", false
    }
    return p.(*packet.PropReasonString).V.String(), true
}


func (v *AuthVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("\t%v\n", v.props[i]))
    }
    return fmt.Sprintf("ReasonCode: %d \nprops:\n%s",
        v.ReasonCode, builder.String())
}

func (v *AuthMessage) String() string {
    return fmt.Sprintf("fixed header: \n%v\nvar header:\n%s\n",
        v.fixedHeader, v.varHeader.String())
}
