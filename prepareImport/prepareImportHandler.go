package prepareImport

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"path/filepath"
	"strings"
	"ts/adapters"
	"ts/config"
)

type Handler struct {
	sourcePath        string
	sentPath          string
	productConverter  *XLSXSheetToCSVConverter
	failuresConverter *XLSXSheetToCSVConverter
	offerConverter    *XLSXSheetToCSVConverter
}

type Deps struct {
	dig.In
	Config *config.Config
}

func NewPrepareImportHandler(deps Deps) *Handler {
	conf := deps.Config
	commonConf := deps.Config.CommonConfig
	return &Handler{
		sourcePath: commonConf.SourcePath,
		sentPath:   commonConf.SentPath,
		productConverter: NewXLSXSheetToCSVConverter(
			commonConf.Sheet.Products,
			conf.ProductCatalog.InProgressPath,
			""),
		failuresConverter: NewXLSXSheetToCSVConverter(
			commonConf.Sheet.Failures,
			conf.ProductCatalog.SourcePath,
			"-failures"),
		offerConverter: NewXLSXSheetToCSVConverter(
			commonConf.Sheet.Offers,
			conf.OfferCatalog.SourcePath,
			""),
	}
}

func (h *Handler) Run() {
	files := getXLSXFiles(h.sourcePath)
	if len(files) == 0 {
		log.Printf("no xlsx files for imports are specified in %v", h.sourcePath)
		return
	}

	for _, fileName := range files {
		filePath := filepath.Join(
			h.sourcePath,
			fileName)
		err := h.convertSheetsData(filePath)
		if err != nil {
			log.Printf("failed to process file %v: %v", filePath, err)
		}
		_, err = adapters.MoveToPath(filePath, h.sentPath)
		if err != nil {
			log.Printf("failed to move %v to %v: %v", filePath, h.sentPath, err)
		}
	}
}

func getXLSXFiles(path string) []string {
	var res []string
	files := adapters.GetFiles(path)
	for _, filePath := range files {
		if isXLSX(filePath) {
			res = append(res, filePath)
		}
	}
	return res
}

func isXLSX(filePath string) bool {
	res := strings.HasSuffix(strings.ToLower(filePath), ".xls") || strings.HasSuffix(strings.ToLower(filePath), ".xlsx")
	return res
}

func (h *Handler) convertSheetsData(filePath string) error {
	err := h.productConverter.Convert(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert Products: %v", err)
	}
	err = h.offerConverter.Convert(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert Offers: %v", err)
	}
	err = h.failuresConverter.Convert(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert Attributes: %v", err)
	}
	return nil
}
