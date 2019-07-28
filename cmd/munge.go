package wotwhb

import (
	"path/filepath"
)

var (
	TotalAmountSpent        = float64(0)
	CurrencyList            = NewUniqueStringList()
	SubProductHumanNameList = NewUniqueStringList()
	ProductCategoryList     = NewUniqueStringList()
	ProductHumanNameList    = NewUniqueStringList()
	PayeeHumanNameList      = NewUniqueStringList()
	PlatformList            = NewUniqueStringList()
	ArchList                = NewUniqueStringList()
	TotalFileSize           = int64(0)
	DownloadStructNameList  = NewUniqueStringList()
)

type MungeComponents struct {
	TotalAmountSpent        float64
	TotalFileSize           int64
	CurrencyList            []string
	SubProductHumanNameList []string
	ProductCategoryList     []string
	ProductHumanNameList    []string
	PayeeHumanNameList      []string
	PlatformList            []string
	ArchList                []string
	DownloadStructNameList  []string
}

func traverseAllOrders(orders []HbOrder) *MungeComponents {
	var mungeComponents MungeComponents
	for _, order := range orders {
		TotalAmountSpent += order.AmountSpent
		CurrencyList.Add(order.Currency)
		ProductCategoryList.Add(order.Product.Category)
		ProductHumanNameList.Add(order.Product.HumanName)
		for _, subproduct := range order.SubProducts {
			SubProductHumanNameList.Add(subproduct.HumanName)
			PayeeHumanNameList.Add(subproduct.Payee.HumanName)
			for _, download := range subproduct.Downloads {
				PlatformList.Add(download.Platform)
				for _, downloadStruct := range download.DownloadStruct {
					ArchList.Add(downloadStruct.Arch)
					TotalFileSize += downloadStruct.FileSize
					DownloadStructNameList.Add(downloadStruct.Name)
				}
			}
		}
	}
	mungeComponents.TotalAmountSpent = TotalAmountSpent
	mungeComponents.TotalFileSize = TotalFileSize
	mungeComponents.CurrencyList = CurrencyList.Contents()
	mungeComponents.ProductCategoryList = ProductCategoryList.Contents()
	mungeComponents.ProductHumanNameList = ProductHumanNameList.Contents()
	mungeComponents.SubProductHumanNameList = SubProductHumanNameList.Contents()
	mungeComponents.PayeeHumanNameList = PayeeHumanNameList.Contents()
	mungeComponents.PlatformList = PlatformList.Contents()
	mungeComponents.ArchList = ArchList.Contents()
	mungeComponents.DownloadStructNameList = DownloadStructNameList.Contents()
	writeJsonToFile(mungeComponents, filepath.Join(ConfigDirectoryFlagValue, mungeListFileBasename))
	return &mungeComponents
}
