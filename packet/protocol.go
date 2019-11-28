// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

//4 - MQTT
var (
    MqttProtocolVersion = byte(5)
    MqttProtocolNameBytes  = []byte{0, 4, 77, 81, 84, 84}
    MqttProtocolName       = "MQTT"
    MqttProtocolNameString = String{
        length: 4,
        data: []byte(MqttProtocolName),
    }
)
