package privatetxn

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PutState saves the data secretly on the ledger.
// How it works:
// 1. Encrypts the data using the secret key passed via PrivateData struct
// 2. Stores the key and encrypted data on public ledger
// 3. Stores the private data (contains secret key) to the implicit private data collection of participating orgs
func PutState(
	ctx contractapi.TransactionContextInterface,
	key string,
	value interface{},
	pvtData *PrivateData,
) error {

	if pvtData.Secret == "" {
		return errors.New("secret is missing in PrivateData")
	}

	pvtBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	outBytes, err := encrypt(pvtData.Secret, pvtBytes)
	if err != nil {
		return err
	}

	if err := ctx.GetStub().PutState(key, outBytes); err != nil {
		return err
	}

	for _, mspId := range pvtData.MSPOrgs {
		if err := ctx.GetStub().PutPrivateData(getImplicitPrivateCollection(mspId), key, pvtBytes); err != nil {
			return err
		}
	}

	return nil
}
