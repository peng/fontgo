package font

// Windows kern table (version 0) test data
func getWindowsKernData() []byte {
	// Kern table version 0 (Windows)
	// version=0, nTables=1
	// subtable: version=0, length=26, coverage=0x0001 (horizontal kerning, format 0)
	// nPairs=2, searchRange=2, entrySelector=1, rangeShift=2
	// Pair 1: left=65('A'), right=86('V'), value=-50
	// Pair 2: left=70('F'), right=46('.'), value=-30
	return []byte{
		0x00, 0x00, // version = 0 (Windows)
		0x00, 0x01, // nTables = 1
		// Subtable header
		0x00, 0x00, // subtable version = 0
		0x00, 0x1A, // length = 26 (header 14 + pairs 12)
		0x00, 0x00, // coverage = 0x0000 (horizontal kerning, format 0)
		0x00, 0x02, // nPairs = 2
		0x00, 0x02, // searchRange = 2
		0x00, 0x01, // entrySelector = 1
		0x00, 0x02, // rangeShift = 2
		// Pair 1: A-V kerning
		0x00, 0x41, // left = 65 ('A')
		0x00, 0x56, // right = 86 ('V')
		0xFF, 0xCE, // value = -50
		// Pair 2: F-. kerning
		0x00, 0x46, // left = 70 ('F')
		0x00, 0x2E, // right = 46 ('.')
		0xFF, 0xE2, // value = -30
	}
}

// Mac kern table (version 1, old format) test data
func getMacKernData() []byte {
	// Kern table version 1 (Mac, old format)
	// version=1, nTables=1
	// subtable: length=32, coverage=0x8000 (horizontal, format 0)
	// tupleIndex=0x0000, nPairs=2
	return []byte{
		0x00, 0x01, // version = 1 (Mac)
		0x00, 0x01, // nTables = 1
		// Subtable header
		0x00, 0x00, 0x00, 0x20, // length = 32
		0x80, 0x00, // coverage = 0x8000 (horizontal, format 0)
		0x00, 0x00, // tupleIndex = 0
		0x00, 0x02, // nPairs = 2
		0x00, 0x02, // searchRange = 2
		0x00, 0x01, // entrySelector = 1
		0x00, 0x02, // rangeShift = 2
		// Pair 1: T-o kerning
		0x00, 0x54, // left = 84 ('T')
		0x00, 0x6F, // right = 111 ('o')
		0xFF, 0xD8, // value = -40
		// Pair 2: W-a kerning
		0x00, 0x57, // left = 87 ('W')
		0x00, 0x61, // right = 97 ('a')
		0xFF, 0xEC, // value = -20
	}
}

// Mac kern table (version 1, new format) test data
func getMacNewKernData() []byte {
	// Kern table version 1 (Mac, new format with nTables=0)
	// version=1, nTables=0 (indicates new format)
	// actual nTables=1 (32-bit)
	return []byte{
		0x00, 0x01, // version = 1 (Mac)
		0x00, 0x00, // nTables = 0 (new format indicator)
		0x00, 0x00, 0x00, 0x01, // actual nTables = 1 (32-bit)
		// Subtable header
		0x00, 0x00, 0x00, 0x20, // length = 32
		0x80, 0x00, // coverage = 0x8000 (horizontal, format 0)
		0x00, 0x00, // tupleIndex = 0
		0x00, 0x01, // nPairs = 1
		0x00, 0x01, // searchRange = 1
		0x00, 0x00, // entrySelector = 0
		0x00, 0x00, // rangeShift = 0
		// Pair 1: L-Y kerning
		0x00, 0x4C, // left = 76 ('L')
		0x00, 0x59, // right = 89 ('Y')
		0xFF, 0xB0, // value = -80
	}
}

