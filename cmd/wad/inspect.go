package wad

import (
	"encoding/binary"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/wii-tools/wadlib"
)

func Inspect(c *cli.Context) error {
	wad, err := wadlib.LoadWADFromFile(c.String("in"))
	if err != nil {
		return err
	}

	fmt.Fprint(c.App.Writer, "=== Title ===\n")
	// We want to separate the high and low halves in order to separate groups.
	upperId := wad.Ticket.TitleID >> 32
	lowerId := (uint32)(wad.Ticket.TitleID)

	// Convert the lower ID to ASCII.
	var lowerBuf [4]byte
	binary.BigEndian.PutUint32(lowerBuf[0:4], lowerId)
	lowerIdAscii := string(lowerBuf[:])

	fmt.Fprintf(c.App.Writer, "App ID: %08x-%08x (%s)\n", upperId, lowerId, lowerIdAscii)
	fmt.Fprintf(c.App.Writer, "Title version: %d\n", wad.TMD.TitleVersion)
	fmt.Fprintf(c.App.Writer, "Requested IOS version: %d\n", wad.TMD.SystemVersionLow)
	fmt.Fprintf(c.App.Writer, "Ticket ID: %d\n", wad.Ticket.TicketID)
	if wad.Ticket.ConsoleID != 0 {
		fmt.Fprintf(c.App.Writer, "Console ID: %d\n", wad.Ticket.ConsoleID)
	}

	// We want time limits to be listed separately if enabled.
	// Relatively redundant, but it does the trick.
	shouldShowLimits := false
	for _, limit := range wad.Ticket.TimeLimits {
		if limit.Code != 0 {
			shouldShowLimits = true
			break
		}
	}
	if shouldShowLimits {
		fmt.Fprint(c.App.Writer, "\n=== Limit ===\n")
		for index, limit := range wad.Ticket.TimeLimits {
			fmt.Fprintf(c.App.Writer, "Index %d has a limit of %d seconds\n", index, limit.Limit)
		}
	}

	fmt.Fprint(c.App.Writer, "\n=== Contents ===\n")

	pluralisation := ""
	if wad.TMD.NumberOfContents != 1 {
		pluralisation = "s"
	}
	fmt.Fprintf(c.App.Writer, "This title contains %d content%s.\n", wad.TMD.NumberOfContents, pluralisation)
	fmt.Fprint(c.App.Writer, "------------------------------------------------------------------------------------\n")
	fmt.Fprint(c.App.Writer, "| ID | Content ID |                   Hash                   |  Type  |    Size    |\n")
	fmt.Fprint(c.App.Writer, "------------------------------------------------------------------------------------\n")
	for _, content := range wad.TMD.Contents {
		contentType := ""
		switch content.Type {
		case wadlib.TitleTypeNormal:
			contentType = "normal"
		case wadlib.TitleTypeShared:
			contentType = "shared"
		default:
			contentType = "other "
		}

		fmt.Fprintf(c.App.Writer, "| %-2d |  %08x  | %20x | %s | %-10d |\n", content.Index, content.ID, content.Hash, contentType, content.Size)
	}
	fmt.Fprint(c.App.Writer, "------------------------------------------------------------------------------------\n")

	return nil
}
