// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "bytes"
    "mqtt/packet"
    "testing"
)

func TestFixedHeader(t *testing.T) {
    header := packet.CreateFixedHeader(packet.PktTypeCONNECT, packet.PktFlagCONNECT, 10)
    buf := bytes.NewBuffer(nil)
    packet.WriteFixedHeader(buf, header)

    rh, _, _ := packet.ReadFixedHeader(buf)
    t.Log(rh)
}
