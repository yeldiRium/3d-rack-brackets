package globals

import "io"

type Globals struct {
	Debug bool
	Stdout io.Writer
}
