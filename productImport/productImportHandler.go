package productImport

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"os"
	"regexp"
	"ts/adapters"
	"ts/config"
	"ts/externalAPI/tradeshiftAPI"
	"ts/productImport/mapping"
	"ts/productImport/ontologyRead"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/ontologyValidator"
	"ts/productImport/reports"
	"ts/utils"
)

type ProductImportHandler struct {
	config        *config.Config
	mapHandler    mapping.MappingHandlerInterface
	rulesHandler  *ontologyRead.RulesHandler
	handler       adapters.HandlerInterface
	validator     ontologyValidator.ValidatorInterface
	reports       *reports.ReportsHandler
	fileManager   *adapters.FileManager
	importHandler *tradeshiftAPI.TradeshiftHandler
}

type Deps struct {
	dig.In
	Config        *config.Config
	MapHandler    mapping.MappingHandlerInterface
	RulesHandler  *ontologyRead.RulesHandler
	Handler       adapters.HandlerInterface
	Validator     ontologyValidator.ValidatorInterface
	Reports       *reports.ReportsHandler
	FileManager   *adapters.FileManager
	ImportHandler *tradeshiftAPI.TradeshiftHandler
}

func NewProductImportHandler(deps Deps) *ProductImportHandler {
	return &ProductImportHandler{
		config:        deps.Config,
		mapHandler:    deps.MapHandler,
		rulesHandler:  deps.RulesHandler,
		handler:       deps.Handler,
		validator:     deps.Validator,
		reports:       deps.Reports,
		fileManager:   deps.FileManager,
		importHandler: deps.ImportHandler,
	}
}

func (ph *ProductImportHandler) Run() {
	//ontology
	var rulesConfig *models.OntologyConfig

	rulesConfig, err := ph.rulesHandler.InitRulesConfig()
	if err != nil {
		log.Fatalf("ontology was not specified: %v", err)
	}

	// mappings
	columnMap := ph.mapHandler.Init(ph.config.Catalog.MappingPath)

	// feed
	err = ph.processProducts(columnMap, rulesConfig)
	if err != nil {
		log.Println(err)
	}
}

func (ph *ProductImportHandler) processProducts(columnMap map[string]string, rulesConfig *models.OntologyConfig) error {
	// if something in progress

	var processedSource []string
	inProgress := adapters.GetFiles(ph.config.Catalog.InProgressPath)
	sources := adapters.GetFiles(ph.config.Catalog.SourcePath)
	// identify fitting report
	if len(inProgress) > 0 {
		for _, processingFile := range inProgress {
			reportFile := findReport(processingFile, utils.SliceDiff(sources, processedSource))
			if reportFile != "" {
				processedSource = append(processedSource, reportFile)
				ph.processFeed(
					ph.config.Catalog.InProgressPath+"/"+processingFile, //feed
					ph.config.Catalog.SourcePath+reportFile,             //report
					columnMap,
					rulesConfig,
					false,
				)
			} else {
				log.Printf("You have the failed feed in progress '%v'. "+
					"Please check the failure report in '%v', "+
					"fill it with the data and appload to the '%v' folder.",
					ph.config.Catalog.InProgressPath+"/"+processingFile,
					ph.config.Catalog.FailResultPath,
					ph.config.Catalog.SourcePath)
			}
		}
	} else if len(sources) == 0 {
		return fmt.Errorf("SOURCE IS NOT FOUND")
	}

	for _, source := range sources {
		if inArr, _ := utils.InArray(source, processedSource); !inArr {
			ph.processFeed(
				ph.config.Catalog.SourcePath+source,
				"",
				columnMap,
				rulesConfig,
				true,
			)
		}
	}
	return nil
}

