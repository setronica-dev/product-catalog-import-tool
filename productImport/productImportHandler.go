package productImport

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"os"
	"regexp"
	"ts/adapters"
	"ts/config"
	"ts/productImport/attribute"
	"ts/productImport/mapping"
	"ts/productImport/ontologyRead"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/ontologyValidator"
	"ts/productImport/product"
	"ts/productImport/reports"
	"ts/productImport/tradeshiftImportHandler"
	"ts/utils"
)

const stageName = "Product Validation Import stage"

type ProductImportHandler struct {
	config           *config.Config
	mapHandler       mapping.MappingHandlerInterface
	rulesHandler     *ontologyRead.RulesHandler
	productHandler   product.ProductHandlerInterface
	attributeHandler attribute.AttributeHandlerInterface
	handler          adapters.HandlerInterface
	validator        ontologyValidator.ValidatorInterface
	reports          *reports.ReportsHandler
	fileManager      *adapters.FileManager
	importHandler    *tradeshiftImportHandler.TradeshiftHandler
	columnMap        *ColumnMap
}

type ColumnMap struct {
	Category  string
	ProductID string
	Name      string
}

type Deps struct {
	dig.In
	Config           *config.Config
	MapHandler       mapping.MappingHandlerInterface
	RulesHandler     *ontologyRead.RulesHandler
	ProductHandler   product.ProductHandlerInterface
	AttributeHandler attribute.AttributeHandlerInterface
	Handler          adapters.HandlerInterface
	Validator        ontologyValidator.ValidatorInterface
	Reports          *reports.ReportsHandler
	FileManager      *adapters.FileManager
	ImportHandler    *tradeshiftImportHandler.TradeshiftHandler
}

func NewProductImportHandler(deps Deps) *ProductImportHandler {
	m := deps.MapHandler.GetColumnMapConfig()
	return &ProductImportHandler{
		config:           deps.Config,
		mapHandler:       deps.MapHandler,
		rulesHandler:     deps.RulesHandler,
		attributeHandler: deps.AttributeHandler,
		productHandler:   deps.ProductHandler,
		handler:          deps.Handler,
		validator:        deps.Validator,
		reports:          deps.Reports,
		fileManager:      deps.FileManager,
		importHandler:    deps.ImportHandler,
		columnMap: &ColumnMap{
			ProductID: m.ProductID,
			Category:  m.Category,
		},
	}
}

func (ph *ProductImportHandler) RunCSV() {
	//ontology
	var rulesConfig *models.OntologyConfig

	rulesConfig, err := ph.rulesHandler.InitRulesConfig()
	if err != nil {
		log.Fatalf("ontology was not specified: %v", err)
	}

	// mappings
	columnMap := ph.mapHandler.Get()

	// feed
	err = ph.runProductValidationImportFlow(columnMap, rulesConfig)
	if err != nil {
		log.Printf("Unexpected error: %v", err)
	}
}

