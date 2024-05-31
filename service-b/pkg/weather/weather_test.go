package weather

import (
	"reflect"
	"testing"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	type args struct {
		celsius float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Test 1",
			args: args{
				celsius: 0,
			},
			want: 32,
		},
		{
			name: "Test 2",
			args: args{
				celsius: 100,
			},
			want: 212,
		},
		{
			name: "Test 3",
			args: args{
				celsius: 37,
			},
			want: 98.60000000000001,
		},
		}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CelsiusToFahrenheit(tt.args.celsius); got != tt.want {
				t.Errorf("CelsiusToFahrenheit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	type args struct {
		celsius float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Test 1",
			args: args{
				celsius: 0,
			},
			want: 273,
		},
		{
			name: "Test 2",
			args: args{
				celsius: 100,
			},
			want: 373,
		},
		{
			name: "Test 3",
			args: args{
				celsius: 37,
			},
			want: 310,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CelsiusToKelvin(tt.args.celsius); got != tt.want {
				t.Errorf("CelsiusToKelvin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWeather(t *testing.T) {

	weatherResponse, _ := GetWeather("São Paulo")

	type args struct {
		city string
	}
	tests := []struct {
		name    string
		args    args
		want    *WeatherResponse
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				city: "São Paulo",
			},
			want: weatherResponse,
			wantErr: false,
		},
		{
			name: "Test 2",
			args: args{
				city: "",
			},
			want:    &WeatherResponse{
				Location: struct {
					Name string `json:"name"`
				}{
					Name: "",
				},
				Current: struct {
					TempC float64 `json:"temp_c"`
				}{
					TempC: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWeather(tt.args.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWeather() = %v, want %v", got, tt.want)
			}
		})
	}
}
