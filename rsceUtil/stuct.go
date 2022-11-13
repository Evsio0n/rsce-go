/*
 *
 * @Author evsio0n
 * @Date 2022/3/27 下午3:53
 * @Email <bbq2001820@gmail.com>
 *
 */

// Package rsce aka Rockchip resources image.
// It's a binary file which contains device tree blob and additional resources
// (like vendor splash screen) and appears as boot.img-second on unpacking.
package RSCEUtil

import (
	"encoding/binary"
)

const (
	HeaderMagic = "RSCE"
	EntryMagic  = "ENTR"
)

type Header struct {
	Magic             [4]byte //"RSCE"
	RSCEver           uint16  //0x0000
	RSCEfileTblVer    uint16  //0x0000
	HdrBlkSize        byte    //0x01
	FileTblBlkOffset  byte    //0x01
	FileTblRecBlkSize byte    //0x01
	Unknown           byte    //0x00
	FileCount         uint32
	//.... All along calculated as 512 bytes.
	//and all in little endian.
	//Reserved [496]byte
}

type FileEntry struct {
	Magic         [4]byte //"ENTR"
	FileName      [256]byte
	FileBlkOffset uint32
	FileSize      uint32

	//.... All along calculated as 512 bytes.
	//Reserved      [244]byte
}
type fileEntry struct {
	FileName      string
	FileBlkOffset uint32
	FileSize      uint32
}

func endAtFirstZeroByte(src []byte) []byte {
	var strBuf []byte
	for _, v := range src {
		if v != 0 {
			strBuf = append(strBuf, v)
		} else {
			break
		}
	}
	return strBuf
}

func (h *Header) ToBytes() [512]byte {
	var headerBuffer [512]byte
	//copy header to buffer
	//string to byte
	HeaderMagicByte := make([]byte, 4)
	HeaderMagicByte = []byte(HeaderMagic)

	copy(headerBuffer[:4], HeaderMagicByte)
	copy(headerBuffer[4:6], toLittleEndianBytes(uint(h.RSCEver), 2, 16))
	copy(headerBuffer[6:8], toLittleEndianBytes(uint(h.RSCEfileTblVer), 2, 16))
	copy(headerBuffer[8:9], []byte{h.HdrBlkSize})
	copy(headerBuffer[9:10], []byte{h.FileTblBlkOffset})
	copy(headerBuffer[10:11], []byte{h.FileTblRecBlkSize})
	copy(headerBuffer[11:12], []byte{h.Unknown})
	copy(headerBuffer[12:16], toLittleEndianBytes(uint(h.FileCount), 4, 16))
	return headerBuffer
}

func (f *fileEntry) ToBytes() [512]byte {
	var fileEntryBuffer [512]byte
	//copy file entry to buffer
	//string to byte
	EntryMagicByte := make([]byte, 4)
	EntryMagicByte = []byte(EntryMagic)
	//string extend to 256 bytes
	fileName := f.FileName
	//null safe if file name is more than 256 bytes
	if len(fileName) > 256 {
		fileName = fileName[:256]
	}
	var fileNameBuffer [256]byte
	copy(fileNameBuffer[:len(fileName)], fileName)
	copy(fileEntryBuffer[:4], EntryMagicByte)
	copy(fileEntryBuffer[4:260], fileNameBuffer[:])
	copy(fileEntryBuffer[260:264], toLittleEndianBytes(uint(f.FileBlkOffset), 4, 32))
	copy(fileEntryBuffer[264:268], toLittleEndianBytes(uint(f.FileSize), 4, 32))
	return fileEntryBuffer
}

func toLittleEndianBytes(num uint, byteLength int, uintByte int) []byte {
	littleEndianBytes := make([]byte, byteLength)
	//check uint type
	switch uintByte {
	case 16:
		binary.LittleEndian.PutUint16(littleEndianBytes[:], uint16(num))
	case 32:
		binary.LittleEndian.PutUint32(littleEndianBytes[:], uint32(num))
		//default. wes should not be here
	default:
		binary.LittleEndian.PutUint16(littleEndianBytes[:], uint16(num))
	}
	return littleEndianBytes
}
