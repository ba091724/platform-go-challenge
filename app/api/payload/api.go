package payload

import "errors"

const (
	CHARACTERISTIC_AGE_GROUP                = "age group"
	CHARACTERISTIC_BIRTH_COUNTRY            = "birth country"
	CHARACTERISTIC_GENDER                   = "gender"
	CHARACTERISTIC_PURCHASES_LAST_MONTH     = "purchases made last month"
	CHARACTERISTIC_SOCIAL_MEDIA_DAILY_HOURS = "hours spent on social media per day" // todo: should be treated like age groups
	GENDER_FEMALE                           = "female"
	GENDER_MALE                             = "male"
)

type AssetUpdateRequest struct {
	Description string `json:"description" binding:"required"`
}

type UserFavoriteRequest struct {
	AssetID int `json:"assetId" binding:"required"`
}

type UserDto struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type AssetDto struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type ChartAssetDto struct {
	Title      string `json:"title"`
	AxesTitles string `json:"axesTitles"`
	PlotData   string `json:"plotData"`
}

type InsightAssetDto struct {
	Text string `json:"text"`
}

type AudienceAssetDto struct {
	Characteristics []string `json:"characteristics"`
}

type AssetDetailsDto struct {
	Asset           AssetDto          `json:"asset"`
	ChartDetails    *ChartAssetDto    `json:"chart,omitempty"`
	InsightDetails  *InsightAssetDto  `json:"insight,omitempty"`
	AudienceDetails *AudienceAssetDto `json:"audience,omitempty"`
}

type UserFavoriteDto struct {
	ID      int             `json:"id"`
	Details AssetDetailsDto `json:"details"`
}

/* error */
type ErrHttp struct {
	Code    int
	Message string
	Err     error
}

func (e ErrHttp) Error() string {
	return e.Message
}

func (e ErrHttp) Unwrap() error {
	return e.Err
}

func (e *ErrHttp) WithMessage(message string) ErrHttp {
	e.Message = message
	return *e
}

var ErrOK = ErrHttp{Code: 200, Err: nil}
var ErrBadRequest = ErrHttp{Code: 400, Err: errors.New("bad request")}
var ErrNotFound = ErrHttp{Code: 404, Err: errors.New("not found")}
var ErrConflict = ErrHttp{Code: 409, Err: errors.New("conflict")}
