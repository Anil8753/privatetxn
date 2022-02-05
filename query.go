package privatetxn

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// GetQueryResult makes the query in private data collection and
// retieves the results from the encrypted data stored on public ledger
// Before returning it decryts the data. Secret key is stored in the
// implicit private data collection
func GetQueryResult(
	ctx contractapi.TransactionContextInterface,
	query string,
) ([][]byte, error) {

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("faled to get the client msp. %w", err)
	}

	pvtColelction := getImplicitPrivateCollection(clientMSP)

	itr, err := ctx.GetStub().GetPrivateDataQueryResult(pvtColelction, query)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateDataQueryResult failed. %w", err)
	}

	entries, err := constructResultFromPrivateDataQueryIterator(ctx, itr)
	if err != nil {
		return nil, fmt.Errorf("constructResultFromPrivateDataQueryIterator failed. %w", err)
	}

	return entries, nil
}

// constructResultFromPrivateDataQueryIterator constructs a slice of assets from the resultsIterator
func constructResultFromPrivateDataQueryIterator(
	ctx contractapi.TransactionContextInterface,
	itr shim.StateQueryIteratorInterface,
) ([][]byte, error) {

	entries := make([][]byte, 0)

	for itr.HasNext() {

		queryResult, err := itr.Next()
		if err != nil {
			return nil, err
		}

		var pvtData PrivateData
		err = json.Unmarshal(queryResult.Value, &pvtData)
		if err != nil {
			return nil, err
		}

		b, err := ctx.GetStub().GetState(queryResult.GetKey())
		if err != nil {
			return nil, err
		}

		entryBytes, err := decrypt(pvtData.Secret, b)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entryBytes)
	}

	return entries, nil
}
