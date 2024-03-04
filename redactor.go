package bishamon

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	ErrInvalidRedactorArgs = errors.New("invalid args")
)

type (
	RedactFieldFunc func(message protoreflect.Message, fd protoreflect.FieldDescriptor) error
	RedactMapFunc   func(valueMap protoreflect.Map, mapKey protoreflect.MapKey, value protoreflect.Value) error
	RedactListFunc  func(list protoreflect.List) error

	FieldsFromMapFunc func(ext any) map[string]struct{}

	RedactorFuncs struct {
		redactFieldsFuncs []RedactFieldFunc
		redactMapFuncs    []RedactMapFunc
		redactListFuncs   []RedactListFunc
	}

	Redactor struct {
		extInfo           *protoimpl.ExtensionInfo
		fieldsFromMapFunc FieldsFromMapFunc
		redactorFuncs     RedactorFuncs
	}

	RedactorOption func(*Redactor)
)

func NewRedactor(
	extInfo *protoimpl.ExtensionInfo,
	opts ...RedactorOption,
) (*Redactor, error) {
	r := &Redactor{
		extInfo: extInfo,
	}

	for _, opt := range opts {
		opt(r)
	}

	if !r.isValid() {
		return nil, ErrInvalidRedactorArgs
	}

	return r, nil
}

func NewClearRedactor(extInfo *protoimpl.ExtensionInfo, opts ...RedactorOption) (*Redactor, error) {
	return NewRedactor(extInfo, append(opts,
		WithRedactFuncs(ClearFieldFunc),
		WithRedactMapFuncs(ClearMapFunc),
		WithRedactListFuncs(ClearListFunc),
	)...)
}

func (r *Redactor) Redact(msg proto.Message) error {
	return r.redact(msg)
}

func (r *Redactor) RedactClone(msg proto.Message) (proto.Message, error) {
	clone := proto.Clone(msg)

	if err := r.redact(clone); err != nil {
		return nil, fmt.Errorf("redact: %w", err)
	}

	return clone, nil
}

func (r *Redactor) redact(msg proto.Message) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic: %v", p)
		}
	}()

	if err = protorange.Range(msg.ProtoReflect(), func(p protopath.Values) error {
		fieldDescriptor := p.Path.Index(-1).FieldDescriptor()
		if fieldDescriptor == nil {
			return nil
		}

		value := p.Index(-1).Value

		if !value.IsValid() {
			return nil
		}

		opts, ok := fieldDescriptor.Options().(*descriptorpb.FieldOptions)
		if !ok || !proto.HasExtension(opts, r.extInfo) {
			return nil
		}

		switch {
		case fieldDescriptor.IsMap():
			return r.processMap(value.Map(), opts, r.extInfo)
		case fieldDescriptor.IsList():
			return r.processList(value.List())
		default:
			return r.processMessage(p, fieldDescriptor)
		}
	}); err != nil {
		return fmt.Errorf("range: %w", err)
	}

	return nil
}

func (r *Redactor) processList(list protoreflect.List) error {
	if len(r.redactorFuncs.redactListFuncs) == 0 {
		return nil
	}

	for _, redactFunc := range r.redactorFuncs.redactListFuncs {
		if err := redactFunc(list); err != nil {
			return fmt.Errorf("redactFunc: %w", err)
		}
	}

	return nil
}

func (r *Redactor) processMessage(
	p protopath.Values,
	fieldDescriptor protoreflect.FieldDescriptor,
) error {
	parent := p.Index(-2)
	if !parent.Value.IsValid() {
		return nil
	}

	for _, redactFunc := range r.redactorFuncs.redactFieldsFuncs {
		if err := redactFunc(parent.Value.Message(), fieldDescriptor); err != nil {
			return fmt.Errorf("redactFunc: %w", err)
		}
	}

	return nil
}

func (r *Redactor) processMap(
	value protoreflect.Map,
	opts *descriptorpb.FieldOptions,
	extInfo *protoimpl.ExtensionInfo,
) error {
	if r.fieldsFromMapFunc == nil ||
		len(r.redactorFuncs.redactMapFuncs) == 0 ||
		!value.IsValid() ||
		opts == nil {
		return nil
	}

	keys := r.fieldsFromMapFunc(proto.GetExtension(opts, extInfo))
	if len(keys) == 0 {
		return nil
	}

	var err error

	value.Range(func(key protoreflect.MapKey, v protoreflect.Value) bool {
		if _, ok := keys[key.String()]; !ok {
			return true
		}

		for _, redactFunc := range r.redactorFuncs.redactMapFuncs {
			if err = redactFunc(value, key, v); err != nil {
				return false
			}
		}

		return true
	})

	return err
}

func (r *Redactor) isValid() bool {
	if (len(r.redactorFuncs.redactFieldsFuncs) == 0 &&
		len(r.redactorFuncs.redactMapFuncs) == 0 &&
		len(r.redactorFuncs.redactListFuncs) == 0) ||
		r.extInfo == nil {
		return false
	}

	return true
}
