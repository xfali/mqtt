// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

import (
    "encoding/binary"
    "fmt"
    "io"
    "mqtt/errcode"
)

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

/*
var (
    HexPayloadFormatIndicator          = []byte{PayloadFormatIndicator}
    HexMessageExpiryInterval           = []byte{MessageExpiryInterval}
    HexContentType                     = []byte{ContentType}
    HexResponseTopic                   = []byte{ResponseTopic}
    HexCorrelationData                 = []byte{CorrelationData}
    HexSubscriptionIdentifier          = []byte{SubscriptionIdentifier}
    HexSessionExpiryInterval           = []byte{SessionExpiryInterval}
    HexAssignedClientIdentifier        = []byte{AssignedClientIdentifier}
    HexServerKeepAlive                 = []byte{ServerKeepAlive}
    HexAuthenticationMethod            = []byte{AuthenticationMethod}
    HexAuthenticationData              = []byte{AuthenticationData}
    HexRequestProblemInformation       = []byte{RequestProblemInformation}
    HexWillDelayInterval               = []byte{WillDelayInterval}
    HexRequestResponseInformation      = []byte{RequestResponseInformation}
    HexResponseInformation             = []byte{ResponseInformation}
    HexServerReference                 = []byte{ServerReference}
    HexReasonString                    = []byte{ReasonString}
    HexReceiveMaximum                  = []byte{ReceiveMaximum}
    HexTopicAliasMaximum               = []byte{TopicAliasMaximum}
    HexTopicAlias                      = []byte{TopicAlias}
    HexMaximumQoS                      = []byte{MaximumQoS}
    HexRetainAvailable                 = []byte{RetainAvailable}
    HexUserProperty                    = []byte{UserProperty}
    HexMaximumPacketSize               = []byte{MaximumPacketSize}
    HexWildcardSubscriptionAvailable   = []byte{WildcardSubscriptionAvailable}
    HexSubscriptionIdentifierAvailable = []byte{SubscriptionIdentifierAvailable}
    HexSharedSubscriptionAvailable     = []byte{SharedSubscriptionAvailable}
)
*/
const (
    PROPERTY_DECHEX_SIZE = 1
)

type Setter interface {
    Set(interface{})
}

type Getter interface {
    Get() interface{}
}

type Property interface {
    Id() int64
    DataLen() int32
    UnmarshalData(io.Reader) (int, error)
    MarshalData(io.Writer) (int, error)

    Setter
    Getter
}

type PropPayloadFormatIndicator struct{ ByteProperty }
type PropMessageExpiryInterval struct{ Uint32Property }
type PropContentType struct{ StringProperty }
type PropResponseTopic struct{ StringProperty }
type PropCorrelationData struct{ StringProperty } //ByteDataProperty
type PropSubscriptionIdentifier struct{ VarIntProperty }
type PropSessionExpiryInterval struct{ Uint32Property }
type PropAssignedClientIdentifier struct{ StringProperty }
type PropServerKeepAlive struct{ Uint16Property }
type PropAuthenticationMethod struct{ StringProperty }
type PropAuthenticationData struct{ StringProperty } //ByteDataProperty
type PropRequestProblemInformation struct{ ByteProperty }
type PropWillDelayInterval struct{ Uint32Property }
type PropRequestResponseInformation struct{ ByteProperty }
type PropResponseInformation struct{ StringProperty }
type PropServerReference struct{ StringProperty }
type PropReasonString struct{ StringProperty }
type PropReceiveMaximum struct{ Uint16Property }
type PropTopicAliasMaximum struct{ Uint16Property }
type PropTopicAlias struct{ Uint16Property }
type PropMaximumQoS struct{ ByteProperty }
type PropRetainAvailable struct{ ByteProperty }
type PropUserProperty struct{ StringPairProperty }
type PropMaximumPacketSize struct{ Uint32Property }
type PropWildcardSubscriptionAvailable struct{ ByteProperty }
type PropSubscriptionIdentifierAvailable struct{ ByteProperty }
type PropSharedSubscriptionAvailable struct{ ByteProperty }

type PropertyCreator func() Property

