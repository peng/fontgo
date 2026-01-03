package font

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type TagItem struct {
	CheckSum uint32 `json:"checkSum"`
	Offset   uint32 `json:"offset"`
	Length   uint32 `json:"length"`
}

type Tables struct {
	Head *Head      `json:"head"`
	Maxp *Maxp      `json:"maxp"`
	Loca []int      `json:"loca"`
	Cmap *Cmap      `json:"cmap,omitempty"`
	Name *NameTable `json:"name,omitempty"`
	Hhea *Hhea      `json:"hhea,omitempty"`
	Hmtx *Hmtx      `json:"hmtx,omitempty"`
	Kern *Kern      `json:"kern,omitempty"`
	Os2  *OS2       `json:"os2"`
	Post *Post      `json:"post"`
	Fvar *Fvar      `json:"fvar"`
	Ltag *Ltag      `json:"ltag,omitempty"`
	Meta *Meta      `json:"meta"`
}

type FontInfo struct {
	OffsetTable  *OffsetTable        `json:"offsetTable"`
	TableContent map[string]*TagItem `json:"tableContent"`
	Tables       *Tables             `json:"tables"`
	Glyphs       *Glyphs             `json:"glyphs"`
}

type Font struct {
	fileByte []byte
	filePath string
	fontInfo *FontInfo
}

func ReadFontFile(filePath string) (f *Font, err error) {
	f = &Font{
		filePath: filePath,
	}
	f.fileByte, err = os.ReadFile(filePath)
	return
}

func (f *Font) GetFontInfo() (fontInfo *FontInfo, err error) {
	fileByte := f.fileByte
	// read offset table
	offsetTable := GetOffsetTable(fileByte)

	// read table content
	numTables := int(offsetTable.NumTables)
	tableContent := GetTableContent(numTables, fileByte)

	headInfo, existHead := tableContent["head"]
	maxpInfo, existMaxp := tableContent["maxp"]
	locaInfo, existLoca := tableContent["loca"]
	cmapInfo, existCmap := tableContent["cmap"]
	nameInfo, existName := tableContent["name"]
	hheaInfo, existHhea := tableContent["hhea"]
	hmtxInfo, existHmtx := tableContent["hmtx"]
	kernInfo, existKern := tableContent["kern"]
	os2Info, existOs2 := tableContent["OS/2"]
	postInfo, existPost := tableContent["post"]
	fvarInfo, existFvar := tableContent["fvar"]
	itagInfo, existLtag := tableContent["Ltag"]
	metaInfo, existMeta := tableContent["meta"]
	// add test

	// tables content
	tables := new(Tables)
	if existHead {
		tables.Head = GetHead(fileByte, int(headInfo.Offset))
	}
	if existMaxp {
		tables.Maxp = GetMaxp(fileByte, int(maxpInfo.Offset))
	}
	if existLoca && tables.Maxp != nil && tables.Head != nil {
		tables.Loca = GetLoca(fileByte, int(locaInfo.Offset), tables.Maxp.NumGlyphs, tables.Head.IndexToLocFormat)
	}
	if existCmap && tables.Maxp != nil {
		tables.Cmap, err = GetCmap(fileByte, int(cmapInfo.Offset), int(tables.Maxp.NumGlyphs))
	}

	if existName {
		tables.Name = GetName(fileByte, int(nameInfo.Offset))
	}

	if existHhea {
		tables.Hhea = GetHhea(fileByte, int(hheaInfo.Offset))
	}

	if existHmtx && existHhea && existMaxp {
		tables.Hmtx = GetHmtx(fileByte, int(hmtxInfo.Offset), int(tables.Hhea.NumOfLongHorMetrics), int(tables.Maxp.NumGlyphs))
	}

	if existKern {
		tables.Kern, err = GetKern(fileByte, int(kernInfo.Offset))
	}

	if existOs2 {
		tables.Os2 = GetOS2(fileByte, int(os2Info.Offset))
	}

	if existPost {
		tables.Post = GetPost(fileByte, int(postInfo.Offset))
	}

	if existFvar {
		tables.Fvar = GetFvar(fileByte, int(fvarInfo.Offset))
	}

	if existLtag {
		ltag, ltagErr := GetLtag(fileByte, int(itagInfo.Offset))
		if ltagErr == nil {
			tables.Ltag = ltag
		}
	}

	if existMeta {
		meta, metaErr := GetMeta(fileByte, int(metaInfo.Offset))

		if metaErr == nil {
			tables.Meta = meta
		}
	}

	fontInfo = new(FontInfo)

	fontInfo.OffsetTable = offsetTable
	fontInfo.TableContent = tableContent
	fontInfo.Tables = tables
	glyfStart := tableContent["glyf"].Offset
	fontInfo.Glyphs = GetGlyphs(fileByte, int(glyfStart), tables.Loca, int(tables.Maxp.NumGlyphs))

	f.fontInfo = fontInfo

	return
}

