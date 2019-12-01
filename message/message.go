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
    //检测message控制报文是否有效
    Valid() bool

    //设置固定头
    SetFixedHeader(header packet.FixedHeader)

    //获得固定头，注意remain length必须已经计算完成
    GetFixedHeader() packet.FixedHeader

    //读取可变头，注意必须在ReadPayload之前调用
    ReadVariableHeader(r io.Reader) (int, error)

    //写可变头，注意必须在WritePayload之前调用
    WriteVariableHeader(w io.Writer) (int, error)

    //读取payload，注意必须在ReadVariableHeader之后调用
    ReadPayload(r io.Reader) (int, error)

    //写payload，注意必须在WriteVariableHeader之后调用
    WritePayload(w io.Writer) (int, error)
}
