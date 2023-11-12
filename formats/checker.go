package formats

import (
	"path"
	"strings"
)

type KnownFormatType string

const (
	UnknownFormatType    KnownFormatType = "unknown"
	GarminFitFormatType  KnownFormatType = "fit"
	GoProVideoFormatType KnownFormatType = "gopro"
)

// TODO: implement
func IsGarminFitFile(filePath string) bool {
	return strings.ToLower(path.Ext(filePath)) == ".fit"
}

// TODO: implement
func IsGoProVideoFile(filePath string) bool {
	// for now we assume that all mp4 files are GoPro videos
	return strings.ToLower(path.Ext(filePath)) == ".mp4"
}

func CheckFormatType(filePath string) KnownFormatType {
	if IsGarminFitFile(filePath) {
		return GarminFitFormatType
	}
	if IsGoProVideoFile(filePath) {
		return GoProVideoFormatType
	}
	return UnknownFormatType
}