func (f *Font) Write(filePath string) (err error) {
	ext := filepath.Ext(filePath)
	if ext != ".ttf" {
		err = errors.New("Not support format!")
		return
	}
	supportTable := []string{"cmap", "fvar", "glyf", "head", "hhea", "hmtx", "kern", "Ltag", "loca", "maxp", "meta", "name", "OS/2", "post"}

	fontInfo := f.fontInfo
	data := []byte{}
	pad4 := func(n int) int { return (n + 3) &^ 3 }
	computeCheckSum := func(b []byte) uint32 {
		padded := make([]byte, pad4(len(b)))
		copy(padded, b)
		var sum uint32
		for i := 0; i < len(padded); i += 4 {
			sum += getUint32(padded[i : i+4])
		}
		return sum
	}

	// prepare table data, skipping missing entries with a warning
	tablesData := map[string][]byte{}
	actualTables := make([]string, 0, len(supportTable))
	var locaFromGlyphs []int
	for _, tag := range supportTable {
		var td []byte
		switch tag {
		case "cmap":
			if fontInfo.Tables.Cmap == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td, err = WriteCmap(fontInfo.Tables.Cmap)
			if err != nil {
				log.Printf("[WARN] table %s write failed: %v", tag, err)
				continue
			}
		case "fvar":
			if fontInfo.Tables.Fvar == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteFvar(fontInfo.Tables.Fvar)
		case "glyf":
			if fontInfo.Glyphs == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			numGlyphs := int(fontInfo.Tables.Maxp.NumGlyphs)
			td, locaFromGlyphs, err = WriteGlyphs(fontInfo.Glyphs, numGlyphs)
			if err != nil {
				log.Printf("[WARN] table %s write failed: %v", tag, err)
				continue
			}
		case "head":
			if fontInfo.Tables.Head == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			head := *fontInfo.Tables.Head
			head.CheckSumAdjustment = 0
			td = WriteHead(&head)
		case "hhea":
			if fontInfo.Tables.Hhea == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteHhea(fontInfo.Tables.Hhea)
		case "hmtx":
			if fontInfo.Tables.Hmtx == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteHmtx(fontInfo.Tables.Hmtx)
		case "kern":
			if fontInfo.Tables.Kern == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteKern(fontInfo.Tables.Kern)
		case "Ltag":
			if fontInfo.Tables.Ltag == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteLtag(fontInfo.Tables.Ltag)
		case "loca":
			if fontInfo.Tables.Head == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			locaData := fontInfo.Tables.Loca
			if len(locaFromGlyphs) > 0 {
				locaData = locaFromGlyphs
			}
			td = WriteLoca(locaData, fontInfo.Tables.Head.IndexToLocFormat)
		case "maxp":
			if fontInfo.Tables.Maxp == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteMaxp(fontInfo.Tables.Maxp)
		case "meta":
			if fontInfo.Tables.Meta == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteMeta(fontInfo.Tables.Meta)
		case "name":
			if fontInfo.Tables.Name == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteName(fontInfo.Tables.Name)
		case "OS/2":
			if fontInfo.Tables.Os2 == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WriteOS2(fontInfo.Tables.Os2)
		case "post":
			if fontInfo.Tables.Post == nil {
				log.Printf("[WARN] table %s data missing, continue", tag)
				continue
			}
			td = WritePost(fontInfo.Tables.Post)
		default:
			log.Printf("[WARN] table %s not handled, continue", tag)
			continue
		}

		tablesData[tag] = td
		actualTables = append(actualTables, tag)
	}

	// Sort actualTables to match WriteTableContent order (TrueType spec requires sorted tags)
	sort.Strings(actualTables)

	// rewrite offset table using only the tables we actually have
	cpOffsetTable := &OffsetTable{}
	cpOffsetTable.ScalerType = fontInfo.OffsetTable.ScalerType
	cpOffsetTable.NumTables = uint16(len(actualTables))

	maxPowerOf2 := uint16(1)
	entrySelector := uint16(0)
	for maxPowerOf2*2 <= cpOffsetTable.NumTables {
		maxPowerOf2 *= 2
		entrySelector++
	}
	cpOffsetTable.SearchRange = maxPowerOf2 * 16
	cpOffsetTable.EntrySelector = entrySelector
	cpOffsetTable.RangeShift = (cpOffsetTable.NumTables*16 - cpOffsetTable.SearchRange)

	offsetTableData := WriteOffsetTable(cpOffsetTable)
	data = append(data, offsetTableData...)

	// table directory size now depends on actual tables only
	tableDirSize := len(actualTables) * 16
	nextOffset := uint32(len(offsetTableData) + tableDirSize)

	// Build the directory entries and output using WriteTableContent
	cpTableContent := make(map[string]*TagItem, len(actualTables))
	for _, tag := range actualTables {
		td := tablesData[tag]
		length := uint32(len(td))
		checkSum := computeCheckSum(td)
		cpTableContent[tag] = &TagItem{
			CheckSum: checkSum,
			Offset:   nextOffset,
			Length:   length,
		}
		nextOffset += uint32(pad4(int(length)))
	}

	data = append(data, WriteTableContent(cpTableContent)...)

	// append table payloads with padding, using cached data from actualTables
	for _, tag := range actualTables {
		td := tablesData[tag]
		data = append(data, td...)
		padLen := pad4(len(td)) - len(td)
		if padLen > 0 {
			data = append(data, make([]byte, padLen)...)
		}
	}

	if err = os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return
	}
	return os.WriteFile(filePath, data, 0o755)
}

