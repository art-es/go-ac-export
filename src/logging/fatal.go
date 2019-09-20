package logging

import (
	"log"
	"reflect"
)

const (
	ErrorInMarshaling      = "error in marshaling"
	ErrorCreatingRequestTo = "error in creating request to"
	ErrorSendingRequestTo  = "error sending request to"
)

func ErrorCreatingRequest(uri string, err error) {
	log.Fatalf("%s %s\nmessage: %s", ErrorCreatingRequestTo, uri, err.Error())
}

func ErrorMarshal(marshalStruct interface{}, err error) {
	t := reflect.TypeOf(marshalStruct)
	var structName string

	if t.Kind() == reflect.String {
		structName = reflect.ValueOf(marshalStruct).String()
	} else {
		structName = t.String()
	}

	log.Fatalf("%s %s\nmessage: %s", ErrorInMarshaling, structName, err.Error())
}

func ErrorSendingRequest(uri string, err error) {
	log.Fatalf("%s %s\nmessage: %s", ErrorSendingRequestTo, uri, err.Error())
}
