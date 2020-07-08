package miniprogram

import (
	"encoding/json"
	"fmt"

	"github.com/antsbean/wechat/util"
)

const (
	bindAccountURL    = "https://api.weixin.qq.com/cgi-bin/express/business/account/bind?access_token=%s"
	getAllAccountURL  = "https://api.weixin.qq.com/cgi-bin/express/business/account/getall?access_token=%s"
	getAllDeliveryURL = "https://api.weixin.qq.com/cgi-bin/express/business/delivery/getall?access_token=%s"
	getPathURL        = "https://api.weixin.qq.com/cgi-bin/express/business/path/get?access_token=%s"
	addOrderURL       = "https://api.weixin.qq.com/cgi-bin/express/local/business/order/add?access_token=%s"
	cancelOrderURL    = "https://api.weixin.qq.com/cgi-bin/express/business/order/cancel?access_token=%s"
	getQuotaURL       = "https://api.weixin.qq.com/cgi-bin/express/business/quota/get?access_token=%s"
)

type BindType string

const (
	// BindAccount 绑定账户
	BindAccount = "bind"
	// UnbindAccount 解除绑定
	UnbindAccount = "unbind"
)

// LogisticsBindAccount Logistics Account info
type LogisticsBindAccount struct {
	Type          BindType `json:"type"`
	BizID         string   `json:"biz_id" validate:"required"`
	DeliveryID    string   `json:"delivery_id" validate:"required"`
	Password      string   `json:"password,omitempty"`
	RemarkContent string   `json:"remark_content,omitempty"`
}

// LogisticsBindAccount bind account
func (wxa *MiniProgram) LogisticsBindAccount(account *LogisticsBindAccount) (commonErr util.CommonError, err error) {
	var accessToken string
	accessToken, err = wxa.GetAccessToken()
	if err != nil {
		return
	}
	urlStr := fmt.Sprintf(bindAccountURL, accessToken)
	resultData, err := util.PostJSON(urlStr, account)
	if err != nil {
		return
	}
	return util.DecodeToCommonError(resultData)
}

type LogisticsServerType struct {
	ServiceType int32  `json:"service_type"` // 服务ID
	ServiceName string `json:"service_name"` // 服务名称
}

type StatusCodeEnum = int32

const (
	BindSuccess StatusCodeEnum = iota // 绑定成功
	Auditing                          // 审核中
	BindFailed                        // 绑定失败
	UnBind                            // 解绑
)

// LogisticsAccount 物流账号信息
type LogisticsAccount struct {
	BizID           string                `json:"biz_id"`      //  快递公司客户编码
	DeliveryID      string                `json:"delivery_id"` // 快递公司ID
	CrateTime       int64                 `json:"create_time"`
	UpdateTime      int64                 `json:"update_time"`
	StatusCode      StatusCodeEnum        `json:"status_code"` // 绑定状态
	Alias           string                `json:"alias"`
	RemarkWrongMsg  string                `json:"remark_wrong_msg"`
	RemarkContent   string                `json:"remark_Content"`
	QuotaNum        int32                 `json:"quota_num"`         // 电子面单余额
	QuotaUpdateTime int64                 `json:"quota_update_time"` // 电子免单更新时间
	ServiceTypes    []LogisticsServerType `json:"service_type"`      // 支持的服务类型
}

// LogisticsGetAllAccount get all account
func (wxa *MiniProgram) LogisticsGetAllAccount() (accounts []*LogisticsAccount, count int32, err error) {
	var accessToken string
	accessToken, err = wxa.GetAccessToken()
	if err != nil {
		return
	}
	urlStr := fmt.Sprintf(getAllAccountURL, accessToken)
	resultData, err := util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	var result struct {
		util.CommonError
		Count    int32               `json:"count"`
		Accounts []*LogisticsAccount `json:"list"`
	}
	if err = json.Unmarshal(resultData, &result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = result.Error("LogisticsGetAllAccount")
		return
	}
	accounts = result.Accounts
	count = result.Count
	return
}

// Delivery delivery info
type Delivery struct {
	DeliveryID   string                `json:"delivery_id"`   // 快递公司ID
	DeliveryName string                `json:"delivery_name"` // 快递公司名称
	CanUseCash   uint8                 `json:"can_use_cash"`  // 是否支持散单,1表示支持
	CanGetQuota  uint8                 `json:"can_get_quota"` // 是否支持查询面单余额,1表示支持
	CashBizID    string                `json:"cash_biz_id"`   // 散单对应的bizid，当can_use_cash=1时有效
	ServiceTypes []LogisticsServerType `json:"service_type"`  // 支持的服务类型
}