// Subset extracts only the specified characters from the font.
// It keeps glyph 0 (.notdef) and adds glyphs for the given characters.
// The font tables (cmap, glyf, loca, hmtx, maxp, hhea) are updated accordingly.
func (f *Font) Subset(chars []string) error {
	if f.fontInfo == nil {
		return errors.New("fontInfo is nil, call GetFontInfo first")
	}
	fontInfo := f.fontInfo

	if fontInfo.Tables.Cmap == nil || fontInfo.Tables.Cmap.WindowsCode == nil {
		return errors.New("cmap table or WindowsCode is nil")
	}

	// Step 1: Collect glyph indices for the requested characters
	// Always include glyph 0 (.notdef)
	neededGlyphSet := make(map[int]bool)
	neededGlyphSet[0] = true

	for _, char := range chars {
		for _, r := range char {
			unicode := int(r)
			if glyphIndex, ok := fontInfo.Tables.Cmap.WindowsCode[unicode]; ok {
				neededGlyphSet[glyphIndex] = true
			}
		}
	}

	// Step 2: Resolve compound glyph dependencies
	// Compound glyphs reference other glyphs, we need to include those too
	compoundMap := make(map[int]*GlyphCompound)
	for i := range fontInfo.Glyphs.Compounds {
		c := &fontInfo.Glyphs.Compounds[i]
		compoundMap[c.GlyphCommon.Index] = c
	}

	// Iteratively resolve dependencies until no new glyphs are added
	changed := true
	for changed {
		changed = false
		for glyphIdx := range neededGlyphSet {
			if compound, ok := compoundMap[glyphIdx]; ok {
				for _, comp := range compound.Component {
					refIdx := int(comp.GlyphIndex)
					if !neededGlyphSet[refIdx] {
						neededGlyphSet[refIdx] = true
						changed = true
					}
				}
			}
		}
	}

	// Step 3: Create sorted list of old glyph indices and mapping to new indices
	oldIndices := make([]int, 0, len(neededGlyphSet))
	for idx := range neededGlyphSet {
		oldIndices = append(oldIndices, idx)
	}
	sort.Ints(oldIndices)

	oldToNew := make(map[int]int)
	for newIdx, oldIdx := range oldIndices {
		oldToNew[oldIdx] = newIdx
	}

	// Step 4: Build new Glyphs with remapped indices
	simpleMap := make(map[int]*GlyphSimple)
	for i := range fontInfo.Glyphs.Simples {
		s := &fontInfo.Glyphs.Simples[i]
		simpleMap[s.GlyphCommon.Index] = s
	}

	newGlyphs := &Glyphs{}
	for newIdx, oldIdx := range oldIndices {
		if simple, ok := simpleMap[oldIdx]; ok {
			newSimple := *simple
			newSimple.GlyphCommon.Index = newIdx
			newGlyphs.Simples = append(newGlyphs.Simples, newSimple)
		} else if compound, ok := compoundMap[oldIdx]; ok {
			newCompound := *compound
			newCompound.GlyphCommon.Index = newIdx
			// Remap component glyph references
			for i := range newCompound.Component {
				oldRef := int(newCompound.Component[i].GlyphIndex)
				newCompound.Component[i].GlyphIndex = uint16(oldToNew[oldRef])
			}
			newGlyphs.Compounds = append(newGlyphs.Compounds, newCompound)
		}
		// If glyph not found (empty glyph), it won't be added, which is correct
	}
	fontInfo.Glyphs = newGlyphs

	// Step 5: Build new cmap with remapped glyph indices
	newWindowsCode := make(map[int]int)
	for unicode, oldGlyphIdx := range fontInfo.Tables.Cmap.WindowsCode {
		if newGlyphIdx, ok := oldToNew[oldGlyphIdx]; ok {
			newWindowsCode[unicode] = newGlyphIdx
		}
	}
	fontInfo.Tables.Cmap.WindowsCode = newWindowsCode

	// Rebuild cmap subtables for format 4 and 12
	if err := rebuildCmapSubtables(fontInfo.Tables.Cmap, newWindowsCode); err != nil {
		return err
	}

	// Step 6: Build new hmtx
	if fontInfo.Tables.Hmtx != nil {
		oldHmtx := fontInfo.Tables.Hmtx
		newHmtx := &Hmtx{}

		for _, oldIdx := range oldIndices {
			if oldIdx < len(oldHmtx.HMetrics) {
				newHmtx.HMetrics = append(newHmtx.HMetrics, oldHmtx.HMetrics[oldIdx])
			} else {
				// Use last advanceWidth with LeftSideBearing from array
				lsbIdx := oldIdx - len(oldHmtx.HMetrics)
				if lsbIdx < len(oldHmtx.LeftSideBearing) {
					lastMetric := oldHmtx.HMetrics[len(oldHmtx.HMetrics)-1]
					newHmtx.HMetrics = append(newHmtx.HMetrics, &LongHorMetric{
						AdvanceWidth:    lastMetric.AdvanceWidth,
						LeftSideBearing: oldHmtx.LeftSideBearing[lsbIdx],
					})
				}
			}
		}
		fontInfo.Tables.Hmtx = newHmtx
		// Clear LeftSideBearing since all metrics are now in HMetrics
		fontInfo.Tables.Hmtx.LeftSideBearing = nil
	}

	// Step 7: Update maxp.NumGlyphs
	if fontInfo.Tables.Maxp != nil {
		fontInfo.Tables.Maxp.NumGlyphs = uint16(len(oldIndices))
	}

	// Step 8: Update hhea.NumOfLongHorMetrics
	if fontInfo.Tables.Hhea != nil && fontInfo.Tables.Hmtx != nil {
		fontInfo.Tables.Hhea.NumOfLongHorMetrics = uint16(len(fontInfo.Tables.Hmtx.HMetrics))
	}

	// Step 9: Rebuild loca (will be done during Write)
	fontInfo.Tables.Loca = nil

	return nil
}

