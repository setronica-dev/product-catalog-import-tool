package offerItemImport


type OfferItemImportHandlerInterface interface {
	Run()
}

type OfferItemMappingHandlerInterface interface {
	Run() error
}
