package main

import (
	"net/http"
	"testing"
)

func Test_getEnvConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		getEnvConfig()
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

func Test_getCityUniqueLink(t *testing.T) {
	type args struct {
		city    string
		country string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := getCityUniqueLink(tt.args.city, tt.args.country)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. getCityUniqueLink() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. getCityUniqueLink() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetWeatherConditions(t *testing.T) {
	type args struct {
		link       string
		dateString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := GetWeatherConditions(tt.args.link, tt.args.dateString)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GetWeatherConditions() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. GetWeatherConditions() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
