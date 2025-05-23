package game

import "io"

type OutputJob interface {
	Run(w io.Writer) error
}
