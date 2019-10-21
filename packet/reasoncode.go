// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

const (
    Success                             = 0   //0x00	成功	CONNACK, PUBACK, PUBREC, PUBREL, PUBCOMP, UNSUBACK, AUTH
    Normaldisconnection                 = 0   //0x00	正常断开	DISCONNECT
    GrantedQoS0                         = 0   //0x00	授权的QoS 0	SUBACK
    GrantedQoS1                         = 1   //0x01	授权的QoS 1	SUBACK
    GrantedQoS2                         = 2   //0x02	授权的QoS 2	SUBACK
    DisconnectwithWillMessage           = 4   //0x04	包含遗嘱的断开	DISCONNECT
    NoMatchingSubscribers               = 16  //0x10	无匹配订阅	PUBACK, PUBREC
    NoSubscriptionExisted               = 17  //0x11	订阅不存在	UNSUBACK
    ContinueAuthentication              = 24  //0x18	继续认证	AUTH
    Reauthenticate                      = 25  //0x19	重新认证	AUTH
    UnspecifiedError                    = 128 //0x80	未指明的错误	CONNACK, PUBACK, PUBREC, SUBACK, UNSUBACK, DISCONNECT
    MalformedPacket                     = 129 //0x81	无效报文	CONNACK, DISCONNECT
    ProtocolError                       = 130 //0x82	协议错误	CONNACK, DISCONNECT
    ImplementationSpecificError         = 131 //0x83	实现错误	CONNACK, PUBACK, PUBREC, SUBACK, UNSUBACK, DISCONNECT
    UnsupportedProtocolVersion          = 132 //0x84	协议版本不支持	CONNACK
    ClientIdentifierNotValid            = 133 //0x85	客户标识符无效	CONNACK
    BadUserNameOrPassword               = 134 //0x86	用户名密码错误	CONNACK
    NotAuthorized                       = 135 //0x87	未授权	CONNACK, PUBACK, PUBREC, SUBACK, UNSUBACK, DISCONNECT
    ServerUnavailable                   = 136 //0x88	服务端不可用	CONNACK
    ServerBusy                          = 137 //0x89	服务端正忙	CONNACK, DISCONNECT
    Banned                              = 138 //0x8A	禁止	CONNACK
    ServerShuttingDown                  = 139 //0x8B	服务端关闭中	DISCONNECT
    BadAuthenticationMethod             = 140 //0x8C	无效的认证方法	CONNACK, DISCONNECT
    KeepAliveTimeout                    = 141 //0x8D	保活超时	DISCONNECT
    SessionTakenOver                    = 142 //0x8E	会话被接管	DISCONNECT
    TopicFilterInvalid                  = 143 //0x8F	主题过滤器无效	SUBACK, UNSUBACK, DISCONNECT
    TopicNameInvalid                    = 144 //0x90	主题名无效	CONNACK, PUBACK, PUBREC, DISCONNECT
    PacketIdentifierInUse               = 145 //0x91	报文标识符已被占用	PUBACK, PUBREC, SUBACK, UNSUBACK
    PacketIdentifierNotFound            = 146 //0x92	报文标识符无效	PUBREL, PUBCOMP
    ReceiveMaximumExceeded              = 147 //0x93	接收超出最大数量	DISCONNECT
    TopicAliasInvalid                   = 148 //0x94	主题别名无效	DISCONNECT
    PacketTooLarge                      = 149 //0x95	报文过长	CONNACK, DISCONNECT
    MessageRateTooHigh                  = 150 //0x96	消息太过频繁	DISCONNECT
    QuotaExceeded                       = 151 //0x97	超出配额	CONNACK, PUBACK, PUBREC, SUBACK, DISCONNECT
    AdministrativeAction                = 152 //0x98	管理行为	DISCONNECT
    PayloadFormatInvalid                = 153 //0x99	载荷格式无效	CONNACK, PUBACK, PUBREC, DISCONNECT
    RetainNotSupported                  = 154 //0x9A	不支持保留	CONNACK, DISCONNECT
    QoSNotSupported                     = 155 //0x9B	不支持的QoS等级	CONNACK, DISCONNECT
    UseAnotherServer                    = 156 //0x9C	（临时）使用其他服务端	CONNACK, DISCONNECT
    ServerMoved                         = 157 //0x9D	服务端已（永久）移动	CONNACK, DISCONNECT
    SharedSubscriptionsNotSupported     = 158 //0x9E	不支持共享订阅	SUBACK, DISCONNECT
    ConnectionRateExceeded              = 159 //0x9F	超出连接速率限制	CONNACK, DISCONNECT
    MaximumConnectTime                  = 160 //0xA0	最大连接时间	DISCONNECT
    SubscriptionIdentifiersNotSupported = 161 //0xA1	不支持订阅标识符	SUBACK, DISCONNECT
    WildcardSubscriptionsNotSupported   = 162 //0xA2	不支持通配符订阅	SUBACK, DISCONNECT
)

type ReasonCode byte
