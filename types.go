package privatetxn

type PrivateData struct {
	// Secret is used to encrypt the data stored on ledger
	Secret string
	// Tags are used for clouch db rich queries
	Tags interface{}
	// Endorsing Orgs, private data will be shared between the endorsing orgs
	MSPOrgs []string
}
