package system

import (
	"io"
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
