package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testdataPath := "testdata"
	inputPath := testdataPath + "/input.txt"
	resultPath := testdataPath + "/result.txt"

	t.Run("offset0_limit0", func(t *testing.T) {
		err := Copy(inputPath, resultPath, 0, 0)
		require.Nil(t, err)
		require.Equal(t, fileGetContent(resultPath), fileGetContent(testdataPath+"/out_offset0_limit0.txt"))
	})

	t.Run("offset0_limit10", func(t *testing.T) {
		err := Copy(inputPath, resultPath, 0, 10)
		require.Nil(t, err)
		require.Equal(t, fileGetContent(resultPath), fileGetContent(testdataPath+"/out_offset0_limit10.txt"))
	})

	t.Run("offset0_limit1000", func(t *testing.T) {
		err := Copy(inputPath, resultPath, 0, 1000)
		require.Nil(t, err)
		require.Equal(t, fileGetContent(resultPath), fileGetContent(testdataPath+"/out_offset0_limit1000.txt"))
	})

	t.Run("offset0_limit10000", func(t *testing.T) {
		err := Copy(inputPath, resultPath, 0, 10000)
		require.Nil(t, err)
		require.Equal(t, fileGetContent(resultPath), fileGetContent(testdataPath+"/out_offset0_limit10000.txt"))
	})

	t.Run("offset100_limit1000", func(t *testing.T) {
		err := Copy(inputPath, resultPath, 100, 1000)
		require.Nil(t, err)
		require.Equal(t, fileGetContent(resultPath), fileGetContent(testdataPath+"/out_offset100_limit1000.txt"))
	})

	t.Run("offset6000_limit1000", func(t *testing.T) {
		err := Copy(inputPath, resultPath, 6000, 1000)
		require.Nil(t, err)
		require.Equal(t, fileGetContent(resultPath), fileGetContent(testdataPath+"/out_offset6000_limit1000.txt"))
	})

	os.Remove(resultPath)
}

func fileGetContent(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
