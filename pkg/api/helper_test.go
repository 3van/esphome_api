package api

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestTypeIDRegistryMatchesProtoIDOptions(t *testing.T) {
	messages := File_api_proto.Messages()
	seen := 0
	for i := 0; i < messages.Len(); i++ {
		desc := messages.Get(i)
		id := testMessageID(t, desc)
		if id == 0 {
			continue
		}
		seen++

		message := NewMessageByTypeID(id)
		if message == nil {
			t.Fatalf("message %s with id %d is not registered", desc.FullName(), id)
		}
		if got := message.ProtoReflect().Descriptor().FullName(); got != desc.FullName() {
			t.Fatalf("id %d resolved to %s, want %s", id, got, desc.FullName())
		}
		if got := TypeID(message); got != id {
			t.Fatalf("TypeID(%s) = %d, want %d", desc.FullName(), got, id)
		}
	}
	if seen == 0 {
		t.Fatal("api proto has no message id options")
	}
}

func testMessageID(t *testing.T, desc protoreflect.MessageDescriptor) uint64 {
	t.Helper()

	options, ok := desc.Options().(*descriptorpb.MessageOptions)
	if !ok || !proto.HasExtension(options, E_Id) {
		return 0
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
	t.Fatalf("message %s has unexpected id extension type %T", desc.FullName(), value)
	return 0
}
