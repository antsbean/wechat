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
	getPath           = "https://api.weixin.qq.com/cgi-bin/express/business/path/get?access_token=%s"
)

// LogisticsBindAccount Logistics Account info
type LogisticsBindAccount struct {
	Type          string `json:"type"`
	BizID         string `json:"biz_id"`
	DeliveryID    string `json:"delivery_id"`
	Password      string `json:"password,omitempty"`
	RemarkContent string `json:"remark_content,omitempty"`
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
	DeliveryID string `json:"delivery_id" validate:"required"` // 快递公司ID
	WaybillID  string `json:"waybill_id" validate:"required"`  // 运单ID
	OpenID     string `json:"open_id"`                         // 用户openid，当add_source=2时无需填写（不发送物流服务通知）
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
	urlStr := fmt.Sprintf(getPath, accessToken)
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
