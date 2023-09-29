package dockerfile

import (
	"reflect"
	"testing"
)

func TestParseImage(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Image
		wantErr bool
	}{
		{
			name: "parse image",
			args: args{
				path: "example.com:8080/myusername/myrepository:mytag",
			},
			want: &Image{
				Registry:   "example.com",
				Port:       "8080",
				Username:   "myusername",
				Repository: "myrepository",
				Tag:        "mytag",
			},
			wantErr: false,
		},
		{
			name: "no matches",
			args: args{
				path: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseImage(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Path(t *testing.T) {
	type fields struct {
		Registry   string
		Port       string
		Username   string
		Repository string
		Tag        string
		ID         string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				ID: "my-id",
			},
			name: "id present",
			want: "my-id",
		},
		{
			name: "all cases",
			fields: fields{
				Registry:   "ghcr.io",
				Username:   "kdwils",
				Repository: "dockerfile",
				Tag:        "1.2.3",
			},
			want: "ghcr.io/kdwils/dockerfile:1.2.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Image{
				Registry:   tt.fields.Registry,
				Port:       tt.fields.Port,
				Username:   tt.fields.Username,
				Repository: tt.fields.Repository,
				Tag:        tt.fields.Tag,
				ID:         tt.fields.ID,
			}
			if got := i.Path(); got != tt.want {
				t.Errorf("Image.Path() = %v, want %v", got, tt.want)
			}
		})
	}
}