// rebuildCmapSubtables rebuilds cmap subtables based on the new unicode->glyph mapping
func rebuildCmapSubtables(cmap *Cmap, unicodeToGlyph map[int]int) error {
	// Collect all unicode code points and sort them
	unicodes := make([]int, 0, len(unicodeToGlyph))
	for u := range unicodeToGlyph {
		unicodes = append(unicodes, u)
	}
	sort.Ints(unicodes)

	// Only keep format 4 and 12 subtables, rebuild them
	newSubTables := make([]map[string]interface{}, 0)
	hasFormat4 := false
	hasFormat12 := false

	for _, subTable := range cmap.SubTables {
		format, ok := subTable["format"].(uint16)
		if !ok {
			continue
		}

		switch format {
		case 4:
			if !hasFormat4 {
				newSubTables = append(newSubTables, buildCmapFormat4(subTable, unicodeToGlyph, unicodes))
				hasFormat4 = true
			}
		case 12:
			if !hasFormat12 {
				newSubTables = append(newSubTables, buildCmapFormat12(subTable, unicodeToGlyph, unicodes))
				hasFormat12 = true
			}
		}
	}

	// If no format 4 exists, create one for BMP characters
	if !hasFormat4 && len(unicodes) > 0 {
		defaultSubTable := map[string]interface{}{
			"platformID":         uint16(3), // Windows
			"platformSpecificID": uint16(1), // Unicode BMP
			"language":           uint16(0),
		}
		newSubTables = append(newSubTables, buildCmapFormat4(defaultSubTable, unicodeToGlyph, unicodes))
	}

	cmap.SubTables = newSubTables
	cmap.NumberSubtables = uint16(len(newSubTables))

	return nil
}

