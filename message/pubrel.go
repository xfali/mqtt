// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "mqtt/packet"
)

//PUBREL报文是对PUBREC报文的响应。它是QoS 2等级协议交换的第三个报文。
type PubRelMessage struct {
    PubAckMessage
}

func NewPubRelMessage() *PubRelMessage {
    ret := &PubRelMessage{
        PubAckMessage{fixedHeader: packet.CreateFixedHeader(packet.PktTypePUBREL, packet.PktFlagPUBREL, 0)},
    }
    return ret
}