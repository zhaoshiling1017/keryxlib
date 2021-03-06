package streams

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"github.com/MediaMath/keryxlib/filters"
	"github.com/MediaMath/keryxlib/message"
	"github.com/MediaMath/keryxlib/pg"
	"github.com/MediaMath/keryxlib/pg/wal"
)

//TxnBuffer is a stream of WAL entries organized by transaction
type TxnBuffer struct {
	Filters          filters.MessageFilter
	WorkingDirectory string
	SchemaReader     *pg.SchemaReader
}

func (b *TxnBuffer) filterRelation(entry *wal.Entry) bool {
	return entry.RelationID > 0 && b.Filters.FilterRelID(entry.RelationID)
}

func (b *TxnBuffer) hasDatabaseConnection(entry *wal.Entry) bool {
	return b.SchemaReader == nil || b.SchemaReader.HaveConnectionToDb(entry.DatabaseID)
}

//Start takes a channel of WAL entries and async selects on it.  As it finds a commit for a transaction it publishes a slice of the entries in that transaction.  Aborted transactions are not published. Rel filtering happens in this stream and not downstream.
func (b *TxnBuffer) Start(entryChan <-chan *wal.Entry) (<-chan []*wal.Entry, error) {
	txns := make(chan []*wal.Entry)

	go func() {
		buffer := message.NewBuffer(b.WorkingDirectory, 10*1024*wal.EntryBytesSize, wal.EntryBytesSize)
		var lastEntry *wal.Entry
		for entry := range entryChan {
			if lastEntry != nil && lastEntry.ReadFrom.Offset() > entry.ReadFrom.Offset() {
				continue
			} else if entry.Type == wal.Unknown {
				continue
			} else if entry.Type != wal.Commit && (!b.hasDatabaseConnection(entry) || b.filterRelation(entry)) {
				continue
			}

			lastEntry = entry

			if entry.Type == wal.Commit {
				entriesBytes := buffer.Remove(entry.TransactionID)
				var entries []*wal.Entry
				for _, entryBytes := range entriesBytes {
					e := wal.EntryFromBytes(entryBytes)
					entries = append(entries, &e)
				}
				if len(entries) != 0 {
					entries = append(entries, entry)
					txns <- entries
				}
			} else if entry.Type == wal.Abort {
				buffer.Remove(entry.TransactionID)
			} else {
				buffer.Add(entry.TransactionID, entry.ToBytes())
			}

		}
		close(txns)
	}()

	return txns, nil
}
