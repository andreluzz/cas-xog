package common

import "os"

func ValidateFolder(folder string) error {
	_, dirErr := os.Stat(folder)
	if os.IsNotExist(dirErr) {
		err := os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}