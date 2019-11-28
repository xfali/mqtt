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

type Message interface {
    Valid() bool
    SetFixedHeader(header packet.FixedHeader)
    GetFixedHeader() packet.FixedHeader
    ReadVariableHeader(r io.Reader) (int, error)
    WriteVariableHeader(w io.Writer) (int, error)
    ReadPayload(r io.Reader) (int, error)
    WritePayload(w io.Writer) (int, error)
}
