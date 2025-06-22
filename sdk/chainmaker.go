/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package sdk

// sdk for user

import (
	vmPb "chainmaker.org/chainmaker/pb-go/v2/vm"
	"fmt"
	"strconv"
	"unsafe"
)

type ResultCode int

const (
	// special parameters passed to contract
	ContractParamCreatorOrgId = "__creator_org_id__"
	ContractParamCreatorRole  = "__creator_role__"
	ContractParamCreatorPk    = "__creator_pk__"
	ContractParamSenderOrgId  = "__sender_org_id__"
	ContractParamSenderRole   = "__sender_role__"
	ContractParamSenderPk     = "__sender_pk__"
	ContractParamBlockHeight  = "__block_height__"
	ContractParamTxId         = "__tx_id__"
	ContractParamContextPtr   = "__context_ptr__"
	ContractParamTxTimeStamp  = "__tx_time_stamp__"

	// method name used by smart contract sdk
	// common
	ContractMethodLogMessage      = "LogMessage"
	ContractMethodSuccessResult   = "SuccessResult"
	ContractMethodErrorResult     = "ErrorResult"
	ContractMethodCallContract    = "CallContract"
	ContractMethodCallContractLen = "CallContractLen"
	ContractMethodEmitEvent       = "EmitEvent"
	ContractMethodArg             = "GetArg"
	ContractMethodArgLen          = "GetArgLen"

	// paillier
	ContractMethodGetPaillierOperationResult    = "GetPaillierOperationResult"
	ContractMethodGetPaillierOperationResultLen = "GetPaillierOperationResultLen"
	// bulletproofs
	ContractMethodGetBulletproofsResult    = "GetBulletproofsResult"
	ContractMethodGetBulletproofsResultLen = "GetBulletproofsResultLen"

	// kv
	ContractMethodGetStateLen      = "GetStateLen"
	ContractMethodGetState         = "GetState"
	ContractMethodPutState         = "PutState"
	ContractMethodDeleteState      = "DeleteState"
	ContractMethodGetBatchStateLen = "GetBatchStateLen"
	ContractMethodGetBatchState    = "GetBatchState"

	//address
	ContractMethodSenderAddress    = "GetSenderAddress"
	ContractMethodSenderAddressLen = "GetSenderAddressLen"
	// kv iterator
	ContractMethodKvIterator        = "KvIterator"
	ContractMethodKvPreIterator     = "KvPreIterator"
	ContractMethodKvIteratorHasNext = "KvIteratorHasNext"
	ContractMethodKvIteratorNextLen = "KvIteratorNextLen"
	ContractMethodKvIteratorNext    = "KvIteratorNext"
	ContractMethodKvIteratorClose   = "KvIteratorClose"

	// keyHistoryKvIterator method
	ContractHistoryKvIterator        = "HistoryKvIterator"
	ContractHistoryKvIteratorHasNext = "HistoryKvIterHasNext"
	ContractHistoryKvIteratorNextLen = "HistoryKvIterNextLen"
	ContractHistoryKvIteratorNext    = "HistoryKvIterNext"
	ContractHistoryKvIteratorClose   = "HistoryKvIterClose"
	// sql
	ContractMethodExecuteQuery       = "ExecuteQuery"
	ContractMethodExecuteQueryOne    = "ExecuteQueryOne"
	ContractMethodExecuteQueryOneLen = "ExecuteQueryOneLen"
	ContractMethodRSNext             = "RSNext"
	ContractMethodRSNextLen          = "RSNextLen"
	ContractMethodRSHasNext          = "RSHasNext"
	ContractMethodRSClose            = "RSClose"
	ContractMethodExecuteUpdate      = "ExecuteUpdate"
	ContractMethodExecuteDdl         = "ExecuteDDL"

	SUCCESS    ResultCode = 0
	ERROR      ResultCode = 1
	DebugLevel            = -1
	InfoLevel             = 0
	WarnLevel             = 1
	ErrorLevel            = 2
	// default batch keys count limit
	defaultLimitKeys = 10000
)

// sysCall provides data interaction with the chain. sysCallReq common param, request var param
//
//go:wasmimport env sys_call
func sysCall(requestHeader string, requestBody string) int32

//go:wasmimport env log_message
func logMessage(msg string)

//go:wasmimport env log_message_with_type
func logMessageWithType(msg string, msgType int32)

