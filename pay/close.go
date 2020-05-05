package pay

import (
	"encoding/xml"
	"fmt"
	"github.com/antsbean/wechat/util"
)

var closeGateway = "https://api.mch.weixin.qq.com/pay/closeorder"

// CloseOrderParams close order params
type CloseOrderParams struct {
	SubAppID   string
	SubMchID   string
	OutTradeNo string
	SignType   string
}

// closeOrderRequest 接口请求参数
type closeOrderRequest struct {
	CommonRequest
	OutTradeNo string `xml:"out_trade_no"` // 商户订单号
}

type CloseResponse struct {
	CommonResponse
	ResultMsg string `xml:"result_msg,omitempty"`
}

// CloseOrder close order
func (pcf *Pay) CloseOrder(c *CloseOrderParams) (rsp CloseResponse, err error) {
	nonceStr := util.RandomStr(32)
	param := make(map[string]interface{})
	param["appid"] = pcf.AppID
	param["mch_id"] = pcf.PayMchID
	param["sub_appid"] = c.SubAppID
	param["sub_mch_id"] = c.SubMchID
	param["nonce_str"] = nonceStr
	param["out_trade_no"] = c.OutTradeNo
	param["sign_type"] = "MD5"

	bizKey := "&key=" + pcf.PayKey
	str := orderParam(param, bizKey)
	sign := util.MD5Sum(str)
	request := closeOrderRequest{
		CommonRequest: CommonRequest{
			AppID:    pcf.AppID,
			MchID:    pcf.PayMchID,
			SubAppID: c.SubAppID,
			SubMchID: c.SubMchID,
			NonceStr: nonceStr,
			Sign:     sign,
			SignType: "MD5",
		},
		OutTradeNo: c.OutTradeNo,
	}
	rawRet, err := util.PostXML(closeGateway, request)
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
