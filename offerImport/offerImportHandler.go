package offerImport

import (
	"fmt"
	"go.uber.org/dig"
	"log"
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

func (o *OfferImportHandler) Run() {
	paths := adapters.GetFiles(o.sourcePath)
	if len(paths) == 0 {
		log.Fatalf("Offer Import failed: please, put file with offers into %v", o.sourcePath)
	}

	for _, path := range paths {
		o.processOffers(path)
	}
}

func (o *OfferImportHandler) processOffers(path string) {
	offers, err := o.uploadOffers(path)
	if err != nil {
		log.Printf("Offer Import failed: source file was not replaced: %v", err)
	}
	o.importHandler.ImportOffers(offers)
}

func (o *OfferImportHandler) uploadOffers(path string) ([]offerReader.RawOffer, error) {
	log.Printf("Offer Import for %v was started", path)

	offers := o.offerReader.UploadOffers(o.getSourcePath(path))
	if len(offers) == 0 {
		return nil, fmt.Errorf("Offer Upload failed: 0 offers were loaded from %v. Please, check file and try again", o.sourcePath)
	}
	err := o.processSourceFile(path)
	return offers, err
}

func (o *OfferImportHandler) processSourceFile(path string) error {
	filePath := o.getSourcePath(path)
	_, err := adapters.MoveToPath(filePath, o.sentPath)
	return err
}

func (o *OfferImportHandler) getSourcePath(path string) string {
	return fmt.Sprintf("%v/%v", o.sourcePath, path)
}
