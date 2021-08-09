package api

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func parseRequest(msgName string, data []byte) (proto.Message, error) {
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(msgName))
	if err != nil {
		return nil, err
	}

	msg := msgType.New().Interface()
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// func parseRequestV1(msgName protoreflect.FullName, data []byte) (proto.Message, error) {
// 	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	msg := proto.MessageV1(msgType.New())
// 	err = proto.Unmarshal(data, msg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return msg, nil
// }
