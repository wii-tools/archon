package wad

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/wii-tools/wadlib"
	"os"
)

func Unpack(c *cli.Context) error {
	wad, err := wadlib.LoadWADFromFile(c.String("in"))
	if err != nil {
		return err
	}

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

	// Start writing!
	// We'll write the ticket, TMD, certificate chain, and meta section (if available).
	dir := directory{
		dir:     out,
		titleId: fmt.Sprintf("%016x", wad.Ticket.TitleID),
	}
	ticket, err := wad.GetTicket()
	if err != nil {
		return err
	}
	err = dir.writeSection("tik", ticket)
	if err != nil {
		return err
	}

	tmd, err := wad.GetTMD()
	if err != nil {
		return err
	}
	err = dir.writeSection("tmd", tmd)
	if err != nil {
		return err
	}

	err = dir.writeSection("certs", wad.CertificateChain)
	if err != nil {
		return err
	}

	// Meta doesn't always exist.
	if len(wad.Meta) != 0 {
		err = dir.writeSection("meta", wad.Meta)
		if err != nil {
			return err
		}
	}

	// Nor does the CRL section, by default.
	if len(wad.Meta) != 0 {
		err = dir.writeSection("crl", wad.CertificateRevocationList)
		if err != nil {
			return err
		}
	}

	// We'll then write all contents listed.
	shouldDecrypt := !c.Bool("no-decrypt")
	for _, content := range wad.Data {
		// First, decrypt.
		if shouldDecrypt {
			err := content.DecryptData(wad.Ticket.TitleKey)
			if err != nil {
				return err
			}
		}

		err = dir.writeContents(content)
		if err != nil {
			return err
		}
	}

	return nil
}
