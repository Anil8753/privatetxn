package privatetxn

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// GetState returns
// 1. Private data stored in the implicit data collection
// 2. Public data which is decypted using the secret key.
//    This secret key is available in the implicit private data collection of the org
func GetState(
	ctx contractapi.TransactionContextInterface,
	key string,
) ([]byte, *PrivateData, error) {

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, nil, fmt.Errorf("faled to get the client msp. %w", err)
	}

	pvtColelction := getImplicitPrivateCollection(clientMSP)

	pvtBytes, err := ctx.GetStub().GetPrivateData(pvtColelction, key)
	if err != nil {
		return nil, nil, fmt.Errorf("GetPrivateData failed. %w", err)
	}

	if len(pvtBytes) == 0 {
		return nil, nil, fmt.Errorf("private data not found for key %s", key)
	}

	var pvtData PrivateData
	if err := json.Unmarshal(pvtBytes, &pvtBytes); err != nil {
		return nil, nil, fmt.Errorf("faled to get parse private data. %w", err)
	}

	pubBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, nil, fmt.Errorf("GetPrivateData failed. %w", err)
	}

	if len(pubBytes) == 0 {
		return nil, nil, fmt.Errorf("public data not found for key %s", key)
	}

	outBytes, err := decrypt(pvtData.Secret, pubBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt failed. %w", err)
	}

	return outBytes, &pvtData, nil
}
