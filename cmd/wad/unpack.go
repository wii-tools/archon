package wad

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func Unpack(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, c.String("in"))
	fmt.Fprintf(c.App.Writer, c.String("out"))
	return nil
}
