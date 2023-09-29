package dockerfile

import (
	"io"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func ParseFromReader(reader io.Reader) (*Dockerfile, error) {
	r, err := parser.Parse(reader)
	if err != nil {
		return nil, err
	}

	dockerfile := new(Dockerfile)
	for _, c := range r.AST.Children {
		cmd := &Command{
			Line:      c.Original,
			Command:   c.Value,
			StartLine: c.StartLine,
			EndLine:   c.EndLine,
			Flags:     c.Flags,
			Values:    make([]string, 0),
		}

		for n := c.Next; n != nil; n = n.Next {
			cmd.Values = append(cmd.Values, n.Value)
		}

		dockerfile.Commands = append(dockerfile.Commands, cmd)
	}

	return dockerfile, nil
}
