package pay

// CommonRequest common request
type CommonRequest struct {
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	SubAppID   string `xml:"sub_appid"`
	SubMchID   string `xml:"sub_mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	SignType   string `xml:"sign_type,omitempty"`
}

// CommonResponse wechat pay common response
type CommonResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid,omitempty"`
	MchID      string `xml:"mch_id,omitempty"`
	SubAppID   string `xml:"sub_appid"`
	SubMchID   string `xml:"sub_mch_id"`
	NonceStr   string `xml:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty"`
	ResultCode string `xml:"result_code,omitempty"`
	ErrCode    string `xml:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty"`
}
