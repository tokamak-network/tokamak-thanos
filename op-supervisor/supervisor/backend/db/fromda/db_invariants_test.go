package fromda

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type statInvariant func(stat os.FileInfo, m *stubMetrics) error
type linkInvariant func(prev, current LinkEntry) error

// checkDBInvariants reads the database log directly and asserts a set of invariants on the data.
func checkDBInvariants(t *testing.T, dbPath string, m *stubMetrics) {
	stat, err := os.Stat(dbPath)
	require.NoError(t, err)

	statInvariants := []statInvariant{
		invariantFileSizeMultipleOfEntrySize,
		invariantFileSizeMatchesEntryCountMetric,
	}
	for _, invariant := range statInvariants {
		require.NoError(t, invariant(stat, m))
	}

	// Read all entries as binary blobs
	file, err := os.OpenFile(dbPath, os.O_RDONLY, 0o644)
	require.NoError(t, err)
	entries := make([]Entry, stat.Size()/EntrySize)
	for i := range entries {
		n, err := io.ReadFull(file, entries[i][:])
		require.NoErrorf(t, err, "failed to read entry %v", i)
		require.EqualValuesf(t, EntrySize, n, "read wrong length for entry %v", i)
	}
	var links []LinkEntry
	for i, e := range entries {
		var v LinkEntry
		require.NoError(t, v.decode(e), "failed to decode entry %d", i)
		links = append(links, v)
	}

	linkInvariants := []linkInvariant{
		invariantDerivedTimestamp,
		invariantSourceTimestamp,
		invariantNumberIncrement,
	}
	for i, link := range links {
		if i == 0 {
			continue
		}
		for _, invariant := range linkInvariants {
			err := invariant(links[i-1], link)
			if err != nil {
				require.NoErrorf(t, err, "Invariant breached: \n%v", fmtEntries(entries))
			}
		}
	}
}

func fmtEntries(entries []Entry) string {
	out := ""
	for i, entry := range entries {
		out += fmt.Sprintf("%v: %x\n", i, entry)
	}
	return out
}

func invariantFileSizeMultipleOfEntrySize(stat os.FileInfo, _ *stubMetrics) error {
	size := stat.Size()
	if size%EntrySize != 0 {
		return fmt.Errorf("expected file size to be a multiple of entry size (%v) but was %v", EntrySize, size)
	}
	return nil
}

func invariantFileSizeMatchesEntryCountMetric(stat os.FileInfo, m *stubMetrics) error {
	size := stat.Size()
	if m.DBDerivedEntryCount*EntrySize != size {
		return fmt.Errorf("expected file size to be entryCount (%v) * entrySize (%v) = %v but was %v", m.DBDerivedEntryCount, EntrySize, m.DBDerivedEntryCount*EntrySize, size)
	}
	return nil
}

func invariantDerivedTimestamp(prev, current LinkEntry) error {
	if current.derived.Timestamp < prev.derived.Timestamp {
		return fmt.Errorf("derived timestamp must be >=, current: %s, prev: %s", current.derived, prev.derived)
	}
	return nil
}

func invariantNumberIncrement(prev, current LinkEntry) error {
	// derived stays the same if the new L1 block is empty.
	derivedSame := current.derived.Number == prev.derived.Number
	// source stays the same if this L2 block is derived from the same L1 block as the last L2 block
	sourceSame := current.source.Number == prev.source.Number
	// At least one of the two must increment, otherwise we are just repeating data in the DB.
	if derivedSame && sourceSame {
		return errors.New("expected at least either source or derived to increment, but both have same number")
	}
	derivedIncrement := current.derived.Number == prev.derived.Number+1
	sourceIncrement := current.source.Number == prev.source.Number+1
	if derivedIncrement == sourceIncrement { // one of the two must be true, the other false, to pass.
		return errors.New("expected source or (excl.) derived to increment")
	}
	return nil
}

func invariantSourceTimestamp(prev, current LinkEntry) error {
	if current.source.Timestamp < prev.source.Timestamp {
		return fmt.Errorf("source timestamp must be >=, current: %s, prev: %s", current.source, prev.source)
	}
	return nil
}
