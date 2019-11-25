// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package errcode

import "errors"

var StringOutOfRange = errors.New("String is out of range:Max is 65535 ")
var ParseVarIntFailed = errors.New("Parse VarInt error ")
var UnknownProperty = errors.New("Unknown Property error ")
