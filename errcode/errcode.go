// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package errcode

type Reason struct {
    Msg  string `json:"msg"`
    Code byte   `json:"code"`
}

func (r *Reason) Error() string {
    return r.Msg
}

var (
    StringOutOfRange        = &Reason{Msg: "String is out of range:Max is 65535 ", Code: ReasonProtocolError}
    ParseVarIntFailed       = &Reason{Msg: "Parse VarInt error", Code: ReasonProtocolError}
    UnknownProperty         = &Reason{Msg: "Unknown Property error", Code: ReasonProtocolError}
    ProtocolNameError       = &Reason{Msg: "Protocol Name error", Code: ReasonProtocolError}
    MessageNotSupport       = &Reason{Msg: "Message Not Support", Code: ReasonProtocolError}
    MessageReadSizeNotMatch = &Reason{Msg: "Message Read Size Not Match", Code: ReasonProtocolError}

    NormalDisconnection                 = &Reason{Msg: "Normal disconnection", Code: ReasonNormalDisconnection}
    GrantedQoS0                         = &Reason{Msg: "Granted QoS 0", Code: ReasonGrantedQoS0}
    GrantedQoS1                         = &Reason{Msg: "Granted QoS 1", Code: ReasonGrantedQoS1}
    GrantedQoS2                         = &Reason{Msg: "Granted QoS 2", Code: ReasonGrantedQoS2}
    DisconnectWithWillMessage           = &Reason{Msg: "Disconnect with WillMessage", Code: ReasonDisconnectWithWillMessage}
    NoMatchingSubscribers               = &Reason{Msg: "No Matching Subscribers", Code: ReasonNoMatchingSubscribers}
    NoSubscriptionExisted               = &Reason{Msg: "No Subscription Existed", Code: ReasonNoSubscriptionExisted}
    ContinueAuthentication              = &Reason{Msg: "Continue Authentication", Code: ReasonContinueAuthentication}
    Reauthenticate                      = &Reason{Msg: "Reauthenticate", Code: ReasonReauthenticate}
    UnspecifiedError                    = &Reason{Msg: "Unspecified Error", Code: ReasonUnspecifiedError}
    MalformedPacket                     = &Reason{Msg: "Malformed Packet", Code: ReasonMalformedPacket}
    ProtocolError                       = &Reason{Msg: "Protocol Error", Code: ReasonProtocolError}
    ImplementationSpecificError         = &Reason{Msg: "Implementation Specific Error", Code: ReasonImplementationSpecificError}
    UnsupportedProtocolVersion          = &Reason{Msg: "Unsupported Protocol Version", Code: ReasonUnsupportedProtocolVersion}
    ClientIdentifierNotValid            = &Reason{Msg: "Client Identifier Not Valid", Code: ReasonClientIdentifierNotValid}
    BadUserNameOrPassword               = &Reason{Msg: "Bad UserName Or Password", Code: ReasonBadUserNameOrPassword}
    NotAuthorized                       = &Reason{Msg: "Not Authorized", Code: ReasonNotAuthorized}
    ServerUnavailable                   = &Reason{Msg: "Server Unavailable", Code: ReasonServerUnavailable}
    ServerBusy                          = &Reason{Msg: "Server Busy", Code: ReasonServerBusy}
    Banned                              = &Reason{Msg: "Banned", Code: ReasonBanned}
    ServerShuttingDown                  = &Reason{Msg: "Server Shutting Down", Code: ReasonServerShuttingDown}
    BadAuthenticationMethod             = &Reason{Msg: "Bad Authentication Method", Code: ReasonBadAuthenticationMethod}
    KeepAliveTimeout                    = &Reason{Msg: "Keep Alive Timeout", Code: ReasonKeepAliveTimeout}
    SessionTakenOver                    = &Reason{Msg: "Session TakenOver", Code: ReasonSessionTakenOver}
    TopicFilterInvalid                  = &Reason{Msg: "Topic Filter Invalid", Code: ReasonTopicFilterInvalid}
    TopicNameInvalid                    = &Reason{Msg: "Topic Name Invalid", Code: ReasonTopicNameInvalid}
    PacketIdentifierInUse               = &Reason{Msg: "Packet Identifier InUse", Code: ReasonPacketIdentifierInUse}
    PacketIdentifierNotFound            = &Reason{Msg: "Packet Identifier Not Found", Code: ReasonPacketIdentifierNotFound}
    ReceiveMaximumExceeded              = &Reason{Msg: "Receive Maximum Exceeded", Code: ReasonReceiveMaximumExceeded}
    TopicAliasInvalid                   = &Reason{Msg: "Topic Alias Invalid", Code: ReasonTopicAliasInvalid}
    PacketTooLarge                      = &Reason{Msg: "Packet Too Large", Code: ReasonPacketTooLarge}
    MessageRateTooHigh                  = &Reason{Msg: "Message Rate Too High", Code: ReasonMessageRateTooHigh}
    QuotaExceeded                       = &Reason{Msg: "Quota Exceeded", Code: ReasonQuotaExceeded}
    AdministrativeAction                = &Reason{Msg: "Administrative Action", Code: ReasonAdministrativeAction}
    PayloadFormatInvalid                = &Reason{Msg: "Payload Format Invalid", Code: ReasonPayloadFormatInvalid}
    RetainNotSupported                  = &Reason{Msg: "Retain Not Supported", Code: ReasonRetainNotSupported}
    QoSNotSupported                     = &Reason{Msg: "QoS Not Supported", Code: ReasonQoSNotSupported}
    UseAnotherServer                    = &Reason{Msg: "Use Another Server", Code: ReasonUseAnotherServer}
    ServerMoved                         = &Reason{Msg: "Server Moved", Code: ReasonServerMoved}
    SharedSubscriptionsNotSupported     = &Reason{Msg: "Shared Subscriptions Not Supported", Code: ReasonSharedSubscriptionsNotSupported}
    ConnectionRateExceeded              = &Reason{Msg: "Connection Rate Exceeded", Code: ReasonConnectionRateExceeded}
    MaximumConnectTime                  = &Reason{Msg: "Maximum Connect Time", Code: ReasonMaximumConnectTime}
    SubscriptionIdentifiersNotSupported = &Reason{Msg: "Subscription Identifiers Not Supported", Code: ReasonSubscriptionIdentifiersNotSupported}
    WildcardSubscriptionsNotSupported   = &Reason{Msg: "Wildcard Subscriptions Not Supported", Code: ReasonWildcardSubscriptionsNotSupported}
)
