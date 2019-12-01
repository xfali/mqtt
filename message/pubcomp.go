// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "mqtt/packet"
)

//PUBCOMP报文是对PUBREL报文的响应。它是QoS 2等级协议交换的第四个也是最后一个报文。
type PubCompMessage struct {
    PubAckMessage
}

func NewPubCompMessage() *PubCompMessage {
    ret := &PubCompMessage{
        PubAckMessage{fixedHeader: packet.CreateFixedHeader(packet.PktTypePUBCOMP, packet.PktFlagPUBCOMP, 0)},
    }
    return ret
}