package gofile

import (
	"github.com/rs/zerolog/log"
	"os"
)

func IsFile(path string) bool {
	stat, statErr := os.Stat(path)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return false
		}
		log.Warn().Msgf("Unknown error while checking is the file exists %v", path)
		return false
	}
	return !stat.IsDir()
}

func IsDir(path string) bool {
	stat, statErr := os.Stat(path)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return false
		}
		log.Warn().Msgf("Unknown error while checking is the dir exists %v", path)
		return false
	}
	return stat.IsDir()
}

func PathExists(path string) bool {
	_, statErr := os.Stat(path)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return false
		}
		log.Warn().Msgf("Unknown error while checking is the dir exists %v", path)
		return false
	}
	return true
}