func PayloadFormatIndicatorCreator() Property              { return &PropPayloadFormatIndicator{} }
func PropMessageExpiryIntervalCreator() Property           { return &PropMessageExpiryInterval{} }
func PropContentTypeCreator() Property                     { return &PropContentType{} }
func PropResponseTopicCreator() Property                   { return &PropResponseTopic{} }
func PropCorrelationDataCreator() Property                 { return &PropCorrelationData{} }
func PropSubscriptionIdentifierCreator() Property          { return &PropSubscriptionIdentifier{} }
func PropSessionExpiryIntervalCreator() Property           { return &PropSessionExpiryInterval{} }
func PropAssignedClientIdentifierCreator() Property        { return &PropAssignedClientIdentifier{} }
func PropServerKeepAliveCreator() Property                 { return &PropServerKeepAlive{} }
func PropAuthenticationMethodCreator() Property            { return &PropAuthenticationMethod{} }
func PropAuthenticationDataCreator() Property              { return &PropAuthenticationData{} }
func PropRequestProblemInformationCreator() Property       { return &PropRequestProblemInformation{} }
func PropWillDelayIntervalCreator() Property               { return &PropWillDelayInterval{} }
func PropRequestResponseInformationCreator() Property      { return &PropRequestResponseInformation{} }
func PropResponseInformationCreator() Property             { return &PropResponseInformation{} }
func PropServerReferenceCreator() Property                 { return &PropServerReference{} }
func PropReasonStringCreator() Property                    { return &PropReasonString{} }
func PropReceiveMaximumCreator() Property                  { return &PropReceiveMaximum{} }
func PropTopicAliasMaximumCreator() Property               { return &PropTopicAliasMaximum{} }
func PropTopicAliasCreator() Property                      { return &PropTopicAlias{} }
func PropMaximumQoSCreator() Property                      { return &PropMaximumQoS{} }
func PropRetainAvailableCreator() Property                 { return &PropRetainAvailable{} }
func PropUserPropertyCreator() Property                    { return &PropUserProperty{} }
func PropMaximumPacketSizeCreator() Property               { return &PropMaximumPacketSize{} }
func PropWildcardSubscriptionAvailableCreator() Property   { return &PropWildcardSubscriptionAvailable{} }
func PropSubscriptionIdentifierAvailableCreator() Property { return &PropSubscriptionIdentifierAvailable{} }
func PropSharedSubscriptionAvailableCreator() Property     { return &PropSharedSubscriptionAvailable{} }

var propFac = map[byte]PropertyCreator{
    PayloadFormatIndicator:          PayloadFormatIndicatorCreator,
    MessageExpiryInterval:           PropMessageExpiryIntervalCreator,
    ContentType:                     PropContentTypeCreator,
    ResponseTopic:                   PropResponseTopicCreator,
    CorrelationData:                 PropCorrelationDataCreator,
    SubscriptionIdentifier:          PropSubscriptionIdentifierCreator,
    SessionExpiryInterval:           PropSessionExpiryIntervalCreator,
    AssignedClientIdentifier:        PropAssignedClientIdentifierCreator,
    ServerKeepAlive:                 PropServerKeepAliveCreator,
    AuthenticationMethod:            PropAuthenticationMethodCreator,
    AuthenticationData:              PropAuthenticationDataCreator,
    RequestProblemInformation:       PropRequestProblemInformationCreator,
    WillDelayInterval:               PropWillDelayIntervalCreator,
    RequestResponseInformation:      PropRequestResponseInformationCreator,
    ResponseInformation:             PropResponseInformationCreator,
    ServerReference:                 PropServerReferenceCreator,
    ReasonString:                    PropReasonStringCreator,
    ReceiveMaximum:                  PropReceiveMaximumCreator,
    TopicAliasMaximum:               PropTopicAliasMaximumCreator,
    TopicAlias:                      PropTopicAliasCreator,
    MaximumQoS:                      PropMaximumQoSCreator,
    RetainAvailable:                 PropRetainAvailableCreator,
    UserProperty:                    PropUserPropertyCreator,
    MaximumPacketSize:               PropMaximumPacketSizeCreator,
    WildcardSubscriptionAvailable:   PropWildcardSubscriptionAvailableCreator,
    SubscriptionIdentifierAvailable: PropSubscriptionIdentifierAvailableCreator,
    SharedSubscriptionAvailable:     PropSharedSubscriptionAvailableCreator,
}

