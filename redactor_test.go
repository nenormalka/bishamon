package bishamon

import (
	"testing"

	"bishamon/data"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestNewClearRedactorClone(t *testing.T) {
	redactor, err := NewClearRedactor(data.E_Sensitive, WithFieldsFromMapFunc(CommonFieldsFromMapFunc))
	if err != nil {
		t.Fatalf("failed to create redactor: %v", err)
	}

	msg := &data.TestMessage{
		Password: "password_test",
		Login:    "login_test",
		Contacts: map[string]string{
			"email": "email_test",
			"phone": "phone_test",
			"addr":  "addr_test",
			"city":  "city_test",
		},
		FollowIds: []string{"1", "2", "3"},
	}

	cloned, err := redactor.RedactClone(msg)
	if err != nil {
		t.Fatalf("failed to redact: %v", err)
	}

	clonedMsg, ok := cloned.(*data.TestMessage)
	if !ok {
		t.Fatalf("failed to cast to TestMessage")
	}

	if password := clonedMsg.GetPassword(); password != "" {
		t.Fatalf("password is not redacted")
	}

	if login := clonedMsg.GetLogin(); login != msg.GetLogin() {
		t.Fatalf("login is redacted")
	}

	if len(clonedMsg.GetContacts()) == len(msg.GetContacts()) {
		t.Fatalf("contacts are not redacted")
	}

	for key := range clonedMsg.GetContacts() {
		if key == "phone" || key == "email" {
			t.Fatalf("%s is not redacted", key)
		}
	}

	if len(clonedMsg.GetFollowIds()) != 0 {
		t.Fatalf("ids are not redacted")
	}
}

func TestNewClearRedactor(t *testing.T) {
	redactor, err := NewClearRedactor(data.E_Sensitive, WithFieldsFromMapFunc(CommonFieldsFromMapFunc))
	if err != nil {
		t.Fatalf("failed to create redactor: %v", err)
	}

	msg := &data.TestMessage{
		Password: "password_test",
		Login:    "login_test",
		Contacts: map[string]string{
			"email": "email_test",
			"phone": "phone_test",
			"addr":  "addr_test",
			"city":  "city_test",
		},
		FollowIds: []string{"1", "2", "3"},
	}

	if err = redactor.Redact(msg); err != nil {
		t.Fatalf("failed to redact: %v", err)
	}

	if password := msg.GetPassword(); password != "" {
		t.Fatalf("password is not redacted")
	}

	if login := msg.GetLogin(); login == "" {
		t.Fatalf("login is redacted")
	}

	for key := range msg.GetContacts() {
		if key == "phone" || key == "email" {
			t.Fatalf("%s is not redacted", key)
		}
	}

	if len(msg.GetFollowIds()) != 0 {
		t.Fatalf("ids are not redacted")
	}
}

func TestNewCustomRedactorClone(t *testing.T) {
	redactor, err := NewRedactor(
		data.E_Sensitive,
		WithFieldsFromMapFunc(CommonFieldsFromMapFunc),
		WithRedactFuncs(func(message protoreflect.Message, fd protoreflect.FieldDescriptor) error {
			if fd.Kind() != protoreflect.StringKind {
				return nil
			}

			message.Set(fd, protoreflect.ValueOfString("*masked*"))

			return nil
		}),
		WithRedactMapFuncs(func(valueMap protoreflect.Map, mapKey protoreflect.MapKey, value protoreflect.Value) error {
			valueMap.Set(mapKey, protoreflect.ValueOfString("*masked*"))
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("failed to create redactor: %v", err)
	}

	msg := &data.TestMessage{
		Password: "password_test",
		Login:    "login_test",
		Contacts: map[string]string{
			"email": "email_test",
			"phone": "phone_test",
			"addr":  "addr_test",
			"city":  "city_test",
		},
		FollowIds: []string{"1", "2", "3"},
	}

	cloned, err := redactor.RedactClone(msg)
	if err != nil {
		t.Fatalf("failed to redact: %v", err)
	}

	clonedMsg, ok := cloned.(*data.TestMessage)
	if !ok {
		t.Fatalf("failed to cast to TestMessage")
	}

	if password := clonedMsg.GetPassword(); password != "*masked*" {
		t.Fatalf("password is not redacted")
	}

	if login := clonedMsg.GetLogin(); login != msg.GetLogin() {
		t.Fatalf("login is redacted")
	}

	for key, value := range clonedMsg.GetContacts() {
		if (key == "phone" || key == "email") && value != "*masked*" {
			t.Fatalf("%s is not redacted", key)
		}
	}
}

func TestNewCustomRedactor(t *testing.T) {
	redactor, err := NewRedactor(
		data.E_Sensitive,
		WithFieldsFromMapFunc(CommonFieldsFromMapFunc),
		WithRedactFuncs(func(message protoreflect.Message, fd protoreflect.FieldDescriptor) error {
			if fd.Kind() != protoreflect.StringKind {
				return nil
			}

			message.Set(fd, protoreflect.ValueOfString("*masked*"))

			return nil
		}),
		WithRedactMapFuncs(func(valueMap protoreflect.Map, mapKey protoreflect.MapKey, value protoreflect.Value) error {
			valueMap.Set(mapKey, protoreflect.ValueOfString("*masked*"))
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("failed to create redactor: %v", err)
	}

	msg := &data.TestMessage{
		Password: "password_test",
		Login:    "login_test",
		Contacts: map[string]string{
			"email": "email_test",
			"phone": "phone_test",
			"addr":  "addr_test",
			"city":  "city_test",
		},
		FollowIds: []string{"1", "2", "3"},
	}

	if err = redactor.Redact(msg); err != nil {
		t.Fatalf("failed to redact: %v", err)
	}

	if password := msg.GetPassword(); password != "*masked*" {
		t.Fatalf("password is not redacted")
	}

	if login := msg.GetLogin(); login != msg.GetLogin() {
		t.Fatalf("login is redacted")
	}

	for key, value := range msg.GetContacts() {
		if (key == "phone" || key == "email") && value != "*masked*" {
			t.Fatalf("%s is not redacted", key)
		}
	}
}
