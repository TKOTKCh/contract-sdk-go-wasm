package main

import (
	"encoding/json"
	"fmt"
	"github.com/TKOTKCh/contract-sdk-go-wasm/sdk"
	"strings"
)

const (
	goodsIdArgKey    = "goodsId"
	nameArgKey       = "name"
	factoryArgKey    = "factory"
	fromArgKey       = "from"
	toArgKey         = "to"
	uploaderArgKey   = "uploader"
	sellerArgKey     = "seller"
	goodsStoreMapKey = "goodsList"
	statusCreate     = iota
	statusTransfer
	statusUpload
	statusSelled
)

type Goods struct {
	GoodsId    string       `json:"GoodsId"`
	Name       string       `json:"Name"`
	Status     uint8        `json:"Status"`
	TraceDatas []*TraceData `json:"TraceDatas"`
}

type TraceData struct {
	Operator     string `json:"Operator"`
	Status       uint8  `json:"Status"`
	OperatorTime string `json:"OperatorTime"`
	Remark       string `json:"Remark"`
}

///////////////////////// 合约入口函数 /////////////////////////

//go:wasmexport init_contract
func InitContract() {
	ctx := sdk.NewSimContext()
	ctx.SuccessResult("Init contract success")
}

//go:wasmexport upgrade
func Upgrade() {
	ctx := sdk.NewSimContext()
	ctx.SuccessResult("Upgrade contract success")
}

///////////////////////// 核心业务函数 /////////////////////////

//go:wasmexport newGoods
func newGoods() {
	ctx := sdk.NewSimContext()

	// 获取参数
	goodsId, _ := ctx.ArgString(goodsIdArgKey)
	name, _ := ctx.ArgString(nameArgKey)
	factory, _ := ctx.ArgString(factoryArgKey)

	// 参数校验
	if isBlank(goodsId) || isBlank(name) || isBlank(factory) {
		ctx.ErrorResult("invalid required parameters")
	}

	// 检查商品是否存在
	if goodsBytes, _ := ctx.GetStateByte(goodsStoreMapKey, goodsId); len(goodsBytes) > 0 {
		ctx.ErrorResult("goodsId is exists!")
	}

	// 构建商品对象
	goods := &Goods{
		GoodsId:    goodsId,
		Name:       name,
		Status:     statusCreate,
		TraceDatas: make([]*TraceData, 0),
	}

	id, _ := ctx.GetSenderPk()
	// 添加初始溯源记录
	goods.TraceDatas = append(goods.TraceDatas, &TraceData{
		Operator: id,
		Status:   statusCreate,
		Remark:   goodsId + ":" + factory + " created"})

	goodsBytes, err1 := json.Marshal(goods)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("newGoods Marshal failed, err: %s", err1))
	}

	err := ctx.PutStateByte(goodsStoreMapKey, goodsId, goodsBytes)
	if err != sdk.SUCCESS {
		ctx.ErrorResult(fmt.Sprintf("newGoods PutStateByte failed, err: %s", err))
	}

	ctx.SuccessResult("newGoods success")
}

//go:wasmexport transferGoods
func transferGoods() {
	ctx := sdk.NewSimContext()

	goodsId, _ := ctx.ArgString(goodsIdArgKey)
	from, _ := ctx.ArgString(fromArgKey)
	to, _ := ctx.ArgString(toArgKey)

	// 参数校验
	if isBlank(goodsId) || isBlank(from) || isBlank(to) {
		ctx.ErrorResult("missing required parameters")
	}
	// 更新状态
	trace := &TraceData{
		Status: statusTransfer,
		Remark: goodsId + ":" + from + "->" + to,
	}

	if !updateGoodsStatus(goodsId, statusTransfer, trace, "transferGoods") {
		ctx.ErrorResult("transfer goods failed")
	}

	ctx.SuccessResult("transferGoods success")
}