func (ph *ProductImportHandler) processFeed(
	sourceFeedPath string,
	validationReportPath string,
	columnMap map[string]string,
	ruleConfig *models.OntologyConfig,
	isInitial bool,
) {
	log.Println("_________________________________")
	log.Printf("PROCESSING SOURCE: %v", sourceFeedPath)
	var er error
	if validationReportPath != "" {
		log.Printf("EDITED REPORT: %v", validationReportPath)
		if validationReportPath, er = adapters.MoveToPath(validationReportPath, ph.config.Catalog.InProgressPath); er != nil {
			log.Printf("ERROR COPYING THE '%v' FILE to the  '%v' folder", validationReportPath, ph.config.Catalog.InProgressPath)
		}
	}
	if isInitial {
		if sourceFeedPath, er = adapters.MoveToPath(sourceFeedPath, ph.config.Catalog.InProgressPath); er != nil {
			log.Printf("ERROR COPYING THE '%v' FILE to the  '%v' folder", sourceFeedPath, ph.config.Catalog.InProgressPath)
		}
	}

	labels := ph.reports.Header
	reportData := make([]*reports.Report, 0)
	if validationReportPath != "" {
		if _, err := os.Stat(validationReportPath); !os.IsNotExist(err) {
			ph.handler.Init(ph.fileManager.GetFileType(validationReportPath))
			reportDataSource := ph.handler.Parse(validationReportPath)
			for _, line := range reportDataSource {
				r := &reports.Report{
					ProductId:    fmt.Sprintf("%v", line[labels.ProductId]),
					Name:         fmt.Sprintf("%v", line[labels.Name]),
					Category:     fmt.Sprintf("%v", line[labels.Category]),
					CategoryName: fmt.Sprintf("%v", line[labels.CategoryName]),
					AttrName:     fmt.Sprintf("%v", line[labels.AttrName]),
					AttrValue:    fmt.Sprintf("%v", line[labels.AttrValue]),
					UoM:          fmt.Sprintf("%v", line[labels.UoM]),
					Errors:       nil,
					Description:  fmt.Sprintf("%v", line[labels.Description]),
					DataType:     fmt.Sprintf("%v", line[labels.DataType]),
					IsMandatory:  fmt.Sprintf("%v", line[labels.IsMandatory]),
					CodedVal:     fmt.Sprintf("%v", line[labels.CodedVal]),
				}
				reportData = append(reportData, r)
			}
		}
	}

	// source
	ph.handler.Init(ph.fileManager.GetFileType(sourceFeedPath))
	parsedData := ph.handler.Parse(sourceFeedPath)

	// validation feed
	feed, hasErrors := ph.validator.Validate(struct {
		Mapping map[string]string
		Rules   *models.OntologyConfig
		Data    []map[string]interface{}
		Report  []*reports.Report
	}{
		Mapping: columnMap,
		Rules:   ruleConfig,
		Data:    parsedData,
		Report:  reportData,
	})


	if !hasErrors {
		log.Printf("SUCCESS: FILE IS VALID. Please check the '%s' folder", ph.config.Catalog.SentPath)
		if _, er = adapters.MoveToPath(sourceFeedPath, ph.config.Catalog.SentPath); er != nil {
			log.Printf("ERROR COPYING THE SOURCE FILE %v to the '%v' folder", sourceFeedPath, ph.config.Catalog.SentPath)
		}

		if validationReportPath != "" {
			if _, er = adapters.MoveToPath(validationReportPath, ph.config.Catalog.SentPath); er != nil {
				log.Printf("ERROR COPYING THE REPORT FILE %v to the '%v' folder", validationReportPath, ph.config.Catalog.SentPath)
			}
		}
	} else {
		log.Printf("FAILURE: check the failure report in '%v', fill it with the data and upload to the '%v' folder.",
			ph.config.Catalog.ReportPath,
			ph.config.Catalog.SourcePath)
		if validationReportPath != "" {
			e := os.Remove(validationReportPath)

			if e != nil {
				log.Println(e)
			}
		}
	}

	cleanUpFailures(sourceFeedPath, ph.config.Catalog.FailResultPath)
	validationReportPath = ph.reports.WriteReport(sourceFeedPath, hasErrors, feed, parsedData, columnMap)
	if !hasErrors {
		log.Println("IMPORT FEED TO TRADESHIFT WAS STARTED")
		er := ph.importHandler.ImportFeedToTradeshift(sourceFeedPath, validationReportPath)
		if er != nil {
			log.Printf("FAILED TO IMPORT VALID FEED TO TRADESHIFT. Reason: %v", er)
		}
	}
}

func findReport(inProgressFile string, sources []string) string {
	report := ""
	pattern := adapters.GetFileName(inProgressFile)

	for _, source := range sources {
		regexp, _ := regexp.Compile(`(-failures)`)
		match := regexp.FindStringIndex(source)
		if len(match) == 2 {
			name := string(source[0:match[0]])
			if name == pattern {
				return source
			}
		}
	}
	return report
}

func cleanUpFailures(sourceFile string, folder string) {
	reports := adapters.GetFiles(folder)
	for _, source := range reports {
		del := findReport(sourceFile, []string{source})
		if del != "" {
			e := os.Remove(folder + "/" + del)
			if e != nil {
				log.Println(e)
			}
		}
	}
}
