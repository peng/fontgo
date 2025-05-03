package font

import (
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

func getUint8(data []byte) uint8 {
	return data[0]
}

func getUint16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}

func getUint24(data []byte) (num uint32, err error) {
	if len(data) < 3 {
		err = errors.New("data to short to read uint24")
		return
	}

	num |= uint32(data[0]) << 16
	num |= uint32(data[1]) << 8
	num |= uint32(data[2])

	return
}

func getUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func getUint64(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}

func getInt8(data []byte) int8 {
	return int8(data[0])
}

func getInt16(data []byte) int16 {
	return int16(getUint16((data)))
}

func getInt32(data []byte) int32 {
	return int32(getUint32(data))
}

func getInt64(data []byte) int64 {
	return int64(getUint64(data))
}

func getString(data []byte) string {
	return string(data)
}

func getFixed(data []byte) float64 {
	return float64(getInt32(data) / 65535)
}

func getFWord(data []byte) int16 {
	return getInt16(data)
}

func getUFWord(data []byte) uint16 {
	return getUint16(data)
}

func get2Dot14(data []byte) float32 {
	return float32(getInt16(data) / 16384)
}

func getLongDateTime(data []byte) int64 {
	// unix second
	longDateTime := getInt64(data)
	starTime := time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()

	unixTime := longDateTime + starTime

	return unixTime
}

func getVersion(data []byte) string {
	// 32 bytes
	major := strconv.Itoa(int(getUint16(data[:2])))
	minor := strconv.Itoa(int(getUint16(data[2:4])))
	return major + "." + minor
}

func DecodeUTF8(data []byte, offset int, numBytes int) string {
	codePoints := make([]rune, numBytes)
	for j := 0; j < numBytes; j++ {
		codePoints[j] = rune(data[offset+j])
	}
	return string(codePoints)
}

func DecodeUTF16(data []byte, offset int, numBytes int) string {
	codePoints := make([]uint16, numBytes/2)
	for j := 0; j < len(codePoints); j++ {
		codePoints[j] = getUint16(data[offset : offset+2])
		offset += 2
	}
	runes := utf16.Decode(codePoints)
	return string(runes)
}

