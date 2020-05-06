package pay

import (
	"encoding/xml"
	"fmt"
	"github.com/antsbean/wechat/util"
)

var queryGateway = "https://api.mch.weixin.qq.com/pay/orderquery"

// CloseOrderParams close order params
type QueryOrderParams struct {
	SubAppID      string
	SubMchID      string
	TransactionID string
	OutTradeNo    string
	SignType      string
}

// closeOrderRequest 接口请求参数
type queryOrderRequest struct {
	CommonRequest
	OutTradeNo    string `xml:"out_trade_no"` // 商户订单号
	TransactionID string `xml:"transaction_id"`
}

type QueryResponse struct {
	CommonResponse
	ResultMsg string `xml:"result_msg,omitempty"`

	OpenID             string `xml:"openid"`
	SubOpenID          string `xml:"sub_openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           int    `xml:"total_fee"`
	SettlementTotalFee int    `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            string `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          int    `xml:"coupon_fee"`
	CouponCount        int    `xml:"coupon_count"`

	// coupon_type_$n 这里只声明 3 个，如果有更多的可以自己组合
	CouponType0 string `xml:"coupon_type_0"`
	CouponType1 string `xml:"coupon_type_1"`
	CouponType2 string `xml:"coupon_type_2"`
	CouponID0   string `xml:"coupon_id_0"`
	CouponID1   string `xml:"coupon_id_1"`
	CouponID2   string `xml:"coupon_id_2"`
	CouponFeed0 string `xml:"coupon_fee_0"`
	CouponFeed1 string `xml:"coupon_fee_1"`
	CouponFeed2 string `xml:"coupon_fee_2"`

	TransactionID  string `xml:"transaction_id"`
	OutTradeNo     string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	TimeEnd        string `xml:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc"`
}

// QueryOrder query order
func (pcf *Pay) QueryOrder(q *QueryOrderParams) (rsp QueryResponse, err error) {
	if q.SignType == "" {
		q.SignType = "MD5"
	}
	nonceStr := util.RandomStr(32)
	param := make(map[string]interface{})
	param["appid"] = pcf.AppID
	param["mch_id"] = pcf.PayMchID
	param["sub_appid"] = q.SubAppID
	param["sub_mch_id"] = q.SubMchID
	param["nonce_str"] = nonceStr
	param["transaction_id"] = q.TransactionID
	param["out_trade_no"] = q.OutTradeNo
	param["sign_type"] = q.SignType

	bizKey := "&key=" + pcf.PayKey
	str := orderParam(param, bizKey)
	sign := util.MD5Sum(str)
	request := queryOrderRequest{
		CommonRequest: CommonRequest{
			AppID:    pcf.AppID,
			MchID:    pcf.PayMchID,
			SubAppID: q.SubAppID,
			SubMchID: q.SubMchID,
			NonceStr: nonceStr,
			Sign:     sign,
			SignType: q.SignType,
		},
		OutTradeNo:    q.OutTradeNo,
		TransactionID: q.TransactionID,
	}
	rawRet, err := util.PostXML(queryGateway, request)
	if err != nil {
		return
	}
	err = xml.Unmarshal(rawRet, &rsp)
	if err != nil {
		return
	}
	if rsp.ReturnCode == "SUCCESS" {
		if rsp.ResultCode == "SUCCESS" {
			err = nil
			return
		}
		err = fmt.Errorf("refund error, errcode=%s,errmsg=%s", rsp.ErrCode, rsp.ErrCodeDes)
		return
	}
	err = fmt.Errorf("[msg : xmlUnmarshalError] [rawReturn : %s] [params : %s] [sign : %s]",
		string(rawRet), str, sign)
	return
}