func (ph *ProductImportHandler) runProductValidationImportFlow(columnMap map[string]string, rulesConfig *models.OntologyConfig) error {
	// if something in progress

	var processedSource []string
	inProgress := adapters.GetFiles(ph.config.ProductCatalog.InProgressPath)
	sources := adapters.GetFiles(ph.config.ProductCatalog.SourcePath)
	// identify fitting report
	if len(inProgress) > 0 {
		for _, processingFile := range inProgress {
			reportFile := findAttributeReport(processingFile, utils.SliceDiff(sources, processedSource))
			if reportFile != "" {
				processedSource = append(processedSource, reportFile)
				ph.processFeed(
					ph.config.ProductCatalog.InProgressPath+"/"+processingFile, //feed
					ph.config.ProductCatalog.SourcePath+reportFile,             //report
					columnMap,
					rulesConfig,
					false,
				)
			} else {
				log.Printf("You have the failed feed in progress '%v'. "+
					"Please check the failure report in '%v', "+
					"fill it with the data and appload to the '%v' folder.",
					ph.config.ProductCatalog.InProgressPath+"/"+processingFile,
					ph.config.ProductCatalog.ReportPath,
					ph.config.ProductCatalog.SourcePath)
			}
		}
	} else if len(sources) == 0 {
		return fmt.Errorf("PRODUCT SOURCE FILES WAS NOT FOUND")
	}

	for _, source := range sources {
		if inArr, _ := utils.InArray(source, processedSource); !inArr {
			ph.processFeed(
				ph.config.ProductCatalog.SourcePath+source,
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
	attributesPath string,
	columnMap map[string]string,
	ruleConfig *models.OntologyConfig,
	isInitial bool,
) {
	var hasErrors bool
	var feed []reports.Report

	// source
	log.Println("_________________________________")
	log.Printf("PROCESSING SOURCE: '%v'", sourceFeedPath)
	var er error

	if sourceFeedPath, er = adapters.MoveToPath(sourceFeedPath, ph.config.ProductCatalog.InProgressPath); er != nil {
		log.Printf("ERROR COPYING THE '%v' FILE to the  '%v' folder", sourceFeedPath, ph.config.ProductCatalog.InProgressPath)
	}
	parsedSourceData, err := ph.productHandler.InitSourceData(sourceFeedPath)
	if err != nil {
		wrongFilePath, _ := adapters.MoveToPath(sourceFeedPath, ph.config.ProductCatalog.SentPath)
		log.Printf("an error occured while reading products data. File was moved to %v.\nReason: %v", wrongFilePath, err)
		return
	}

	if isInitial {
		feed, hasErrors, err = ph.runInitialOntologyValidation(sourceFeedPath, parsedSourceData, columnMap, ruleConfig)
		if err != nil {
			log.Printf("ontology validation for products has been failed: %v", err)
		}
	} else {
		feed, hasErrors, err = ph.runSecondaryOntologyValidation(sourceFeedPath, attributesPath, parsedSourceData, columnMap, ruleConfig)
		if err != nil {
			log.Printf("ontology validation for attributes and products has been failed: %v", err)
		}
	}

	attributesReportPath := ph.reports.WriteReport(sourceFeedPath, hasErrors, feed, parsedSourceData)

	if !hasErrors {
		log.Println("Product import to Tradeshift has been started")
		er := ph.importHandler.ImportFeedToTradeshift(attributesReportPath)
		if er != nil {
			log.Printf("Product import to Tradeshift has been failed, report was not built. Reason: %v", er)
		}
	}
}

func (ph *ProductImportHandler) runInitialOntologyValidation(
	sourceFeedPath string,
	parsedSourceData []map[string]interface{},
	columnMap map[string]string,
	ruleConfig *models.OntologyConfig) ([]reports.Report, bool, error) {

	// validation feed
	feed, hasErrors := ph.validator.InitialValidation(
		columnMap,
		ruleConfig,
		parsedSourceData,
	)

	if !hasErrors {
		log.Printf("DATA IS VALID. Please check the result here '%s'", ph.config.ProductCatalog.SentPath)
		if _, er := adapters.MoveToPath(sourceFeedPath, ph.config.ProductCatalog.SentPath); er != nil {
			log.Printf("ERROR COPYING THE SOURCE FILE %v to the '%v' folder", sourceFeedPath, ph.config.ProductCatalog.SentPath)
		}

	} else {
		log.Printf("The validation has found inconsistency in your attributes based on rules. "+
			"You can find the report here '%v'. You can apply corrections right into this report and upload it "+
			"into the source folder %v to continue the process.",
			ph.config.ProductCatalog.ReportPath,
			ph.config.ProductCatalog.SourcePath)
	}
	return feed, hasErrors, nil
}

func (ph *ProductImportHandler) runSecondaryOntologyValidation(
	sourceFeedPath string,
	attributesPath string,
	parsedSourceData []map[string]interface{},
	columnMap map[string]string,
	ruleConfig *models.OntologyConfig,
) ([]reports.Report, bool, error) {

	var er error
	var attributeReportData []*attribute.Attribute
	if attributesPath != "" {
		log.Printf("PROCESSING REPORT: '%v'", attributesPath)
		if attributesPath, er = adapters.MoveToPath(attributesPath, ph.config.ProductCatalog.InProgressPath); er != nil {
			log.Printf("ERROR COPYING THE '%v' FILE to the  '%v' folder", attributesPath, ph.config.ProductCatalog.InProgressPath)
		}
	} else {
		return nil, true, fmt.Errorf("empty attributes filename has been detected: %v", attributesPath)
	}
	// fixed attributes
	attributeReportData, err := ph.attributeHandler.InitAttributeData(attributesPath)
	if err != nil {
		wrongAttributesPath, _ := adapters.MoveToPath(attributesPath, ph.config.ProductCatalog.SentPath)
		return nil, true, fmt.Errorf("an error occured while reading attributes report. File was moved to %v.\n"+
			"Reason: %v", wrongAttributesPath, err)
	}

	// validation feed
	feed, hasErrors := ph.validator.SecondaryValidation(
		columnMap,
		ruleConfig,
		parsedSourceData,
		attributeReportData,
	)

	if !hasErrors {
		log.Printf("DATA IS VALID. Please check the result here '%s'", ph.config.ProductCatalog.SentPath)
		if _, er = adapters.MoveToPath(sourceFeedPath, ph.config.ProductCatalog.SentPath); er != nil {
			log.Printf("ERROR COPYING THE SOURCE FILE %v to the '%v' folder", sourceFeedPath, ph.config.ProductCatalog.SentPath)
		}

		if _, er = adapters.MoveToPath(attributesPath, ph.config.ProductCatalog.SentPath); er != nil {
			log.Printf("ERROR COPYING THE REPORT FILE %v to the '%v' folder", attributesPath, ph.config.ProductCatalog.SentPath)
		}

	} else {
		log.Printf("The validation has found inconsistency in your attributes based on rules. "+
			"You can find the report here '%v'. You can apply corrections right into this report and upload it "+
			"into the source folder %v to continue the process.",
			ph.config.ProductCatalog.ReportPath,
			ph.config.ProductCatalog.SourcePath)
		if attributesPath != "" {
			e := os.Remove(attributesPath)

			if e != nil {
				log.Println(e)
			}
		}
	}

	cleanUpAttributeReports(sourceFeedPath, ph.config.ProductCatalog.ReportPath)
	return feed, hasErrors, nil
}

func findAttributeReport(inProgressFile string, sources []string) string {
	report := ""
	pattern := adapters.GetFileName(inProgressFile)

	for _, source := range sources {
		regexp, _ := regexp.Compile(`(_attributes)`)
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

func cleanUpAttributeReports(sourceFile string, folder string) {
	reportsList := adapters.GetFiles(folder)
	for _, source := range reportsList {
		del := findAttributeReport(sourceFile, []string{source})
		if del != "" {
			e := os.Remove(folder + "/" + del)
			if e != nil {
				log.Println(e)
			}
		}
	}
}
