// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "mqtt/container/binlist"
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

func init() {
    connectMessageVarHeaderMeta = binlist.NewMetaList()
    connectMessageVarHeaderMeta.Put(ConnectMessageProtocolName, 0, 6)
    connectMessageVarHeaderMeta.Put(ConnectMessageProtocolLevel, 6, 1)
    connectMessageVarHeaderMeta.Put(ConnectMessageConnectFlags, 7, 1)
    connectMessageVarHeaderMeta.Put(ConnectMessageKeepAlive, 8, 2)
}

type ConnectMessage struct {
    fixedHeader packet.FixedHeader
    varHeader   []byte
}

func NewConnectMessage() *ConnectMessage {
    return &ConnectMessage{

    }
}

func (m *ConnectMessage) SetWillQos(v byte) {
    m.varHeader[7] |= (v & 0x3) << 3
}

func (m *ConnectMessage) SetVersion(v byte) {
    m.varHeader[6] = v
}

func (m *ConnectMessage) SetCleanStart(v bool) {
    if v {
        m.varHeader[7] |= 1 << 1
    } else {
        m.varHeader[7] &= 0xFF & ^(1 << 1)
    }
}

func (m *ConnectMessage) SetWillEnable() {
    m.varHeader[7] |= 1 << 2
}

func (m *ConnectMessage) SetClientId(v []byte) {

}

func (m *ConnectMessage) SetKeepAlive(v uint16) {
    m.varHeader[8] = byte(v >> 8)
    m.varHeader[9] = byte(v & 0xFF)
}

func (m *ConnectMessage) SetWillRetain(v bool) {
    m.varHeader[7] |= 1 << 5
}

func (m *ConnectMessage) SetWillTopic(v []byte) {

}

func (m *ConnectMessage) SetWillMessage(v []byte) {

}

func (m *ConnectMessage) SetUsername(v []byte) {
    m.varHeader[7] |= 1 << 7
}

func (m *ConnectMessage) SetPassword(v []byte) {
    m.varHeader[7] |= 1 << 6
}
