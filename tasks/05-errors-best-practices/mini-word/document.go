package miniword

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const maxPages = 3

var (
	errInvalidPageNumber = errors.New("invalid page number")
	errNoMorePages       = errors.New("no more pages")
	errEmptyText         = errors.New("empty text")
)

type Document struct {
	totalPages       int
	activePageNumber int
	pagesContent     []string
	err              error
}

func NewDocument() *Document {
	d := &Document{
		pagesContent: make([]string, 0),
	}

	d.AddPage()
	d.SetActivePage(1)

	return d
}

func (d *Document) AddPage() {
	if d.err != nil {
		return
	}

	if d.totalPages == maxPages {
		d.err = errNoMorePages

		return
	}

	d.totalPages++
	d.pagesContent = append(d.pagesContent, "")
}

func (d *Document) SetActivePage(number int) {
	if d.err != nil {
		return
	}

	if number < 1 || number > d.totalPages {
		d.err = errInvalidPageNumber

		return
	}

	d.activePageNumber = number
}

func (d *Document) WriteText(s string) {
	if d.err != nil {
		return
	}

	if s == "" {
		d.err = errEmptyText

		return
	}

	d.pagesContent[d.activePageNumber-1] += s
}

func (d *Document) WriteTo(w io.Writer) (n int64, err error) {
	if d.err != nil {
		return 0, d.err
	}

	sb := strings.Builder{}

	for i, pageContent := range d.pagesContent {
		sb.WriteString(fmt.Sprintf("--- Page %d ---", i+1))
		sb.WriteRune('\n')
		sb.WriteString(pageContent)
		sb.WriteRune('\n')
	}

	writtenBytes, err := w.Write([]byte(sb.String()))
	if err != nil {
		return int64(writtenBytes), nil
	}

	return int64(writtenBytes), nil
}