// SimContextCommon common context
type SimContextCommon interface {
	// Arg get arg from transaction parameters, as:  arg1, code := ctx.Arg("arg1")
	Arg(key string) ([]byte, ResultCode)
	// Arg get arg from transaction parameters, as:  arg1, code := ctx.ArgString("arg1")
	ArgString(key string) (string, ResultCode)
	// Args return args
	Args() []*EasyCodecItem

	// SuccessResult record the execution result of the transaction, multiple calls will override
	SuccessResult(msg string)
	// SuccessResultByte record the execution result of the transaction, multiple calls will override
	SuccessResultByte(msg []byte)
	// ErrorResult record the execution result of the transaction. multiple calls will append. Once there is an error, it cannot be called success method
	ErrorResult(msg string)
	// CallContract cross contract call
	CallContract(contractName string, method string, param map[string][]byte) ([]byte, ResultCode)
	// GetCreatorOrgId get tx creator org id
	GetCreatorOrgId() (string, ResultCode)

	// Log record log to chain server
	Log(msg string)
	// Debugf record log to chain server
	// @param format: 日志格式化模板
	// @param a: 模板参数
	Debugf(format string, a ...interface{})
	// Infof record log to chain server
	// @param format: 日志格式化模板
	// @param a: 模板参数
	Infof(format string, a ...interface{})
	// Warnf record log to chain server
	// @param format: 日志格式化模板
	// @param a: 模板参数
	Warnf(format string, a ...interface{})
	// Errorf record log to chain server
	// @param format: 日志格式化模板
	// @param a: 模板参数
	Errorf(format string, a ...interface{})

	// GetCreatorRole get tx creator role
	GetCreatorRole() (string, ResultCode)
	// GetCreatorPk get tx creator pk
	GetCreatorPk() (string, ResultCode)
	// GetSenderOrgId get tx sender org id
	GetSenderOrgId() (string, ResultCode)
	// GetSenderOrgId get tx sender role
	GetSenderRole() (string, ResultCode)
	// GetSenderOrgId get tx sender pk
	GetSenderPk() (string, ResultCode)
	// GetBlockHeight get tx block height
	GetBlockHeight() (string, ResultCode)
	// GetTxTimeStamp get tx timestamp
	// @return1: 交易timestamp
	// @return2: 获取错误信息
	GetTxTimeStamp() (string, ResultCode)
	// GetTxId get current tx id
	GetTxId() (string, ResultCode)
	// GetTxInfo get tx info
	// @param txId :合约交易ID
	GetTxInfo(txId string) ([]byte, ResultCode)
	// EmitEvent emit event, you can subscribe to the event using the SDK
	EmitEvent(topic string, data ...string) ResultCode

	// GetSenderAddr Get the address of the origin caller address, same with Origin()
	// @return1: origin caller address
	// @return2: 获取错误信息
	GetSenderAddr() (string, ResultCode)
	// Sender Get the address of the sender address, if the contract is called by another contract, the result will be
	// the caller contract's address.
	// Sender will return system contract address when executing the init or upgrade method (If you need to return the
	// user address, we recommend using Origin method here), because the init and upgrade methods are cross-contract
	// txs (system contract -> common contract).
	// @return1: sender address
	// @return2: 获取错误信息
	Sender() (string, ResultCode)

	// Origin Get the address of the tx origin caller address
	// @return1: origin caller address
	// @return2: 获取错误信息
	Origin() (string, ResultCode)
}

