package bishamon

import (
	lilith "github.com/nenormalka/lilith/methods"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ClearFieldFunc = func(message protoreflect.Message, fieldDescriptor protoreflect.FieldDescriptor) error {
		message.Clear(fieldDescriptor)
		return nil
	}

	ClearMapFunc = func(valueMap protoreflect.Map, mapKey protoreflect.MapKey, _ protoreflect.Value) error {
		valueMap.Clear(mapKey)
		return nil
	}

	ClearListFunc = func(list protoreflect.List) error {
		list.Truncate(0)
		return nil
	}

	CommonFieldsFromMapFunc = func(ext any) map[string]struct{} {
		extNeeded, ok := ext.(interface {
			GetMapKeysToRedact() []string
		})
		if !ok || len(extNeeded.GetMapKeysToRedact()) == 0 {
			return nil
		}

		return lilith.ArrayToMapValues(extNeeded.GetMapKeysToRedact())
	}
)
