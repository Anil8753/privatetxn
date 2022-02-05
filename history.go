package privatetxn

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// GetHistoryForKey returns the history of a key
// It gets the encrypted history of a key and decrypts the history
//  using the secret key stored in the implicit private data colelction of the org
func GetHistoryForKey(
	ctx contractapi.TransactionContextInterface,
	key string,
) ([][]byte, error) {

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("faled to get the client msp. %w", err)
	}

	pvtColelction := getImplicitPrivateCollection(clientMSP)

	pvtBytes, err := ctx.GetStub().GetPrivateData(pvtColelction, key)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateData failed. %w", err)
	}

	if len(pvtBytes) == 0 {
		return nil, fmt.Errorf("private data not found for key %s", key)
	}

	var pvtData PrivateData
	if err := json.Unmarshal(pvtBytes, &pvtBytes); err != nil {
		return nil, fmt.Errorf("failed to get parse private data. %w", err)
	}

	itr, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("GetHistoryForKey failed. %w", err)
	}

	results, err := constructSecretHistoryFromIterator(itr, pvtData.Secret)
	if err != nil {
		return nil, fmt.Errorf("constructSecretHistoryFromIterator failed. %w", err)
	}

	return results, nil
}

// constructSecretHistoryFromIterator constructs a slice of assets from the HistoryQueryIteratorInterface
func constructSecretHistoryFromIterator(
	itr shim.HistoryQueryIteratorInterface,
	secret string,
) ([][]byte, error) {

	entries := make([][]byte, 0)

	for itr.HasNext() {

		result, err := itr.Next()
		if err != nil {
			return nil, err
		}

		entryBytes, err := decrypt(secret, result.Value)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entryBytes)
	}

	return entries, nil
}
