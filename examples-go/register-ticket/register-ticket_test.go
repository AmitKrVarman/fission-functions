package main

import (
	"net/http"
	"testing"
)

func Test_getEnvDatabaseConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		getEnvDatabaseConfig()
	}
}

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
