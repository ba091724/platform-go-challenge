package models

type Asset struct {
	ID          int
	Description string
	Type        int
}

type Chart struct {
	Asset
	Title      string
	AxesTitles string
	PlotData   string
}

type Insight struct {
	Asset
	Text string
}

type Audience struct {
	Asset
	ID int
}

type AudienceCharacteristic struct {
	ID                  int
	AudienceID          int
	CharacteristicID    int
	CharacteristicValue int
}

type User struct {
	ID   int
	Name string
}

type UserFavorite struct {
	ID int
	// UserID  int
	// AssetID int
	User  *User
	Asset *Asset
}
