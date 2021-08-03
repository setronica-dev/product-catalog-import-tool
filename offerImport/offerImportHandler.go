package offerImport

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"path/filepath"
	"ts/adapters"
	"ts/config"
	"ts/offerImport/importHandler"
	"ts/offerImport/offerReader"
)

type OfferImportHandler struct {
	sourcePath    string
	sentPath      string
	offerReader   *offerReader.OfferReader
	importHandler importHandler.ImportOfferInterface
}

type Deps struct {
	dig.In
	Config        *config.Config
	OfferReader   *offerReader.OfferReader
	ImportHandler importHandler.ImportOfferInterface
}

func NewOfferImportHandler(deps Deps) *OfferImportHandler {
	return &OfferImportHandler{
		sourcePath:    deps.Config.OfferCatalog.SourcePath,
		sentPath:      deps.Config.OfferCatalog.SentPath,
		offerReader:   deps.OfferReader,
		importHandler: deps.ImportHandler,
	}
}

func (o *OfferImportHandler) RunCSV() {
	sourceFileNames := adapters.GetFiles(o.sourcePath)
	if len(sourceFileNames) == 0 {
		log.Printf("Offer Import failed: please, put file with offers into %v", o.sourcePath)
		return
	}

	for _, fileName := range sourceFileNames {
		o.runOfferImportFlow(fileName)
	}
}

func (o *OfferImportHandler) runOfferImportFlow(fileName string) {
	offers, err := o.uploadOffers(fileName)
	if err != nil {
		log.Printf("Offer Import failed: source file was not replaced: %v", err)
	}

	o.importHandler.ImportOffers(offers)
}

func (o *OfferImportHandler) uploadOffers(fileName string) ([]offerReader.RawOffer, error) {
	log.Printf("Offer Import for %v was started", fileName)

	offers := o.offerReader.UploadOffers(filepath.Join(o.sourcePath, fileName))
	if len(offers) == 0 {
		return nil, fmt.Errorf(
			"Offer Upload failed: 0 offers were loaded from %v. Please, check file and try again",
			o.sourcePath)
	}
	err := o.processSourceFile(fileName)
	return offers, err
}

func (o *OfferImportHandler) processSourceFile(fileName string) error {
	_, err := adapters.MoveToPath(filepath.Join(o.sourcePath, fileName), o.sentPath)
	return err
}
