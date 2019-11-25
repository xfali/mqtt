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

//type Property struct {
//    DecHex []byte
//    Data   []byte
//}

type Property interface {
    DecHex() []byte
    DataLen() int32
    UnmarshalData(r io.Reader) error
    MarshalData(w io.Writer) error
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

func UnmarshalProp(r io.Reader) (prop Property, err error) {
    b := make([]byte, 1)
    _, err = io.ReadFull(r, b)
    if err != nil {
        return nil, err
    }
    prop = CreateProperty(b[0])
    if prop == nil {
        return nil, errcode.UnknownProperty
    }

    err = prop.UnmarshalData(r)
    if err != nil {
        return nil, err
    }
    return
}

func MarshalProp(w io.Writer, prop Property) error {
    _, err := w.Write(prop.DecHex())
    if err != nil {
        return err
    }

    return prop.MarshalData(w)
}

func CreateProperty(t byte) Property {
    if v, ok := propFac[t]; ok {
        return v()
    }
    return nil
}

type Uint16Property struct {
    decHex []byte
    v      uint16
}

func (prop *Uint16Property) DecHex() []byte {
    return prop.decHex
}

func (prop *Uint16Property) DataLen() int32 {
    return 2
}

func (prop *Uint16Property) UnmarshalData(r io.Reader) error {
    buf := make([]byte, 2)
    _, err := io.ReadFull(r, buf)
    if err != nil {
        return err
    }

    prop.v = binary.BigEndian.Uint16(buf)
    return nil
}

func (prop *Uint16Property) MarshalData(w io.Writer) error {
    buf := make([]byte, 2)
    binary.BigEndian.PutUint16(buf, prop.v)
    _, err := w.Write(buf)
    return err
}

type Uint32Property struct {
    decHex []byte
    v      uint32
}

func (prop *Uint32Property) DecHex() []byte {
    return prop.decHex
}

func (prop *Uint32Property) DataLen() int32 {
    return 4
}

func (prop *Uint32Property) UnmarshalData(r io.Reader) error {
    buf := make([]byte, 4)
    _, err := io.ReadFull(r, buf)
    if err != nil {
        return err
    }

    prop.v = binary.BigEndian.Uint32(buf)
    return nil
}

func (prop *Uint32Property) MarshalData(w io.Writer) error {
    buf := make([]byte, 4)
    binary.BigEndian.PutUint32(buf, prop.v)
    _, err := w.Write(buf)
    return err
}

type ByteProperty struct {
    decHex []byte
    v      byte
}

func (prop *ByteProperty) DecHex() []byte {
    return prop.decHex
}

func (prop *ByteProperty) DataLen() int32 {
    return 1
}

func (prop *ByteProperty) UnmarshalData(r io.Reader) error {
    buf := make([]byte, 1)
    _, err := io.ReadFull(r, buf)
    if err != nil {
        return err
    }

    prop.v = buf[0]
    return nil
}

func (prop *ByteProperty) MarshalData(w io.Writer) error {
    _, err := w.Write([]byte{prop.v})
    return err
}

type VarIntProperty struct {
    decHex []byte
    v      VarInt
}

func (prop *VarIntProperty) DecHex() []byte {
    return prop.decHex
}

func (prop *VarIntProperty) DataLen() int32 {
    return int32(prop.v.Length())
}

func (prop *VarIntProperty) UnmarshalData(r io.Reader) error {
    _, err := prop.v.LoadFromReader(r)
    return err
}

func (prop *VarIntProperty) MarshalData(w io.Writer) error {
    _, err := w.Write(prop.v.Bytes())
    return err
}

type StringProperty struct {
    decHex []byte
    v      String
}

func (prop *StringProperty) DecHex() []byte {
    return prop.decHex
}

func (prop *StringProperty) DataLen() int32 {
    return int32(prop.v.AllLength())
}

func (prop *StringProperty) UnmarshalData(r io.Reader) error {
    s, err := ParseString(r)
    if err != nil {
        return err
    }
    prop.v = *s
    return nil
}

func (prop *StringProperty) MarshalData(w io.Writer) error {
    return WriteString(w, prop.v)
}

type StringPairProperty struct {
    decHex []byte
    v      [2]String
}

func (prop *StringPairProperty) DecHex() []byte {
    return prop.decHex
}

func (prop *StringPairProperty) DataLen() int32 {
    return int32(prop.v[0].AllLength() + prop.v[1].AllLength())
}

func (prop *StringPairProperty) UnmarshalData(r io.Reader) error {
    s1, err := ParseString(r)
    if err != nil {
        return err
    }
    prop.v[0] = *s1

    s2, err := ParseString(r)
    if err != nil {
        return err
    }
    prop.v[1] = *s2
    return nil
}

func (prop *StringPairProperty) MarshalData(w io.Writer) error {
    err := WriteString(w, prop.v[0])
    if err != nil {
        return err
    }
    return WriteString(w, prop.v[1])
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

func ReadProperty(r io.Reader) error {
    v := NewFromReader(r)
    if v == nil {
        return errcode.ParseVarIntFailed
    }
    length := v.ToInt()
    buf := make([]byte, length)
    _, err := r.Read(buf)
    if err != nil {
        return err
    }

    return nil
}
