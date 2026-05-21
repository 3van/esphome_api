package api

import (
	"fmt"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Request and response types from/to esphome
const UndefinedTypeID uint64 = 0

var (
	apiMessageRegistryOnce sync.Once
	apiMessageTypesByID    map[uint64]protoreflect.MessageType
	apiMessageIDsByName    map[protoreflect.FullName]uint64
)

func TypeID(message interface{}) uint64 {
	if message == nil {
		return UndefinedTypeID
	}

	value := reflect.ValueOf(message)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return UndefinedTypeID
	}

	_, idsByName := apiMessageRegistry()
	if protoMessage, ok := message.(proto.Message); ok {
		if id := idsByName[protoMessage.ProtoReflect().Descriptor().FullName()]; id != 0 {
			return id
		}
		return UndefinedTypeID
	}

	messageType := value.Type()
	for messageType.Kind() == reflect.Ptr {
		messageType = messageType.Elem()
	}
	if messageType.PkgPath() != reflect.TypeOf(HelloRequest{}).PkgPath() {
		return UndefinedTypeID
	}

	if id := idsByName[apiMessageFullName(messageType.Name())]; id != 0 {
		return id
	}
	return UndefinedTypeID
}

func NewMessageByTypeID(typeID uint64) proto.Message {
	messagesByID, _ := apiMessageRegistry()
	messageType, ok := messagesByID[typeID]
	if !ok {
		return nil
	}
	return messageType.New().Interface()
}

func apiMessageRegistry() (map[uint64]protoreflect.MessageType, map[protoreflect.FullName]uint64) {
	apiMessageRegistryOnce.Do(func() {
		apiMessageTypesByID = make(map[uint64]protoreflect.MessageType)
		apiMessageIDsByName = make(map[protoreflect.FullName]uint64)

		messages := File_api_proto.Messages()
		for i := 0; i < messages.Len(); i++ {
			registerAPIMessage(messages.Get(i))
		}
	})
	return apiMessageTypesByID, apiMessageIDsByName
}

func registerAPIMessage(desc protoreflect.MessageDescriptor) {
	if id := apiMessageID(desc); id != 0 {
		messageType, err := protoregistry.GlobalTypes.FindMessageByName(desc.FullName())
		if err != nil {
			panic(fmt.Sprintf("api message %s with id %d is not registered: %v", desc.FullName(), id, err))
		}
		if existing, ok := apiMessageTypesByID[id]; ok && existing.Descriptor().FullName() != desc.FullName() {
			panic(fmt.Sprintf("api message id %d is used by both %s and %s", id, existing.Descriptor().FullName(), desc.FullName()))
		}
		apiMessageTypesByID[id] = messageType
		apiMessageIDsByName[desc.FullName()] = id
	}

	nested := desc.Messages()
	for i := 0; i < nested.Len(); i++ {
		registerAPIMessage(nested.Get(i))
	}
}

func apiMessageID(desc protoreflect.MessageDescriptor) uint64 {
	options, ok := desc.Options().(*descriptorpb.MessageOptions)
	if !ok || !proto.HasExtension(options, E_Id) {
		return UndefinedTypeID
	}

	value := proto.GetExtension(options, E_Id)
	switch id := value.(type) {
	case uint32:
		return uint64(id)
	case *uint32:
		if id != nil {
			return uint64(*id)
		}
	}
	return UndefinedTypeID
}

func apiMessageFullName(name string) protoreflect.FullName {
	if pkg := File_api_proto.Package(); pkg != "" {
		return pkg.Append(protoreflect.Name(name))
	}
	return protoreflect.FullName(name)
}
