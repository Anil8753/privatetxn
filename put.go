package privatetxn

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PutState stores the private data in the orgs' implicit private collection and corresponding private data key on the public ledger.
func PutState(ctx contractapi.TransactionContextInterface, orgsMsp []string, key string, data []byte) error {

	// normalize the JSON structure keys in sorted order.
	// Otherwise couchdb will reorder the keys and returns the sorted keys on Get call.
	// This causes GetPrivateDataHash hash vaidation failure
	normalizedData, err := jsonRemarshal(data)
	if err != nil {
		return fmt.Errorf("failed to JsonRemarshal. %w", err)
	}

	dataHash := getHash(normalizedData)

	// This key must be unique in private date collection
	pvtCollectionKey, err := getPrivateCollectionKey(ctx, dataHash)
	if err != nil {
		return fmt.Errorf("failed to private colleection key. %w", err)
	}

	for _, orgMsp := range orgsMsp {

		if err := putPDC(ctx, orgMsp, pvtCollectionKey, normalizedData); err != nil {
			return err
		}
	}

	if err := ctx.GetStub().PutState(key, []byte(pvtCollectionKey)); err != nil {
		return fmt.Errorf("failed to save public state. %w", err)
	}

	return nil
}
