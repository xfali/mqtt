// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "io"
    "mqtt/packet"
)

//客户端发送PINGREQ报文给服务端，可被用于：
//1、在没有任何其他MQTT控制报文从客户端发给服务端时，告知服务端客户端还活着。
//2、请求服务端发送响应以确认服务端还活着。
//3、使用网络已确认网络连接没有断开。
type PingReqMessage struct {
    fixedHeader packet.FixedHeader
}

func NewPingReqMessage() *PingReqMessage {
    ret := &PingReqMessage{
        fixedHeader: packet.CreateFixedHeader(packet.PktTypePINGREQ, packet.PktFlagPINGREQ, 0),
    }
    return ret
}

func (m *PingReqMessage) SetFixedHeader(header packet.FixedHeader) {
    m.fixedHeader = header
}

func (m *PingReqMessage) GetFixedHeader() packet.FixedHeader {
    return m.fixedHeader
}

func (msg *PingReqMessage) ReadVariableHeader(r io.Reader) (int, error) {
    return 0, nil
}

func (msg *PingReqMessage) WriteVariableHeader(w io.Writer) (int, error) {
    return 0, nil
}

func (msg *PingReqMessage) ReadPayload(r io.Reader) (int, error) {
    return 0, nil
}

func (msg *PingReqMessage) WritePayload(w io.Writer) (int, error) {
    return 0, nil
}

func (msg *PingReqMessage) Valid() bool {
    return true
}
