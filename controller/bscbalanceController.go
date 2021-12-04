package controller

import (
	"fmt"
	Model "github.com/Adrenesis/crypto-curry-print/model"
	View "github.com/Adrenesis/crypto-curry-print/view"
	"github.com/tyler-sommer/stick"
	"log"
	"net/http"
	"strings"
)

func HandleBSCBalance(w http.ResponseWriter, r *http.Request) {
	contractsString := r.URL.Query()["contracts"]
	submit := r.URL.Query()["submit"]
	fmt.Println(contractsString)
	//fmt.Println("confirm", confirm)
	//fmt.Println("refresh", refreshAll)
	var bscBalances Model.BSCBalances
	var bscContracts Model.BSCContracts
	fmt.Println(nil)
	if (len(contractsString) > 0) && (len(submit) > 0) {
		bscContracts.Contracts = strings.Split(contractsString[0], "\n")
		bscContracts1 := Model.ReadBSCContractsQLDB()
		for i := 0; i < len(bscContracts1.Contracts); i++ {

			bscContracts.Contracts = append(bscContracts.Contracts, bscContracts1.Contracts[i])
			//fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[i].CoinDatum))
		}
		bscBalances = Model.ReadBSCBalancesSQLDB()
		bscBalances1 := Model.ReadBSCBalancesFromBSCScan(bscContracts)
		for i := 0; i < len(bscBalances1.Balances); i++ {

			bscBalances.Balances = append(bscBalances.Balances, bscBalances1.Balances[i])
			//fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[i].CoinDatum))
		}
		for i := 0; i < len(bscBalances.Balances); i++ {
			fmt.Println("test")
			bscBalances.Balances[i].CoinDatum = Model.ReadCryptoByBSCContractSQLDB(bscBalances.Balances[i].Contract)
			//fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[i].CoinDatum))
		}
		Model.WriteBSCBalancesSQLDB(bscBalances)
		Model.WriteBSCContractsSQLDB(bscContracts)
		//bscBalances = Model.ReadBSCBalancesSQLDB()
		bscContracts = Model.ReadBSCContractsQLDB()
	}
	fmt.Println(fmt.Sprintf("%v", bscBalances.Balances[0].CoinDatum))
	env := View.GetEnv()
	fmt.Println(fmt.Sprintf("%v", bscBalances))
	fmt.Println(fmt.Sprintf("%v", bscContracts))
	p := map[string]stick.Value{"balances": bscBalances, "contracts": bscContracts}
	var err = env.Execute("balances.html.twig", w, p) // Loads "bar.html.twig" relative to fsRoot.
	if err != nil {
		log.Fatal(err)
	}

}
