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
	addOrderURL       = "https://api.weixin.qq.com/cgi-bin/express/local/business/order/add?access_token=%s"
	cancelOrderURL    = "https://api.weixin.qq.com/cgi-bin/express/business/order/cancel?access_token=%s"
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

// SenderReceiver sender receiver info
type SenderReceiver struct {
	Name           string  `json:"name"`            // 是	姓名，最长不超过256个字符
	City           string  `json:"city"`            // 是	城市名称，如广州市
	Address        string  `json:"address"`         // 是	地址(街道、小区、大厦等，用于定位)
	AddressDetail  string  `json:"addess_detail"`   // 是	地址详情(楼号、单元号、层号)
	Phone          string  `json:"phone"`           // 是	电话/手机号，最长不超过64个字符
	Lng            float32 `json:"lng"`             // 是	经度（火星坐标或百度坐标，和 coordinate_type 字段配合使用，确到小数点后6位
	Lat            float32 `json:"lat"`             // 是	纬度（火星坐标或百度坐标，和 coordinate_type 字段配合使用，精确到小数点后6位）
	CoordinateType int32   `json:"coordinate_type"` // 0	否	坐标类型，0：火星坐标（高德，腾讯地图均采用火星坐标） 1：百度坐标
}

// CargoGoods cargo goods
type CargoGoods struct {
	GoodCount int32   `json:"good_count"` // 是	货物数量
	GoodName  string  `json:"good_name"`  // 是	货品名称
	GoodPrice float32 `json:"good_price"` // 否	货品单价，精确到小数点后两位（如果小数点后位数多于两位，则四舍五入保留两位小数）
	GoodUnit  string  `json:"good_unit"`  // 否	货品单位，最长不超过20个字符
}

// CargoGoodsDetail cargo goods detail
type CargoGoodsDetail struct {
	Goods []*CargoGoods `json:"goods"` // 否 货物详情，最长不超过10240个字符
}

// LogisticsCargo cargo
type LogisticsCargo struct {
	GoodsValue        float32           `json:"goods_value"`         // 是	货物价格，单位为元，精确到小数点后两位（如果小数点后位数多于两位，则四舍五入保留两位小数），范围为(0-5000]
	GoodsHeight       float32           `json:"goods_height"`        // 否	货物高度，单位为CM，精确到小数点后两位（如果小数点后位数多于两位，则四舍五入保留两位小数），范围为(0-45]
	GoodsLength       float32           `json:"goods_length"`        // 否	货物长度，单位为CM，精确到小数点后两位（如果小数点后位数多于两位，则四舍五入保留两位小数），范围为(0-65]
	GoodsWidth        float32           `json:"goods_width"`         // 否	货物宽度，单位为cm，精确到小数点后两位（如果小数点后位数多于两位，则四舍五入保留两位小数），范围为(0-50]
	GoodsWeight       float32           `json:"goods_weight"`        // 是	货物重量，单位为kg，精确到小数点后两位（如果小数点后位数多于两位，则四舍五入保留两位小数），范围为(0-50]
	GoodsDetail       *CargoGoodsDetail `json:"goods_detail"`        // 否	货物详情，最长不超过10240个字符
	GoodsPickupInfo   string            `json:"goods_pickup_info"`   // 否	货物取货信息，用于骑手到店取货，最长不超过100个字符
	GoodsDeliveryInfo string            `json:"goods_delivery_info"` // 否	货物交付信息，最长不超过100个字符
	CargoFirstClass   string            `json:"cargo_first_class"`   // 是	品类一级类目, 详见品类表
	CargoSecondClass  string            `json:"cargo_second_class"`  // 是	品类二级类目
}