//go:wasmexport uploadGoods
func uploadGoods() {
	ctx := sdk.NewSimContext()

	goodsId, _ := ctx.ArgString(goodsIdArgKey)
	if len(goodsId) == 0 {
		ctx.ErrorResult("invalid goodsId")
	}

	uploader, _ := ctx.ArgString(uploaderArgKey)
	if len(uploader) == 0 {
		ctx.ErrorResult("invalid uploader")
	}

	traceData := &TraceData{Status: statusUpload, Remark: goodsId + ":" + uploader + " upload"}

	if !updateGoodsStatus(goodsId, statusUpload, traceData, "uploadGoods") {
		ctx.ErrorResult("upload goods failed")
	}
	ctx.SuccessResult("uploadGoods success")
}

//go:wasmexport sellGoods
func sellGoods() {
	ctx := sdk.NewSimContext()

	goodsId, _ := ctx.ArgString(goodsIdArgKey)
	if len(goodsId) == 0 {
		ctx.ErrorResult("invalid goodsId")
	}

	seller, _ := ctx.ArgString(sellerArgKey)
	if len(seller) == 0 {
		ctx.ErrorResult("invalid seller")
	}

	traceData := &TraceData{Status: statusSelled, Remark: goodsId + ":" + "selled by " + seller}

	if !updateGoodsStatus(goodsId, statusSelled, traceData, "sellGoods") {
		ctx.ErrorResult("sellGoods failed")
	} else {
		ctx.SuccessResult("sellGoods success")
	}
}

//go:wasmexport goodsStatus
func getGoodsStatus() {
	ctx := sdk.NewSimContext()
	goods := getGoodsByArgs("getGoodsStatus")
	if goods == nil {
		ctx.ErrorResult("getGoodsStatus failed")
	}

	var statusBytes []byte
	statusBytes, err1 := json.Marshal(goods.Status)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("getGoodsStatus Marshal goods.Status failed, err: %s", err1))
	}
	ctx.SuccessResult(string(statusBytes))
}

//go:wasmexport traceGoods
func getTraceInfo() {
	ctx := sdk.NewSimContext()
	goods := getGoodsByArgs("getTraceInfo")
	if goods == nil {
		ctx.ErrorResult("getGoodsStatus failed")
	}

	var traceDataBytes []byte
	traceDataBytes, err1 := json.Marshal(goods.TraceDatas)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("getTraceInfo Marshal goods.Status failed, err: %s", err1))
	}
	ctx.SuccessResult(string(traceDataBytes))
}

// /////////////////////// 辅助函数 /////////////////////////
func updateGoodsStatus(goodsId string, status uint8, trace *TraceData, method string) bool {
	ctx := sdk.NewSimContext()

	goodsBytes, err := ctx.GetStateByte(goodsStoreMapKey, goodsId)
	if err != sdk.SUCCESS {
		return false
	}

	var goods Goods
	err1 := json.Unmarshal(goodsBytes, &goods)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf(method+" Unmarshal error : %s", err))
		return false
	}
	goods.Status = status

	id, err := ctx.GetSenderPk()
	if err != sdk.SUCCESS {
		ctx.ErrorResult(fmt.Sprintf(method+" GetSenderOrgId failed, err: %s", err))
		return false
	}

	trace.Operator = id
	goods.TraceDatas = append(goods.TraceDatas, trace)

	//重新存储进去
	goodsBytes, err1 = json.Marshal(goods)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf(method+" Marshal goods failed, err: %s", err))
		return false
	}
	err = ctx.PutStateByte(goodsStoreMapKey, goodsId, goodsBytes)

	if err != sdk.SUCCESS {
		ctx.ErrorResult(fmt.Sprintf(method+" PutStateByte failed, err: %s", err))
		return false
	}
	ctx.SuccessResult(method + " success")
	return true
}

func isBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func getGoodsByArgs(method string) *Goods {
	ctx := sdk.NewSimContext()

	goodsId, _ := ctx.ArgString(goodsIdArgKey)
	if len(goodsId) == 0 {
		return nil
	}

	goodsBytes, err := ctx.GetStateByte(goodsStoreMapKey, goodsId)
	if err != sdk.SUCCESS {
		return nil
	}
	if len(goodsBytes) == 0 {
		return nil
	}

	var goods *Goods
	err1 := json.Unmarshal(goodsBytes, &goods)
	if err1 != nil {
		return nil
	}
	return goods
}

func main() {

}
