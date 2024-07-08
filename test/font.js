const fs = require("fs");
const assert = require("assert");
const http = require("http");
const path = require("path");

const server = http.createServer((req, res) => {
  res.writeHead(200, { "Content-Type": "application/json" });

  if (req.url === "/getfontjson") {
    fs.readFile(
      path.resolve(__dirname, "./HanyiSentyCrayon.ttf"),
      (err, data) => {
        if (err) throw err;
        //   console.log(data);
        // const fontData = new Uint8Array(data);
        const ttfInfo = new TrueTypeFont(data);
        delete ttfInfo.file.data
        // console.log(ttfInfo)
        let resData = "json stringify error";
        try {
          resData = JSON.stringify(ttfInfo);
        } catch (e) {
          console.log(e);
        }
        console.log("hhhh", resData);
        res.end(resData);
      }
    );
    return;
  }
  res.end("nothing");
});

server.listen(3111, () => {
  console.log("Server is running on http://localhost:3111");
});

function BinaryReader(arrayBuffer) {
  // assert(arrayBuffer instanceof ArrayBuffer);
  this.pos = 0;
  this.data = new Uint8Array(arrayBuffer);
}

BinaryReader.prototype = {
  seek: function (pos) {
    assert(pos >= 0 && pos <= this.data.length);
    var oldPos = this.pos;
    this.pos = pos;
    return oldPos;
  },

  tell: function () {
    return this.pos;
  },

  getUint8: function () {
    //读取单字节无符号整型
    assert(this.pos < this.data.length);
    return this.data[this.pos++];
  },

  getUint16: function () {
    //读取双字节无符号整型
    return ((this.getUint8() << 8) | this.getUint8()) >>> 0;
  },

  getUint32: function () {
    //读取四字节无符号整型
    return this.getInt32() >>> 0;
  },

  getInt16: function () {
    //读取双字节有符号整型
    var result = this.getUint16();
    if (result & 0x8000) {
      result -= 1 << 16;
    }
    return result;
  },

  getInt32: function () {
    //读取四字节有符号整型
    /*
     * 位移运算解析
     *
     * 例如 1
     * 8 位计数法：00000001
     *  00000001000000000000000000000000  左位移24位
     *          000000010000000000000000  左位移16位
     *                  0000000100000000  左位移8位
     *                          00000001  没有位移
     *  00000001000000010000000100000001  进行或运算，结果，转换为十进制值为16843009
     */
    return (
      (this.getUint8() << 24) |
      (this.getUint8() << 16) |
      (this.getUint8() << 8) |
      this.getUint8()
    );
  },

  getFword: function () {
    return this.getInt16();
  },

  get2Dot14: function () {
    //读取定点数，00.00000000000000
    return this.getInt16() / (1 << 14);
  },

  getFixed: function () {
    //读取定点数，00.00
    return this.getInt32() / (1 << 16);
  },

  getString: function (length) {
    //由arraybuffer转字符串（ascii编码）
    var result = "";
    for (var i = 0; i < length; i++) {
      result += String.fromCharCode(this.getUint8());
    }
    return result;
  },

  getDate: function () {
    //读取日期
    var macTime = this.getUint32() * 0x100000000 + this.getUint32();
    var utcTime = macTime * 1000 + Date.UTC(1904, 1, 1);
    return new Date(utcTime);
  },
};

function TrueTypeFont(arrayBuffer) {
  this.file = new BinaryReader(arrayBuffer);
  this.tables = this.readOffsetTables(this.file);
  this.readHeadTable(this.file);
  this.length = this.glyphCount();
}

