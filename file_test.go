package dockerfile

import (
	"reflect"
	"testing"
)

func TestCommand_String(t *testing.T) {
	type fields struct {
		Line      string
		Command   string
		Flags     []string
		StartLine int
		EndLine   int
		Values    []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "format COPY command",
			fields: fields{
				Command: "COPY",
				Flags:   []string{"--from=builder"},
				Values:  []string{"/some/path/test.txt", "."},
			},
			want: "COPY --from=builder /some/path/test.txt .",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				Line:      tt.fields.Line,
				Command:   tt.fields.Command,
				Flags:     tt.fields.Flags,
				StartLine: tt.fields.StartLine,
				EndLine:   tt.fields.EndLine,
				Values:    tt.fields.Values,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Command.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatLine(t *testing.T) {
	type args struct {
		cmd    string
		values []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "entrypoint",
			args: args{
				cmd:    "ENTRYPOINT",
				values: []string{"./my-binary", "run"},
			},
			want: `ENTRYPOINT ["./my-binary", "run"]`,
		},
		{
			name: "default",
			args: args{
				cmd:    "COPY",
				values: []string{"/my/path/test.txt", "."},
			},
			want: "COPY /my/path/test.txt .",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatLine(tt.args.cmd, tt.args.values); got != tt.want {
				t.Errorf("formatLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerfile_SetBaseImageTag(t *testing.T) {
	type fields struct {
		Commands []*Command
	}
	type args struct {
		tag string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantDockerfile *Dockerfile
	}{
		{
			name:    "no commands",
			wantErr: true,
			fields: fields{
				Commands: nil,
			},
			wantDockerfile: &Dockerfile{},
		},
		{
			name:    "no base image",
			wantErr: true,
			fields: fields{
				Commands: []*Command{
					{
						Command: "COPY",
					},
				},
			},
			args: args{
				tag: "1.2.3",
			},
			wantDockerfile: &Dockerfile{
				Commands: []*Command{
					{
						Command: "COPY",
					},
				},
			},
		},
		{
			name:    "base image present",
			wantErr: false,
			fields: fields{
				Commands: []*Command{
					{
						Command: "COPY",
					},
					{
						Command: "FROM",
						Values:  []string{"my-image:1.2.3 as something"},
					},
					{
						Command: "FROM",
						Values:  []string{"base-image:1.2.3"},
					},
				},
			},
			args: args{
				tag: "1.2.4",
			},
			wantDockerfile: &Dockerfile{
				Commands: []*Command{
					{
						Command: "COPY",
					},
					{
						Command: "FROM",
						Values:  []string{"my-image:1.2.3 as something"},
					},
					{
						Command: "FROM",
						Values:  []string{"base-image:1.2.4"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dockerfile{
				Commands: tt.fields.Commands,
			}
			err := d.SetBaseImageTag(tt.args.tag)
			if !reflect.DeepEqual(d, tt.wantDockerfile) {
				t.Errorf("Dockerfile.SetBaseImageTag() dockerfile = %v, want %v", d, tt.wantDockerfile)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Dockerfile.SetBaseImageTag() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}
