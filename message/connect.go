// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "io"
    "mqtt/container/binlist"
    "mqtt/errcode"
    "mqtt/packet"
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

func (msg *ConnectMessage) ReadVarHeader(r io.Reader) (int, error) {
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

func (msg *ConnectMessage) WriteVarHeader(w io.Writer) (int, error) {
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
    ret := &ConnectMessage{}
    ret.varHeader.ProtocolName = packet.MqttProtocolNameString
    ret.varHeader.ProtocolVersion = packet.MqttProtocolVersion
    return ret
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

func (m *ConnectMessage) SetWillMessage(v []byte) {

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
