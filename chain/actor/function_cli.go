// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package actor

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FunctionClientMetaData contains all meta data concerning the FunctionClient contract.
var FunctionClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"EmptyRequestData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderIsNotRegistry\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"score\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"functionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"func\",\"type\":\"address\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"score\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// FunctionClientABI is the input ABI used to generate the binding from.
// Deprecated: Use FunctionClientMetaData.ABI instead.
var FunctionClientABI = FunctionClientMetaData.ABI

// FunctionClient is an auto generated Go binding around an Ethereum contract.
type FunctionClient struct {
	FunctionClientCaller     // Read-only binding to the contract
	FunctionClientTransactor // Write-only binding to the contract
	FunctionClientFilterer   // Log filterer for contract events
}

// FunctionClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type FunctionClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunctionClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FunctionClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunctionClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FunctionClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunctionClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FunctionClientSession struct {
	Contract     *FunctionClient   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FunctionClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FunctionClientCallerSession struct {
	Contract *FunctionClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// FunctionClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FunctionClientTransactorSession struct {
	Contract     *FunctionClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// FunctionClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type FunctionClientRaw struct {
	Contract *FunctionClient // Generic contract binding to access the raw methods on
}

// FunctionClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FunctionClientCallerRaw struct {
	Contract *FunctionClientCaller // Generic read-only contract binding to access the raw methods on
}

// FunctionClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FunctionClientTransactorRaw struct {
	Contract *FunctionClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFunctionClient creates a new instance of FunctionClient, bound to a specific deployed contract.
func NewFunctionClient(address common.Address, backend bind.ContractBackend) (*FunctionClient, error) {
	contract, err := bindFunctionClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionClient{FunctionClientCaller: FunctionClientCaller{contract: contract}, FunctionClientTransactor: FunctionClientTransactor{contract: contract}, FunctionClientFilterer: FunctionClientFilterer{contract: contract}}, nil
}

// NewFunctionClientCaller creates a new read-only instance of FunctionClient, bound to a specific deployed contract.
func NewFunctionClientCaller(address common.Address, caller bind.ContractCaller) (*FunctionClientCaller, error) {
	contract, err := bindFunctionClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionClientCaller{contract: contract}, nil
}

// NewFunctionClientTransactor creates a new write-only instance of FunctionClient, bound to a specific deployed contract.
func NewFunctionClientTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionClientTransactor, error) {
	contract, err := bindFunctionClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionClientTransactor{contract: contract}, nil
}

// NewFunctionClientFilterer creates a new log filterer instance of FunctionClient, bound to a specific deployed contract.
func NewFunctionClientFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionClientFilterer, error) {
	contract, err := bindFunctionClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionClientFilterer{contract: contract}, nil
}

// bindFunctionClient binds a generic wrapper to an already deployed contract.
func bindFunctionClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FunctionClient *FunctionClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionClient.Contract.FunctionClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FunctionClient *FunctionClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionClient.Contract.FunctionClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FunctionClient *FunctionClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionClient.Contract.FunctionClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FunctionClient *FunctionClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FunctionClient *FunctionClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FunctionClient *FunctionClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionClient.Contract.contract.Transact(opts, method, params...)
}

// HandleOracleFulfillment is a paid mutator transaction binding the contract method 0x42ba6d28.
//
// Solidity: function handleOracleFulfillment(bytes32 requestId, address node, uint256 score, bytes response, bytes err) returns()
func (_FunctionClient *FunctionClientTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, node common.Address, score *big.Int, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionClient.contract.Transact(opts, "handleOracleFulfillment", requestId, node, score, response, err)
}

// HandleOracleFulfillment is a paid mutator transaction binding the contract method 0x42ba6d28.
//
// Solidity: function handleOracleFulfillment(bytes32 requestId, address node, uint256 score, bytes response, bytes err) returns()
func (_FunctionClient *FunctionClientSession) HandleOracleFulfillment(requestId [32]byte, node common.Address, score *big.Int, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionClient.Contract.HandleOracleFulfillment(&_FunctionClient.TransactOpts, requestId, node, score, response, err)
}

// HandleOracleFulfillment is a paid mutator transaction binding the contract method 0x42ba6d28.
//
// Solidity: function handleOracleFulfillment(bytes32 requestId, address node, uint256 score, bytes response, bytes err) returns()
func (_FunctionClient *FunctionClientTransactorSession) HandleOracleFulfillment(requestId [32]byte, node common.Address, score *big.Int, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionClient.Contract.HandleOracleFulfillment(&_FunctionClient.TransactOpts, requestId, node, score, response, err)
}

