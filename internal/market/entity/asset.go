package entity

type Asset struct {
	id           string
	name         string
	marketVolume int
}

func NewAsset(id string, name string, marketVolume int) *Asset {
	return &Asset{
		id,
		name,
		marketVolume,
	}
}
