package main

import "testing"

func Test_main(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Hello World",
			args: args{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