// LogisticsOrderInfo order info
type LogisticsOrderInfo struct {
	DeliveryServiceCode  string `json:"delivery_service_code"`  // 	否	配送服务代码 不同配送公司自定义, 顺丰和达达不填
	OrderType            int32  `json:"order_type"`             // 	0	否	订单类型, 0: 即时单 1 预约单，如预约单，需要设置expected_delivery_time或expected_finish_time或expected_pick_time
	ExpectedDeliveryTime int64  `json:"expected_delivery_time"` //	0	否	期望派单时间(达达支持，表示达达系统调度时间, 到那个时间才会有状态更新的回调通知)，unix-timestamp, 比如1586342180
	ExpectedFinishTime   int64  `json:"expected_finish_time"`   //0	否	期望送达时间(美团、顺丰同城急送支持），unix-timestamp, 比如1586342180
	PectedPickTime       int64  `json:"pected_pick_time"`       //	0	否	期望取件时间（闪送、顺丰同城急送支持，闪送需要设置两个小时后的时间，顺丰同城急送只需传expected_finish_time或expected_pick_time其中之一即可，同时都传则以expected_finish_time为准），unix-timestamp, 比如1586342180
	PoiSeq               string `json:"poi_seq"`                //	否	门店订单流水号，建议提供，方便骑手门店取货，最长不超过32个字符
	Note                 string `json:"note"`                   // 否	备注，最长不超过200个字符
	OrderTime            int64  `json:"order_time"`             // 否	用户下单付款时间, 顺丰必填, 比如1555220757
	IsInsured            int32  `json:"is_insured"`             // 0	否	是否保价，0，非保价，1.保价
	DeclaredValue        int32  `json:"declared_value"`         // 否	保价金额，单位为元，精确到分
	Tips                 int32  `json:"tips"`                   // 否	小费，单位为元, 下单一般不加小费
	IsDirectDelivery     int32  `json:"is_direct_delivery"`     //		否	是否选择直拿直送（0：不需要；1：需要。选择直拿直送后，同一时间骑手只能配送此订单至完成，配送费用也相应高一些，闪送必须选1，达达可选0或1，其余配送公司不支持直拿直送）
	CashOnDelivery       int32  `json:"cash_on_delivery"`       // 否	骑手应付金额，单位为元，精确到分
	CashOnPickup         int32  `json:"cash_on_pickup"`         //否	骑手应收金额，单位为元，精确到分
	RiderPickMethod      int32  `json:"rider_pick_method"`      //否	物流流向，1：从门店取件送至用户；2：从用户取件送至门店
	IsFinishCodeNeeded   int32  `json:"is_finish_code_needed"`  //否	收货码（0：不需要；1：需要。收货码的作用是：骑手必须输入收货码才能完成订单妥投）
	IsPickupCodeNeeded   int32  `json:"is_pickup_code_needed"`  // 否	取货码（0：不需要；1：需要。取货码的作用是：骑手必须输入取货码才能从商家取货）
}

// LogisticsShop shop
type LogisticsShop struct {
	WxaPath    string `json:"wxa_path"`
	ImgURL     string `json:"img_url"`
	GoodsName  string `json:"goods_name"`
	GoodsCount string `json:"goods_count"`
	WxaAppID   string `json:"wxa_appid"`
}

// LogisticsOrderRequest order request info
type LogisticsOrderRequest struct {
	DeliveryToken string             `json:"delivery_token"` // 否	预下单接口返回的参数，配送公司可保证在一段时间内运费不变
	ShopID        string             `json:"shopid"`         // 是	商家id，由配送公司分配的appkey
	ShopOrderID   string             `json:"shop_order_id"`  // 是	唯一标识订单的 ID，由商户生成, 不超过128字节
	ShopNO        string             `json:"shop_no"`        // 是	商家门店编号，在配送公司登记，如果只有一个门店，美团闪送必填, 值为店铺id
	DeliverySign  string             `json:"delivery_sign"`  // 是	用配送公司提供的appSecret加密的校验串说明
	DeliveryID    string             `json:"delivery_id"`    // 是	配送公司ID
	Openid        string             `json:"openid"`         // 是	下单用户的openid
	Sender        SenderReceiver     `json:"sender"`         // 是	发件人信息，顺丰同城急送必须填写，美团配送、达达、闪送，若传了shop_no的值可不填该字段
	Receiver      SenderReceiver     `json:"receiver"`       // 是	收件人信息
	Cargo         LogisticsCargo     `json:"cargo"`          // 是	货物信息
	OrderInfo     LogisticsOrderInfo `json:"order_info"`     // 是	订单信息
	Shop          LogisticsShop      `json:"shop"`           // 是	商品信息，会展示到物流通知消息中
	SubBizID      string             `json:"sub_biz_id"`     // 否	子商户id，区分小程序内部多个子商户
}

type DeliveryCommonError struct {
	util.CommonError
	ResultCode int32  `json:"resultcode"` // 运力返回的错误码
	ResultMsg  string `json:"resultmsg"`  // 运力返回的错误描述
}

// LogisticsOrderResponse order response
type LogisticsOrderResponse struct {
	DeliveryCommonError
	Fee              int32  `json:"fee"`               // 实际运费(单位：元)，运费减去优惠券费用
	DeliverFee       int32  `json:"deliverfee"`        //	运费(单位：元)
	CouponFee        int32  `json:"couponfee"`         //	优惠券费用(单位：元)
	Tips             int32  `json:"tips"`              // 小费(单位：元)
	Insurancefee     int32  `json:"insurancefee"`      //	保价费(单位：元)
	Distance         int32  `json:"distance"`          // 配送距离(整数单位：米)
	WayBillID        string `json:"waybill_id"`        // 配送单号
	OrderStatus      int32  `json:"order_status"`      // 配送状态
	FinishCode       int32  `json:"finish_code"`       // 收货码
	PickupCode       int32  `json:"pickup_code"`       // 取货码
	DispatchDuration int64  `json:"dispatch_duration"` //预计骑手接单时间，单位秒，比如5分钟，就填300, 无法预计填0
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
	DeliveryID string `json:"delivery_id" validate:"required"` // 快递公司ID
	WaybillID  string `json:"waybill_id" validate:"required"`  // 运单ID
	OpenID     string `json:"open_id"`                         // 用户openid，当add_source=2时无需填写（不发送物流服务通知）
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