// SimContext kv context
type SimContext interface {
	SimContextCommon
	// GetState get [key+"#"+field] from chain and db
	GetState(key string, field string) (string, ResultCode)
	// GetStateByte get [key+"#"+field] from chain and db
	GetStateByte(key string, field string) ([]byte, ResultCode)

	// GetStateWithExists get [key, field] from chain
	// @param key: 获取的参数名
	// @param field: 获取的参数名
	// @return1: 获取结果，格式为string
	// @return2: 是否存在，bool, 字符串长度为0不代表不存在
	// @return3: ResultCode
	GetStateWithExists(key, field string) (string, bool, ResultCode)
	// GetBatchState get [BatchKeys] from chain
	// @param batchKey: 获取的参数名
	// @return1: 获取结果
	// @return2: 获取错误信息
	GetBatchState(batchKeys []*vmPb.BatchKey) ([]*vmPb.BatchKey, ResultCode)
	// GetStateFromKeyByte get [key] from chain
	// @param key: 获取的参数名
	// @return1: 获取结果，格式为[]byte, nil表示不存在
	GetStateFromKeyByte(key string) ([]byte, ResultCode)

	// GetStateByte get [key] from chain and db
	GetStateFromKey(key string) ([]byte, ResultCode)

	// GetStateFromKeyWithExists get [key] from chain
	// @param key: 获取的参数名
	// @return1: 获取结果，格式为string
	// @return2: 是否存在，bool, 字符串长度为0不代表不存在
	// @return3: ResultCode
	GetStateFromKeyWithExists(key string) (string, bool, ResultCode)

	// PutState put [key+"#"+field, value] to chain
	PutState(key string, field string, value string) ResultCode
	// PutStateByte put [key+"#"+field, value] to chain
	PutStateByte(key string, field string, value []byte) ResultCode
	// PutStateFromKey put [key, value] to chain
	PutStateFromKey(key string, value string) ResultCode
	// PutStateFromKeyByte put [key, value] to chain
	PutStateFromKeyByte(key string, value []byte) ResultCode
	// DeleteState delete [key+"#"+field] to chain
	DeleteState(key string, field string) ResultCode
	// DeleteStateFromKey delete [key] to chain
	DeleteStateFromKey(key string) ResultCode
	// NewIterator range of [startKey, limitKey), front closed back open
	NewIterator(startKey string, limitKey string) (ResultSetKV, ResultCode)
	// NewIteratorWithField range of [key+"#"+startField, key+"#"+limitField), front closed back open
	NewIteratorWithField(key string, startField string, limitField string) (ResultSetKV, ResultCode)
	// NewIteratorPrefixWithKeyField range of [key+"#"+field, key+"#"+field], front closed back closed
	NewIteratorPrefixWithKeyField(key string, field string) (ResultSetKV, ResultCode)
	// NewIteratorPrefixWithKey range of [key, key], front closed back closed
	NewIteratorPrefixWithKey(key string) (ResultSetKV, ResultCode)
	// NewHistoryKvIterForKey query all historical data of key, field
	// @param1: 查询历史的key
	// @param2: 查询历史的field
	// @return1: 根据key, field 生成的历史迭代器
	// @return2: 获取错误信息
	NewHistoryKvIterForKey(startKey string, startField string) (KeyHistoryKvIter, ResultCode)
}

type SimContextCommonImpl struct {
	origin string
}

type SimContextImpl struct {
	SimContextCommonImpl
}

func NewSimContext() SimContext {
	return &SimContextImpl{}
}
func (s *SimContextImpl) GetState(key string, field string) (string, ResultCode) {
	return GetState(key, field)
}
func (s *SimContextImpl) GetStateByte(key string, field string) ([]byte, ResultCode) {
	return GetStateByte(key, field)
}
func (s *SimContextImpl) GetStateWithExists(key, field string) (string, bool, ResultCode) {
	value, err := GetStateByte(key, field)
	if err != SUCCESS || value == nil {
		return "", false, ERROR
	}
	return string(value), true, ERROR
}
func (s *SimContextImpl) GetStateFromKey(key string) ([]byte, ResultCode) {
	return GetStateByte(key, "")
}

func (s *SimContextImpl) GetStateFromKeyWithExists(key string) (string, bool, ResultCode) {
	return s.GetStateWithExists(key, "")
}

func (s *SimContextImpl) GetStateFromKeyByte(key string) ([]byte, ResultCode) {
	return s.GetStateByte(key, "")
}

