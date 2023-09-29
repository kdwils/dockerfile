package dockerfile

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	FROM        = "FROM"
	WORKDIR     = "WORKDIR"
	ARG         = "ARG"
	MAINTAINER  = "MAINTAINER"
	RUN         = "RUN"
	CMD         = "CMD"
	LABEL       = "LABEL"
	EXPOSE      = "EXPOSE"
	ENV         = "ENV"
	ADD         = "ADD"
	COPY        = "COPY"
	ENTRYPOINT  = "ENTRYPOINT"
	HEALTHCHECK = "HEALTHCHECK"
	SHELL       = "SHELL"
	STOPSIGNAL  = "STOPSIGNAL"
	ONBUILD     = "ONBUILD"
)

// Dockerfile represents the commands that make up the contents of a dockerfile
type Dockerfile struct {
	Commands []*Command
}

// Command represents a single line in a dockerfile
type Command struct {
	// Line is the original dockerfile line
	Line string
	// Command is the dockerfile command such as RUN, COPY, or ADD
	Command string
	// Flags represent the flags of the command, such as --from in a COPY statement
	Flags []string
	// StartLine is the line number the command begins on
	StartLine int
	// EndLine is the line number the command ends on
	EndLine int
	// Values are the arguments to a command following the Commmand and Flags
	Values []string
}

// String formats the docker command back to a string
func (c *Command) String() string {
	if c == nil {
		return ""
	}

	cmd := c.Command
	if len(c.Flags) != 0 {
		cmd = formatLine(cmd, c.Flags)
	}

	if len(c.Values) != 0 {
		cmd = formatLine(cmd, c.Values)
	}

	return cmd
}

func formatLine(cmd string, values []string) string {
	cmd = strings.ToUpper(cmd)

	if len(values) == 0 {
		return cmd
	}

	switch cmd {
	case ENTRYPOINT:
		args := make([]string, len(values))
		for i, v := range values {
			args[i] = fmt.Sprintf(`"%s"`, v)
		}
		return fmt.Sprintf("%s [%s]", cmd, strings.Join(args, ", "))
	default:
		return fmt.Sprintf("%s %s", cmd, strings.Join(values, " "))
	}
}

// SetBaseImageTag sets the tag to the image in the last FROM statement. It is assumed that the last FROM statement is the base image.
func (d *Dockerfile) SetBaseImageTag(tag string) error {
	for i := len(d.Commands) - 1; i > 0; i-- {
		if !strings.EqualFold(FROM, d.Commands[i].Command) {
			continue
		}

		if len(d.Commands[i].Values) == 0 {
			return errors.New("no values supplied to last FROM statement")
		}

		image, err := ParseImage(d.Commands[i].Values[0])
		if err != nil {
			return err
		}

		image.Tag = tag
		d.Commands[i].Values[0] = image.Path()
		return nil
	}

	return errors.New("no FROM statements found")
}

// WriteContents writes the commands of the dockerfile to the provided writer
func (d *Dockerfile) WriteContents(w io.Writer) error {
	var line = 1
	contents := make([]string, 0)

	for _, c := range d.Commands {
		for line < c.StartLine {
			contents = append(contents, "\n")
			line++
		}
		contents = append(contents, c.String())

		// handle cases where lines extended by a '\' are removed due to how the buildkit parser works
		// this will remove the extra new lines that are left over when the command is shortened into a single line
		if (c.EndLine - c.StartLine) > 1 {
			line = line + c.EndLine - c.StartLine + 1
		}
	}

	_, err := w.Write([]byte(strings.Join(contents, "")))
	return err
}