// buildCmapFormat4 builds a format 4 cmap subtable
func buildCmapFormat4(oldSubTable map[string]interface{}, unicodeToGlyph map[int]int, unicodes []int) map[string]interface{} {
	// Filter to BMP only (0-0xFFFF) for format 4
	bmpUnicodes := make([]int, 0)
	for _, u := range unicodes {
		if u <= 0xFFFF {
			bmpUnicodes = append(bmpUnicodes, u)
		}
	}

	// Build segments - use idRangeOffset method for non-contiguous glyph mappings
	type segment struct {
		startCode     uint16
		endCode       uint16
		idDelta       int16
		idRangeOffset uint16
		glyphIdArray  []uint16 // only used when idRangeOffset != 0
	}
	var segments []segment

	i := 0
	for i < len(bmpUnicodes) {
		segStart := bmpUnicodes[i]
		segEnd := bmpUnicodes[i]
		firstGlyph := unicodeToGlyph[segStart]

		// Try to extend segment with consecutive unicodes that have consecutive glyph indices
		canUseDelta := true
		j := i + 1
		for j < len(bmpUnicodes) {
			u := bmpUnicodes[j]
			g := unicodeToGlyph[u]

			if u != segEnd+1 {
				// Unicode not consecutive, end segment
				break
			}

			expectedGlyph := firstGlyph + (u - segStart)
			if g != expectedGlyph {
				// Glyph not consecutive, can't use simple delta
				canUseDelta = false
				break
			}

			segEnd = u
			j++
		}

		if canUseDelta {
			// Use idDelta method
			delta := int16(firstGlyph - segStart)
			segments = append(segments, segment{
				startCode:     uint16(segStart),
				endCode:       uint16(segEnd),
				idDelta:       delta,
				idRangeOffset: 0,
			})
			i = j
		} else {
			// Single character segment with idDelta
			delta := int16(firstGlyph - segStart)
			segments = append(segments, segment{
				startCode:     uint16(segStart),
				endCode:       uint16(segStart),
				idDelta:       delta,
				idRangeOffset: 0,
			})
			i++
		}
	}

	// Add terminating segment
	segments = append(segments, segment{
		startCode:     0xFFFF,
		endCode:       0xFFFF,
		idDelta:       1,
		idRangeOffset: 0,
	})

	segCount := len(segments)
	segCountX2 := uint16(segCount * 2)

	// Calculate searchRange, entrySelector, rangeShift
	maxPow2 := 1
	entrySelector := uint16(0)
	for maxPow2*2 <= segCount {
		maxPow2 *= 2
		entrySelector++
	}
	searchRange := uint16(maxPow2 * 2)
	rangeShift := segCountX2 - searchRange

	// Build arrays
	endCode := make([]uint16, segCount)
	startCode := make([]uint16, segCount)
	idDelta := make([]uint16, segCount)
	idRangeOffset := make([]uint16, segCount)

	for idx, seg := range segments {
		endCode[idx] = seg.endCode
		startCode[idx] = seg.startCode
		idDelta[idx] = uint16(seg.idDelta)
		idRangeOffset[idx] = seg.idRangeOffset
	}

	// Calculate length
	length := uint16(16 + segCount*8) // header + 4 arrays * segCount * 2 bytes

	newSubTable := make(map[string]interface{})
	newSubTable["platformID"] = oldSubTable["platformID"]
	newSubTable["platformSpecificID"] = oldSubTable["platformSpecificID"]
	newSubTable["format"] = uint16(4)
	newSubTable["length"] = length
	newSubTable["language"] = oldSubTable["language"]
	newSubTable["segCountX2"] = segCountX2
	newSubTable["searchRange"] = searchRange
	newSubTable["entrySelector"] = entrySelector
	newSubTable["rangeShift"] = rangeShift
	newSubTable["endCode"] = endCode
	newSubTable["reservedPad"] = uint16(0)
	newSubTable["startCode"] = startCode
	newSubTable["idDelta"] = idDelta
	newSubTable["idRangeOffset"] = idRangeOffset
	newSubTable["glyphIndexArray"] = []uint16{}

	return newSubTable
}

