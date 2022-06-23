package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

type Batch struct {
	BPS    []*big.Int `abi:"recipientBPS_"`
	Values []string   `abi:"recipients_"`
}

func ABIDecodeString(data []byte) (bps []*big.Int, address []string, err error) {
	uint256ArrTy, err := abi.NewType("uint256[]", "", nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "abi.NewType-uint256ArrTy error")
	}

	addressArrTy, err := abi.NewType("address[]", "", nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "abi.NewType-addressArrTy error")
	}
	// stringTy, err := abi.NewType("string", "", nil)
	// if err != nil {
	// 	return "", errors.Wrap(err, "abi.NewType error")
	// }
	args := abi.Arguments{
		{
			Type: uint256ArrTy,
			Name: "recipientBPS_",
		},
		{
			Type: addressArrTy,
			Name: "recipients_",
		},
	}

	tem := Batch{}
	unpacked, err := args.Unpack(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "args.Unpack error")
	}
	tem, ok := unpacked[0].(Batch)
	if !ok {
		return nil, nil, errors.Wrap(err, "ivs[0] type error")
	}
	bps = tem.BPS
	address = tem.Values
	return bps, address, err
	// var (
	// 	s  string
	// 	ok bool
	// )
	// fmt.Printf("unpacked: %v\n", unpacked...)
	// if len(unpacked) != 0 {
	// 	s, ok = unpacked[0].(string)
	// 	if !ok {
	// 		return "", errors.New("arg is not string.")
	// 	}
	// }

	// return s, nil
}
