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
	log.Println("_________________________________")
	sourceFileNames := adapters.GetFiles(o.sourcePath)
	if len(sourceFileNames) == 0 {
		log.Println("Source to import offers isn’t found. Skip step.")
		return
	}

	for _, fileName := range sourceFileNames {
		o.runOfferImportFlow(fileName)
	}
}

func (o *OfferImportHandler) runOfferImportFlow(fileName string) {
	offers, err := o.uploadOffers(fileName)
	if err != nil {
		_, _ = adapters.MoveToPath(filepath.Join(o.sourcePath, fileName), o.sentPath)
		log.Printf("An error occurred while uploading the offer: %v. Skip step, invalid file was moved to %v", err, o.sentPath)
		return
	}

	o.importHandler.ImportOffers(offers)
}

func (o *OfferImportHandler) uploadOffers(fileName string) ([]offerReader.RawOffer, error) {
	log.Printf("Offers file processing '%v' has been started", fileName)

	offers := o.offerReader.UploadOffers(filepath.Join(o.sourcePath, fileName))
	if len(offers) == 0 {
		return nil, fmt.Errorf(
			"0 offers were loaded from %v. Please, check file and try again",
			o.sourcePath)
	}
	err := o.processSourceFile(fileName)
	return offers, err
}

func (o *OfferImportHandler) processSourceFile(fileName string) error {
	_, err := adapters.MoveToPath(filepath.Join(o.sourcePath, fileName), o.sentPath)
	return err
}
