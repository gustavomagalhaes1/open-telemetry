package cep

import (
	"reflect"
	"testing"
)

func TestGetAddressFromViaCEP(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name    string
		args    args
		want    *AddressResponse
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				cep: "01001000",
			},
			want: &AddressResponse{
				CEP:         "01001-000",
				Logradouro:  "Praça da Sé",
				Complemento: "lado ímpar",
				Bairro:      "Sé",
				Localidade:  "São Paulo",
				UF:          "SP",
				Erro:        false,
			},
			wantErr: false,
		},
		{
			name: "Test 2",
			args: args{
				cep: "99999999",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 3",
			args: args{
				cep: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAddressFromViaCEP(tt.args.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressFromViaCEP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressFromViaCEP() = %v, want %v", got, tt.want)
			}
		})
	}
}
