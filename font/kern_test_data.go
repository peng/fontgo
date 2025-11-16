package font

// Windows kern table (version 0) test data
func getWindowsKernData() []byte {
	// Kern table version 0 (Windows)
	// version=0, nTables=1
	// subtable: version=0, length=26, coverage=0x0001 (horizontal kerning)
	// nPairs=2, searchRange=2, entrySelector=1, rangeShift=2
	// Pair 1: left=65('A'), right=86('V'), value=-50
	// Pair 2: left=70('F'), right=46('.'), value=-30
	return []byte{
		0x00, 0x00, // version = 0 (Windows)
		0x00, 0x01, // nTables = 1
		// Subtable header
		0x00, 0x00, // subtable version = 0
		0x00, 0x1A, // length = 26 (header 14 + pairs 12)
		0x00, 0x01, // coverage = 0x0001 (horizontal kerning, format 0)
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
