package entity

type Investor struct {
	id            string
	name          string
	assetPosition []*InvestorAssetPostion
}

func NewInvestor(id string) *Investor {
	return &Investor{
		id:            id,
		assetPosition: []*InvestorAssetPostion{},
	}
}

func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPostion) {
	i.assetPosition = append(i.assetPosition, assetPosition)
}

func (i *Investor) UpdateAssetPosition(assetId string, qtShares int) {
	assetPosition := i.GetAssetPosition(assetId)
	if assetPosition == nil {
		i.assetPosition = append(i.assetPosition, NewInvestorAssetPosition(assetId, qtShares))
	} else {
		assetPosition.shares += qtShares
	}
}

func (i *Investor) GetAssetPosition(assetId string) *InvestorAssetPostion {
	for _, assetPosition := range i.assetPosition {
		if assetPosition.assetId == assetId {
			return assetPosition
		}
	}
	return nil
}

type InvestorAssetPostion struct {
	assetId string
	shares  int
}

func NewInvestorAssetPosition(assetId string, shares int) *InvestorAssetPostion {
	return &InvestorAssetPostion{
		assetId: assetId,
		shares:  shares,
	}
}
