package controller

import (
	"fmt"
	Model "github.com/Adrenesis/crypto-curry-print/model"
	View "github.com/Adrenesis/crypto-curry-print/view"
	"github.com/tyler-sommer/stick"
	"log"
	"net/http"
	"sort"
)

func sortBSCbalancesPriceDecrease(data Model.BSCBalances) {
	sort.Slice(data.Balances[:], func(i, j int) bool {
		return data.Balances[i].USDConvert > data.Balances[j].USDConvert
	})
}

func HandleBSCBalance(w http.ResponseWriter, r *http.Request) {
	addressParam := r.URL.Query()["address"]
	submit := r.URL.Query()["submit"]
	//submit := r.URL.Query()["refresh"]
	//fmt.Println(contractsString)
	//fmt.Println("confirm", confirm)
	//fmt.Println("refresh", refreshAll)
	var bscBalances Model.BSCBalances
	var bscContracts Model.BSCContracts
	var addresses Model.BSCaddresses
	var addresses1 Model.BSCaddresses
	//fmt.Println(nil)
	bscContracts = Model.ReadBSCContractsQLDB("ram")
	if (len(addressParam) > 0) && (len(submit) > 0) {
		fmt.Println(addressParam[0])
		if addressParam[0] != "" {
			addresses1.Addresses = append(addresses1.Addresses, addressParam[0])
		}

		//var bscContracts1 Model.BSCContracts
		//bscContracts1.Contracts = strings.Split(contractsString[0], "\n")
		//for i := 0; i < len(bscContracts1.Contracts); i++ {
		//
		//	bscContracts.Contracts = append(bscContracts.Contracts, bscContracts1.Contracts[i])
		//	//fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[i].CoinDatum))
		//}
		//bscBalances1 := Model.ReadBSCBalancesFromBSCScan(bscContracts)
		//for i := 0; i < len(bscBalances1.Balances); i++ {
		//
		//	bscBalances.Balances = append(bscBalances.Balances, bscBalances1.Balances[i])
		//	//fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[i].CoinDatum))
		//}

		//bscContracts = Model.ReadBSCContractsQLDB()
		//bscBalances = Model.ReadBSCBalancesSQLDB()
		//bscContracts = Model.ReadBSCContractsQLDB()
	}
	Model.CreateBSCaddressesTable("hdd")
	Model.CreateBSCaddressesTable("ram")
	Model.WriteBSCaddressesSQLDB(addresses1, "ram")
	addresses = Model.ReadBSCaddressesSQLDB("ram")
	fmt.Println(fmt.Sprintf("%v", addresses))
	//for i := 0; i<len(addresses1.Addresses); i++ {
	//	addresses.Addresses = append(addresses.Addresses, addresses1.Addresses[i])
	//}

	for j := 0; j < len(addresses.Addresses); j++ {
		bscBalances1 := Model.ReadBSCBalancesFromUnmarshal(addresses.Addresses[j])
		for i := 0; i < len(bscBalances1.Balances); i++ {
			bscBalances.Balances = append(bscBalances.Balances, bscBalances1.Balances[i])
		}
	}
	for i := 0; i < len(bscBalances.Balances); i++ {
		//fmt.Println("test")
		bscBalances.Balances[i].CoinDatum = Model.ReadCryptoByBSCContractSQLDB(bscBalances.Balances[i].Contract, "ram")
		bscBalances.Balances[i].USDConvert = bscBalances.Balances[i].Amount * bscBalances.Balances[i].CoinDatum.Properties.Dollar.Price
		//fmt.Println(bscBalances.Balances[i].CoinDatum.Name)
		//fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[i].CoinDatum.Properties.Dollar.Price))
		//fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[i].CoinDatum.Properties.Dollar.Price))
	}
	//contractsString = ""
	//for i := 0; i < len(bscBalances.Balances); i++ {
	//	if i == 0 {
	//		contractsString = "\"" + bscBalances.Balances[i].Contract + "\""
	//	} else {
	//		contractsString += ",\"" + bscBalances.Balances[i].Contract + "\""
	//	}
	//}
	//coinData := Model.ReadBSCPricesFromBitQuery("0x55d398326f99059ff775485246999027b3197955", contractsString)
	//Model.WriteCryptosByBSCContract(coinData, "ram")
	Model.WriteBSCBalancesSQLDB(bscBalances, "ram")
	//Model.WriteBSCContractsSQLDB(bscContracts, "ram")
	sortBSCbalancesPriceDecrease(bscBalances)
	//fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[0].CoinDatum))
	env := View.GetEnv()
	//fmt.Println(fmt.Sprintf("%v", bscBalances))
	//fmt.Println(fmt.Sprintf("%v", bscContracts))
	p := map[string]stick.Value{"balances": bscBalances, "contracts": bscContracts}
	var err = env.Execute("balances.html.twig", w, p) // Loads "bar.html.twig" relative to fsRoot.
	if err != nil {
		log.Fatal(err)
	}
	//time.Sleep(20 * time.Second)
	//fmt.Println("####################################### 20s after refresh")

}
