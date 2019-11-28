// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "fmt"
    "io"
    "mqtt/container/binlist"
    "mqtt/errcode"
    "mqtt/packet"
    "mqtt/util"
    "strings"
)

const (
    ConnectMessageProtocolName  = "ProtocolName"
    ConnectMessageProtocolLevel = "ProtocolLevel"
    ConnectMessageConnectFlags  = "ConnectFlags"
    ConnectMessageKeepAlive     = "KeepAlive"
    ConnectMessageProperties    = "Properties"
)

var connectMessageVarHeaderMeta *binlist.BinMetaList

//func init() {
//    connectMessageVarHeaderMeta = binlist.NewMetaList()
//    connectMessageVarHeaderMeta.Put(ConnectMessageProtocolName, 0, 6)
//    connectMessageVarHeaderMeta.Put(ConnectMessageProtocolLevel, 6, 1)
//    connectMessageVarHeaderMeta.Put(ConnectMessageConnectFlags, 7, 1)
//    connectMessageVarHeaderMeta.Put(ConnectMessageKeepAlive, 8, 2)
//}

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

type ConnectMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   ConnectVarHeader
    payload     ConnectPayload
}

func (v *ConnectVarHeader) String() string {
    builder := strings.Builder{}
    for i := range v.props {
        builder.WriteString(fmt.Sprintf("%v\n", v.props[i]))
    }
    return fmt.Sprintf("protocal name: %s version: %d flag: %b keepAlive: %d \nprops: %s",
        v.ProtocolName.String(), v.ProtocolVersion, v.Flag, v.KeepAlive, builder.String())
}

func (v *ConnectPayload) String() string {
    builder := strings.Builder{}
    for i := range v.WillProps {
        builder.WriteString(fmt.Sprintf("%v\n", v.WillProps[i]))
    }
    return fmt.Sprintf("ClientId: %s WillTopic: %s WillPayload: %v Username: %s Password: %v \nprops: %s",
        v.ClientId.String(), v.WillTopic.String(), v.WillPayload, v.Username.String(), v.Password, builder.String())
}

