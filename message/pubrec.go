// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import "mqtt/packet"

//PUBREC报文是对QoS等级2的PUBLISH报文的响应。它是QoS 2等级协议交换的第二个报文。
type PubRecMessage struct {
    PubAckMessage
}

func NewPubRecMessage() *PubRecMessage {
    ret := &PubRecMessage{
        PubAckMessage{fixedHeader: packet.CreateFixedHeader(packet.PktTypePUBREC, packet.PktFlagPUBREC, 0)},
    }
    return ret
}