// buildCmapFormat12 builds a format 12 cmap subtable
func buildCmapFormat12(oldSubTable map[string]interface{}, unicodeToGlyph map[int]int, unicodes []int) map[string]interface{} {
	// Build groups
	var groups []*CmapFormat8nGroup

	if len(unicodes) > 0 {
		segStart := uint32(unicodes[0])
		segEnd := uint32(unicodes[0])
		firstGlyph := uint32(unicodeToGlyph[unicodes[0]])

		for i := 1; i < len(unicodes); i++ {
			u := uint32(unicodes[i])
			g := uint32(unicodeToGlyph[unicodes[i]])
			expectedGlyph := firstGlyph + (u - segStart)

			if u == segEnd+1 && g == expectedGlyph {
				segEnd = u
			} else {
				groups = append(groups, &CmapFormat8nGroup{
					StartCharCode:  segStart,
					EndCharCode:    segEnd,
					StartGlyphCode: firstGlyph,
				})
				segStart = u
				segEnd = u
				firstGlyph = g
			}
		}
		groups = append(groups, &CmapFormat8nGroup{
			StartCharCode:  segStart,
			EndCharCode:    segEnd,
			StartGlyphCode: firstGlyph,
		})
	}

	// Calculate length: 16 (header) + nGroups * 12
	length := uint32(16 + len(groups)*12)

	newSubTable := make(map[string]interface{})
	newSubTable["platformID"] = oldSubTable["platformID"]
	newSubTable["platformSpecificID"] = oldSubTable["platformSpecificID"]
	newSubTable["format"] = uint16(12)
	newSubTable["reserved"] = uint16(0)
	newSubTable["length"] = length
	newSubTable["language"] = oldSubTable["language"]
	newSubTable["nGroups"] = uint32(len(groups))
	newSubTable["groups"] = groups

	return newSubTable
}
