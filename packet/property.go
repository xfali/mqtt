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
    Unmarshal(r io.Reader, v interface{}) error
    Marshal(w io.Writer, v interface{}) error
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
