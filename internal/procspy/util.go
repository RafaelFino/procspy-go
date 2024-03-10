package procspy

import (
	"fmt"
	"os"
)

func LoadFile(path string) (string, error) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("[util] file not found")
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("[util] failed to open file")
	}

	defer file.Close()

	buf := make([]byte, stat.Size())
	n, err := file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("[util] failed to read file")
	}

	return string(buf[:n]), nil
}

func WriteFile(path, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("[util] failed to create file")
	}

	defer file.Close()

	wrote, err := file.WriteString(data)
	if err != nil {
		return fmt.Errorf("[util] failed to write file")
	}

	if wrote != len(data) {
		return fmt.Errorf("[util] failed to write all data")
	}

	return nil
}
