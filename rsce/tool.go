/*
 *
 * @Author evsio0n
 * @Date 2022/3/27 下午3:53
 * @Email <bbq2001820@gmail.com>
 *
 */

package rsce

import (
	"encoding/binary"
	"io/ioutil"
	"math"
	"os"
	"path"

	"github.com/evsio0n/log"
)

func UnPackRSCE(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Error("Open file  failed: ", err.Error())
		return
	}
	//Check if Header is valid, Header first four bytes should be "RSCE"
	buffer := [512]byte{}
	file.Read(buffer[:])
	if string(buffer[:4]) != HeaderMagic {
		log.Panic("Invalid RSCE header")
	}
	//Find File count in header and parse uint32 little endian offset=0x12 with 4 bytes
	fileToalCount := int(binary.LittleEndian.Uint32(buffer[12:16]))
	log.Info("File count:", fileToalCount)
	//make struct
	//Parse file entries
	//Started after 500 bytes
	var fileEntryOffset int
	fileEntries := make([]fileEntry, fileToalCount)
	fileEntryOffset = 512*fileToalCount + 512
	fileEntryBuffer := make([]byte, fileEntryOffset)
	//read from top again
	file.Seek(0, 0)
	file.Read(fileEntryBuffer[:])
	//cut header buffer 512 bytes
	fileEntryBuffer = fileEntryBuffer[512:]
	for i := 0; i < fileToalCount; i++ {
		//Parse file entry
		//First 4 bytes should be "ENTR"
		if string(fileEntryBuffer[i*512:i*512+4]) != EntryMagic {
			log.Error("Invalid RSCE entry at offset:", i*512)
		}
		//Parse file name
		//File name is 256 bytes and remove null bytes
		strWithZeroByte := fileEntryBuffer[i*512+4 : i*512+260]
		fileName := string(endAtFirstZeroByte(strWithZeroByte))

		log.Info("Found File, name:", fileName)
		//Parse file offset
		//File offset is 4 bytes
		fileOffset := binary.LittleEndian.Uint32(fileEntryBuffer[i*512+260 : i*512+264])
		log.Info("File offset:", fileOffset)
		//Parse file size
		//File size is 4 bytes
		fileSize := binary.LittleEndian.Uint32(fileEntryBuffer[i*512+264 : i*512+268])
		log.Info("File size:", fileSize)
		fileEntries[i] = fileEntry{
			FileName:      fileName,
			FileBlkOffset: fileOffset,
			FileSize:      fileSize,
		}
	}
	//release os.File

	//Save files from file entries
	for _, fileEntry := range fileEntries {
		//Create files
		//Check if files exists
		if _, err := os.Stat(fileEntry.FileName); err == nil {
			log.Info("File \"", fileEntry.FileName, "\" exists, removing...")
			os.Remove(fileEntry.FileName)
		}
		log.Info("Creating files:", fileEntry.FileName)
		files, err := os.Create("./" + fileEntry.FileName)
		if err != nil {
			log.Error("Failed to create files:", fileEntry.FileName, " error:\"", err.Error(), "\"")
		}
		//Seek to file offset
		file.Seek(int64(fileEntry.FileBlkOffset)*512, 0)
		//Read files
		fileBuffer := make([]byte, fileEntry.FileSize)
		file.Read(fileBuffer[:])
		//Write files
		files.Write(fileBuffer[:])
		//Close files
		files.Close()
	}
	file.Close()

}

func GenerateRSCE(contentFilePaths []string, outputFile string) {
	//if output file exists, remove it
	if _, err := os.Stat(outputFile); err == nil {
		log.Info("Output file exists, removing...")
		err := os.Remove(outputFile)
		if err != nil {
			log.Panic("Failed to remove exist output file:", outputFile, " error:", err.Error())
		}
	}
	//Create RSCE file
	rsceFile, err := os.Create(outputFile)
	if err != nil {
		log.Error("Failed to create RSCE file:", outputFile, " error:\"", err.Error(), "\"")
	}
	var allBuffer []byte
	//create little endian uint32

	header := &Header{
		RSCEver:           uint16(0),
		RSCEfileTblVer:    uint16(0),
		HdrBlkSize:        byte(1),
		FileTblBlkOffset:  byte(1),
		FileTblRecBlkSize: byte(1),
		Unknown:           byte(0),
		FileCount:         uint32(len(contentFilePaths)),
	}
	//Write header
	var headerBuffer [512]byte
	//Store all file buffer after rsce buffer
	var allFileBuffer []byte
	//copy header to buffer
	headerBuffer = header.ToBytes()
	//read file from contentFilePaths
	var fileEntries []fileEntry
	for _, contentFilePath := range contentFilePaths {
		//open file
		file, err := os.Open(contentFilePath)
		defer file.Close()
		if err != nil {
			log.Error("Failed to open file:", contentFilePath, " error:\"", err.Error(), "\" Skip this file")
		} else {
			//read file
			var fileBuffer []byte
			if err == nil {
				fileBuffer, err = ioutil.ReadAll(file)
				if err == nil {
					//parse file name without path
					fileName := path.Base(contentFilePath)
					//calculate file size / 512
					fileSize := uint32(len(fileBuffer))
					//append all file buffer to file buffer
					//check if file is multiple of 512
					if fileSize%512 != 0 {
						//add zero bytes
						fileBuffer = append(fileBuffer, make([]byte, 512-fileSize%512)...)
					}

					//append file entry to file entries
					allFileBuffer = append(allFileBuffer, fileBuffer...)
					fileEntries = append(fileEntries, fileEntry{
						FileName: fileName,
						FileSize: fileSize,
					})
				}

			}

		}

	}
	if len(fileEntries) == 0 {
		log.Panic("No file available")
	} else {
		fileEntriesBuffer := make([]byte, 512*len(fileEntries))
		var lastOffset uint32
		//write file entries
		//append to header buffer
		fistFileOffset := len(fileEntries) + 1
		for i, fileEntry := range fileEntries {
			//write file entry
			//copy file entry to buffer
			//calculate file offset fist all entries and block size
			if i == 0 {
				fileEntry.FileBlkOffset = uint32(fistFileOffset)

			} else {
				fileEntry.FileBlkOffset = uint32Plus(lastOffset, fileEntries[i-1].FileSize/512) + 1
			}

			lastOffset = fileEntry.FileBlkOffset
			//append to file entries buffer
			fE := fileEntry.ToBytes()
			copy(fileEntriesBuffer[i*512:], fE[:])
			log.Info("Adding file:", fileEntry.FileName, " file size:", fileEntry.FileSize, " file offset:", fileEntry.FileBlkOffset)

		}
		//Finally, write all buffer together
		allBuffer = append(allBuffer, headerBuffer[:]...)
		allBuffer = append(allBuffer, fileEntriesBuffer[:]...)
		allBuffer = append(allBuffer, allFileBuffer[:]...)
		//write all buffer to rsce file
		rsceFile.Write(allBuffer[:])
		rsceFile.Close()
	}

}

func uint32Plus(a uint32, b uint32) uint32 {
	//using math to add two uint32
	return uint32(math.Pow(2, 32)) + a + b
}
