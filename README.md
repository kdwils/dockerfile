# Dockerfile

Dockerfile is a package that provides a methods on parsing and updating dockerfile commands. The initial intent of this package was to provide an interface to update base image tags of dockerfiles in a programmatic way. This could be used for automated updates to fix image vulnerabilites, mass image updates, etc.

Behind the scenes, this package is wrapped around buildkit's dockerfile parser, and provides methods to make modifying a dockerfile easier.

# Usage

Example of how to parse the contents of a dockerfile
```golang
import (
    "log"
    "os"

    "github.com/kdwils/dockerfile"
)

f, err := os.ReadFile("Dockerfile")
if err != nil {
    // handle err
}

d, err := dockerfile.ParseFromReader(f)
if err != nil {
    // handle err
}

err = d.WriteContents(os.Stdout)
if err != nil {
    // handler err
}
```