func (s *SimContextImpl) PutState(key string, field string, value string) ResultCode {
	return PutState(key, field, value)
}
func (s *SimContextImpl) PutStateByte(key string, field string, value []byte) ResultCode {
	return PutState(key, field, string(value))
}
func (s *SimContextImpl) PutStateFromKey(key string, value string) ResultCode {
	return PutState(key, "", value)
}
func (s *SimContextImpl) PutStateFromKeyByte(key string, value []byte) ResultCode {
	return PutStateByte(key, "", value)
}
func (s *SimContextImpl) DeleteState(key string, field string) ResultCode {
	return DeleteState(key, field)
}
func (s *SimContextImpl) DeleteStateFromKey(key string) ResultCode {
	return DeleteState(key, "")
}
func (s *SimContextImpl) GetBatchState(batchKeys []*vmPb.BatchKey) ([]*vmPb.BatchKey, ResultCode) {
	if err := s.batchKeysLimit(batchKeys); err != nil {
		return nil, ERROR
	}
	getBatchStateKeys := vmPb.BatchKeys{Keys: batchKeys}
	getBatchStateKeysByte, err := getBatchStateKeys.Marshal()
	if err != nil {
		return nil, ERROR
	}
	ec := NewEasyCodec()
	ec.AddBytes("BatchKeys", getBatchStateKeysByte)
	value, code := GetBytesFromChain(ec, ContractMethodSenderAddressLen, ContractMethodSenderAddress)
	if code != SUCCESS {
		return nil, code
	}
	keys := &vmPb.BatchKeys{}
	if err = keys.Unmarshal(value); err != nil {
		return nil, ERROR
	}
	return keys.Keys, SUCCESS
}
func (s *SimContextImpl) batchKeysLimit(keys []*vmPb.BatchKey) error {
	if len(keys) > defaultLimitKeys {
		return fmt.Errorf("over batch keys count limit %d", defaultLimitKeys)
	}
	return nil
}

// common
func (s *SimContextCommonImpl) Arg(key string) ([]byte, ResultCode) {
	return Arg(key)
}
func (s *SimContextCommonImpl) ArgString(key string) (string, ResultCode) {
	val, code := Arg(key)
	return string(val), code
}
func (s *SimContextCommonImpl) Args() []*EasyCodecItem {
	return Args()
}
func (s *SimContextCommonImpl) Log(msg string) {
	LogMessage(msg)
}
func (s *SimContextCommonImpl) CallContract(contractName string, method string, param map[string][]byte) ([]byte, ResultCode) {
	return CallContract(contractName, method, param)
}
func (s *SimContextCommonImpl) SuccessResult(msg string) {
	sysCall(getRequestHeader(ContractMethodSuccessResult), msg)
}
func (s *SimContextCommonImpl) SuccessResultByte(msg []byte) {
	sysCall(getRequestHeader(ContractMethodSuccessResult), string(msg))
}
func (s *SimContextCommonImpl) ErrorResult(msg string) {
	sysCall(getRequestHeader(ContractMethodErrorResult), string(msg))
}
func (s *SimContextCommonImpl) GetCreatorOrgId() (string, ResultCode) {
	return stringArg(ContractParamCreatorOrgId)
}
func (s *SimContextCommonImpl) GetCreatorRole() (string, ResultCode) {
	return stringArg(ContractParamCreatorRole)
}
func (s *SimContextCommonImpl) GetCreatorPk() (string, ResultCode) {
	return stringArg(ContractParamCreatorPk)
}
func (s *SimContextCommonImpl) GetSenderOrgId() (string, ResultCode) {
	return stringArg(ContractParamSenderOrgId)
}
func (s *SimContextCommonImpl) GetSenderRole() (string, ResultCode) {
	return stringArg(ContractParamSenderRole)
}
func (s *SimContextCommonImpl) GetSenderPk() (string, ResultCode) {
	return stringArg(ContractParamSenderPk)
}
func (s *SimContextCommonImpl) GetBlockHeight() (string, ResultCode) {
	return stringArg(ContractParamBlockHeight)
}
func (s *SimContextCommonImpl) GetTxId() (string, ResultCode) {
	return stringArg(ContractParamTxId)
}
func (s *SimContextCommonImpl) GetTxTimeStamp() (string, ResultCode) {
	return stringArg(ContractParamTxTimeStamp)
}
func (s *SimContextCommonImpl) GetTxInfo(txId string) ([]byte, ResultCode) {
	paramTxId := "txId"
	paramMethod := "method"

	contractName := "CHAIN_QUERY"
	method := "GET_TX_BY_TX_ID"
	args := map[string][]byte{
		paramTxId:   []byte(txId),
		paramMethod: []byte(method),
	}
	return s.CallContract(contractName, method, args)
}

func (s *SimContextCommonImpl) EmitEvent(topic string, data ...string) ResultCode {
	return EmitEvent(topic, data...)
}
func (s *SimContextCommonImpl) Debugf(format string, a ...interface{}) {
	LogMessageWityType(getMessage(format, a), DebugLevel)
}

func (s *SimContextCommonImpl) Infof(format string, a ...interface{}) {
	LogMessageWityType(getMessage(format, a), InfoLevel)
}

