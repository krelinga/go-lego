package golden

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type GoldenEntry struct {
	Path string
	Line int
	Text string
}

var byteOrder = binary.LittleEndian

// magicNumber is written as the first 8 bytes of every golden-entry file so that readers can detect
// files that are corrupt, truncated, or of the wrong format before attempting to decode entries.
var magicNumber = [8]byte{'G', 'L', 'D', 'N', 0, 1, 0, 0}

func writeString(w io.Writer, s string) error {
	b := []byte(s)
	if err := binary.Write(w, byteOrder, uint32(len(b))); err != nil {
		return err
	}
	_, err := w.Write(b)
	return err
}

func readString(r io.Reader) (string, error) {
	var length uint32
	if err := binary.Read(r, byteOrder, &length); err != nil {
		return "", err
	}
	b := make([]byte, length)
	if _, err := io.ReadFull(r, b); err != nil {
		return "", err
	}
	return string(b), nil
}

// WriteGoldenEntry appends a GoldenEntry to the file at path, creating it if necessary. When
// creating a new file the magic number is written first.
func WriteGoldenEntry(path string, entry GoldenEntry) error {
	var buf bytes.Buffer
	if err := writeString(&buf, entry.Path); err != nil {
		return fmt.Errorf("encoding Path: %w", err)
	}
	if err := binary.Write(&buf, byteOrder, int64(entry.Line)); err != nil {
		return fmt.Errorf("encoding Line: %w", err)
	}
	if err := writeString(&buf, entry.Text); err != nil {
		return fmt.Errorf("encoding Text: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stating file: %w", err)
	}
	if info.Size() == 0 {
		if _, err := f.Write(magicNumber[:]); err != nil {
			return fmt.Errorf("writing magic number: %w", err)
		}
	}

	if err := binary.Write(f, byteOrder, uint64(buf.Len())); err != nil {
		return fmt.Errorf("writing entry size: %w", err)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("writing entry data: %w", err)
	}
	return nil
}

// ReadGoldenEntries reads all GoldenEntry records from the file at path. It returns an error if the
// file does not begin with the expected magic number.
func ReadGoldenEntries(path string) ([]GoldenEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var gotMagic [8]byte
	if _, err := io.ReadFull(f, gotMagic[:]); err != nil {
		return nil, fmt.Errorf("reading magic number: %w", err)
	}
	if gotMagic != magicNumber {
		return nil, fmt.Errorf("invalid magic number: got %x, want %x", gotMagic, magicNumber)
	}

	var entries []GoldenEntry
	for {
		var entrySize uint64
		if err := binary.Read(f, byteOrder, &entrySize); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("reading entry size: %w", err)
		}

		payload := make([]byte, entrySize)
		if _, err := io.ReadFull(f, payload); err != nil {
			return nil, fmt.Errorf("reading entry payload: %w", err)
		}
		r := bytes.NewReader(payload)

		var entry GoldenEntry
		entry.Path, err = readString(r)
		if err != nil {
			return nil, fmt.Errorf("decoding Path: %w", err)
		}
		var line int64
		if err := binary.Read(r, byteOrder, &line); err != nil {
			return nil, fmt.Errorf("decoding Line: %w", err)
		}
		entry.Line = int(line)
		entry.Text, err = readString(r)
		if err != nil {
			return nil, fmt.Errorf("decoding Text: %w", err)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
