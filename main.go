package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"ts/adapters"
	"ts/config"
	"ts/di"
	"ts/externalAPI/rest"
	"ts/externalAPI/tradeshiftAPI"
	"ts/mapping"
	"ts/models"
	"ts/ontology"
	"ts/ontologyValidator"
	"ts/reports"
	"ts/utils"
)

func main() {
	config.Init()

	// init di build container
	bc := di.BuildContainer()

	// inject stuff and start service
	if err := bc.Invoke(start); err != nil {
		log.Fatalf("instantiation error\n%s", err)
	}
}

func start(
	config *config.Config,
	mapHandler mapping.HandlerInterface,
	httpHandler http.Handler,
	rulesHandler *ontology.RulesHandler,
	handler adapters.HandlerInterface,
	validator ontologyValidator.ValidatorInterface,
	reports *reports.ReportsHandler,
	fileManager *adapters.FileManager,
	rest rest.RestClientInterface,
	tradeshiftAPI *tradeshiftAPI.TradeshiftAPI,
	importHandler *tradeshiftAPI.TradeshiftHandler,
) {
	// TODO: rethink a storage for mappings and rules, move to routines

	//ontology
	var rulesConfig *models.OntologyConfig

	rulesConfig,err := rulesHandler.InitRulesConfig()
	if err!=nil {
		log.Fatalf("ontology was not specified: %v", err)
	}

	// mappings
	var columnMap map[string]string
	if config.Catalog.MappingPath != "" {
		if _, err := os.Stat(config.Catalog.MappingPath); !os.IsNotExist(err) {
			mapHandler.Init(config.Catalog.MappingPath)
			columnMap = mapHandler.Get()
		}
	}

	// if something in progress
	var processedSource []string
	inProgress := adapters.GetFiles(config.Catalog.InProgressPath)
	sources := adapters.GetFiles(config.Catalog.SourcePath)
	// identify fitting report
	if len(inProgress) > 0 {
		for _, processingFile := range inProgress {
			reportFile := findReport(processingFile, utils.SliceDiff(sources, processedSource))
			if reportFile != "" {
				processedSource = append(processedSource, reportFile)
				processFeed(
					config.Catalog.InProgressPath+"/"+processingFile, //feed
					config.Catalog.SourcePath+reportFile,             //report
					columnMap,
					rulesConfig,
					config,
					fileManager,
					handler,
					importHandler,
					reports,
					validator,
					false,
				)
			} else {
				log.Printf("You have the failed feed in progress '%v'. "+
					"Please check the failure report in '%v', "+
					"fill it with the data and appload to the '%v' folder.",
					config.Catalog.InProgressPath+"/"+processingFile,
					config.Catalog.FailResultPath,
					config.Catalog.SourcePath)
			}
		}
	} else if len(sources) == 0 {
		log.Println("SOURCE IS NOT FOUND")
		return
	}

	for _, source := range sources {
		if inArr, _ := utils.InArray(source, processedSource); !inArr {
			processFeed(
				config.Catalog.SourcePath+source,
				"",
				columnMap,
				rulesConfig,
				config,
				fileManager,
				handler,
				importHandler,
				reports,
				validator,
				true,
			)
		}
	}

	return
}

func processFeed(
	sourceFeedPath string,
	validationReportPath string,
	columnMap map[string]string,
	ruleConfig *models.OntologyConfig,
	config *config.Config,
	fileManager *adapters.FileManager,
	handler adapters.HandlerInterface,
	importHandler *tradeshiftAPI.TradeshiftHandler,
	reports *reports.ReportsHandler,
	validator ontologyValidator.ValidatorInterface,
	isInitial bool,
) {
	log.Println("_________________________________")
	log.Printf("PROCESSING SOURCE: %v", sourceFeedPath)
	var er error
	if validationReportPath != "" {
		log.Printf("EDITED REPORT: %v", validationReportPath)
		if validationReportPath, er = adapters.MoveToPath(validationReportPath, config.Catalog.InProgressPath); er != nil {
			log.Printf("ERROR COPYING THE '%v' FILE to the  '%v' folder", validationReportPath, config.Catalog.InProgressPath)
		}
	}
	if isInitial {
		if sourceFeedPath, er = adapters.MoveToPath(sourceFeedPath, config.Catalog.InProgressPath); er != nil {
			log.Printf("ERROR COPYING THE '%v' FILE to the  '%v' folder", sourceFeedPath, config.Catalog.InProgressPath)
		}
	}

	labels := reports.Header
	reportData := make([]*models.Report, 0)
	if validationReportPath != "" {
		if _, err := os.Stat(validationReportPath); !os.IsNotExist(err) {
			handler.Init(fileManager.GetFileType(validationReportPath))
			reportDataSource := handler.Parse(validationReportPath)
			for _, line := range reportDataSource {
				r := &models.Report{
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
	handler.Init(fileManager.GetFileType(sourceFeedPath))
	parsedData := handler.Parse(sourceFeedPath)

	// validation feed
	feed, hasErrors := validator.Validate(struct {
		Mapping map[string]string
		Rules   *models.OntologyConfig
		Data    []map[string]interface{}
		Report  []*models.Report
	}{
		Mapping: columnMap,
		Rules:   ruleConfig,
		Data:    parsedData,
		Report:  reportData,
	})

	if !hasErrors {
		log.Printf("SUCCESS: FILE IS VALID. Please check the '%s' folder", config.Catalog.SentPath)
		if _, er = adapters.MoveToPath(sourceFeedPath, config.Catalog.SentPath); er != nil {
			log.Printf("ERROR COPYING THE SOURCE FILE %v to the '%v' folder", sourceFeedPath, config.Catalog.SentPath)
		}

		if validationReportPath != "" {
			if _, er = adapters.MoveToPath(validationReportPath, config.Catalog.SentPath); er != nil {
				log.Printf("ERROR COPYING THE REPORT FILE %v to the '%v' folder", validationReportPath, config.Catalog.SentPath)
			}
		}
	} else {
		log.Printf("FAILURE: check the failure report in '%v', fill it with the data and upload to the '%v' folder.",
			config.Catalog.ReportPath,
			config.Catalog.SourcePath)
		if validationReportPath != "" {
			e := os.Remove(validationReportPath)

			if e != nil {
				log.Println(e)
			}
		}
	}

	cleanUpFailures(sourceFeedPath, config.Catalog.FailResultPath)
	validationReportPath = reports.WriteReport(sourceFeedPath, hasErrors, feed, parsedData, columnMap)
	if !hasErrors {
		log.Println("IMPORT FEED TO TRADESHIFT WAS STARTED")
		er := importHandler.ImportFeedToTradeshift(sourceFeedPath, validationReportPath)
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
