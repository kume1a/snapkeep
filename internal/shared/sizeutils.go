package shared

import "os"

type FileSizeInUnits struct {
	InBytes     int64
	InKilobytes float64
	InMegabytes float64
	InGigabytes float64
}

func ConvertBytes(sizeBytes int64) FileSizeInUnits {
	return FileSizeInUnits{
		InBytes:     sizeBytes,
		InKilobytes: float64(sizeBytes) / 1024.0,
		InMegabytes: float64(sizeBytes) / (1024.0 * 1024.0),
		InGigabytes: float64(sizeBytes) / (1024.0 * 1024.0 * 1024.0),
	}
}

func GetFileSize(filePath string) (FileSizeInUnits, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return FileSizeInUnits{}, err
	}

	return ConvertBytes(fileInfo.Size()), nil
}
