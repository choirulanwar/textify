package file

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/choirulanwar/textify/backend/service"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	App   *service.App
	exDir string
}

func New(a *service.App) *Service {
	return &Service{
		App: a,
	}
}

func (s *Service) SaveJson(fileName string, jsonData any) error {
	text, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.exDir+fileName, text, 0644); err != nil {
		return err
	}
	return nil
}

func (s *Service) ReadJson(fileName string) (any, error) {
	file, err := os.ReadFile(s.exDir + fileName)
	if err != nil {
		return nil, err
	}

	var data any
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) FileExists(fileName string) bool {
	_, err := os.Stat(s.exDir + fileName)
	return err == nil
}

type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime string `json:"modTime"`
}

func (s *Service) ReadFileInfo(fileName string) (FileInfo, error) {
	info, err := os.Stat(s.exDir + fileName)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime().Format(time.RFC3339),
	}, nil
}

func (s *Service) ListDirFiles(dirPath string) ([]FileInfo, error) {
	files, err := os.ReadDir(s.exDir + dirPath)
	if err != nil {
		return nil, err
	}

	var filesInfo []FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		filesInfo = append(filesInfo, FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime().Format(time.RFC3339),
		})
	}
	return filesInfo, nil
}

func (s *Service) DeleteFile(path string) error {
	err := os.Remove(s.exDir + path)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CopyFile(src string, dst string) error {
	sourceFile, err := os.Open(s.exDir + src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	err = os.MkdirAll(s.exDir+dst[:strings.LastIndex(dst, "/")], 0755)
	if err != nil {
		return err
	}

	destFile, err := os.Create(s.exDir + dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OpenSaveFileDialog(filterPattern string, defaultFileName string, savedContent string) (string, error) {
	return s.OpenSaveFileDialogBytes(filterPattern, defaultFileName, []byte(savedContent))
}

func (s *Service) OpenSaveFileDialogBytes(filterPattern string, defaultFileName string, savedContent []byte) (string, error) {
	path, err := wruntime.SaveFileDialog(s.App.Ctx, wruntime.SaveDialogOptions{
		DefaultFilename: defaultFileName,
		Filters: []wruntime.FileFilter{{
			Pattern: filterPattern,
		}},
		CanCreateDirectories: true,
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	if err := os.WriteFile(path, savedContent, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func (s *Service) OpenFileFolder(path string, relative bool) error {
	var absPath string
	var err error
	if relative {
		absPath, err = filepath.Abs(s.exDir + path)
	} else {
		absPath, err = filepath.Abs(path)
	}
	if err != nil {
		return err
	}
	switch os := runtime.GOOS; os {
	case "windows":
		cmd := exec.Command("explorer", "/select,", absPath)
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	case "darwin":
		cmd := exec.Command("open", "-R", absPath)
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	case "linux":
		cmd := exec.Command("xdg-open", absPath)
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("unsupported OS")
}
