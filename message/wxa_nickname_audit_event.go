package message

// WXANickNameAuditEvent wxa_nickname_audit event
// see https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Mini_Programs/wxa_nickname_audit.html
type WXANickNameAuditEvent struct {
	Ret      int32  `xml:"ret"`
	NickName string `xml:"nick_name"`
	Reason   string `xml:"reason"`
}