func UnmarshalProp(r io.Reader) (Property, int, error) {
    b := make([]byte, 1)
    n1, err := io.ReadFull(r, b)
    if err != nil {
        return nil, n1, err
    }
    prop := CreateProperty(b[0])
    if prop == nil {
        return nil, n1, errcode.UnknownProperty
    }

    n2, err := prop.UnmarshalData(r)
    if err != nil {
        return nil, n1 + n2, err
    }
    return prop, n1 + n2, nil
}

func MarshalProp(w io.Writer, prop Property) (int, error) {
    n1, err := w.Write([]byte{byte(prop.Id())})
    if err != nil {
        return n1, err
    }

    n2, err := prop.MarshalData(w)
    return n1 + n2, err
}

func CreateProperty(t byte) Property {
    if v, ok := propFac[t]; ok {
        return v()
    }
    return nil
}

type Uint16Property struct {
    V uint16
}

func (prop *Uint16Property) Set(v interface{}) {
    prop.V = v.(uint16)
}

func (prop *Uint16Property) Get() interface{} {
    return prop.V
}

func (prop *Uint16Property) DataLen() int32 {
    return 2
}

func (prop *Uint16Property) UnmarshalData(r io.Reader) (int, error) {
    buf := make([]byte, 2)
    n, err := io.ReadFull(r, buf)
    if err != nil {
        return n, err
    }

    prop.V = binary.BigEndian.Uint16(buf)
    return n, nil
}

func (prop *Uint16Property) MarshalData(w io.Writer) (int, error) {
    buf := make([]byte, 2)
    binary.BigEndian.PutUint16(buf, prop.V)
    return w.Write(buf)
}

type Uint32Property struct {
    V uint32
}

func (prop *Uint32Property) Set(v interface{}) {
    prop.V = v.(uint32)
}

func (prop *Uint32Property) Get() interface{} {
    return prop.V
}

func (prop *Uint32Property) DataLen() int32 {
    return 4
}

func (prop *Uint32Property) UnmarshalData(r io.Reader) (int, error) {
    buf := make([]byte, 4)
    n, err := io.ReadFull(r, buf)
    if err != nil {
        return n, err
    }

    prop.V = binary.BigEndian.Uint32(buf)
    return n, nil
}

func (prop *Uint32Property) MarshalData(w io.Writer) (int, error) {
    buf := make([]byte, 4)
    binary.BigEndian.PutUint32(buf, prop.V)
    return w.Write(buf)
}

type ByteProperty struct {
    V byte
}

func (prop *ByteProperty) Set(v interface{}) {
    prop.V = v.(byte)
}

func (prop *ByteProperty) Get() interface{} {
    return prop.V
}

func (prop *ByteProperty) DataLen() int32 {
    return 1
}

func (prop *ByteProperty) UnmarshalData(r io.Reader) (int, error) {
    buf := make([]byte, 1)
    n, err := io.ReadFull(r, buf)
    if err != nil {
        return n, err
    }

    prop.V = buf[0]
    return n, nil
}

func (prop *ByteProperty) MarshalData(w io.Writer) (int, error) {
    return w.Write([]byte{prop.V})
}

type VarIntProperty struct {
    V VarInt
}

func (prop *VarIntProperty) Set(v interface{}) {
    prop.V = v.(VarInt)
}

func (prop *VarIntProperty) Get() interface{} {
    return prop.V
}

func (prop *VarIntProperty) DataLen() int32 {
    return int32(prop.V.Length())
}

func (prop *VarIntProperty) UnmarshalData(r io.Reader) (int, error) {
    size := 0
    for {
        b, n, err := prop.V.LoadFromReader(r)
        if err != nil {
            return size + n, err
        }
        if b {
            break
        }
        size += n
    }
    return size, nil
}

func (prop *VarIntProperty) MarshalData(w io.Writer) (int, error) {
    return w.Write(prop.V.Bytes())
}

type StringProperty struct {
    V String
}

func (prop *StringProperty) Set(v interface{}) {
    prop.V = v.(String)
}

func (prop *StringProperty) Get() interface{} {
    return prop.V
}

func (prop *StringProperty) DataLen() int32 {
    return int32(prop.V.AllLength())
}

func (prop *StringProperty) UnmarshalData(r io.Reader) (int, error) {
    s, n, err := ParseString(r)
    if err != nil {
        return n, err
    }
    prop.V = *s
    return n, nil
}

func (prop *StringProperty) MarshalData(w io.Writer) (int, error) {
    return WriteString(w, prop.V)
}

type StringPairProperty struct {
    V StringPair
}