func (s *SimContextCommonImpl) Warnf(format string, a ...interface{}) {
	LogMessageWityType(getMessage(format, a), WarnLevel)
}

func (s *SimContextCommonImpl) Errorf(format string, a ...interface{}) {
	LogMessageWityType(getMessage(format, a), ErrorLevel)
}

func (s *SimContextCommonImpl) Origin() (string, ResultCode) {
	if s.origin != "" {
		return s.origin, SUCCESS
	}
	result, code := GetSenderAddress()
	if code == SUCCESS {
		s.origin = result
	}
	return result, code
}

func (s *SimContextCommonImpl) GetSenderAddr() (string, ResultCode) {
	return s.Origin()
}

func (s *SimContextCommonImpl) Sender() (string, ResultCode) {
	//TODO Detail Implement
	return s.Origin()
}

// TODO 原本sdk对于参数的实现方法是在实例内用全局变量，由于vm-wasmer用的vm-pool同一个合约用的一个实例，在并发情况下可能有问题
// TODO 现在修改成从链上获取，待检查,暂时注释掉allocate、deallocate两个跟参数相关的接口
//var argsBytes []byte
//var argsMap []*EasyCodecItem
//var argsFlag bool

//go:wasmexport runtime_type
func runtimeType() int32 {
	var ContractRuntimeGoSdkType int32 = 4
	//argsFlag = false
	return ContractRuntimeGoSdkType
}

////go:wasmexport deallocate
//func deallocate(size int32) {
//	argsBytes = make([]byte, size)
//	argsMap = make([]*EasyCodecItem, 0)
//	argsFlag = false
//}
//
////go:wasmexport allocate
//func allocate(size int32) uintptr {
//	argsBytes = make([]byte, size)
//	argsMap = make([]*EasyCodecItem, 0)
//	argsFlag = false
//
//	return uintptr(unsafe.Pointer(&argsBytes[0]))
//}

func getRequestHeader(method string) string {
	ec := NewEasyCodec()
	ec.AddValue(EasyKeyType_SYSTEM, "ctx_ptr", EasyValueType_INT32, getCtxPtr())
	ec.AddValue(EasyKeyType_SYSTEM, "version", EasyValueType_STRING, "v1.2.0")
	ec.AddValue(EasyKeyType_SYSTEM, "method", EasyValueType_STRING, method)
	return string(ec.Marshal())
}

// LogMessage
func LogMessage(msg string) {
	logMessage(msg)
}

func LogMessageWityType(msg string, msgType int32) {
	logMessageWithType(msg, msgType)
}

// getMessage
func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

// GetSenderAddress get senderAddress from chain
func GetSenderAddress() (string, ResultCode) {
	ec := NewEasyCodec()
	result, code := GetBytesFromChain(ec, ContractMethodSenderAddressLen, ContractMethodSenderAddress)
	if code != SUCCESS {
		return "", code
	}
	return string(result), code
}

// GetState get state from chain
func GetState(key string, field string) (string, ResultCode) {
	result, code := GetStateByte(key, field)
	if code != SUCCESS {
		return "", code
	}
	return string(result), code
}

// GetState get state from chain
func GetStateByte(key string, field string) ([]byte, ResultCode) {
	ec := NewEasyCodec()
	ec.AddString("key", key)
	ec.AddString("field", field)
	return GetBytesFromChain(ec, ContractMethodGetStateLen, ContractMethodGetState)
}

func GetBytesFromChain(ec *EasyCodec, methodLen string, method string) ([]byte, ResultCode) {
	// # get len
	// ## prepare param
	var valueLen int32 = 0
	valuePtr := int32(uintptr(unsafe.Pointer(&valueLen)))
	ec.AddInt32("value_ptr", valuePtr)
	b := ec.Marshal()
	// ## send req get len
	code := sysCall(getRequestHeader(methodLen), string(b))
	// ## verify
	if code != int32(SUCCESS) {
		return nil, ERROR
	}
	if valueLen == 0 {
		return nil, SUCCESS
	}
	// # get data
	// ## prepare param
	valueByte := make([]byte, valueLen)
	ec.RemoveKey("value_ptr")
	valuePtr = int32(uintptr(unsafe.Pointer(&valueByte[0])))
	ec.AddInt32("value_ptr", valuePtr)
	b = ec.Marshal()
	// ## send req get value
	code2 := sysCall(getRequestHeader(method), string(b))
	if code2 != int32(SUCCESS) {
		return nil, ERROR
	}
	return valueByte, SUCCESS
}

