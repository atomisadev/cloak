package store

import "io/fs"

type Store interface {
	Write(path string, data []byte) error
	Read(path string) ([]byte, error)
	List() ([]fs.FileInfo, error)
}
