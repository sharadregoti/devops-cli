package utils

import (
	structpb "github.com/golang/protobuf/ptypes/struct"
	structpbt "google.golang.org/protobuf/types/known/structpb"
)

func GetMap(data map[string]interface{}) map[string]*structpb.Value {
	m := make(map[string]*structpb.Value)
	for k, v := range data {
		newVal, err := structpbt.NewValue(v)
		if err != nil {
			return nil
		}
		m[k] = newVal
	}

	return m
}

func GetMapInterface(data map[string]*structpb.Value) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range data {
		m[k] = v.AsInterface()
	}

	return m
}