// GetInt32FromChain get i32 from chain
func GetInt32FromChain(ec *EasyCodec, method string) (int32, ResultCode) {
	// # get len
	// ## prepare param
	var valueLen int32 = 0
	valuePtr := int32(uintptr(unsafe.Pointer(&valueLen)))
	ec.AddInt32("value_ptr", valuePtr)
	b := ec.Marshal()
	// ## send req get len
	code := sysCall(getRequestHeader(method), string(b))
	return valueLen, ResultCode(code)
}

// GetStateFromKey get state from chain
func GetStateFromKey(key string) ([]byte, ResultCode) {
	return GetStateByte(key, "")
}

// EmitEvent emit Event to chain
func EmitEvent(topic string, data ...string) ResultCode {
	// prepare param
	var items []*EasyCodecItem
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "topic",
		ValueType: EasyValueType_STRING,
		Value:     topic,
	})
	for index, value := range data {
		items = append(items, &EasyCodecItem{
			KeyType:   EasyKeyType_USER,
			Key:       "data" + strconv.FormatInt(int64(index), 10),
			ValueType: EasyValueType_STRING,
			Value:     value,
		})
	}
	b := EasyMarshal(items)
	reqBody := string(b)
	// send req put value
	code := sysCall(getRequestHeader(ContractMethodEmitEvent), reqBody)
	if code != int32(SUCCESS) {
		return ERROR
	}
	return SUCCESS
}

// PutState put state to chain
func PutState(key string, field string, value string) ResultCode {
	// prepare param
	ec := NewEasyCodec()
	ec.AddString("key", key)
	ec.AddString("field", field)
	ec.AddBytes("value", []byte(value))
	b := ec.Marshal()
	// send req put value
	code := sysCall(getRequestHeader(ContractMethodPutState), string(b))
	if code != int32(SUCCESS) {
		return ERROR
	}
	return SUCCESS
}

// PutState put state to chain
func PutStateByte(key string, field string, value []byte) ResultCode {
	return PutState(key, field, string(value))
}

// PutStateFromKey put state to chain
func PutStateFromKey(key string, value string) ResultCode {
	return PutState(key, "", value)
}

// PutStateFromKey put state to chain
func PutStateFromKeyByte(key string, value []byte) ResultCode {
	return PutStateByte(key, "", value)
}

// DeleteState delete state to chain
func DeleteState(key string, field string) ResultCode {
	// prepare param
	ec := NewEasyCodec()
	ec.AddString("key", key)
	ec.AddString("field", field)
	b := ec.Marshal()
	// send req put value
	code := sysCall(getRequestHeader(ContractMethodDeleteState), string(b))
	if code != int32(SUCCESS) {
		return ERROR
	}
	return SUCCESS
}

// CallContract call other contract from chain
func CallContract(contractName string, method string, param map[string][]byte) ([]byte, ResultCode) {
	// # get len
	// ## prepare param
	var valueLen int32 = 0
	valuePtr := int32(uintptr(unsafe.Pointer(&valueLen)))

	ec := NewEasyCodec()
	ecMap := NewEasyCodecWithMap(param)
	paramBytes := ecMap.Marshal()
	ec.AddBytes("param", paramBytes)
	ec.AddInt32("value_ptr", valuePtr)
	ec.AddString("contract_name", contractName)
	ec.AddString("method", method)
	b := ec.Marshal()
	// ## send req get call len
	code := sysCall(getRequestHeader(ContractMethodCallContractLen), string(b))
	if code != int32(SUCCESS) {
		return nil, ERROR
	}
	if valueLen == 0 {
		return nil, SUCCESS
	}

	// # get data
	// ## prepare param
	valueByte := make([]byte, valueLen)
	valuePtr = int32(uintptr(unsafe.Pointer(&valueByte[0])))
	ec.RemoveKey("value_ptr")
	ec.AddInt32("value_ptr", valuePtr)
	b = ec.Marshal()
	// ## send req get value
	code2 := sysCall(getRequestHeader(ContractMethodCallContract), string(b))
	if code2 != int32(SUCCESS) {
		return nil, ERROR
	}
	return valueByte, SUCCESS
}

func DeleteStateFromKey(key string) ResultCode {
	return DeleteState(key, "")
}

