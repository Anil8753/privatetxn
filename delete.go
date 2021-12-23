package privatetxn

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func DelState(ctx contractapi.TransactionContextInterface, orgsMsp string, key string) error {

	if err := ctx.GetStub().DelState(key); err != nil {
		return fmt.Errorf("failed to delete key '%s' from the public ledger. %w", key, err)
	}

	return nil
}
