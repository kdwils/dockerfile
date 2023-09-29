package dockerfile

import (
	"errors"
	"regexp"
)

const (
	ImagePathRegex = `^(?:(?P<hostname>[^:/]+)(?::(?P<port>\d+))\/)?(?:(?P<username>[^/]+)\/)?(?:(?P<repository>[^:]+)(?::(?P<tag>[^/]+))?|(?P<ID>[^/]+))$`
)

// Image represents the components of a docker image path
type Image struct {
	Registry   string
	Port       string
	Username   string
	Repository string
	Tag        string
	ID         string
}

// ParseImage parses a full Docker image path into its individual components
func ParseImage(path string) (*Image, error) {
	regex, err := regexp.Compile(ImagePathRegex)
	if err != nil {
		return nil, err
	}

	match := regex.FindStringSubmatch(path)
	if match == nil {
		return nil, errors.New("invalid image path: match is nil")
	}
	if len(match) == 0 {
		return nil, errors.New("invalid image path: no matches found")
	}

	return &Image{
		Registry:   match[1],
		Port:       match[2],
		Username:   match[3],
		Repository: match[4],
		Tag:        match[5],
		ID:         match[6],
	}, nil
}

// Path constructs a full Docker image path from its components.
func (i *Image) Path() string {
	if i.ID != "" {
		return i.ID
	}

	path := ""
	if i.Registry != "" {
		path += i.Registry
		if i.Port != "" {
			path += ":" + i.Port
		}
		path += "/"
	}

	if i.Username != "" {
		path += i.Username + "/"
	}

	path += i.Repository
	if i.Tag != "" {
		path += ":" + i.Tag
	}

	return path
}
