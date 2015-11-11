package wal

import (
	"fmt"

	"github.com/MediaMath/keryxlib/pg/control"
)

// NewCursorAtCheckpoint creates a new cursor pointing at the current checkpoint
func NewCursorAtCheckpoint(path string) (cursor *Cursor, err error) {
	control, err := control.NewControlFromDataDir(path)
	if err == nil {
		blockReader := blockReader{path, control.XlogBlcksz, control.MaxAlign}
		checkPointLocation := LocationFromUint32s(control.CheckPointLogID, control.CheckPointRecordOffset)

		cursor = &Cursor{checkPointLocation, blockReader}
	}

	return
}

// Cursor models a position in the WAL of a PostgreSQL system
type Cursor struct {
	location Location
	reader   blockReader
}

func (c Cursor) String() string {
	return fmt.Sprintf("%.8X%.16X", c.location.timelineID, c.location.Offset())
}

// MoveTo sets the cursor to point at the specified location in the WAL even if its invalid
func (c Cursor) MoveTo(location Location) Cursor {
	return Cursor{location, c.reader}
}

// ReadEntry will read a tuple at the current location and if successful advance to the next tuple
func (c Cursor) ReadEntry() (entry *Entry, cur Cursor, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			}
		}
	}()

	cur = c
	block := cur.reader.readBlock(cur.location)
	page := &Page{block}

	recordHeader := NewRecordHeader(block, cur.location)
	afterRecordHeader := cur.MoveTo(cur.location.Add(RecordHeaderSize).Aligned())
	recordBody := NewRecordBody(recordHeader)

	var bytesRead uint64

	if cur.location.IsOnSamePageAs(afterRecordHeader.location) {
		bytesRead = recordBody.AppendBodyAfterHeader(block, afterRecordHeader.location)
		cur = afterRecordHeader
	}

	for !recordBody.IsComplete() {
		cur = cur.MoveTo(cur.location.StartOfNextPage())

		nextBlock := cur.reader.readBlock(cur.location)
		nextPage := Page{nextBlock}

		cur = cur.MoveTo(cur.location.Add(nextPage.HeaderLength()))

		bytesRead = recordBody.AppendContinuation(nextPage)
		if bytesRead == 0 {
			return nil, c, nil
		}
	}

	entry = NewEntry(page, recordHeader, recordBody)
	cur = cur.MoveTo(cur.location.Add(bytesRead).Aligned())

	nextRecord := scanForRecordWithPrevious(c, c.MoveTo(c.location.Add(1).Aligned()))
	if nextRecord != nil {
		cur = *nextRecord
	} else {
		cur = c
		if entry.Type != Commit && entry.Type != Abort {
			entry = nil
		}
	}

	return
}

func scanForRecordWithPrevious(previous, startAt Cursor) *Cursor {
	out := samePageScanForRecordWithPrevious(previous, startAt)
	if out == nil {
		out = multiPageScanForRecordWithPrevious(previous, startAt)
	}
	return out
}

func samePageScanForRecordWithPrevious(previous, startAt Cursor) *Cursor {
	block := startAt.reader.readBlock(startAt.location)

	cur := startAt.MoveTo(startAt.location.Aligned())

	for cur.location.IsOnSamePageAs(startAt.location) {
		maybeHeader := NewRecordHeader(block, cur.location)
		if maybeHeader != nil && maybeHeader.Previous().Offset() == previous.location.Offset() {
			return &cur
		}

		cur = cur.MoveTo(cur.location.Add(1).Aligned())
	}

	return nil
}

func multiPageScanForRecordWithPrevious(previous, startAt Cursor) (out *Cursor) {
	var cur *Cursor
	for cur == nil {
		startAt = startAt.MoveTo(startAt.location.StartOfNextPage())
		cur = cursorAtFirstRecordOnPage(startAt)
	}

	block := cur.reader.readBlock(cur.location)
	maybeHeader := NewRecordHeader(block, cur.location)
	if maybeHeader != nil && maybeHeader.Previous().Offset() == previous.location.Offset() {
		out = cur
	}

	return
}

func cursorAtFirstRecordOnPage(startAt Cursor) (out *Cursor) {
	block := startAt.reader.readBlock(startAt.location)
	page := Page{block}

	cur := startAt.MoveTo(startAt.location.StartOfPage().Add(page.HeaderLength()))

	if cont := page.Continuation(); cont != nil {
		afterCont := cur.location.Add(uint64(len(cont) + 4)).Aligned()
		if afterCont.IsOnSamePageAs(cur.location) && afterCont.ToEndOfPage() >= RecordHeaderSize {
			curAfterCont := cur.MoveTo(afterCont)
			out = &curAfterCont
		}
	} else {
		out = &cur
	}

	return
}