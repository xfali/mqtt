// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

import (
    "bytes"
    "encoding/binary"
    "io"
    "mqtt/errcode"
)

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
    DecHex() []byte
    DataLen() int32
    UnmarshalData(io.Reader) (int, error)
    MarshalData(io.Writer) (int, error)
}

type PropPayloadFormatIndicator ByteProperty
type PropMessageExpiryInterval Uint32Property
type PropContentType StringProperty
type PropResponseTopic StringProperty
type PropCorrelationData StringProperty //ByteDataProperty
type PropSubscriptionIdentifier VarIntProperty
type PropSessionExpiryInterval Uint32Property
type PropAssignedClientIdentifier StringProperty
type PropServerKeepAlive Uint16Property
type PropAuthenticationMethod StringProperty
type PropAuthenticationData StringProperty //ByteDataProperty
type PropRequestProblemInformation ByteProperty
type PropWillDelayInterval Uint32Property
type PropRequestResponseInformation ByteProperty
type PropResponseInformation StringProperty
type PropServerReference StringProperty
type PropReasonString StringProperty
type PropReceiveMaximum Uint16Property
type PropTopicAliasMaximum Uint16Property
type PropTopicAlias Uint16Property
type PropMaximumQoS ByteProperty
type PropRetainAvailable ByteProperty
type PropUserProperty StringPairProperty
type PropMaximumPacketSize Uint32Property
type PropWildcardSubscriptionAvailable ByteProperty
type PropSubscriptionIdentifierAvailable ByteProperty
type PropSharedSubscriptionAvailable ByteProperty

type PropertyCreator func() Property

func PayloadFormatIndicatorCreator() Property              { return &ByteProperty{decHex: []byte{PayloadFormatIndicator}} }
func PropMessageExpiryIntervalCreator() Property           { return &Uint32Property{decHex: []byte{MessageExpiryInterval}} }
func PropContentTypeCreator() Property                     { return &StringProperty{decHex: []byte{ContentType}} }
func PropResponseTopicCreator() Property                   { return &StringProperty{decHex: []byte{ResponseTopic}} }
func PropCorrelationDataCreator() Property                 { return &StringProperty{decHex: []byte{CorrelationData}} }
func PropSubscriptionIdentifierCreator() Property          { return &VarIntProperty{decHex: []byte{SubscriptionIdentifier}} }
func PropSessionExpiryIntervalCreator() Property           { return &Uint32Property{decHex: []byte{SessionExpiryInterval}} }
func PropAssignedClientIdentifierCreator() Property        { return &StringProperty{decHex: []byte{AssignedClientIdentifier}} }
func PropServerKeepAliveCreator() Property                 { return &Uint16Property{decHex: []byte{ServerKeepAlive}} }
func PropAuthenticationMethodCreator() Property            { return &StringProperty{decHex: []byte{AuthenticationMethod}} }
func PropAuthenticationDataCreator() Property              { return &StringProperty{decHex: []byte{AuthenticationData}} }
func PropRequestProblemInformationCreator() Property       { return &ByteProperty{decHex: []byte{RequestProblemInformation}} }
func PropWillDelayIntervalCreator() Property               { return &Uint32Property{decHex: []byte{WillDelayInterval}} }
func PropRequestResponseInformationCreator() Property      { return &ByteProperty{decHex: []byte{RequestResponseInformation}} }
func PropResponseInformationCreator() Property             { return &StringProperty{decHex: []byte{ResponseInformation}} }
func PropServerReferenceCreator() Property                 { return &StringProperty{decHex: []byte{ServerReference}} }
func PropReasonStringCreator() Property                    { return &StringProperty{decHex: []byte{ReasonString}} }
func PropReceiveMaximumCreator() Property                  { return &Uint16Property{decHex: []byte{ReceiveMaximum}} }
func PropTopicAliasMaximumCreator() Property               { return &Uint16Property{decHex: []byte{TopicAliasMaximum}} }
func PropTopicAliasCreator() Property                      { return &Uint16Property{decHex: []byte{TopicAlias}} }
func PropMaximumQoSCreator() Property                      { return &ByteProperty{decHex: []byte{MaximumQoS}} }
func PropRetainAvailableCreator() Property                 { return &ByteProperty{decHex: []byte{RetainAvailable}} }
func PropUserPropertyCreator() Property                    { return &StringPairProperty{decHex: []byte{UserProperty}} }
func PropMaximumPacketSizeCreator() Property               { return &Uint32Property{decHex: []byte{MaximumPacketSize}} }
func PropWildcardSubscriptionAvailableCreator() Property   { return &ByteProperty{decHex: []byte{WildcardSubscriptionAvailable}} }
func PropSubscriptionIdentifierAvailableCreator() Property { return &ByteProperty{decHex: []byte{SubscriptionIdentifierAvailable}} }
func PropSharedSubscriptionAvailableCreator() Property     { return &ByteProperty{decHex: []byte{SharedSubscriptionAvailable}} }

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
    n1, err := w.Write(prop.DecHex())
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
    decHex []byte
    V      uint16
}