// LogisticsGetAllDelivery get all delivery
func (wxa *MiniProgram) LogisticsGetAllDelivery() (accounts []*Delivery, count int32, err error) {
	var accessToken string
	accessToken, err = wxa.GetAccessToken()
	if err != nil {
		return
	}
	urlStr := fmt.Sprintf(getAllDeliveryURL, accessToken)
	resultData, err := util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	var result struct {
		util.CommonError
		Count      int32       `json:"count"`
		Deliveries []*Delivery `json:"data"`
	}
	if err = json.Unmarshal(resultData, &result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = result.Error("LogisticsGetAllAccount")
		return
	}
	accounts = result.Deliveries
	count = result.Count
	return
}

// PathRequest get path request
type LogisticsPathRequest struct {
	OrderID    string `json:"order_id" validate:"required"`    // 订单ID
	OpenID     string `json:"open_id"`                         // 用户openid，当add_source=2时无需填写（不发送物流服务通知）
	DeliveryID string `json:"delivery_id" validate:"required"` // 快递公司ID
	WaybillID  string `json:"waybill_id" validate:"required"`  // 运单ID
}

type ActionTypeEnum int32

const (
	CargoCanvassingSuccess    ActionTypeEnum = iota + 100001 // 揽件
	CargoCanvassingFailed                                    // 揽件失败
	CargoCanvassingAssigning                                 // 揽件分配阶段
	CarriagePathUpdated       ActionTypeEnum = 200001        // 运输轨迹更新
	DistributeExpressDelivery ActionTypeEnum = 300002        // 开发派发快递
	SignExpressDelivery                                      // 签收快递
	SignFailedExpressDelivery                                // 签收快递
	CancelOrder               ActionTypeEnum = 400001        // 取消订单
	StrandedDelivery                                         // 滞留快递

)

// PathRequest path item info
type LogisticsPathItem struct {
	ActionTime int64          `json:"action_time" xml:"ActionTime"` // 轨迹节点 Unix 时间戳
	ActionType ActionTypeEnum `json:"action_type" xml:"ActionType"` // 轨迹节点类型
	ActionMsg  string         `json:"action_msg" xml:"ActionMsg"`   // 轨迹节点详情
}

// LogisticsGetPath get path
func (wxa *MiniProgram) LogisticsGetPath(request *LogisticsPathRequest) (items []*LogisticsPathItem, count int32, err error) {
	var accessToken string
	accessToken, err = wxa.GetAccessToken()
	if err != nil {
		return
	}
	urlStr := fmt.Sprintf(getPathURL, accessToken)
	resultData, err := util.PostJSON(urlStr, request)
	if err != nil {
		return
	}
	var result struct {
		util.CommonError
		PathItemNum int32                `json:"path_item_num"`
		Items       []*LogisticsPathItem `json:"path_item_list"`
	}
	if err = json.Unmarshal(resultData, &result); err != nil {
		return
	}
	count = result.PathItemNum
	items = result.Items
	return

}

// SenderReceiver sender receiver info
type SenderReceiver struct {
	Name     string `json:"name"`      // 是	姓名，最长不超过256个字符
	Tel      string `json:"tel"`       // 否 发件人座机号码，若不填写则必须填写 mobile，不超过32字节
	Mobile   string `json:"mobile"`    // 是	电话/手机号，最长不超过64个字符
	Company  string `json:"company"`   // 否 公司名称
	PostCode string `json:"post_code"` // 否 邮政编码
	Country  string `json:"country"`   // 否 国家
	Province string `json:"province"`  // 否 国家
	City     string `json:"city"`      // 是	城市名称，如广州市
	Area     string `json:"area"`      // 区域
	Address  string `json:"address"`   // 是	地址(街道、小区、大厦等，用于定位)
}

// CargoGoods cargo goods
type CargoGoods struct {
	Count int32  `json:"count"` // 是	货物数量
	Name  string `json:"name"`  // 是	货品名称
}

// LogisticsCargo cargo
type LogisticsCargo struct {
	Count      int32         `json:"count"`       // 是	包裹数量, 需要和detail_list size保持一致
	Weight     int32         `json:"weight"`      // 是	包裹总重量，单位是千克(kg)
	SpaceX     int32         `json:"space_x"`     // 是	包裹长度，单位厘米(cm)
	SpaceY     int32         `json:"space_y"`     // 是	包裹宽度，单位厘米(cm)
	SpaceZ     int32         `json:"space_z"`     // 是	包裹高度，单位厘米(cm)
	DetailList []*CargoGoods `json:"detail_list"` // 是	包裹中商品详情列表
}

// LogisticsShop shop
type LogisticsShop struct {
	WxaPath    string `json:"wxa_path"`
	ImgURL     string `json:"img_url"`
	GoodsName  string `json:"goods_name"`
	GoodsCount int32  `json:"goods_count"`
}

type LogisticsInsured struct {
	UseInsured   int32 `json:"use_insured"`
	InsuredValue int32 `json:"insured_value"`
}

