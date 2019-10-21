// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

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
    vAUTH              = 0 //Reserved	0	0	0	0
)

type FixedHeader [11]byte

func (h FixedHeader) Type() byte {
    return h[0] >> 4
}

func (h FixedHeader) Flag() byte {
    return h[0] & 0x0F
}

//return dup, QoS, retain
func (h FixedHeader) PubFlag() (bool, uint8, bool) {
    flag := h[0] & 0x0F
    return flag >> 3 == 1, (flag & 0x6) >> 1, flag & 0x1 == 1
}

func (h FixedHeader) RemainLength() uint64 {
    return DecodeVaruint(h[1:])
}

