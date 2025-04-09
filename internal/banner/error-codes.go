package banner

import "yadwy-backend/internal/common"

const (
	FailedToGetAllBanners common.ErrorCode = "failed_to_get_all_banners"
	FailedToCreateBanner  common.ErrorCode = "failed_to_create_banner"
	FailedToUploadImage   common.ErrorCode = "failed_to_upload_image"
	FailedToParseImage    common.ErrorCode = "failed_to_parse_image"
)