// FunctionClientRequestFulfilledIterator is returned from FilterRequestFulfilled and is used to iterate over the raw logs and unpacked data for RequestFulfilled events raised by the FunctionClient contract.
type FunctionClientRequestFulfilledIterator struct {
	Event *FunctionClientRequestFulfilled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FunctionClientRequestFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionClientRequestFulfilled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FunctionClientRequestFulfilled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FunctionClientRequestFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FunctionClientRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FunctionClientRequestFulfilled represents a RequestFulfilled event raised by the FunctionClient contract.
type FunctionClientRequestFulfilled struct {
	Id     [32]byte
	Node   common.Address
	Score  *big.Int
	Result []byte
	Err    []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRequestFulfilled is a free log retrieval operation binding the contract event 0xd0ceebb26e373797d34f277e1ddb02d72d53e53eb23ebeef0fcbb06191c76ea1.
//
// Solidity: event RequestFulfilled(bytes32 indexed id, address indexed node, uint256 score, bytes result, bytes err)
func (_FunctionClient *FunctionClientFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte, node []common.Address) (*FunctionClientRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _FunctionClient.contract.FilterLogs(opts, "RequestFulfilled", idRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return &FunctionClientRequestFulfilledIterator{contract: _FunctionClient.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

// WatchRequestFulfilled is a free log subscription operation binding the contract event 0xd0ceebb26e373797d34f277e1ddb02d72d53e53eb23ebeef0fcbb06191c76ea1.
//
// Solidity: event RequestFulfilled(bytes32 indexed id, address indexed node, uint256 score, bytes result, bytes err)
func (_FunctionClient *FunctionClientFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionClientRequestFulfilled, id [][32]byte, node []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _FunctionClient.contract.WatchLogs(opts, "RequestFulfilled", idRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FunctionClientRequestFulfilled)
				if err := _FunctionClient.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestFulfilled is a log parse operation binding the contract event 0xd0ceebb26e373797d34f277e1ddb02d72d53e53eb23ebeef0fcbb06191c76ea1.
//
// Solidity: event RequestFulfilled(bytes32 indexed id, address indexed node, uint256 score, bytes result, bytes err)
func (_FunctionClient *FunctionClientFilterer) ParseRequestFulfilled(log types.Log) (*FunctionClientRequestFulfilled, error) {
	event := new(FunctionClientRequestFulfilled)
	if err := _FunctionClient.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FunctionClientRequestSentIterator is returned from FilterRequestSent and is used to iterate over the raw logs and unpacked data for RequestSent events raised by the FunctionClient contract.
type FunctionClientRequestSentIterator struct {
	Event *FunctionClientRequestSent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FunctionClientRequestSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionClientRequestSent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FunctionClientRequestSent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FunctionClientRequestSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FunctionClientRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FunctionClientRequestSent represents a RequestSent event raised by the FunctionClient contract.
type FunctionClientRequestSent struct {
	Id         [32]byte
	FunctionId [32]byte
	Arg2       common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRequestSent is a free log retrieval operation binding the contract event 0xf66c42a7b945272fe61591062073bc8ec78eef6116c496d08b5cf0e54b7e885f.
//
// Solidity: event RequestSent(bytes32 indexed id, bytes32 indexed functionId, address func)
func (_FunctionClient *FunctionClientFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte, functionId [][32]byte) (*FunctionClientRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var functionIdRule []interface{}
	for _, functionIdItem := range functionId {
		functionIdRule = append(functionIdRule, functionIdItem)
	}

	logs, sub, err := _FunctionClient.contract.FilterLogs(opts, "RequestSent", idRule, functionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionClientRequestSentIterator{contract: _FunctionClient.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

// WatchRequestSent is a free log subscription operation binding the contract event 0xf66c42a7b945272fe61591062073bc8ec78eef6116c496d08b5cf0e54b7e885f.
//
// Solidity: event RequestSent(bytes32 indexed id, bytes32 indexed functionId, address func)
func (_FunctionClient *FunctionClientFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionClientRequestSent, id [][32]byte, functionId [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var functionIdRule []interface{}
	for _, functionIdItem := range functionId {
		functionIdRule = append(functionIdRule, functionIdItem)
	}

	logs, sub, err := _FunctionClient.contract.WatchLogs(opts, "RequestSent", idRule, functionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FunctionClientRequestSent)
				if err := _FunctionClient.contract.UnpackLog(event, "RequestSent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestSent is a log parse operation binding the contract event 0xf66c42a7b945272fe61591062073bc8ec78eef6116c496d08b5cf0e54b7e885f.
//
// Solidity: event RequestSent(bytes32 indexed id, bytes32 indexed functionId, address func)
func (_FunctionClient *FunctionClientFilterer) ParseRequestSent(log types.Log) (*FunctionClientRequestSent, error) {
	event := new(FunctionClientRequestSent)
	if err := _FunctionClient.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