// LogisticsOrderRequest order request info
type LogisticsOrderRequest struct {
	AddSource       int32               `json:"add_source"`    // 是 订单来源，0为小程序订单，2为App或H5订单，填2则不发送物流服务通知
	WXAppID         string              `json:"wx_appid"`      // 否 App或H5的appid，add_source=2时必填，需和开通了物流助手的小程序绑定同一open帐号
	OrderID         string              `json:"order_id"`      // 是 订单ID，须保证全局唯一，不超过512字节
	OpenID          string              `json:"openid"`        // 否 用户openid，当add_source=2时无需填写（不发送物流服务通知）
	DeliveryID      string              `json:"delivery_id"`   // 是 配送公司ID
	BizID           string              `json:"biz_id"`        // 是 快递客户编码或者现付编码
	CustomRemark    string              `json:"custom_remark"` // 否 快递备注信息，比如"易碎物品"，不超过1024字节
	TagID           int32               `json:"tagid"`         // 否 订单标签id，用于平台型小程序区分平台上的入驻方，tagid须与入驻方账号一一对应，非平台型小程序无需填写该字段
	Sender          SenderReceiver      `json:"sender"`        // 是 发件人信息，顺丰同城急送必须填写，美团配送、达达、闪送，若传了shop_no的值可不填该字段
	Receiver        SenderReceiver      `json:"receiver"`      // 是 收件人信息
	Cargo           LogisticsCargo      `json:"cargo"`         // 是 货物信息
	Shop            LogisticsShop       `json:"shop"`          // 是 商品信息，会展示到物流通知消息中
	Insured         LogisticsInsured    `json:"insured"`       // 是 保价信息
	LogisticsServer LogisticsServerType `json:"service"`       // 是 服务类型
	ExpectTime      int64               `json:"expect_time"`   // 否 Unix 时间戳, 单位秒，顺丰必须传。 预期的上门揽件时间，0表示已事先约定取件时间；否则请传预期揽件时间戳，需大于当前时间，收件员会在预期时间附近上门。例如expect_time为“1557989929”，表示希望收件员将在2019年05月16日14:58:49-15:58:49内上门取货。说明：若选择 了预期揽件时间，请不要自己打单，由上门揽件的时候打印。如果是下顺丰散单，则必传此字段，否则不会有收件员上门揽件。

}

type DeliveryCommonError struct {
	util.CommonError
	DeliveryResultCode int32  `json:"delivery_resultcode"` // 运力返回的错误码
	DeliveryResultMsg  string `json:"delivery_resultmsg"`  // 运力返回的错误描述
}

type WayBillData struct {
	Key   string `json:"key"`   //运单信息key
	Value string `json:"value"` //运单信息value
}

// LogisticsOrderResponse order response
type LogisticsOrderResponse struct {
	DeliveryCommonError
	OrderID      string         `json:"order_id"`     // 订单ID
	WayBillID    string         `json:"waybill_id"`   // 运单ID
	WayBillInfos []*WayBillData `json:"waybill_data"` // 运单信息
}

// LogisticsAddOrder add order
func (wxa *MiniProgram) LogisticsAddOrder(request *LogisticsOrderRequest) (*LogisticsOrderResponse, error) {
	var accessToken string
	accessToken, err := wxa.GetAccessToken()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf(addOrderURL, accessToken)
	resultData, err := util.PostJSON(urlStr, request)
	if err != nil {
		return nil, err
	}
	var resp LogisticsOrderResponse
	if err := json.Unmarshal(resultData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelOrderRequest cancel order request
type CancelOrderRequest struct {
	OrderID    string `json:"order_id" validate:"required"`    // 订单ID
	OpenID     string `json:"open_id"`                         // 用户openid，当add_source=2时无需填写（不发送物流服务通知）
	DeliveryID string `json:"delivery_id" validate:"required"` // 快递公司ID
	WaybillID  string `json:"waybill_id" validate:"required"`  // 运单ID
}

// LogisticsCancelOrder cancel order
func (wxa *MiniProgram) LogisticsCancelOrder(request *CancelOrderRequest) (*DeliveryCommonError, error) {
	var accessToken string
	accessToken, err := wxa.GetAccessToken()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf(cancelOrderURL, accessToken)
	resultData, err := util.PostJSON(urlStr, request)
	if err != nil {
		return nil, err
	}
	var resp DeliveryCommonError
	if err := json.Unmarshal(resultData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetQuota get quota
func (wxa *MiniProgram) GetQuota(deliveryID, bizID string) (int32, error) {
	var accessToken string
	accessToken, err := wxa.GetAccessToken()
	if err != nil {
		return 0, err
	}
	urlStr := fmt.Sprintf(getQuotaURL, accessToken)
	resultData, err := util.PostJSON(urlStr, map[string]string{"delivery_id": deliveryID, "biz_id": bizID})
	if err != nil {
		return 0, err
	}
	var resp struct {
		QuotaNum int32 `json:"quota_num"`
	}
	if err := json.Unmarshal(resultData, &resp); err != nil {
		return 0, err
	}
	return resp.QuotaNum, nil
}
