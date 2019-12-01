// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package message

import (
    "io"
    "mqtt/errcode"
    "mqtt/packet"
)

type Creator func() Message

func ConnectMessageCreator() Message     { return NewConnectMessage() }
func ConnackMessageCreator() Message     { return NewConnackMessage() }
func PublishMessageCreator() Message     { return NewPublishMessage() }
func PubAckMessageCreator() Message      { return NewPubAckMessage() }
func PubRecMessageCreator() Message      { return NewPubRecMessage() }
func PubRelMessageCreator() Message      { return NewPubRelMessage() }
func PubCompMessageCreator() Message     { return NewPubCompMessage() }
func SubscribeMessageCreator() Message   { return NewSubscribeMessage() }
func SubAckMessageCreator() Message      { return NewSubAckMessage() }
func UnsubscribeMessageCreator() Message { return NewUnsubscribeMessage() }
func UnsubAckMessageCreator() Message    { return NewUnsubAckMessage() }
func PingReqMessageCreator() Message     { return NewPingReqMessage() }
func PingRespMessageCreator() Message    { return NewPingRespMessage() }

var creatorMap = map[byte]Creator{
    packet.PktTypeCONNECT:     ConnectMessageCreator,
    packet.PktTypeCONNACK:     ConnackMessageCreator,
    packet.PktTypePUBLISH:     PublishMessageCreator,
    packet.PktTypePUBACK:      PubAckMessageCreator,
    packet.PktTypePUBREC:      PubRecMessageCreator,
    packet.PktTypePUBREL:      PubRelMessageCreator,
    packet.PktTypePUBCOMP:     PubCompMessageCreator,
    packet.PktTypeSUBSCRIBE:   SubscribeMessageCreator,
    packet.PktTypeSUBACK:      SubAckMessageCreator,
    packet.PktTypeUNSUBSCRIBE: UnsubscribeMessageCreator,
    packet.PktTypeUNSUBACK:    UnsubAckMessageCreator,
    packet.PktTypePINGREQ:     PingReqMessageCreator,
    packet.PktTypePINGRESP:    PingRespMessageCreator,
}

func ReadMessage(r io.Reader) (Message, int, error) {
    f, n, err := packet.ReadFixedHeader(r)
    if err != nil {
        return nil, n, err
    }

    creator := creatorMap[f.Type()]
    if creator == nil {
        return nil, n, errcode.MessageNotSupport
    }

    msg := creator()
    msg.SetFixedHeader(f)

    n2, err2 := msg.ReadVariableHeader(r)
    n += n2
    if err2 != nil {
        return nil, n, err2
    }

    n3, err3 := msg.ReadPayload(r)
    n += n3

    if f.RemainLength() != int64(n2+n3) {
        return nil, n, errcode.MessageReadSizeNotMatch
    }

    return msg, n, err3
}

func WriteMessage(w io.Writer, m Message) (int, error) {
    fixedHeader := m.GetFixedHeader()
    n, err := packet.WriteFixedHeader(w, fixedHeader)
    if err != nil {
        return n, err
    }

    n2, err2 := m.WriteVariableHeader(w)
    n += n2
    if err2 != nil {
        return n, err2
    }

    n3, err3 := m.WritePayload(w)
    n += n3
    return n, err3
}
