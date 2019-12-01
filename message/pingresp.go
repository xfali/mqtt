// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import "mqtt/packet"

//服务端发送PINGRESP报文响应客户端的PINGREQ报文。表示服务端还活着。
//保持连接（Keep Alive）处理中用到这个报文
type PingRespMessage struct {
    PingReqMessage
}

func NewPingRespMessage() *PingRespMessage {
    ret := &PingRespMessage{
        PingReqMessage{fixedHeader: packet.CreateFixedHeader(packet.PktTypePINGRESP, packet.PktFlagPINGRESP, 0),},
    }
    return ret
}
