package fileex

import "os"

// MkdirIfNotExist
/**
 * @Description: path if not exist then create this path (dir)
 * @param path
 * @return error
 */
func MkdirIfNotExist(path string) error {
	b, err := PathExistOrNot(path)
	if err != nil {
		return err
	}
	if !b {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}
