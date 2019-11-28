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

func ConnectMessageCreator() Message { return NewConnectMessage() }
func ConnackMessageCreator() Message { return NewConnackMessage() }

var creatorMap = map[byte]Creator{
    packet.PktTypeCONNECT: ConnectMessageCreator,
    packet.PktTypeCONNACK: ConnackMessageCreator,
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