func (prop *StringPairProperty) Set(v interface{}) {
    prop.V = v.(StringPair)
}

func (prop *StringPairProperty) Get() interface{} {
    return prop.V
}

func (prop *StringPairProperty) DataLen() int32 {
    return int32(prop.V[0].AllLength() + prop.V[1].AllLength())
}

func (prop *StringPairProperty) UnmarshalData(r io.Reader) (int, error) {
    s1, n1, err := ParseString(r)
    if err != nil {
        return n1, err
    }
    prop.V[0] = *s1

    s2, n2, err := ParseString(r)
    if err != nil {
        return n1 + n2, err
    }
    prop.V[1] = *s2
    return n1 + n2, nil
}

func (prop *StringPairProperty) MarshalData(w io.Writer) (int, error) {
    n1, err := WriteString(w, prop.V[0])
    if err != nil {
        return n1, err
    }
    n2, err := WriteString(w, prop.V[1])
    return n1 + n2, err
}

func FindPropValue(t int64, props []Property) interface{} {
    for i := range props {
        if props[i].Id() == t {
            return props[i]
        }
    }
    return nil
}

//f return true, stop find
func FindPropValues(t int64, props []Property, f func(Property) bool) {
    for i := range props {
        if props[i].Id() == t {
            if f(props[i]) {
                return
            }
        }
    }
}

func FindAndSetPropValue(prop Property, props []Property) bool {
    for i := range props {
        if props[i].Id() == prop.Id() {
            prop.Set(props[i].Get())
            return true
        }
    }
    return false
}

func ReadProperties(r io.Reader) ([]Property, int, error) {
    v := NewFromReader(r)
    if v == nil {
        return nil, 0, errcode.ParseVarIntFailed
    }
    length := int(v.ToInt())
    size := v.Length()
    var propList []Property
    for size < length {
        p, n, err := UnmarshalProp(r)
        if err != nil {
            return nil, size, err
        }
        propList = append(propList, p)
        size += n
    }

    return propList, size, nil
}

func ReadPropertyMap(r io.Reader) (map[int64]Property, int, error) {
    v := NewFromReader(r)
    if v == nil {
        return nil, 0, errcode.ParseVarIntFailed
    }
    length := int(v.ToInt())
    size := v.Length()
    propMap := map[int64]Property{}
    for size < length {
        p, n, err := UnmarshalProp(r)
        if err != nil {
            return nil, size, err
        }
        propMap[p.Id()] = p
        size += n
    }

    return propMap, size, nil
}

func WriteProperties(w io.Writer, props []Property) (int, error) {
    v := VarInt{}
    propLen := 0
    for _, p := range props {
        propLen += PROPERTY_DECHEX_SIZE + int(p.DataLen())
    }
    v.InitFromUInt64(uint64(propLen))
    size := 0

    n, err := w.Write(v.Bytes())
    if err != nil {
        return n, err
    }
    size += n
    for _, v := range props {
        n, err := MarshalProp(w, v)
        if err != nil {
            return size + n, err
        }
        size += n
    }
    return size, nil
}

func (p *PropPayloadFormatIndicator) Id() int64          { return PayloadFormatIndicator }
func (p *PropMessageExpiryInterval) Id() int64           { return MessageExpiryInterval }
func (p *PropContentType) Id() int64                     { return ContentType }
func (p *PropResponseTopic) Id() int64                   { return ResponseTopic }
func (p *PropCorrelationData) Id() int64                 { return CorrelationData }
func (p *PropSubscriptionIdentifier) Id() int64          { return SubscriptionIdentifier }
func (p *PropSessionExpiryInterval) Id() int64           { return SessionExpiryInterval }
func (p *PropAssignedClientIdentifier) Id() int64        { return AssignedClientIdentifier }
func (p *PropServerKeepAlive) Id() int64                 { return ServerKeepAlive }
func (p *PropAuthenticationMethod) Id() int64            { return AuthenticationMethod }
func (p *PropAuthenticationData) Id() int64              { return AuthenticationData }
func (p *PropRequestProblemInformation) Id() int64       { return RequestProblemInformation }
func (p *PropWillDelayInterval) Id() int64               { return WillDelayInterval }
func (p *PropRequestResponseInformation) Id() int64      { return RequestResponseInformation }
func (p *PropResponseInformation) Id() int64             { return ResponseInformation }
func (p *PropServerReference) Id() int64                 { return ServerReference }
func (p *PropReasonString) Id() int64                    { return ReasonString }
func (p *PropReceiveMaximum) Id() int64                  { return ReceiveMaximum }
func (p *PropTopicAliasMaximum) Id() int64               { return TopicAliasMaximum }
func (p *PropTopicAlias) Id() int64                      { return TopicAlias }
func (p *PropMaximumQoS) Id() int64                      { return MaximumQoS }
func (p *PropRetainAvailable) Id() int64                 { return RetainAvailable }
func (p *PropUserProperty) Id() int64                    { return UserProperty }
func (p *PropMaximumPacketSize) Id() int64               { return MaximumPacketSize }
func (p *PropWildcardSubscriptionAvailable) Id() int64   { return WildcardSubscriptionAvailable }
func (p *PropSubscriptionIdentifierAvailable) Id() int64 { return SubscriptionIdentifierAvailable }
func (p *PropSharedSubscriptionAvailable) Id() int64     { return SharedSubscriptionAvailable }

