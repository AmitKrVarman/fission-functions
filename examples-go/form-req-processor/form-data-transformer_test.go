package main

import (
	"net/http"
	"reflect"
	"testing"
)

func TestHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		Handler(tt.args.w, tt.args.r)
	}
}

func Test_transformData(t *testing.T) {
	type args struct {
		typeFormdata TypeFormData
	}
	tests := []struct {
		name string
		args args
		want TranformedData
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := transformData(tt.args.typeFormdata); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. transformData() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
