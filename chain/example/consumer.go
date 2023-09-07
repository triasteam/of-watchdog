// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

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

// FunctionConsumerMetaData contains all meta data concerning the FunctionConsumer contract.
var FunctionConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"contractRegistry\",\"name\":\"_reg\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyRequestData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderIsNotRegistry\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"FuncResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"score\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"functionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"funcWorker\",\"type\":\"address\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"secrets\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"}],\"name\":\"executeRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"score\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestError\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestResponse\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// FunctionConsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use FunctionConsumerMetaData.ABI instead.
var FunctionConsumerABI = FunctionConsumerMetaData.ABI

// FunctionConsumer is an auto generated Go binding around an Ethereum contract.
type FunctionConsumer struct {
	FunctionConsumerCaller     // Read-only binding to the contract
	FunctionConsumerTransactor // Write-only binding to the contract
	FunctionConsumerFilterer   // Log filterer for contract events
}

// FunctionConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type FunctionConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunctionConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FunctionConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunctionConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FunctionConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunctionConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FunctionConsumerSession struct {
	Contract     *FunctionConsumer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FunctionConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FunctionConsumerCallerSession struct {
	Contract *FunctionConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// FunctionConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FunctionConsumerTransactorSession struct {
	Contract     *FunctionConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// FunctionConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type FunctionConsumerRaw struct {
	Contract *FunctionConsumer // Generic contract binding to access the raw methods on
}

// FunctionConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FunctionConsumerCallerRaw struct {
	Contract *FunctionConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// FunctionConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FunctionConsumerTransactorRaw struct {
	Contract *FunctionConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFunctionConsumer creates a new instance of FunctionConsumer, bound to a specific deployed contract.
func NewFunctionConsumer(address common.Address, backend bind.ContractBackend) (*FunctionConsumer, error) {
	contract, err := bindFunctionConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionConsumer{FunctionConsumerCaller: FunctionConsumerCaller{contract: contract}, FunctionConsumerTransactor: FunctionConsumerTransactor{contract: contract}, FunctionConsumerFilterer: FunctionConsumerFilterer{contract: contract}}, nil
}

// NewFunctionConsumerCaller creates a new read-only instance of FunctionConsumer, bound to a specific deployed contract.
func NewFunctionConsumerCaller(address common.Address, caller bind.ContractCaller) (*FunctionConsumerCaller, error) {
	contract, err := bindFunctionConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionConsumerCaller{contract: contract}, nil
}

// NewFunctionConsumerTransactor creates a new write-only instance of FunctionConsumer, bound to a specific deployed contract.
func NewFunctionConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionConsumerTransactor, error) {
	contract, err := bindFunctionConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionConsumerTransactor{contract: contract}, nil
}

// NewFunctionConsumerFilterer creates a new log filterer instance of FunctionConsumer, bound to a specific deployed contract.
func NewFunctionConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionConsumerFilterer, error) {
	contract, err := bindFunctionConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionConsumerFilterer{contract: contract}, nil
}

// bindFunctionConsumer binds a generic wrapper to an already deployed contract.
func bindFunctionConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FunctionConsumer *FunctionConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionConsumer.Contract.FunctionConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FunctionConsumer *FunctionConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.FunctionConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FunctionConsumer *FunctionConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.FunctionConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FunctionConsumer *FunctionConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FunctionConsumer *FunctionConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FunctionConsumer *FunctionConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.contract.Transact(opts, method, params...)
}