func (p *PropPayloadFormatIndicator) String() string          { return fmt.Sprintf("PayloadFormatIndicator %v", p.V) }
func (p *PropMessageExpiryInterval) String() string           { return fmt.Sprintf("MessageExpiryInterval %v", p.V) }
func (p *PropContentType) String() string                     { return fmt.Sprintf("ContentType %v", p.V) }
func (p *PropResponseTopic) String() string                   { return fmt.Sprintf("ResponseTopic %v", p.V) }
func (p *PropCorrelationData) String() string                 { return fmt.Sprintf("CorrelationData %v", p.V) }
func (p *PropSubscriptionIdentifier) String() string          { return fmt.Sprintf("SubscriptionIdentifier %v", p.V) }
func (p *PropSessionExpiryInterval) String() string           { return fmt.Sprintf("SessionExpiryInterval %v", p.V) }
func (p *PropAssignedClientIdentifier) String() string        { return fmt.Sprintf("AssignedClientIdentifier %v", p.V) }
func (p *PropServerKeepAlive) String() string                 { return fmt.Sprintf("ServerKeepAlive %v", p.V) }
func (p *PropAuthenticationMethod) String() string            { return fmt.Sprintf("AuthenticationMethod %v", p.V) }
func (p *PropAuthenticationData) String() string              { return fmt.Sprintf("AuthenticationData %v", p.V) }
func (p *PropRequestProblemInformation) String() string       { return fmt.Sprintf("RequestProblemInformation %v", p.V) }
func (p *PropWillDelayInterval) String() string               { return fmt.Sprintf("WillDelayInterval %v", p.V) }
func (p *PropRequestResponseInformation) String() string      { return fmt.Sprintf("RequestResponseInformation %v", p.V) }
func (p *PropResponseInformation) String() string             { return fmt.Sprintf("ResponseInformation %v", p.V) }
func (p *PropServerReference) String() string                 { return fmt.Sprintf("ServerReference %v", p.V) }
func (p *PropReasonString) String() string                    { return fmt.Sprintf("ReasonString %v", p.V) }
func (p *PropReceiveMaximum) String() string                  { return fmt.Sprintf("ReceiveMaximum %v", p.V) }
func (p *PropTopicAliasMaximum) String() string               { return fmt.Sprintf("TopicAliasMaximum %v", p.V) }
func (p *PropTopicAlias) String() string                      { return fmt.Sprintf("TopicAlias %v", p.V) }
func (p *PropMaximumQoS) String() string                      { return fmt.Sprintf("MaximumQoS %v", p.V) }
func (p *PropRetainAvailable) String() string                 { return fmt.Sprintf("RetainAvailable %v", p.V) }
func (p *PropUserProperty) String() string                    { return fmt.Sprintf("UserProperty %s", p.V.String()) }
func (p *PropMaximumPacketSize) String() string               { return fmt.Sprintf("MaximumPacketSize %v", p.V) }
func (p *PropWildcardSubscriptionAvailable) String() string   { return fmt.Sprintf("WildcardSubscriptionAvailable %v", p.V) }
func (p *PropSubscriptionIdentifierAvailable) String() string { return fmt.Sprintf("SubscriptionIdentifierAvailable %v", p.V) }
func (p *PropSharedSubscriptionAvailable) String() string     { return fmt.Sprintf("SharedSubscriptionAvailable %v", p.V) }
