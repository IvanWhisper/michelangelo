package fileex

import "os"

// PathExistOrNot
/**
 * @Description: path or file exist or not
 * @param path
 * @return bool
 * @return error
 */
func PathExistOrNot(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// PathDelete
/**
 * @Description: del path or file
 * @param path
 * @return error
 */
func PathDelete(path string) error {
	return os.RemoveAll(path)
}
