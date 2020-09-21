package wad

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/wii-tools/wadlib"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Unpack(c *cli.Context) error {
	wad, err := wadlib.LoadWADFromFile(c.String("in"))
	if err != nil {
		return err
	}

	// TODO: more descriptive error if the output directory is an existing file?
	out := c.String("out")
	err = os.Mkdir(out, 0755)
	if os.IsExist(err) {
		// Check if this is a file.
		// We are okay with overwriting existing files.
		stat, err := os.Stat(out)
		if err != nil {
			return err
		}

		if !stat.IsDir() {
			return errors.New(fmt.Sprintf("%s is a directory", out))
		}
	} else if err != nil {
		// This isn't an error we know how to cope with.
		return err
	}

	titleId := fmt.Sprintf("%016x", wad.Ticket.TitleID)
	// Ensure that the given ticket has its title key encrypted.
	err = wad.Ticket.EncryptKey()
	if err != nil {
		return err
	}

	// Start writing!
	// We'll write the ticket, TMD, certificate chain, and meta section (if available).
	dir := newDirectory(out)
	err = dir.writeStruct(fmt.Sprintf("%s.tik", titleId), wad.Ticket)
	if err != nil {
		return err
	}

	err = dir.writeTMD(fmt.Sprintf("%s.tmd", titleId), wad.TMD)
	if err != nil {
		return err
	}

	err = dir.writeFile(fmt.Sprintf("%s.certs", titleId), wad.CertificateChain)
	if err != nil {
		return err
	}

	// Meta doesn't always exist.
	if len(wad.Meta) != 0 {
		err = dir.writeFile(fmt.Sprintf("%s.meta", titleId), wad.Meta)
		if err != nil {
			return err
		}
	}

	// We'll then write all contents listed.
	for _, content := range wad.Data {
		err = dir.writeContents(content)
		if err != nil {
			return err
		}
	}

	return nil
}

// Type directory helps us avoid a ton of nested functions.
type directory struct {
	dir string
}

func newDirectory(dir string) directory {
	return directory{dir: dir}
}

// writeFile writes binary contents of a struct to the named file within a directory.
func (d *directory) writeStruct(name string, given interface{}) error {
	// Write the struct to a buffer.
	var tmp bytes.Buffer
	err := binary.Write(&tmp, binary.BigEndian, given)
	if err != nil {
		return err
	}

	// Read the buffer's contents.
	contents, err := ioutil.ReadAll(&tmp)
	if err != nil {
		return err
	}

	// Handle as a byte array.
	return d.writeFile(name, contents)
}

// writeTMD writes a title's TMD within a directory.
func (d *directory) writeTMD(name string, given *wadlib.TMD) error {
	// First, handle the fixed-length BinaryTMD.
	var tmp bytes.Buffer
	err := binary.Write(&tmp, binary.BigEndian, given.BinaryTMD)
	if err != nil {
		return err
	}

	// Then, write all individual content records.
	for _, content := range given.Contents {
		err := binary.Write(&tmp, binary.BigEndian, content)
		if err != nil {
			return err
		}
	}

	// Read the buffer's contents.
	contents, err := ioutil.ReadAll(&tmp)
	if err != nil {
		return err
	}

	// Handle as a byte array.
	return d.writeFile(name, contents)
}

// writeContents writes contents to the content's index within a directory.
func (d *directory) writeContents(content wadlib.WADFile) error {
	name := fmt.Sprintf("%08x.app", content.Index)
	return d.writeFile(name, content.RawData)
}

// writeFile writes contents to the named file within a directory.
func (d *directory) writeFile(name string, contents []byte) error {
	path := filepath.Join(d.dir, name)
	return ioutil.WriteFile(path, contents, os.ModePerm)
}
