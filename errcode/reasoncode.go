// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package errcode

const (
    ReasonSuccess                             = 0   //0x00	成功	CONNACK, PUBACK, PUBREC, PUBREL, PUBCOMP, UNSUBACK, AUTH
    ReasonNormalDisconnection                 = 0   //0x00	正常断开	DISCONNECT
    ReasonGrantedQoS0                         = 0   //0x00	授权的QoS 0	SUBACK
    ReasonGrantedQoS1                         = 1   //0x01	授权的QoS 1	SUBACK
    ReasonGrantedQoS2                         = 2   //0x02	授权的QoS 2	SUBACK
    ReasonDisconnectWithWillMessage           = 4   //0x04	包含遗嘱的断开	DISCONNECT
    ReasonNoMatchingSubscribers               = 16  //0x10	无匹配订阅	PUBACK, PUBREC
    ReasonNoSubscriptionExisted               = 17  //0x11	订阅不存在	UNSUBACK
    ReasonContinueAuthentication              = 24  //0x18	继续认证	AUTH
    ReasonReauthenticate                      = 25  //0x19	重新认证	AUTH
    ReasonUnspecifiedError                    = 128 //0x80	未指明的错误	CONNACK, PUBACK, PUBREC, SUBACK, UNSUBACK, DISCONNECT
    ReasonMalformedPacket                     = 129 //0x81	无效报文	CONNACK, DISCONNECT
    ReasonProtocolError                       = 130 //0x82	协议错误	CONNACK, DISCONNECT
    ReasonImplementationSpecificError         = 131 //0x83	实现错误	CONNACK, PUBACK, PUBREC, SUBACK, UNSUBACK, DISCONNECT
    ReasonUnsupportedProtocolVersion          = 132 //0x84	协议版本不支持	CONNACK
    ReasonClientIdentifierNotValid            = 133 //0x85	客户标识符无效	CONNACK
    ReasonBadUserNameOrPassword               = 134 //0x86	用户名密码错误	CONNACK
    ReasonNotAuthorized                       = 135 //0x87	未授权	CONNACK, PUBACK, PUBREC, SUBACK, UNSUBACK, DISCONNECT
    ReasonServerUnavailable                   = 136 //0x88	服务端不可用	CONNACK
    ReasonServerBusy                          = 137 //0x89	服务端正忙	CONNACK, DISCONNECT
    ReasonBanned                              = 138 //0x8A	禁止	CONNACK
    ReasonServerShuttingDown                  = 139 //0x8B	服务端关闭中	DISCONNECT
    ReasonBadAuthenticationMethod             = 140 //0x8C	无效的认证方法	CONNACK, DISCONNECT
    ReasonKeepAliveTimeout                    = 141 //0x8D	保活超时	DISCONNECT
    ReasonSessionTakenOver                    = 142 //0x8E	会话被接管	DISCONNECT
    ReasonTopicFilterInvalid                  = 143 //0x8F	主题过滤器无效	SUBACK, UNSUBACK, DISCONNECT
    ReasonTopicNameInvalid                    = 144 //0x90	主题名无效	CONNACK, PUBACK, PUBREC, DISCONNECT
    ReasonPacketIdentifierInUse               = 145 //0x91	报文标识符已被占用	PUBACK, PUBREC, SUBACK, UNSUBACK
    ReasonPacketIdentifierNotFound            = 146 //0x92	报文标识符无效	PUBREL, PUBCOMP
    ReasonReceiveMaximumExceeded              = 147 //0x93	接收超出最大数量	DISCONNECT
    ReasonTopicAliasInvalid                   = 148 //0x94	主题别名无效	DISCONNECT
    ReasonPacketTooLarge                      = 149 //0x95	报文过长	CONNACK, DISCONNECT
    ReasonMessageRateTooHigh                  = 150 //0x96	消息太过频繁	DISCONNECT
    ReasonQuotaExceeded                       = 151 //0x97	超出配额	CONNACK, PUBACK, PUBREC, SUBACK, DISCONNECT
    ReasonAdministrativeAction                = 152 //0x98	管理行为	DISCONNECT
    ReasonPayloadFormatInvalid                = 153 //0x99	载荷格式无效	CONNACK, PUBACK, PUBREC, DISCONNECT
    ReasonRetainNotSupported                  = 154 //0x9A	不支持保留	CONNACK, DISCONNECT
    ReasonQoSNotSupported                     = 155 //0x9B	不支持的QoS等级	CONNACK, DISCONNECT
    ReasonUseAnotherServer                    = 156 //0x9C	（临时）使用其他服务端	CONNACK, DISCONNECT
    ReasonServerMoved                         = 157 //0x9D	服务端已（永久）移动	CONNACK, DISCONNECT
    ReasonSharedSubscriptionsNotSupported     = 158 //0x9E	不支持共享订阅	SUBACK, DISCONNECT
    ReasonConnectionRateExceeded              = 159 //0x9F	超出连接速率限制	CONNACK, DISCONNECT
    ReasonMaximumConnectTime                  = 160 //0xA0	最大连接时间	DISCONNECT
    ReasonSubscriptionIdentifiersNotSupported = 161 //0xA1	不支持订阅标识符	SUBACK, DISCONNECT
    ReasonWildcardSubscriptionsNotSupported   = 162 //0xA2	不支持通配符订阅	SUBACK, DISCONNECT
)

type ReasonCode byte
