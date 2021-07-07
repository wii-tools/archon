package arc

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/wii-tools/arclib"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Extract(c *cli.Context) error{
	// Create an empty U8 archive
	arc := arclib.ARC{}

	err := arc.LoadFromFile(c.String("in"))
	if err != nil {
		return err
	}

	// Create output directory if already not existing
	out := c.String("out")
	err = os.Mkdir(out, 0755)
	if os.IsExist(err) {
		// Check if this is a file.
		// We are okay with overwriting existing files within folders.
		stat, err := os.Stat(out)
		if err != nil {
			return err
		}

		if !stat.IsDir() {
			return errors.New(fmt.Sprintf("%s is not a directory", out))
		}
	} else if err != nil {
		// This isn't an error we know how to cope with.
		return err
	}


	for _, file := range arc.Contents() {
		data, err := arc.Read(file)
		if err != nil {
			return err
		}

		path := filepath.Join(out, file)
		last := strings.LastIndex(path, "/")
		mkdirPath := path[:last]

		err = os.MkdirAll(mkdirPath, 0700)
		if os.IsExist(err) {
			// We need to continue to create other parent directories
			continue
		}

		err = ioutil.WriteFile(path, data, 0777)
		if err != nil {
			return err
		}
	}

	fmt.Fprint(c.App.Writer, "Successfully extracted U8 archive!\n")

	return nil
}
