// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

import (
    "io"
    "mqtt/errcode"
)

const (
    PktTypeReserved    = 0  //禁止    保留
    PktTypeCONNECT     = 1  //客户端到服务端    客户端请求连接服务端
    PktTypeCONNACK     = 2  //服务端到客户端    连接报文确认
    PktTypePUBLISH     = 3  //两个方向都允许    发布消息
    PktTypePUBACK      = 4  //两个方向都允许    QoS 1消息发布收到确认
    PktTypePUBREC      = 5  //两个方向都允许    发布收到（保证交付第一步）
    PktTypePUBREL      = 6  //两个方向都允许    发布释放（保证交付第二步）
    PktTypePUBCOMP     = 7  //两个方向都允许    QoS 2消息发布完成（保证交互第三步）
    PktTypeSUBSCRIBE   = 8  //客户端到服务端    客户端订阅请求
    PktTypeSUBACK      = 9  //服务端到客户端    订阅请求报文确认
    PktTypeUNSUBSCRIBE = 10 //客户端到服务端    客户端取消订阅请求
    PktTypeUNSUBACK    = 11 //服务端到客户端    取消订阅报文确认
    PktTypePINGREQ     = 12 //客户端到服务端    心跳请求
    PktTypePINGRESP    = 13 //服务端到客户端    心跳响应
    PktTypeDISCONNECT  = 14 //两个方向都允许    断开连接通知
    PktTypeAUTH        = 15 //两个方向都允许    认证信息交换
)

const (
    PktFlagCONNECT     = 0 //Reserved	0	0	0	0
    PktFlagCONNACK     = 0 //Reserved	0	0	0	0
    PktFlagPUBLISH     = 0 //Used in MQTT v5.0	DUP	QoS	RETAIN
    PktFlagPUBACK      = 0 //Reserved	0	0	0	0
    PktFlagPUBREC      = 0 //Reserved	0	0	0	0
    PktFlagPUBREL      = 2 //Reserved	0	0	1	0
    PktFlagPUBCOMP     = 0 //Reserved	0	0	0	0
    PktFlagSUBSCRIBE   = 2 //Reserved	0	0	1	0
    PktFlagSUBACK      = 0 //Reserved	0	0	0	0
    PktFlagUNSUBSCRIBE = 2 //Reserved	0	0	1	0
    PktFlagUNSUBACK    = 0 //Reserved	0	0	0	0
    PktFlagPINGREQ     = 0 //Reserved	0	0	0	0
    PktFlagPINGRESP    = 0 //Reserved	0	0	0	0
    PktFlagDISCONNECT  = 0 //Reserved	0	0	0	0
    PktFlagAUTH        = 0 //Reserved	0	0	0	0
)

type FixedHeader struct {
    TypeFlag byte
    Len      int64
}

func ReadFixedHeader(r io.Reader) (FixedHeader, int, error) {
    fh := FixedHeader{}
    buf := make([]byte, 1)
    size := 0
    n, err := r.Read(buf)
    if err != nil {
        return fh, n, err
    }
    fh.TypeFlag = buf[0]
    size += n

    v := VarInt{}
    _, n2, err := v.LoadFromReader(r)
    if err != nil {
        if err != nil {
            return fh, size + n2, err
        }
    }
    fh.Len = v.ToInt()

    return fh, size + n2, nil
}

func WriteFixedHeader(w io.Writer, header FixedHeader) (int, error) {
    size := 0
    n1, err1 := w.Write([]byte{header.TypeFlag})
    if err1 != nil {
        return n1, err1
    }
    size += n1

    v := VarInt{}
    v.InitFromUInt64(uint64(header.Len))

    n2, err2 := w.Write(v.Bytes())
    if err2 != nil {
        return size + n2, err2
    }
    return size + n2, nil
}

func CreateFixedHeader(t, f byte, len int64) FixedHeader {
    return FixedHeader{
        TypeFlag: (t << 4) | (f & 0x0F),
        Len:      len,
    }
}

func (h FixedHeader) Type() byte {
    return h.TypeFlag >> 4
}

func (h *FixedHeader) CheckLen(n int) error {
    if h.Len >= int64(n) {
        return nil
    } else {
        return errcode.ProtocolError
    }
}

func (h FixedHeader) Flag() byte {
    return h.TypeFlag & 0x0F
}

//return dup, QoS, retain
func (h FixedHeader) PubFlag() (bool, uint8, bool) {
    flag := h.TypeFlag & 0x0F
    return flag>>3 == 1, (flag & 0x6) >> 1, flag&0x1 == 1
}

func (h *FixedHeader) SetDup(v bool) {
    if v {
        h.TypeFlag |= 1 << 3
    } else {
        h.TypeFlag &= 0xFF & ^(1 << 3)
    }
}

func (h *FixedHeader) SetQos(v byte) {
    v &= 0x3
    h.TypeFlag &= 0xFF & ^(0x3 << 1)
    h.TypeFlag |= v << 1
}

func (h *FixedHeader) SetRetain(v bool) {
    if v {
        h.TypeFlag |= 1
    } else {
        h.TypeFlag &= 0xFF & ^1
    }
}

func (h FixedHeader) RemainLength() int64 {
    return h.Len
}
