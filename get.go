package privatetxn

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// GetState fetchs the private data key from the public ledger and retrive the data from the private data collection of the org.
// At the same time it checks the integrity of the private data by validating the hash of the data with public hash on the ledger.
func GetState(ctx contractapi.TransactionContextInterface, orgMsp string, key string) ([]byte, error) {

	pvtKeyBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get public state. %w", err)
	}

	pvtCollectionKey := string(pvtKeyBytes)

	return getPrivateData(ctx, orgsMsp, pvtCollectionKey)
}