// Windows format 2 kern table test data
func getWindowsFormat2KernData() []byte {
	// Kern table version 0 (Windows), format 2
	// Simple n×m array of kerning values
	return []byte{
		0x00, 0x00, // version = 0 (Windows)
		0x00, 0x01, // nTables = 1
		// Subtable header (starts at byte 4)
		0x00, 0x00, // subtable version = 0
		0x00, 0x26, // length = 38 bytes
		0x00, 0x02, // coverage = 0x0002 (horizontal, format 2)
		0x00, 0x04, // rowWidth = 4 bytes per row
		0x00, 0x0E, // leftOffsetTable = 14 (from subtable start)
		0x00, 0x16, // rightOffsetTable = 22 (from subtable start)
		0x00, 0x1E, // array = 30 (from subtable start)
		// Left class table (at byte 4+14=18)
		0x00, 0x41, // firstGlyph = 65 ('A')
		0x00, 0x02, // nGlyphs = 2
		0x00, 0x1E, // offset for glyph 65 ('A')
		0x00, 0x20, // offset for glyph 66 ('B')
		// Right class table (at byte 4+22=26)
		0x00, 0x56, // firstGlyph = 86 ('V')
		0x00, 0x02, // nGlyphs = 2
		0x00, 0x00, // offset for glyph 86 ('V')
		0x00, 0x02, // offset for glyph 87 ('W')
		// Kern value array (at byte 4+30=34)
		0xFF, 0xCE, // A-V = -50
		0xFF, 0xD8, // A-W = -40
		0xFF, 0xE2, // B-V = -30
		0xFF, 0xEC, // B-W = -20
	}
}

// Mac format 2 kern table test data
func getMacFormat2KernData() []byte {
	// Kern table version 1 (Mac), format 2
	return []byte{
		0x00, 0x01, // version = 1 (Mac)
		0x00, 0x01, // nTables = 1
		// Subtable header (starts here)
		0x00, 0x00, 0x00, 0x28, // length = 40 bytes
		0x80, 0x02, // coverage = 0x8002 (horizontal, format 2)
		0x00, 0x00, // tupleIndex = 0
		0x00, 0x04, // rowWidth = 4 bytes per row
		0x00, 0x10, // leftOffsetTable = offset 16 from subtable start
		0x00, 0x18, // rightOffsetTable = offset 24 from subtable start
		0x00, 0x20, // array = offset 32 from subtable start
		// Left class table (at offset 16)
		0x00, 0x54, // firstGlyph = 84 ('T')
		0x00, 0x02, // nGlyphs = 2
		0x00, 0x20, // offset for glyph 84 ('T')
		0x00, 0x22, // offset for glyph 85 ('U')
		// Right class table (at offset 24)
		0x00, 0x6F, // firstGlyph = 111 ('o')
		0x00, 0x02, // nGlyphs = 2
		0x00, 0x00, // offset for glyph 111 ('o')
		0x00, 0x02, // offset for glyph 112 ('p')
		// Kern value array (at offset 32)
		0xFF, 0xD8, // T-o = -40
		0xFF, 0xEC, // T-p = -20
		0xFF, 0xF0, // U-o = -16
		0xFF, 0xF6, // U-p = -10
	}
}

// Mac format 3 kern table test data
func getMacFormat3KernData() []byte {
	// Kern table version 1 (Mac), format 3
	// Simple n×m index array
	return []byte{
		0x00, 0x01, // version = 1 (Mac)
		0x00, 0x01, // nTables = 1
		// Subtable header
		0x00, 0x00, 0x00, 0x20, // length = 32 bytes
		0x80, 0x03, // coverage = 0x8003 (horizontal, format 3)
		0x00, 0x00, // tupleIndex = 0
		0x00, 0x04, // glyphCount = 4
		0x02, // kernValueCount = 2
		0x02, // leftClassCount = 2
		0x02, // rightClassCount = 2
		0x00, // flags = 0
		// Kern values (2 values)
		0xFF, 0xCE, // kernValue[0] = -50
		0xFF, 0xEC, // kernValue[1] = -20
		// Left class array (4 glyphs)
		0x00, 0x00, 0x01, 0x01, // leftClass[0..3]
		// Right class array (4 glyphs)
		0x00, 0x01, 0x00, 0x01, // rightClass[0..3]
		// Kern index array (2×2=4 entries)
		0x00, 0x01, 0x01, 0x00, // kernIndex[0..3]
	}
}
