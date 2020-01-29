package node

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/INFURA/go-ethlibs/eth"
	"github.com/INFURA/go-ethlibs/jsonrpc"
	"github.com/INFURA/go-ethlibs/node"
	"github.com/pkg/errors"
)

// CustomClient ddd
type CustomClient struct {
	node.Client
}

// CallParams parameters for eth_call
type CallParams struct {
	From     eth.Data
	To       eth.Data
	Gas      eth.Quantity
	Value    eth.Quantity
	Data     eth.Data
	GasPrice eth.Quantity
}

// GetNewCustomClient returns a new ethereum client
func GetNewCustomClient(target string) (CustomClient, error) {
	ctx := context.Background()
	client, err := node.NewClient(ctx, target)
	if err != nil {
		return CustomClient{nil}, err
	}
	return CustomClient{client}, nil
}

// TransactionByBlockNumberAndIndex get transaction based on its ID in a block designed by its hash
func (c *CustomClient) TransactionByBlockNumberAndIndex(ctx context.Context, number uint64, index uint64) (*eth.Transaction, error) {
	n := eth.QuantityFromUInt64(number)
	i := eth.QuantityFromUInt64(index)
	request := jsonrpc.Request{
		ID:     jsonrpc.ID{Num: 1},
		Method: "eth_getTransactionByBlockNumberAndIndex",
		Params: jsonrpc.MustParams(&n, &i),
	}

	response, err := c.Request(ctx, &request)
	if err != nil {
		return nil, errors.Wrap(err, "could not make  request")
	}

	if len(response.Result) == 0 || bytes.Equal(response.Result, json.RawMessage(`null`)) {
		// Then the transaction isn't recognized
		return nil, node.ErrTransactionNotFound
	}

	tx := eth.Transaction{}
	err = tx.UnmarshalJSON(response.Result)
	return &tx, err
}

// GetBalance balance of an address from latest state
func (c *CustomClient) GetBalance(ctx context.Context, address string) (uint64, error) {
	request := jsonrpc.Request{
		ID:     jsonrpc.ID{Num: 1},
		Method: "eth_getBalance",
		Params: jsonrpc.MustParams(address, "latest"),
	}

	response, err := c.Request(ctx, &request)
	if err != nil {
		return 0, errors.Wrap(err, "could not make  request")
	}

	if response.Error != nil {
		return 0, errors.New(string(*response.Error))
	}

	tx := eth.Quantity{}
	err = tx.UnmarshalJSON(response.Result)
	return tx.UInt64(), err
}

// GetGasPrice get current gas price
func (c *CustomClient) GetGasPrice(ctx context.Context) (uint64, error) {
	request := jsonrpc.Request{
		ID:     jsonrpc.ID{Num: 1},
		Method: "eth_gasPrice",
	}

	response, err := c.Request(ctx, &request)
	if err != nil {
		return 0, errors.Wrap(err, "could not make  request")
	}

	if response.Error != nil {
		return 0, errors.New(string(*response.Error))
	}

	tx := eth.Quantity{}
	err = tx.UnmarshalJSON(response.Result)
	return tx.UInt64(), err
}

// CallContract from latest state
func (c *CustomClient) CallContract(ctx context.Context, param CallParams) (string, error) {
	request := jsonrpc.Request{
		ID:     jsonrpc.ID{Num: 1},
		Method: "eth_call",
		Params: jsonrpc.MustParams(&param, "latest"),
	}

	response, err := c.Request(ctx, &request)
	if err != nil {
		return "", errors.Wrap(err, "could not make  request")
	}

	if response.Error != nil {
		return "", errors.New(string(*response.Error))
	}

	var tx eth.Data
	err = tx.UnmarshalJSON(response.Result)
	return string(tx), err
}
