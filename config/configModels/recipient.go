package configModels

type Recipients struct {
	collection []*recipient
}

type recipient struct {
	id   string
	name string
}

func (r *Recipients) GetRecipientIDByName(name string) string {
	for _, item := range r.collection {
		if item.name == name {
			return item.id
		}
	}
	return ""
}
