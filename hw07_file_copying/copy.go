package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}

		return err
	}

	defer closeFile(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	size := fileInfo.Size()

	if size <= 0 {
		return nil
	}

	if size < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit < size && limit != 0 {
		size = limit
	}

	if offset > 0 {
		file.Seek(offset, io.SeekStart)
	}
	// N := int(size) + int(offset) // мы заранее знаем сколько хотим прочитать
	N := int(size)
	// create and start new bar
	bar := pb.StartNew(N)
	buf := make([]byte, N) // подготавливаем буфер нужного размера

	file2, err2 := os.Create(toPath) // открываем файл (не забыть про err!)
	if err2 != nil {
		return err2
	}
	defer closeFile(file2)

	offsetT := 0
	for offsetT < N {
		read, err := file.Read(buf[offsetT:])
		offsetT += read
		for range read {
			bar.Increment()
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panicf("failed to read: %v", err)
		}
	}

	for range int(size) - offsetT {
		bar.Increment()
	}

	_, err3 := file2.Write(buf[:offsetT])

	// finish bar
	bar.Finish()

	if err3 != nil {
		log.Panicf("failed to write: %v", err3)
	}

	return nil
}

// func fileGetContent(path string) {

// 	N := int(size)                 // мы заранее знаем сколько хотим прочитать
// 	buf := make([]byte, N)         // подготавливаем буфер нужного размера
// 	file, err := os.Open(fromPath) // открываем файл (не забыть про err!)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return err
// 		}
// 		return err
// 	}

// }

func closeFile(file io.Closer) {
	err := file.Close()
	if err != nil {
		log.Print("file from impossible close")
	}
}
