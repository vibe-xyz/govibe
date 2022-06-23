package models

import "github.com/ethereum/go-ethereum/crypto"

var (
	Topic_ERC20Created   = crypto.Keccak256Hash([]byte("ERC20Created(address,address)"))
	Topic_ERC721Created  = crypto.Keccak256Hash([]byte("ERC721Created(address,address)"))
	Topic_ERC1155Created = crypto.Keccak256Hash([]byte("ERC1155Created(address,address)"))
	Topic_LogDeploy      = crypto.Keccak256Hash([]byte("LogDeploy(address,bytes,address)"))
)