func (prop *Uint16Property) DecHex() []byte {
    return prop.decHex
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

func (prop *Uint16Property)Set(v interface{}) {
    prop.V = v.(uint16)
}

func (prop *Uint16Property)Get() interface{} {
    return prop.V
}

type Uint32Property struct {
    decHex []byte
    V      uint32
}

func (prop *Uint32Property) DecHex() []byte {
    return prop.decHex
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
    decHex []byte
    V      byte
}

func (prop *ByteProperty) DecHex() []byte {
    return prop.decHex
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
    decHex []byte
    V      VarInt
}

func (prop *VarIntProperty) DecHex() []byte {
    return prop.decHex
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
    decHex []byte
    V      String
}

func (prop *StringProperty) DecHex() []byte {
    return prop.decHex
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
    decHex []byte
    V      [2]String
}

func (prop *StringPairProperty) DecHex() []byte {
    return prop.decHex
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

type PropOpt func(io.Writer) error

func BuildProperty(w io.Writer, opts ...PropOpt) (err error) {
    buf := bytes.NewBuffer(make([]byte, 1024))
    for _, opt := range opts {
        err = opt(buf)
        if err != nil {
            return err
        }
    }
    length := buf.Len()
    data := buf.Bytes()
    varBuf := make([]byte, MaxVarUintBufSize)
    n := EncodeVaruint(varBuf, uint64(length))
    _, err = w.Write(varBuf[:n])
    if err != nil {
        return err
    }
    _, err = w.Write(data)
    return
}

func SetSessionExpiryInterval(v uint32) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 5)
        b[0] = SessionExpiryInterval
        binary.BigEndian.PutUint32(b[1:], v)
        _, err := w.Write(b)
        return err
    }
}

func SetReceiveMaximum(v uint16) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 3)
        b[0] = ReceiveMaximum
        binary.BigEndian.PutUint16(b[1:], v)
        _, err := w.Write(b)
        return err
    }
}

func SetMaximumPacketSize(v uint32) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 5)
        b[0] = MaximumPacketSize
        binary.BigEndian.PutUint32(b[1:], v)
        _, err := w.Write(b)
        return err
    }
}

func SetTopicAliasMaximum(v uint16) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 3)
        b[0] = TopicAliasMaximum
        binary.BigEndian.PutUint16(b[1:], v)
        _, err := w.Write(b)
        return err
    }
}

func SetRequestResponseInformation(v byte) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 2)
        b[0] = RequestResponseInformation
        b[1] = v
        _, err := w.Write(b)
        return err
    }
}

func SetRequestProblemInformation(v byte) PropOpt {
    return func(w io.Writer) error {
        b := make([]byte, 2)
        b[0] = RequestProblemInformation
        b[1] = v
        _, err := w.Write(b)
        return err
    }
}

func SetUserProperty(v map[string]string) PropOpt {
    return func(w io.Writer) error {
        _, err := w.Write([]byte{UserProperty})
        if err != nil {
            return err
        }
        for k, v := range v {
            _, err1 := EncodeString(w, k)
            if err1 != nil {
                return err1
            }
            _, err2 := EncodeString(w, v)
            if err2 != nil {
                return err2
            }
        }
        return nil
    }
}

func SetAuthenticationMethod(v string) PropOpt {
    return func(w io.Writer) error {
        _, err := w.Write([]byte{AuthenticationMethod})
        if err != nil {
            return err
        }
        _, err1 := EncodeString(w, v)
        return err1
    }
}

func SetAuthenticationData(v []byte) PropOpt {
    return func(w io.Writer) error {
        _, err := w.Write([]byte{AuthenticationData})
        if err != nil {
            return err
        }
        _, err1 := w.Write(v)
        return err1
    }
}

func ReadProperty(r io.Reader) ([]Property, int, error) {
    v := NewFromReader(r)
    if v == nil {
        return nil, 0, errcode.ParseVarIntFailed
    }
    length := int(v.ToInt())
    var size int = 0
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

func WriteProperty(w io.Writer, props []Property) (int, error) {
    v := VarInt{}
    propLen := 0
    for _, v := range props {
        propLen += len(v.DecHex()) + int(v.DataLen())
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
