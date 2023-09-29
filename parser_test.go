package dockerfile

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFromReader(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Dockerfile
		wantErr bool
	}{
		{
			name: "parse from reader",
			args: args{
				reader: bytes.NewReader([]byte(`COPY --from=builder my-binary .`)),
			},
			want: &Dockerfile{
				Commands: []*Command{
					{
						Command:   "COPY",
						Line:      "COPY --from=builder my-binary .",
						Flags:     []string{"--from=builder"},
						Values:    []string{"my-binary", "."},
						EndLine:   1,
						StartLine: 1,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFromReader(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("ParseFromReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
