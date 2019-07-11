package payment

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/medivhzhan/weapp/util"
)

const (
	refundQueryAPI = "/pay/refundquery"
)

// 退款查询条件
type RefundQuery struct {
	AppID    string `xml:"appid"`               // 小程序ID
	MchID    string `xml:"mch_id"`              // 商户号
	SignType string `xml:"sign_type,omitempty"` // 签名类型: 目前支持HMAC-SHA256和MD5，默认为MD5
	NonceStr string `xml:"nonce_str"`           // 随机字符串
	Sign     string `xml:"sign"`                // 签名

	// 官方文档为四选一
	// 但是以订单号作为查询会返回多次退款的数据
	// 根据退款单号作为查询则只会返回本次的退款数据
	// OutTradeNo    string `xml:"out_trade_no,omitempty"`   // 商户订单号
	// TransactionID string `xml:"transaction_id,omitempty"` // 微信支付订单号
	OutRefundNo string `xml:"out_refund_no,omitempty"` // 商户退款单号
	RefundID    string `xml:"refund_id,omitempty"`     // 微信退款单号
}

type refundQuery struct {
	RefundQuery

	XMLName  xml.Name `xml:"xml"`
	Sign     string   `xml:"sign"`                // 签名
	NonceStr string   `xml:"nonce_str"`           // 随机字符串
	SignType string   `xml:"sign_type,omitempty"` // 签名类型: 目前支持HMAC-SHA256和MD5，默认为MD5
}

/*

以订单查询返回信息

<?xml version="1.0" encoding="utf-8"?>

<xml>
  <appid>000000</appid>
  <cash_fee>6</cash_fee>
  <mch_id>000000</mch_id>
  <nonce_str>GwUbVYPcPZz0T8Lg</nonce_str>
  <out_refund_no_0>20190705171118247841</out_refund_no_0>
  <out_refund_no_1>20190705171048982976</out_refund_no_1>
  <out_trade_no>20190705170959530641</out_trade_no>
  <refund_account_0>REFUND_SOURCE_UNSETTLED_FUNDS</refund_account_0>
  <refund_account_1>REFUND_SOURCE_UNSETTLED_FUNDS</refund_account_1>
  <refund_channel_0>ORIGINAL</refund_channel_0>
  <refund_channel_1>ORIGINAL</refund_channel_1>
  <refund_count>2</refund_count>
  <refund_fee>4</refund_fee>
  <refund_fee_0>1</refund_fee_0>
  <refund_fee_1>3</refund_fee_1>
  <refund_id_0>50000201152019070510323842871</refund_id_0>
  <refund_id_1>50000201152019070510395660782</refund_id_1>
  <refund_recv_accout_0>工商银行借记卡3234</refund_recv_accout_0>
  <refund_recv_accout_1>工商银行借记卡3234</refund_recv_accout_1>
  <refund_status_0>SUCCESS</refund_status_0>
  <refund_status_1>SUCCESS</refund_status_1>
  <refund_success_time_0>2019-07-05 17:15:51</refund_success_time_0>
  <refund_success_time_1>2019-07-05 17:15:19</refund_success_time_1>
  <result_code>SUCCESS</result_code>
  <return_code>SUCCESS</return_code>
  <return_msg>OK</return_msg>
  <sign>000000</sign>
  <total_fee>6</total_fee>
  <transaction_id>4200000289201907058815070594</transaction_id>
</xml>
*/

/*

以退款单号查询返回信息

<?xml version="1.0" encoding="utf-8"?>

<xml>
  <appid>000000</appid>
  <cash_fee>6</cash_fee>
  <mch_id>000000</mch_id>
  <nonce_str>ORXm2pdxYtc2VoUd</nonce_str>
  <out_refund_no_0>20190705171118247841</out_refund_no_0>
  <out_trade_no>20190705170959530641</out_trade_no>
  <refund_account_0>REFUND_SOURCE_UNSETTLED_FUNDS</refund_account_0>
  <refund_channel_0>ORIGINAL</refund_channel_0>
  <refund_count>1</refund_count>
  <refund_fee>1</refund_fee>
  <refund_fee_0>1</refund_fee_0>
  <refund_id_0>50000201152019070510323842871</refund_id_0>
  <refund_recv_accout_0>工商银行借记卡3234</refund_recv_accout_0>
  <refund_status_0>SUCCESS</refund_status_0>
  <refund_success_time_0>2019-07-05 17:15:51</refund_success_time_0>
  <result_code>SUCCESS</result_code>
  <return_code>SUCCESS</return_code>
  <return_msg>OK</return_msg>
  <sign>000000</sign>
  <total_fee>6</total_fee>
  <transaction_id>4200000289201907058815070594</transaction_id>
</xml>
*/

// 退款查询响应
type RefundQueryResponse struct {
	AppID    string `xml:"appid"`     // 小程序ID
	MchID    string `xml:"mch_id"`    // 商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	Sign     string `xml:"sign"`      // 签名

	TransactionID string `xml:"transaction_id"`  // 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`    // 商户订单号
	RefundID      string `xml:"refund_id_0"`     // 微信退款单号
	OutRefundNo   string `xml:"out_refund_no_0"` // 商户退款单号

	TotalFee  int `xml:"total_fee"`  // 订单金额
	RefundFee int `xml:"refund_fee"` // 退款总金额,单位为分

	RefundStatus RefundStatusType `xml:"refund_status_0"` // 退款状态

	SuccessTime    string `xml:"refund_success_time_0,omitempty"` // 退款成功时间2017-12-15 09:46:01
	ReceiveAccount string `xml:"refund_recv_accout_0"`            // 退款入账账户:取当前退款单的退款入账方
}

type refundQueryResponse struct {
	response
	RefundQueryResponse
}

type RefundStatusType string

const (
	RefundStatusSuccess    RefundStatusType = "SUCCESS"
	RefundStatusChange     RefundStatusType = "CHANGE"
	RefundStatusClose      RefundStatusType = "REFUNDCLOSE"
	RefundStatusProcessing RefundStatusType = "PROCESSING"
)

// 查询退款订单
func (s RefundQuery) Query(key string) (*RefundQueryResponse, error) {
	data, err := s.prepare(key)
	if err != nil {
		return nil, err
	}

	resData, err := util.PostXML(baseURL+refundQueryAPI, data)
	fmt.Println(string(resData))
	if err != nil {
		return nil, err
	}

	result := &refundQueryResponse{}
	if err = xml.Unmarshal(resData, &result); err != nil {
		return nil, err
	}
	err = result.Check()
	if err != nil {
		return nil, err
	}

	return &result.RefundQueryResponse, nil
}

// 预处理
func (s RefundQuery) prepare(key string) (refundQuery, error) {
	result := refundQuery{
		RefundQuery: s,
		SignType:    "MD5",
		NonceStr:    util.RandomString(32),
	}

	signData := map[string]string{
		"appid":     result.AppID,
		"mch_id":    result.MchID,
		"nonce_str": result.NonceStr,
		"sign_type": result.SignType,
	}

	if result.OutRefundNo == "" && result.RefundID == "" {
		return result, errors.New("out_refund_no与refund_id必须二选一")
	}
	if result.OutRefundNo != "" {
		signData["out_refund_no"] = result.OutRefundNo
	}
	if result.RefundID != "" {
		signData["refund_id"] = result.RefundID
	}

	sign, err := util.SignByMD5(signData, key)
	result.Sign = sign
	return result, err
}
