package srm

import (
	"io/fs"
	"path/filepath"
  "os"
)

func CheckExtension(path string, ext string) bool {
  return filepath.Ext(path) == ext
}

// Walks the directory recursively and returns a list of paths to executables under it.
func GetExecsIn(path string) ([]string, error) {
	var execs []string

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && CheckExtension(path, ".exe") {
			execs = append(execs, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return execs, nil
}

func WriteToFile(bytes []byte, path string) error {
  file, err := os.Create("manifest.json")
  if err != nil {
    return err
  }

  defer file.Close()
  file.Write(bytes)
  return nil
}
