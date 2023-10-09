package models

type RawData struct {
	Number         int64  `json:"number"`
	TxTrieRoot     string `json:"txTrieRoot"`
	WitnessAddress string `json:"witness_address"`
	ParentHash     string `json:"parentHash"`
	Version        int    `json:"version"`
	Timestamp      int64  `json:"timestamp"`
}

type TransactionContract struct {
	Parameter struct {
		Value struct {
			Data            string `json:"data"`
			OwnerAddress    string `json:"owner_address"`
			ToAddress       string `json:"to_address"`
			ContractAddress string `json:"contract_address"`
		} `json:"value"`
	} `json:"parameter"`
	TypeUrl string `json:"type_url"`
	Type    string `json:"type"`
}

type Transaction struct {
	Ret []struct {
		ContractRet string `json:"contractRet"`
	} `json:"ret"`
	Signature []string `json:"signature"`
	TxID      string   `json:"txID"`
	RawData   struct {
		TransactionRawData []TransactionContract `json:"contract"`
	} `json:"raw_data"`
	RawDataHex string `json:"raw_data_hex"`
}

type Block struct {
	BlockID     string `json:"blockID"`
	BlockHeader struct {
		RawData          RawData `json:"raw_data"`
		WitnessSignature string  `json:"witness_signature"`
	} `json:"block_header"`
	Transactions []Transaction `json:"transactions"`
}
