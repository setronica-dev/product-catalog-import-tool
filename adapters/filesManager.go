package adapters

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//TODO: needs cleanup
type FileManager struct {
	SourcePath              string
	MappingPath             string
	OntologyPath            string
	ReportPath              string
	SentPath                string
	InProgressPath          string
	SuccessResultFolderPath string
	FailResultFolderPath    string
}

func NewFileManager(deps Deps) *FileManager {
	conf := deps.Config.ProductCatalog
	return &FileManager{
		SourcePath:   conf.SourcePath,
		MappingPath:  conf.MappingPath,
		OntologyPath: conf.OntologyPath,
		ReportPath:   conf.ReportPath,

		SentPath:                conf.SentPath,
		InProgressPath:          conf.InProgressPath,
		SuccessResultFolderPath: conf.SuccessResultPath,
		FailResultFolderPath:    conf.ReportPath,
	}
}

func GetFileType(filePath string) FileType {
	ext := strings.TrimLeft(filepath.Ext(filePath), ".")
	switch ext {
	case "csv":
		return CSV
	case "xls":
		return XLSX
	case "xlsx":
		return XLSX
	default:
		return FileType(ext)
	}
}

func (m *FileManager) BuildTradeshiftImportResultsPath(feedPath string) string {
	return fmt.Sprintf("%v/%v", m.ReportPath, m.buildImportResultsFileName(feedPath))
}

func (m *FileManager) buildImportResultsFileName(feedPath string) string {
	sourceFileName := GetFileName(feedPath)
	return fmt.Sprintf("%v_tradeshift-import-results.txt", sourceFileName)
}

//-----------utils-----------
func GetFileName(path string) string {
	r := strings.Split(path, "/")
	l := len(r)
	res := strings.Split(r[l-1], ".")
	return strings.Join(res[0:(len(res)-1)], ".")
}

func MoveFile(oldLocation string, newLocation string) error {
	err := os.Rename(oldLocation, newLocation)
	return err
}

func CopyFile(oldLocation string, newLocation string) error {
	sourceFile, err := os.Open(oldLocation)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create new file
	newFile, err := os.Create(newLocation)
	if err != nil {
		return err
	}
	defer newFile.Close()

	bytesCopied, err := io.Copy(newFile, sourceFile)
	if err != nil {
		return err
	}
	log.Printf("Copied %d bytes.", bytesCopied)
	return nil
}

func GetFiles(folder string) []string {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}
	var file []string
	for _, f := range files {
		if !f.IsDir() && f.Name() != ".git" && f.Name() != ".gitkeep" {
			name := f.Name()
			file = append(file, name)
		}
	}
	return file
}

func MoveToPath(filePath string, destination string) (string, error) {
	fileName := GetFileName(filePath) + filepath.Ext(filePath)
	destinationFile := destination + "/" + fileName
	if err := MoveFile(filePath, destinationFile); err != nil {
		return filePath, err
	}
	return destinationFile, nil
}
