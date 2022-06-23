package models

import "github.com/everFinance/goar"

func InitArweave(key_file string, ar_node string) (*goar.Wallet, error) {
	wallet, err := goar.NewWalletFromPath(key_file, ar_node)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