TrueTypeFont.prototype = {
  glyphCount: function () {
    assert("maxp" in this.tables);
    var old = this.file.seek(this.tables["maxp"].offset + 4);
    var count = this.file.getUint16();
    this.file.seek(old);
    return count;
  },
  readOffsetTables: function (file) {
    var tables = {};
    const offsetTable = {};
    offsetTable.scalerType = file.getUint32();
    var numTables = file.getUint16();
    offsetTable.numTables = numTables;
    offsetTable.searchRange = file.getUint16();
    offsetTable.entrySelector = file.getUint16();
    offsetTable.rangeShift = file.getUint16();
    this.offsetTable = offsetTable;

    for (var i = 0; i < numTables; i++) {
      var tag = file.getString(4);
      tables[tag] = {
        checksum: file.getUint32(),
        offset: file.getUint32(),
        length: file.getUint32(),
      };

      if (tag !== "head") {
        assert(
          this.calculateTableChecksum(
            file,
            tables[tag].offset,
            tables[tag].length
          ) === tables[tag].checksum
        );
      }
    }

    return tables;
  },

  calculateTableChecksum: function (file, offset, length) {
    var old = file.seek(offset);
    var sum = 0;
    var nlongs = ((length + 3) / 4) | 0;
    while (nlongs--) {
      sum = ((sum + file.getUint32()) & 0xffffffff) >>> 0;
    }

    file.seek(old);
    return sum;
  },

  readHeadTable: function (file) {
    assert("head" in this.tables);
    file.seek(this.tables["head"].offset);

    this.version = file.getFixed();
    this.fontRevision = file.getFixed();
    this.checksumAdjustment = file.getUint32();
    this.magicNumber = file.getUint32();
    assert(this.magicNumber === 0x5f0f3cf5);
    this.flags = file.getUint16();
    this.unitsPerEm = file.getUint16();
    this.created = file.getDate();
    this.modified = file.getDate();
    this.xMin = file.getFword();
    this.yMin = file.getFword();
    this.xMax = file.getFword();
    this.yMax = file.getFword();
    this.macStyle = file.getUint16();
    this.lowestRecPPEM = file.getUint16();
    this.fontDirectionHint = file.getInt16();
    this.indexToLocFormat = file.getInt16();
    this.glyphDataFormat = file.getInt16();
  },

  getGlyphOffset: function (index) {
    assert("loca" in this.tables);
    var table = this.tables["loca"];
    var file = this.file;
    var offset, old;

    if (this.indexToLocFormat === 1) {
      old = file.seek(table.offset + index * 4);
      offset = file.getUint32();
    } else {
      old = file.seek(table.offset + index * 2);
      offset = file.getUint16() * 2;
    }

    file.seek(old);

    return offset + this.tables["glyf"].offset;
  },

  readGlyph: function (index) {
    var offset = this.getGlyphOffset(index);
    var file = this.file;

    if (offset >= this.tables["glyf"].offset + this.tables["glyf"].length) {
      return null;
    }

    assert(offset >= this.tables["glyf"].offset);
    assert(offset < this.tables["glyf"].offset + this.tables["glyf"].length);

    file.seek(offset);

    var glyph = {
      numberOfContours: file.getInt16(),
      xMin: file.getFword(),
      yMin: file.getFword(),
      xMax: file.getFword(),
      yMax: file.getFword(),
    };

    assert(glyph.numberOfContours >= -1);

    if (glyph.numberOfContours === -1) {
      this.readCompoundGlyph(file, glyph);
    } else {
      this.readSimpleGlyph(file, glyph);
    }

    return glyph;
  },

  readSimpleGlyph: function (file, glyph) {
    var ON_CURVE = 1,
      X_IS_BYTE = 2,
      Y_IS_BYTE = 4,
      REPEAT = 8,
      X_DELTA = 16,
      Y_DELTA = 32;

    glyph.type = "simple";
    glyph.contourEnds = [];
    var points = (glyph.points = []);

    for (var i = 0; i < glyph.numberOfContours; i++) {
      glyph.contourEnds.push(file.getUint16());
    }

    // skip over intructions
    file.seek(file.getUint16() + file.tell());

    if (glyph.numberOfContours === 0) {
      return;
    }

    var numPoints = Math.max.apply(null, glyph.contourEnds) + 1;

    var flags = [];

    for (i = 0; i < numPoints; i++) {
      var flag = file.getUint8();
      flags.push(flag);
      points.push({
        onCurve: (flag & ON_CURVE) > 0,
      });

      if (flag & REPEAT) {
        var repeatCount = file.getUint8();
        assert(repeatCount > 0);
        i += repeatCount;
        while (repeatCount--) {
          flags.push(flag);
          points.push({
            onCurve: (flag & ON_CURVE) > 0,
          });
        }
      }
    }

    function readCoords(name, byteFlag, deltaFlag, min, max) {
      var value = 0;

      for (var i = 0; i < numPoints; i++) {
        var flag = flags[i];
        if (flag & byteFlag) {
          if (flag & deltaFlag) {
            value += file.getUint8();
          } else {
            value -= file.getUint8();
          }
        } else if (~flag & deltaFlag) {
          value += file.getInt16();
        } else {
          // value is unchanged.
        }

        points[i][name] = value;
      }
    }

    readCoords("x", X_IS_BYTE, X_DELTA, glyph.xMin, glyph.xMax);
    readCoords("y", Y_IS_BYTE, Y_DELTA, glyph.yMin, glyph.yMax);
  },
  drawGlyph: function (index, ctx) {
    var glyph = this.readGlyph(index);

    if (glyph === null || glyph.type !== "simple") {
      return false;
    }

    var p = 0,
      c = 0,
      first = 1;

    while (p < glyph.points.length) {
      var point = glyph.points[p];
      if (first === 1) {
        ctx.moveTo(point.x, point.y);
        first = 0;
      } else {
        ctx.lineTo(point.x, point.y);
      }

      if (p === glyph.contourEnds[c]) {
        c += 1;
        first = 1;
      }

      p += 1;
    }

    return true;
  },
};
