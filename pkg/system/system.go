package system

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
	"k8s.io/klog/v2"
)

var FileSystem afero.Fs = afero.NewOsFs()

var Os = runtime.GOOS

func ReadFile(path string) ([]byte, error) {
	f, err := FileSystem.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			klog.Errorf("Error closing file %s: %v", path, err)
		}
	}()
	return io.ReadAll(f)
}

func WriteFile(name string, data []byte, perm os.FileMode) error {
	dirname := filepath.Dir(name)
	err := FileSystem.MkdirAll(dirname, perm)
	if err != nil {
		return err
	}
	f, err := FileSystem.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			klog.Errorf("Error closing file %s: %v", name, err)
		}
	}()
	_, err = f.Write(data)
	return err
}