// LatestError is a free data retrieval call binding the contract method 0xfffeb84e.
//
// Solidity: function latestError() view returns(bytes)
func (_FunctionConsumer *FunctionConsumerCaller) LatestError(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FunctionConsumer.contract.Call(opts, &out, "latestError")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// LatestError is a free data retrieval call binding the contract method 0xfffeb84e.
//
// Solidity: function latestError() view returns(bytes)
func (_FunctionConsumer *FunctionConsumerSession) LatestError() ([]byte, error) {
	return _FunctionConsumer.Contract.LatestError(&_FunctionConsumer.CallOpts)
}

// LatestError is a free data retrieval call binding the contract method 0xfffeb84e.
//
// Solidity: function latestError() view returns(bytes)
func (_FunctionConsumer *FunctionConsumerCallerSession) LatestError() ([]byte, error) {
	return _FunctionConsumer.Contract.LatestError(&_FunctionConsumer.CallOpts)
}

// LatestRequestId is a free data retrieval call binding the contract method 0x1aa46f59.
//
// Solidity: function latestRequestId() view returns(bytes32)
func (_FunctionConsumer *FunctionConsumerCaller) LatestRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionConsumer.contract.Call(opts, &out, "latestRequestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LatestRequestId is a free data retrieval call binding the contract method 0x1aa46f59.
//
// Solidity: function latestRequestId() view returns(bytes32)
func (_FunctionConsumer *FunctionConsumerSession) LatestRequestId() ([32]byte, error) {
	return _FunctionConsumer.Contract.LatestRequestId(&_FunctionConsumer.CallOpts)
}

// LatestRequestId is a free data retrieval call binding the contract method 0x1aa46f59.
//
// Solidity: function latestRequestId() view returns(bytes32)
func (_FunctionConsumer *FunctionConsumerCallerSession) LatestRequestId() ([32]byte, error) {
	return _FunctionConsumer.Contract.LatestRequestId(&_FunctionConsumer.CallOpts)
}

// LatestResponse is a free data retrieval call binding the contract method 0xbef3a2f0.
//
// Solidity: function latestResponse() view returns(bytes)
func (_FunctionConsumer *FunctionConsumerCaller) LatestResponse(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FunctionConsumer.contract.Call(opts, &out, "latestResponse")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// LatestResponse is a free data retrieval call binding the contract method 0xbef3a2f0.
//
// Solidity: function latestResponse() view returns(bytes)
func (_FunctionConsumer *FunctionConsumerSession) LatestResponse() ([]byte, error) {
	return _FunctionConsumer.Contract.LatestResponse(&_FunctionConsumer.CallOpts)
}

// LatestResponse is a free data retrieval call binding the contract method 0xbef3a2f0.
//
// Solidity: function latestResponse() view returns(bytes)
func (_FunctionConsumer *FunctionConsumerCallerSession) LatestResponse() ([]byte, error) {
	return _FunctionConsumer.Contract.LatestResponse(&_FunctionConsumer.CallOpts)
}

// ExecuteRequest is a paid mutator transaction binding the contract method 0x2912557a.
//
// Solidity: function executeRequest(bytes32 name, string source, bytes secrets, string[] args) returns(bytes32)
func (_FunctionConsumer *FunctionConsumerTransactor) ExecuteRequest(opts *bind.TransactOpts, name [32]byte, source string, secrets []byte, args []string) (*types.Transaction, error) {
	return _FunctionConsumer.contract.Transact(opts, "executeRequest", name, source, secrets, args)
}

// ExecuteRequest is a paid mutator transaction binding the contract method 0x2912557a.
//
// Solidity: function executeRequest(bytes32 name, string source, bytes secrets, string[] args) returns(bytes32)
func (_FunctionConsumer *FunctionConsumerSession) ExecuteRequest(name [32]byte, source string, secrets []byte, args []string) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.ExecuteRequest(&_FunctionConsumer.TransactOpts, name, source, secrets, args)
}

// ExecuteRequest is a paid mutator transaction binding the contract method 0x2912557a.
//
// Solidity: function executeRequest(bytes32 name, string source, bytes secrets, string[] args) returns(bytes32)
func (_FunctionConsumer *FunctionConsumerTransactorSession) ExecuteRequest(name [32]byte, source string, secrets []byte, args []string) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.ExecuteRequest(&_FunctionConsumer.TransactOpts, name, source, secrets, args)
}

// HandleOracleFulfillment is a paid mutator transaction binding the contract method 0x42ba6d28.
//
// Solidity: function handleOracleFulfillment(bytes32 requestId, address node, uint256 score, bytes response, bytes err) returns()
func (_FunctionConsumer *FunctionConsumerTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, node common.Address, score *big.Int, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionConsumer.contract.Transact(opts, "handleOracleFulfillment", requestId, node, score, response, err)
}

// HandleOracleFulfillment is a paid mutator transaction binding the contract method 0x42ba6d28.
//
// Solidity: function handleOracleFulfillment(bytes32 requestId, address node, uint256 score, bytes response, bytes err) returns()
func (_FunctionConsumer *FunctionConsumerSession) HandleOracleFulfillment(requestId [32]byte, node common.Address, score *big.Int, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.HandleOracleFulfillment(&_FunctionConsumer.TransactOpts, requestId, node, score, response, err)
}

// HandleOracleFulfillment is a paid mutator transaction binding the contract method 0x42ba6d28.
//
// Solidity: function handleOracleFulfillment(bytes32 requestId, address node, uint256 score, bytes response, bytes err) returns()
func (_FunctionConsumer *FunctionConsumerTransactorSession) HandleOracleFulfillment(requestId [32]byte, node common.Address, score *big.Int, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.HandleOracleFulfillment(&_FunctionConsumer.TransactOpts, requestId, node, score, response, err)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_FunctionConsumer *FunctionConsumerTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _FunctionConsumer.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_FunctionConsumer *FunctionConsumerSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.SetOracle(&_FunctionConsumer.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_FunctionConsumer *FunctionConsumerTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _FunctionConsumer.Contract.SetOracle(&_FunctionConsumer.TransactOpts, oracle)
}

// FunctionConsumerFuncResponseIterator is returned from FilterFuncResponse and is used to iterate over the raw logs and unpacked data for FuncResponse events raised by the FunctionConsumer contract.
type FunctionConsumerFuncResponseIterator struct {
	Event *FunctionConsumerFuncResponse // Event containing the contract specifics and raw log

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
func (it *FunctionConsumerFuncResponseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionConsumerFuncResponse)
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
		it.Event = new(FunctionConsumerFuncResponse)
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
func (it *FunctionConsumerFuncResponseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FunctionConsumerFuncResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FunctionConsumerFuncResponse represents a FuncResponse event raised by the FunctionConsumer contract.
type FunctionConsumerFuncResponse struct {
	RequestId [32]byte
	Result    []byte
	Err       []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFuncResponse is a free log retrieval operation binding the contract event 0x51732bd592d531e70a33ba1a99ea425d26fa7cb95aaf81189a1b2a0e7649edc9.
//
// Solidity: event FuncResponse(bytes32 indexed requestId, bytes result, bytes err)
func (_FunctionConsumer *FunctionConsumerFilterer) FilterFuncResponse(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionConsumerFuncResponseIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionConsumer.contract.FilterLogs(opts, "FuncResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionConsumerFuncResponseIterator{contract: _FunctionConsumer.contract, event: "FuncResponse", logs: logs, sub: sub}, nil
}

// WatchFuncResponse is a free log subscription operation binding the contract event 0x51732bd592d531e70a33ba1a99ea425d26fa7cb95aaf81189a1b2a0e7649edc9.
//
// Solidity: event FuncResponse(bytes32 indexed requestId, bytes result, bytes err)
func (_FunctionConsumer *FunctionConsumerFilterer) WatchFuncResponse(opts *bind.WatchOpts, sink chan<- *FunctionConsumerFuncResponse, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionConsumer.contract.WatchLogs(opts, "FuncResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FunctionConsumerFuncResponse)
				if err := _FunctionConsumer.contract.UnpackLog(event, "FuncResponse", log); err != nil {
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

// ParseFuncResponse is a log parse operation binding the contract event 0x51732bd592d531e70a33ba1a99ea425d26fa7cb95aaf81189a1b2a0e7649edc9.
//
// Solidity: event FuncResponse(bytes32 indexed requestId, bytes result, bytes err)
func (_FunctionConsumer *FunctionConsumerFilterer) ParseFuncResponse(log types.Log) (*FunctionConsumerFuncResponse, error) {
	event := new(FunctionConsumerFuncResponse)
	if err := _FunctionConsumer.contract.UnpackLog(event, "FuncResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FunctionConsumerRequestFulfilledIterator is returned from FilterRequestFulfilled and is used to iterate over the raw logs and unpacked data for RequestFulfilled events raised by the FunctionConsumer contract.
type FunctionConsumerRequestFulfilledIterator struct {
	Event *FunctionConsumerRequestFulfilled // Event containing the contract specifics and raw log

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
func (it *FunctionConsumerRequestFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionConsumerRequestFulfilled)
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
		it.Event = new(FunctionConsumerRequestFulfilled)
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
func (it *FunctionConsumerRequestFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FunctionConsumerRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FunctionConsumerRequestFulfilled represents a RequestFulfilled event raised by the FunctionConsumer contract.
type FunctionConsumerRequestFulfilled struct {
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
func (_FunctionConsumer *FunctionConsumerFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte, node []common.Address) (*FunctionConsumerRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _FunctionConsumer.contract.FilterLogs(opts, "RequestFulfilled", idRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return &FunctionConsumerRequestFulfilledIterator{contract: _FunctionConsumer.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

// WatchRequestFulfilled is a free log subscription operation binding the contract event 0xd0ceebb26e373797d34f277e1ddb02d72d53e53eb23ebeef0fcbb06191c76ea1.
//
// Solidity: event RequestFulfilled(bytes32 indexed id, address indexed node, uint256 score, bytes result, bytes err)
func (_FunctionConsumer *FunctionConsumerFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionConsumerRequestFulfilled, id [][32]byte, node []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _FunctionConsumer.contract.WatchLogs(opts, "RequestFulfilled", idRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FunctionConsumerRequestFulfilled)
				if err := _FunctionConsumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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
func (_FunctionConsumer *FunctionConsumerFilterer) ParseRequestFulfilled(log types.Log) (*FunctionConsumerRequestFulfilled, error) {
	event := new(FunctionConsumerRequestFulfilled)
	if err := _FunctionConsumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FunctionConsumerRequestSentIterator is returned from FilterRequestSent and is used to iterate over the raw logs and unpacked data for RequestSent events raised by the FunctionConsumer contract.
type FunctionConsumerRequestSentIterator struct {
	Event *FunctionConsumerRequestSent // Event containing the contract specifics and raw log

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
func (it *FunctionConsumerRequestSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionConsumerRequestSent)
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
		it.Event = new(FunctionConsumerRequestSent)
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
func (it *FunctionConsumerRequestSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FunctionConsumerRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FunctionConsumerRequestSent represents a RequestSent event raised by the FunctionConsumer contract.
type FunctionConsumerRequestSent struct {
	Id         [32]byte
	FunctionId [32]byte
	FuncWorker common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRequestSent is a free log retrieval operation binding the contract event 0xf66c42a7b945272fe61591062073bc8ec78eef6116c496d08b5cf0e54b7e885f.
//
// Solidity: event RequestSent(bytes32 indexed id, bytes32 indexed functionId, address funcWorker)
func (_FunctionConsumer *FunctionConsumerFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte, functionId [][32]byte) (*FunctionConsumerRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var functionIdRule []interface{}
	for _, functionIdItem := range functionId {
		functionIdRule = append(functionIdRule, functionIdItem)
	}

	logs, sub, err := _FunctionConsumer.contract.FilterLogs(opts, "RequestSent", idRule, functionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionConsumerRequestSentIterator{contract: _FunctionConsumer.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

// WatchRequestSent is a free log subscription operation binding the contract event 0xf66c42a7b945272fe61591062073bc8ec78eef6116c496d08b5cf0e54b7e885f.
//
// Solidity: event RequestSent(bytes32 indexed id, bytes32 indexed functionId, address funcWorker)
func (_FunctionConsumer *FunctionConsumerFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionConsumerRequestSent, id [][32]byte, functionId [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var functionIdRule []interface{}
	for _, functionIdItem := range functionId {
		functionIdRule = append(functionIdRule, functionIdItem)
	}

	logs, sub, err := _FunctionConsumer.contract.WatchLogs(opts, "RequestSent", idRule, functionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FunctionConsumerRequestSent)
				if err := _FunctionConsumer.contract.UnpackLog(event, "RequestSent", log); err != nil {
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
// Solidity: event RequestSent(bytes32 indexed id, bytes32 indexed functionId, address funcWorker)
func (_FunctionConsumer *FunctionConsumerFilterer) ParseRequestSent(log types.Log) (*FunctionConsumerRequestSent, error) {
	event := new(FunctionConsumerRequestSent)
	if err := _FunctionConsumer.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
