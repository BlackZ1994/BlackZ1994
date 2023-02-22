package httputil

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		url    string
		method string
		param  interface{}
		header map[string]string
	}
	tests := []struct {
		name string
		args args
		want *RestClient
	}{
		{"restClient1", args{"https://doSomething1", http.MethodGet, nil, nil}, &RestClient{url: "https://doSomething1", method: http.MethodGet, param: nil, header: nil}},
		{"restClient1", args{"https://doSomething2", http.MethodPost, nil, nil}, &RestClient{url: "https://doSomething2", method: http.MethodPost, param: nil, header: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.url, tt.args.method, tt.args.param, tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