var eightBitMacEncodings = map[string]string{
	"x-mac-croatian": "ÄÅÇÉÑÖÜáàâäãåçéèêëíìîïñóòôöõúùûü†°¢£§•¶ß®Š™´¨≠ŽØ∞±≤≥∆µ∂∑∏š∫ªºΩžø" + "¿¡¬√ƒ≈Ć«Č… ÀÃÕŒœĐ—“”‘’÷◊©⁄€‹›Æ»–·‚„‰ÂćÁčÈÍÎÏÌÓÔđÒÚÛÙıˆ˜¯πË˚¸Êæˇ",
	"x-mac-cyrillic": "АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ†°Ґ£§•¶І®©™Ђђ≠Ѓѓ∞±≤≥іµґЈЄєЇїЉљЊњ" + "јЅ¬√ƒ≈∆«»… ЋћЌќѕ–—“”‘’÷„ЎўЏџ№Ёёяабвгдежзийклмнопрстуфхцчшщъыьэю",
	"x-mac-gaelic": // http://unicode.org/Public/MAPPINGS/VENDORS/APPLE/GAELIC.TXT
	"ÄÅÇÉÑÖÜáàâäãåçéèêëíìîïñóòôöõúùûü†°¢£§•¶ß®©™´¨≠ÆØḂ±≤≥ḃĊċḊḋḞḟĠġṀæø" +
		"ṁṖṗɼƒſṠ«»… ÀÃÕŒœ–—“”‘’ṡẛÿŸṪ€‹›Ŷŷṫ·Ỳỳ⁊ÂÊÁËÈÍÎÏÌÓÔ♣ÒÚÛÙıÝýŴŵẄẅẀẁẂẃ",
	"x-mac-greek": // mac_greek
	"Ä¹²É³ÖÜ΅àâä΄¨çéèêë£™îï•½‰ôö¦€ùûü†ΓΔΘΛΞΠß®©ΣΪ§≠°·Α±≤≥¥ΒΕΖΗΙΚΜΦΫΨΩ" +
		"άΝ¬ΟΡ≈Τ«»… ΥΧΆΈœ–―“”‘’÷ΉΊΌΎέήίόΏύαβψδεφγηιξκλμνοπώρστθωςχυζϊϋΐΰ\u00AD",
	"x-mac-icelandic": // mac_iceland
	"ÄÅÇÉÑÖÜáàâäãåçéèêëíìîïñóòôöõúùûüÝ°¢£§•¶ß®©™´¨≠ÆØ∞±≤≥¥µ∂∑∏π∫ªºΩæø" +
		"¿¡¬√ƒ≈∆«»… ÀÃÕŒœ–—“”‘’÷◊ÿŸ⁄€ÐðÞþý·‚„‰ÂÊÁËÈÍÎÏÌÓÔÒÚÛÙıˆ˜¯˘˙˚¸˝˛ˇ",
	"x-mac-inuit": // http://unicode.org/Public/MAPPINGS/VENDORS/APPLE/INUIT.TXT
	"ᐃᐄᐅᐆᐊᐋᐱᐲᐳᐴᐸᐹᑉᑎᑏᑐᑑᑕᑖᑦᑭᑮᑯᑰᑲᑳᒃᒋᒌᒍᒎᒐᒑ°ᒡᒥᒦ•¶ᒧ®©™ᒨᒪᒫᒻᓂᓃᓄᓅᓇᓈᓐᓯᓰᓱᓲᓴᓵᔅᓕᓖᓗ" +
		"ᓘᓚᓛᓪᔨᔩᔪᔫᔭ… ᔮᔾᕕᕖᕗ–—“”‘’ᕘᕙᕚᕝᕆᕇᕈᕉᕋᕌᕐᕿᖀᖁᖂᖃᖄᖅᖏᖐᖑᖒᖓᖔᖕᙱᙲᙳᙴᙵᙶᖖᖠᖡᖢᖣᖤᖥᖦᕼŁł",
	"x-mac-ce": // mac_latin2
	"ÄĀāÉĄÖÜáąČäčĆćéŹźĎíďĒēĖóėôöõúĚěü†°Ę£§•¶ß®©™ę¨≠ģĮįĪ≤≥īĶ∂∑łĻļĽľĹĺŅ" +
		"ņŃ¬√ńŇ∆«»… ňŐÕőŌ–—“”‘’÷◊ōŔŕŘ‹›řŖŗŠ‚„šŚśÁŤťÍŽžŪÓÔūŮÚůŰűŲųÝýķŻŁżĢˇ",
	"macintosh": // mac_roman
	"ÄÅÇÉÑÖÜáàâäãåçéèêëíìîïñóòôöõúùûü†°¢£§•¶ß®©™´¨≠ÆØ∞±≤≥¥µ∂∑∏π∫ªºΩæø" +
		"¿¡¬√ƒ≈∆«»… ÀÃÕŒœ–—“”‘’÷◊ÿŸ⁄€‹›ﬁﬂ‡·‚„‰ÂÊÁËÈÍÎÏÌÓÔÒÚÛÙıˆ˜¯˘˙˚¸˝˛ˇ",
	"x-mac-romanian": // mac_romanian
	"ÄÅÇÉÑÖÜáàâäãåçéèêëíìîïñóòôöõúùûü†°¢£§•¶ß®©™´¨≠ĂȘ∞±≤≥¥µ∂∑∏π∫ªºΩăș" +
		"¿¡¬√ƒ≈∆«»… ÀÃÕŒœ–—“”‘’÷◊ÿŸ⁄€‹›Țț‡·‚„‰ÂÊÁËÈÍÎÏÌÓÔÒÚÛÙıˆ˜¯˘˙˚¸˝˛ˇ",
	"x-mac-turkish": // mac_turkish
	"ÄÅÇÉÑÖÜáàâäãåçéèêëíìîïñóòôöõúùûü†°¢£§•¶ß®©™´¨≠ÆØ∞±≤≥¥µ∂∑∏π∫ªºΩæø" +
		"¿¡¬√ƒ≈∆«»… ÀÃÕŒœ–—“”‘’÷◊ÿŸĞğİıŞş‡·‚„‰ÂÊÁËÈÍÎÏÌÓÔÒÚÛÙˆ˜¯˘˙˚¸˝˛ˇ",
}

func DecodeMACSTRING(data []byte, offset int, dataLength int, platformSpecifi string) string {
	table, exists := eightBitMacEncodings[platformSpecifi]
	if !exists {
		return ""
	}

	var result string
	for i := 0; i < dataLength; i++ {
		c := data[offset+i]
		// In all eight-bit Mac encodings, the characters 0x00..0x7F are
		// mapped to U+0000..U+007F; we only need to look up the others.
		if c <= 0x7F {
			result += string(rune(c))
		} else {
			result += string(table[c&0x7F])
		}
	}

	return result
}

func FromCharCode(data []int) string {
	var str strings.Builder
	for _, val := range data {
		str.WriteRune(rune(val))
	}
	return str.String()
}

func FromCharCodeByte(data []byte) string {
	var toInt []int
	for _, val := range data {
		toInt = append(toInt, int(val))
	}

	return FromCharCode(toInt)
}
