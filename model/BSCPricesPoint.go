package model

type BSCPricePoints []BSCPricePoint

type BSCPricePoint struct {
	Price float64
	Block BlockchainBlock
	CMCId int64
}