// SuccessResult record success data
func SuccessResult(msg string) {
	sysCall(getRequestHeader(ContractMethodSuccessResult), msg)
}

// SuccessResult record success data
func SuccessResultByte(msg []byte) {
	sysCall(getRequestHeader(ContractMethodSuccessResult), string(msg))
}

// ErrorResult record error msg
func ErrorResult(msg string) {
	sysCall(getRequestHeader(ContractMethodErrorResult), string(msg))
}

func GetCreatorOrgId() (string, ResultCode) {
	return stringArg(ContractParamCreatorOrgId)
}
func GetCreatorRole() (string, ResultCode) {
	return stringArg(ContractParamCreatorRole)
}
func GetCreatorPk() (string, ResultCode) {
	return stringArg(ContractParamCreatorPk)
}
func GetSenderOrgId() (string, ResultCode) {
	return stringArg(ContractParamSenderOrgId)
}
func GetSenderRole() (string, ResultCode) {
	return stringArg(ContractParamSenderRole)
}
func GetSenderPk() (string, ResultCode) {
	return stringArg(ContractParamSenderPk)
}
func GetBlockHeight() (string, ResultCode) {
	return stringArg(ContractParamBlockHeight)
}
func GetTxId() (string, ResultCode) {
	return stringArg(ContractParamTxId)
}
func getCtxPtr() int32 {
	if str, resultCode := stringArg(ContractParamContextPtr); resultCode != SUCCESS {
		LogMessage("failed to get ctx ptr")
		return 0
	} else {
		ptr, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			LogMessage("get ptr err: " + err.Error())
		}
		return int32(ptr)
	}
}

//func getArgsMap() error {
//	if !argsFlag {
//		argsMap = EasyUnmarshal(argsBytes)
//		argsFlag = true
//	}
//	return nil
//}

// TODO 新的从链上读取参数的实现，待检查
func getArgsMap() ([]*EasyCodecItem, error) {
	// # get len
	// ## prepare param
	var argsMap []*EasyCodecItem
	ec := NewEasyCodec()
	argsBytes, code := GetBytesFromChain(ec, ContractMethodArgLen, ContractMethodArg)
	if code != SUCCESS {
		return nil, fmt.Errorf("get arg from chain error")
	}
	if len(argsBytes) > 0 {
		return nil, fmt.Errorf("args len is 0")
	}
	argsMap = EasyUnmarshal(argsBytes)
	return argsMap, nil
}

func stringArg(key string) (string, ResultCode) {
	result, code := Arg(key)
	return string(result), code
}
func Arg(key string) ([]byte, ResultCode) {
	argsMap, err := getArgsMap()
	if err != nil {
		LogMessage("get Arg error:" + err.Error())
		return nil, ERROR
	}
	for _, v := range argsMap {
		if v.Key == key {
			return v.Value.([]byte), SUCCESS
		}
	}
	return nil, ERROR
}
func ArgString(key string) (string, ResultCode) {
	argsMap, err := getArgsMap()
	if err != nil {
		LogMessage("get Arg error:" + err.Error())
		return "", ERROR
	}
	for _, v := range argsMap {
		if v.Key == key {
			return string(v.Value.([]byte)), SUCCESS
		}
	}
	return "", ERROR
}

func Args() []*EasyCodecItem {
	argsMap, err := getArgsMap()
	if err != nil {
		LogMessage("get Args error:" + err.Error())
	}
	return argsMap
}

func (s *SimContextImpl) newIterator(startKey string, startField string, limitKey string, limitField string) (ResultSetKV, ResultCode) { //main.go中调用
	ec := NewEasyCodec()

	ec.AddString("start_key", startKey)
	ec.AddString("start_field", startField)
	ec.AddString("limit_key", limitKey)
	ec.AddString("limit_field", limitField)
	index, code := GetInt32FromChain(ec, ContractMethodKvIterator)
	return &ResultSetKvImpl{index}, code
}

func (s *SimContextImpl) NewIteratorWithField(key string, startField string, limitField string) (ResultSetKV, ResultCode) {
	return s.newIterator(key, startField, key, limitField)
}

// NewIterator
func (s *SimContextImpl) NewIterator(key string, limit string) (ResultSetKV, ResultCode) {
	return s.newIterator(key, "", limit, "")
}

