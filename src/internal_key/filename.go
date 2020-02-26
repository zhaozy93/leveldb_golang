package internal_key

import (
	"fmt"
)

func makeFileName(path string, number uint64, suffix string) string {
	return fmt.Sprintf("%s/%06d.%s", path, number, suffix)
}

func TableFileName(path string, number uint64) string {
	return makeFileName(path, number, "ldb")
}

func DescriptorFileName(path string, number uint64) string {
	return fmt.Sprintf("%s/MANIFEST-%.12d", path, number)
}

func CurrentFileName(path string) string {
	return path + "/CURRENT"
}
func TempFileName(path string, number uint64) string {
	return makeFileName(path, number, "dbtmp")
}
