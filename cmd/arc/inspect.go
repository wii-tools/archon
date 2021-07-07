package arc

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/wii-tools/arclib"
)

func Inspect(c *cli.Context) error{
	// Create an empty U8 Archive
	arc := arclib.ARC{}

	err := arc.LoadFromFile(c.String("in"))
	if err != nil {
		return err
	}

	data := arc.Contents()
	num := len(data)

	fmt.Fprintf(c.App.Writer, "There are %d files in this U8 archive.\n\n", num)
	for _, contents := range data{
		size, err := arc.Read(contents)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.App.Writer, "File Path: %s | Size: %d\n",  contents, len(size))
	}

	return nil
}