func (s *SimContextImpl) NewIteratorPrefixWithKeyField(startKey string, startField string) (ResultSetKV, ResultCode) {
	ec := NewEasyCodec()
	ec.AddString("start_key", startKey)
	ec.AddString("start_field", startField)
	index, code := GetInt32FromChain(ec, ContractMethodKvPreIterator)
	return &ResultSetKvImpl{index}, code
}
func (s *SimContextImpl) NewHistoryKvIterForKey(startKey string, startField string) (KeyHistoryKvIter, ResultCode) {
	ec := NewEasyCodec()
	ec.AddString("start_key", startKey)
	ec.AddString("start_field", startField)
	index, code := GetInt32FromChain(ec, ContractHistoryKvIterator)
	return &KeyHistoryKvIterImpl{
		key:   startKey,
		field: startField,
		index: index,
	}, code
}
func (s *SimContextImpl) NewIteratorPrefixWithKey(key string) (ResultSetKV, ResultCode) {
	return s.NewIteratorPrefixWithKeyField(key, "")
}

// ResultSet iterator query result KVdb
type ResultSetKvImpl struct { //为kv查询后的上下文
	index int32 // 链的句柄的index
}

func (r *ResultSetKvImpl) HasNext() bool {
	ec := NewEasyCodec()
	ec.AddInt32("rs_index", r.index)
	data, _ := GetInt32FromChain(ec, ContractMethodKvIteratorHasNext)
	return data != 0
}

func (r *ResultSetKvImpl) NextRow() (*EasyCodec, ResultCode) {
	ec := NewEasyCodec()
	ec.AddInt32("rs_index", r.index)
	bytes, code := GetBytesFromChain(ec, ContractMethodKvIteratorNextLen, ContractMethodKvIteratorNext)
	if code != SUCCESS {
		return nil, ERROR
	}
	ec = NewEasyCodecWithBytes(bytes)
	return ec, code
}

func (r *ResultSetKvImpl) Close() (bool, ResultCode) {
	ec := NewEasyCodec()
	ec.AddInt32("rs_index", r.index)
	data, code := GetInt32FromChain(ec, ContractMethodKvIteratorClose)
	return data != 0, code
}

func (r *ResultSetKvImpl) Next() (string, string, []byte, ResultCode) {
	ec, code := r.NextRow()
	if code != SUCCESS {
		return "", "", nil, ERROR
	}
	k, _ := ec.GetString("key")
	field, _ := ec.GetString("field")
	v, _ := ec.GetBytes("value")
	return k, field, v, code
}

type KeyHistoryKvIterImpl struct {
	key   string
	field string
	index int32
}

func (k *KeyHistoryKvIterImpl) HasNext() bool {
	ec := NewEasyCodec()
	ec.AddInt32("ks_index", k.index)
	data, _ := GetInt32FromChain(ec, ContractHistoryKvIteratorHasNext)
	return data != 0
}

func (k *KeyHistoryKvIterImpl) NextRow() (*EasyCodec, ResultCode) {
	ec := NewEasyCodec()
	ec.AddInt32("ks_index", k.index)
	bytes, code := GetBytesFromChain(ec, ContractHistoryKvIteratorNextLen, ContractHistoryKvIteratorNext)
	if code != SUCCESS {
		return nil, ERROR
	}
	ec = NewEasyCodecWithBytes(bytes)
	return ec, code
}

func (k *KeyHistoryKvIterImpl) Close() (bool, ResultCode) {
	ec := NewEasyCodec()
	ec.AddInt32("ks_index", k.index)
	data, code := GetInt32FromChain(ec, ContractHistoryKvIteratorClose)
	return data != 0, code
}

func (k *KeyHistoryKvIterImpl) Next() (*KeyModification, ResultCode) {
	ec, code := k.NextRow()
	if code != SUCCESS {
		return nil, ERROR
	}
	value, _ := ec.GetBytes("value")
	txId, _ := ec.GetString("txId")
	blockHeight, _ := ec.GetInt32("blockHeight")
	isDeleteBool, _ := ec.GetInt32("isDelete")
	timestamp, _ := ec.GetString("timestamp")
	isDelete := false
	if isDeleteBool == 1 {
		isDelete = true
	}

	return &KeyModification{
		Key:         k.key,
		Field:       k.field,
		Value:       value,
		TxId:        txId,
		BlockHeight: int(blockHeight),
		IsDelete:    isDelete,
		Timestamp:   timestamp,
	}, code
}
