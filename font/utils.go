package font

import (
	"encoding/binary"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

func getUint8(data []byte) uint8 {
	return data[0]
}

func writeUint8(num uint8) []byte {
	return []byte{num}
}

func getUint16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}

func writeUint16(num uint16) []byte {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, num)
	return data
}

func getUint24(data []byte) (num int, err error) {
	if len(data) < 3 {
		err = errors.New("data to short to read uint24")
		return
	}
	var n uint32
	n |= uint32(data[0]) << 16
	n |= uint32(data[1]) << 8
	n |= uint32(data[2])
	num = int(n)
	return
}

func writeUint24(num int) []byte {
	n := uint32(num)
	return []byte{
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}
}

func getUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func writeUint32(num uint32) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, num)
	return data
}

func getUint64(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}

func writeUint64(num uint64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, num)
	return data
}

func getInt8(data []byte) int8 {
	return int8(data[0])
}

func writeInt8(num int8) []byte {
	return writeUint8(uint8(num))
}

func getInt16(data []byte) int16 {
	return int16(getUint16((data)))
}

func writeInt16(num int16) []byte {
	return writeUint16(uint16(num))
}

func getInt32(data []byte) int32 {
	return int32(getUint32(data))
}

func writeInt32(num int32) []byte {
	return writeUint32(uint32(num))
}

func getInt64(data []byte) int64 {
	return int64(getUint64(data))
}

func writeInt64(num int64) []byte {
	return writeUint64(uint64(num))
}

func getString(data []byte) string {
	return string(data)
}

func writeString(str string) []byte {
	return []byte(str)
}

func getFixed(data []byte) float64 {
	return float64(getInt32(data)) / 65536.0
}

func writeFixed(num float64) []byte {
	return writeInt32(int32(num * 65536.0))
}

func getFixed32(data []byte) uint32 {
	// get fixed32 data
	return getUint32(data)
}

func writeFixed32(num uint32) []byte {
	return writeUint32(num)
}

func getFWord(data []byte) int16 {
	return getInt16(data)
}

func writeFWord(num int16) []byte {
	return writeInt16(num)
}

func getUFWord(data []byte) uint16 {
	return getUint16(data)
}

func writeUFWord(num uint16) []byte {
	return writeUint16(num)
}

func get2Dot14(data []byte) float32 {
	return float32(getInt16(data)) / 16384.0
}

func write2Dot14(num float32) []byte {
	return writeInt16(int16(num * 16384.0))
}

func getLongDateTime(data []byte) int64 {
	// unix second
	longDateTime := getInt64(data)
	starTime := time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()

	unixTime := longDateTime + starTime

	return unixTime
}

func writeLongDateTime(unixTime int64) []byte {
	starTime := time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	longDateTime := unixTime - starTime
	return writeInt64(longDateTime)
}

func getVersion(data []byte) string {
	// 32 bytes
	major := strconv.Itoa(int(getUint16(data[:2])))
	minor := strconv.Itoa(int(getUint16(data[2:4])))
	return major + "." + minor
}

func writeVersion(version string) []byte {
	parts := strings.SplitN(version, ".", 2)
	var major, minor uint16
	if len(parts) > 0 {
		if v, err := strconv.Atoi(parts[0]); err == nil {
			major = uint16(v)
		}
	}
	if len(parts) == 2 {
		if v, err := strconv.Atoi(parts[1]); err == nil {
			minor = uint16(v)
		}
	}
	data := make([]byte, 4)
	binary.BigEndian.PutUint16(data[:2], major)
	binary.BigEndian.PutUint16(data[2:], minor)
	return data
}

func DecodeUTF8(data []byte, offset int, numBytes int) string {
	codePoints := make([]rune, numBytes)
	for j := 0; j < numBytes; j++ {
		codePoints[j] = rune(data[offset+j])
	}
	return string(codePoints)
}

func DecodeUTF16(data []byte, offset int, numBytes int) string {
	// Bounds check
	if offset < 0 || numBytes < 0 || offset+numBytes > len(data) {
		log.Printf("[WARNING] DecodeUTF16: out of bounds, offset=%d numBytes=%d dataLen=%d", offset, numBytes, len(data))
		return ""
	}
	// UTF-16 requires pairs of bytes
	if numBytes%2 != 0 {
		numBytes-- // truncate to even
	}
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
		log.Printf("[WARNING] DecodeMACSTRING: encoding not found for platformSpecific=%s", platformSpecifi)
		return ""
	}

	// Bounds check
	if offset < 0 || dataLength < 0 || offset+dataLength > len(data) {
		log.Printf("[WARNING] DecodeMACSTRING: out of bounds, offset=%d dataLength=%d dataLen=%d", offset, dataLength, len(data))
		return ""
	}

	var result strings.Builder
	result.Grow(dataLength)
	for i := 0; i < dataLength; i++ {
		c := data[offset+i]
		// In all eight-bit Mac encodings, the characters 0x00..0x7F are
		// mapped to U+0000..U+007F; we only need to look up the others.
		if c <= 0x7F {
			result.WriteByte(c)
		} else {
			result.WriteRune(rune(table[c&0x7F]))
		}
	}

	return result.String()
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