func (v *ConnectMessage) String() string {
    return fmt.Sprintf("fixed header: %v, var header: %s, payload: %s",
        v.fixedHeader, v.varHeader.String(), v.payload.String())
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
    if msg.GetWillEnable() {
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

    if msg.GetWillEnable() {
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

func (m *ConnectMessage) SetWillQos(v byte) {
    m.varHeader.Flag |= (v & 0x3) << 3
}

func (m *ConnectMessage) SetVersion(v byte) {
    m.varHeader.ProtocolVersion = v
}

func (m *ConnectMessage) SetCleanStart(v bool) {
    if v {
        m.varHeader.Flag |= 1 << 1
    } else {
        m.varHeader.Flag &= 0xFF & ^(1 << 1)
    }
}

func (m *ConnectMessage) SetWillEnable(v bool) {
    if v {
        m.varHeader.Flag |= 1 << 2
    } else {
        m.varHeader.Flag &= 0xFF & ^(1 << 2)
    }
}

func (m *ConnectMessage) GetWillEnable() bool {
    return int(m.varHeader.Flag) & ^(1 << 2) != 0
}

func (m *ConnectMessage) SetClientId(v string) {
    m.payload.ClientId.Reset(v)
}

func (m *ConnectMessage) SetKeepAlive(v uint16) {
    m.varHeader.KeepAlive = v
}

func (m *ConnectMessage) SetWillRetain(v bool) {
    m.varHeader.Flag |= 1 << 5
}

func (m *ConnectMessage) SetWillTopic(v string) {
    m.payload.WillTopic.Reset(v)
}

func (m *ConnectMessage) SetWillPayload(v []byte) {
    m.payload.WillPayload.Reset(v)
}

func (m *ConnectMessage) SetUsername(v string) {
    m.varHeader.Flag |= 1 << 7
    m.payload.Username.Reset(v)
}

func (m *ConnectMessage) haveUsername() bool {
    return int(m.varHeader.Flag) & ^(1 << 7) != 0
}

func (m *ConnectMessage) SetPassword(v []byte) {
    m.varHeader.Flag |= 1 << 6
    m.payload.Password.Reset(v)
}

func (m *ConnectMessage) havePassword() bool {
    return int(m.varHeader.Flag) & ^(1 << 6) != 0
}

func (m *ConnectMessage) SetSessionExpiryInterval(v uint32) {
    p := &packet.PropSessionExpiryInterval{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

func (m *ConnectMessage) SetReceiveMaximum(v uint16) {
    p := &packet.PropReceiveMaximum{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

func (m *ConnectMessage) SetMaximumPacketSize(v uint32) {
    p := &packet.PropMaximumPacketSize{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
}

func (m *ConnectMessage) SetTopicAliasMaximum(v uint16) {
    p := &packet.PropTopicAliasMaximum{}
    p.V = v
    m.varHeader.props = append(m.varHeader.props, p)
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

func (m *ConnectMessage) SetAuthenticationMethod(v string) {
    p := &packet.PropAuthenticationMethod{}
    s, err := packet.FromString(v)
    if err != nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

func (m *ConnectMessage) SetAuthenticationData(v []byte) {
    p := &packet.PropAuthenticationData{}
    s, err := packet.FromString(string(v))
    if err != nil {
        p.V = s
        m.varHeader.props = append(m.varHeader.props, p)
    }
}

func (m *ConnectMessage) GetSessionExpiryInterval() (uint32, bool) {
    p := packet.FindPropValue(packet.SessionExpiryInterval, m.varHeader.props).(*packet.PropSessionExpiryInterval)
    if p == nil {
        return 0, false
    }
    return p.V, true
}

func (m *ConnectMessage) GetReceiveMaximum() (uint16, bool) {
    p := packet.FindPropValue(packet.ReceiveMaximum, m.varHeader.props).(*packet.PropReceiveMaximum)
    if p == nil {
        return 0, false
    }
    return p.V, true
}

func (m *ConnectMessage) GetMaximumPacketSize() (uint32, bool) {
    p := packet.FindPropValue(packet.MaximumPacketSize, m.varHeader.props).(*packet.PropMaximumPacketSize)
    if p == nil {
        return 0, false
    }
    return p.V, true
}

func (m *ConnectMessage) GtTopicAliasMaximum() (uint16, bool) {
    p := packet.FindPropValue(packet.TopicAliasMaximum, m.varHeader.props).(*packet.PropTopicAliasMaximum)
    if p == nil {
        return 0, false
    }
    return p.V, true
}

func (m *ConnectMessage) GetRequestResponseInformation() (byte, bool) {
    p := packet.FindPropValue(packet.RequestResponseInformation, m.varHeader.props).(*packet.PropRequestResponseInformation)
    if p == nil {
        return 0, false
    }
    return p.V, true
}

func (m *ConnectMessage) GetRequestProblemInformation() (byte, bool) {
    p := packet.FindPropValue(packet.RequestProblemInformation, m.varHeader.props).(*packet.PropRequestProblemInformation)
    if p == nil {
        return 0, false
    }
    return p.V, true
}

func (m *ConnectMessage) GetUserProperty() (map[string]string, bool) {
    ret := map[string]string{}
    packet.FindPropValues(packet.UserProperty, m.varHeader.props, func(property packet.Property) bool {
        p := property.(*packet.PropUserProperty)
        if p != nil {
            ret[p.V[0].String()] = p.V[1].String()
        }
        return false
    })
    return ret, len(ret) > 0
}

func (m *ConnectMessage) GetAuthenticationMethod() (string, bool) {
    p := packet.FindPropValue(packet.AuthenticationMethod, m.varHeader.props).(*packet.PropAuthenticationMethod)
    if p == nil {
        return "", false
    }
    return p.V.String(), true
}

func (m *ConnectMessage) GetAuthenticationData() ([]byte, bool) {
    p := packet.FindPropValue(packet.AuthenticationData, m.varHeader.props).(*packet.PropAuthenticationData)
    if p == nil {
        return nil, false
    }
    return []byte(p.V.String()), true
}
