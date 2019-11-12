// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

const (
    PayloadFormatIndicator          = 1  //0x01	载荷格式说明	字节	PUBLISH, Will Properties
    MessageExpiryInterval           = 2  //0x02	消息过期时间	四字节整数	PUBLISH, Will Properties
    ContentType                     = 3  //0x03	内容类型	UTF-8编码字符串	PUBLISH, Will Properties
    ResponseTopic                   = 8  //0x08	响应主题	UTF-8编码字符串	PUBLISH, Will Properties
    CorrelationData                 = 9  //0x09	相关数据	二进制数据	PUBLISH, Will Properties
    SubscriptionIdentifier          = 11 //0x0B	定义标识符	变长字节整数	PUBLISH, SUBSCRIBE
    SessionExpiryInterval           = 17 //0x11	会话过期间隔	四字节整数	CONNECT, CONNACK, DISCONNECT
    AssignedClientIdentifier        = 18 //0x12	分配客户标识符	UTF-8编码字符串	CONNACK
    ServerKeepAlive                 = 19 //0x13	服务端保活时间	双字节整数	CONNACK
    AuthenticationMethod            = 21 //0x15	认证方法	UTF-8编码字符串	CONNECT, CONNACK, AUTH
    AuthenticationData              = 22 //0x16	认证数据	二进制数据	CONNECT, CONNACK, AUTH
    RequestProblemInformation       = 23 //0x17	请求问题信息	字节	CONNECT
    WillDelayInterval               = 24 //0x18	遗嘱延时间隔	四字节整数	Will Properties
    RequestResponseInformation      = 25 //0x19	请求响应信息	字节	CONNECT
    ResponseInformation             = 26 //0x1A	请求信息	UTF-8编码字符串	CONNACK
    ServerReference                 = 28 //0x1C	服务端参考	UTF-8编码字符串	CONNACK, DISCONNECT
    ReasonString                    = 31 //0x1F	原因字符串	UTF-8编码字符串	CONNACK, PUBACK, PUBREC, PUBREL, PUBCOMP, SUBACK, UNSUBACK, DISCONNECT, AUTH
    ReceiveMaximum                  = 33 //0x21	接收最大数量	双字节整数	CONNECT, CONNACK
    TopicAliasMaximum               = 34 //0x22	主题别名最大长度	双字节整数	CONNECT, CONNACK
    TopicAlias                      = 35 //0x23	主题别名	双字节整数	PUBLISH
    MaximumQoS                      = 36 //0x24	最大QoS	字节	CONNACK
    RetainAvailable                 = 37 //0x25	保留属性可用性	字节	CONNACK
    UserProperty                    = 38 //0x26	用户属性	UTF-8字符串对	CONNECT, CONNACK, PUBLISH, Will Properties, PUBACK, PUBREC, PUBREL, PUBCOMP, SUBSCRIBE, SUBACK, UNSUBSCRIBE, UNSUBACK, DISCONNECT, AUTH
    MaximumPacketSize               = 39 //0x27	最大报文长度	四字节整数	CONNECT, CONNACK
    WildcardSubscriptionAvailable   = 40 //0x28	通配符订阅可用性	字节	CONNACK
    SubscriptionIdentifierAvailable = 41 //0x29	订阅标识符可用性	字节	CONNACK
    SharedSubscriptionAvailable     = 42 //0x2A	共享订阅可用性	字节	CONNACK
)

type Identifier uint16
