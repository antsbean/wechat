package material

import (
	"encoding/json"
	"fmt"

	"github.com/antsbean/wechat/util"
)

//MediaType 媒体文件类型
type MediaType string

const (
	//MediaTypeImage 媒体文件:图片
	MediaTypeImage MediaType = "image"
	//MediaTypeVoice 媒体文件:声音
	MediaTypeVoice = "voice"
	//MediaTypeVideo 媒体文件:视频
	MediaTypeVideo = "video"
	//MediaTypeThumb 媒体文件:缩略图
	MediaTypeThumb = "thumb"
)

const (
	mediaUploadURL      = "https://api.weixin.qq.com/cgi-bin/media/upload"
	mediaUploadImageURL = "https://api.weixin.qq.com/cgi-bin/media/uploadimg"
	mediaGetURL         = "https://api.weixin.qq.com/cgi-bin/media/get"
)

//Media 临时素材上传返回信息
type Media struct {
	util.CommonError

	Type         MediaType `json:"type"`
	MediaID      string    `json:"media_id"`
	ThumbMediaID string    `json:"thumb_media_id"`
	CreatedAt    int64     `json:"created_at"`
}

// MediaUploadWithAK 临时素材上传
func (material *Material) MediaUploadWithAK(ak string, mediaType MediaType, filename string) (media Media, err error) {
	uri := fmt.Sprintf("%s?access_token=%s&type=%s", mediaUploadURL, ak, mediaType)
	var response []byte
	response, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &media)
	if err != nil {
		return
	}
	if media.ErrCode != 0 {
		err = fmt.Errorf("MediaUpload error : errcode=%v , errmsg=%v", media.ErrCode, media.ErrMsg)
		return
	}
	return
}

//MediaUpload 临时素材上传
func (material *Material) MediaUpload(mediaType MediaType, filename string) (media Media, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}
	return material.MediaUploadWithAK(accessToken, mediaType, filename)
}

//GetMediaURLWithAK 返回临时素材的下载地址供用户自己处理
//NOTICE: URL 不可公开，因为含access_token 需要立即另存文件
func (material *Material) GetMediaURLWithAK(ak, mediaID string) (mediaURL string, err error) {
	mediaURL = fmt.Sprintf("%s?access_token=%s&media_id=%s", mediaGetURL, ak, mediaID)
	return
}

//GetMediaURL 返回临时素材的下载地址供用户自己处理
//NOTICE: URL 不可公开，因为含access_token 需要立即另存文件
func (material *Material) GetMediaURL(mediaID string) (mediaURL string, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}
	return material.GetMediaURLWithAK(accessToken, mediaID)
}

//resMediaImage 图片上传返回结果
type resMediaImage struct {
	util.CommonError

	URL string `json:"url"`
}

//ImageUpload 图片上传
func (material *Material) ImageUpload(filename string) (url string, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", mediaUploadImageURL, accessToken)
	var response []byte
	response, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	var image resMediaImage
	err = json.Unmarshal(response, &image)
	if err != nil {
		return
	}
	if image.ErrCode != 0 {
		err = fmt.Errorf("UploadImage error : errcode=%v , errmsg=%v", image.ErrCode, image.ErrMsg)
		return
	}
	url = image.URL
	return

}
