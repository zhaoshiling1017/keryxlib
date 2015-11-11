package wal

import "testing"

func TestPageExpectations(t *testing.T) {
	for _, exp := range pageExpectations {
		failIfPageNotMatching(t, exp, Page{exp.bs})
	}
}

var pageExpectations = []pageExpectation{
	{
		3, 1, NewLocationWithDefaults(0x0000000013000000), 0x55085653550bf6f3, 0x01000000, 0x00002000, []byte{0x0e, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc4, 0x09, 0x6f, 0x00, 0x82, 0x01, 0x00}, true, true, 32,
		[]byte{
			0x66, 0xd0, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x13,
			0xf3, 0xf6, 0x0b, 0x55, 0x53, 0x56, 0x08, 0x55, 0x00, 0x00, 0x00, 0x01, 0x00, 0x20, 0x00, 0x00,
			0x0d, 0x00, 0x00, 0x00, 0x0e, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc4, 0x09, 0x6f, 0x00, 0x82, 0x01,
			0x00, 0x00, 0x00, 0x00},
	},
	{
		1, 1, NewLocationWithDefaults(0x0000000013002000), 0, 0, 0, []byte{0x64, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xea, 0xf5}, true, false, 16,
		[]byte{
			0x66, 0xd0, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x13,
			0x0b, 0x00, 0x00, 0x00, 0x64, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xea, 0xf5, 0x07},
	},
	{
		0, 1, NewLocationWithDefaults(0x0000000013004000), 0, 0, 0, []byte{}, false, false, 16,
		[]byte{
			0x66, 0xd0, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x13,
			0x62, 0x94, 0x8e, 0x7e, 0x00, 0x00, 0x00, 0x00, 0xc8, 0x3f, 0x00, 0x13, 0xbe, 0xf6, 0x07, 0x00},
	},
	{
		1, 1, NewLocationWithDefaults(0x0000000013006000), 0, 0, 0, []byte{0x64, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xed, 0xf5}, true, false, 16,
		[]byte{
			0x66, 0xd0, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x60, 0x00, 0x13,
			0x0b, 0x00, 0x00, 0x00, 0x64, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xed, 0xf5, 0x07},
	},
	{
		1, 1, NewLocationWithDefaults(0x0000000013008000), 0, 0, 0, []byte{0x0e, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc6, 0x09, 0x20, 0x00, 0x82}, true, false, 16,
		[]byte{
			0x66, 0xd0, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x13,
			0x0b, 0x00, 0x00, 0x00, 0x0e, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc6, 0x09, 0x20, 0x00, 0x82, 0x01},
	},
	{
		1, 1, NewLocationWithDefaults(0x000000001300a000), 0, 0, 0, []byte{0x7f, 0x06, 0x00, 0x00, 0x0b, 0x40, 0x00, 0x00, 0x0e, 0x40, 0x00}, true, false, 16,
		[]byte{
			0x66, 0xd0, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x13,
			0x0b, 0x00, 0x00, 0x00, 0x7f, 0x06, 0x00, 0x00, 0x0b, 0x40, 0x00, 0x00, 0x0e, 0x40, 0x00, 0x00},
	},
	{
		1, 1, NewLocationWithDefaults(0x000000001300c000), 0, 0, 0, []byte{0x64, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xed, 0xf5}, true, false, 16,
		[]byte{
			0x66, 0xd0, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x13,
			0x0b, 0x00, 0x00, 0x00, 0x64, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xed, 0xf5, 0x07},
	},
	{
		0, 1, NewLocationWithDefaults(0x000000001300e000), 0, 0, 0, []byte{}, false, false, 16,
		[]byte{
			0x66, 0xd0, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x13,
			0x8b, 0x8e, 0x84, 0x1f, 0x00, 0x00, 0x00, 0x00, 0xc0, 0xdf, 0x00, 0x13, 0xbe, 0xf6, 0x07, 0x00},
	},
	{
		1, 1, NewLocationWithDefaults(0x0000000013010000), 0, 0, 0, []byte{0x00, 0x7b, 0xe4, 0x0e, 0x00, 0x0b, 0x61, 0x73, 0x64, 0x66, 0x00}, true, false, 16,
		[]byte{
			0x66, 0xd0, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x13,
			0x0b, 0x00, 0x00, 0x00, 0x00, 0x7b, 0xe4, 0x0e, 0x00, 0x0b, 0x61, 0x73, 0x64, 0x66, 0x00, 0x00},
	},
	{
		1, 1, NewLocationWithDefaults(0x0000000013012000), 0, 0, 0, []byte{0x0e, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc8, 0x09, 0x63, 0x00, 0x82}, true, false, 16,
		[]byte{
			0x66, 0xd0, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x01, 0x13,
			0x0b, 0x00, 0x00, 0x00, 0x0e, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc8, 0x09, 0x63, 0x00, 0x82, 0x01},
	},
}

type pageExpectation struct {
	info         uint16
	timelineID   uint32
	location     Location
	systemID     uint64
	segmentSize  uint32
	blockSize    uint32
	continuation []byte
	isCont       bool
	isLong       bool
	headerLength uint64
	bs           []byte
}

func failIfPageNotMatching(t *testing.T, exp pageExpectation, act Page) {
	if !act.MagicValueIsValid() {
		t.Error("invalid page")
	}
	if exp.info != act.Info() {
		t.Errorf("expected %.4x but got %.4x for info", exp.info, act.Info())
	}
	if exp.timelineID != act.TimelineID() {
		t.Errorf("expected %.8x but got %.8x for timeline id", exp.timelineID, act.TimelineID())
	}
	if exp.location.offset != act.Location().offset {
		t.Errorf("expected %.16x but got %.16x for location", exp.location.offset, act.Location().offset)
	}
	if exp.systemID != act.SystemID() {
		t.Errorf("expected %.16x but got %.16x for system id", exp.systemID, act.SystemID())
	}
	if exp.segmentSize != act.SegmentSize() {
		t.Errorf("expected %.8x but got %.8x for segment size", exp.segmentSize, act.SegmentSize())
	}
	if exp.blockSize != act.BlockSize() {
		t.Errorf("expected %.8x but got %.8x for block size", exp.blockSize, act.BlockSize())
	}
	if !continuationsMatch(exp.continuation, act.Continuation()) {
		t.Error("continuation didn't match", exp.continuation, act.Continuation())
	}
	if exp.isCont != act.IsCont() {
		t.Errorf("expected %v but got %v for isCont", exp.isCont, act.IsCont())
	}
	if exp.isLong != act.IsLong() {
		t.Errorf("expected %v but got %v for isLong", exp.isLong, act.IsLong())
	}
	if exp.headerLength != act.HeaderLength() {
		t.Errorf("expected %v but got %v for header length", exp.headerLength, act.HeaderLength())
	}
}

func continuationsMatch(a, b []byte) bool {
	if len(a) == len(b) {
		for i := range b {
			if a[i] != b[i] {
				return false
			}
		}
	} else {
		return false
	}

	return true
}