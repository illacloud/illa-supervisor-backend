package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const CUSTOMIZATION_LANGUAGE_EN_US = "en-US"
const CUSTOMIZATION_LANGUAGE_ZH_CN = "zh-CN"
const CUSTOMIZATION_LANGUAGE_JA_JP = "ja-JP"
const CUSTOMIZATION_LANGUAGE_KO_KR = "ko-KR"
const CUSTOMIZATION_LANGUAGE_CS_CZ = "cs-CZ"
const CUSTOMIZATION_LANGUAGE_DA_DK = "da-DK"
const CUSTOMIZATION_LANGUAGE_DE_DE = "de-DE"
const CUSTOMIZATION_LANGUAGE_EL_GR = "el-GR"
const CUSTOMIZATION_LANGUAGE_ES_ES = "es-ES"
const CUSTOMIZATION_LANGUAGE_FI_FI = "fi-FI"
const CUSTOMIZATION_LANGUAGE_FR_FR = "fr-FR"
const CUSTOMIZATION_LANGUAGE_IT_IT = "it-IT"
const CUSTOMIZATION_LANGUAGE_NL_NL = "nl-NL"
const CUSTOMIZATION_LANGUAGE_NO_NO = "no-NO"
const CUSTOMIZATION_LANGUAGE_PL_PL = "pl-PL"
const CUSTOMIZATION_LANGUAGE_PT_PT = "pt-PT"
const CUSTOMIZATION_LANGUAGE_RU_RU = "ru-RU"
const CUSTOMIZATION_LANGUAGE_RO_RO = "ro-RO"
const CUSTOMIZATION_LANGUAGE_SV_SE = "sv-SE"
const CUSTOMIZATION_LANGUAGE_UK_UA = "uk-UA"

const PENDING_USER_ID = 0
const PENDING_USER_NICKNAME = "pending"
const PENDING_USER_PASSWORDDIGEST = "pending"
const PENDING_USER_AVATAR = ""

type User struct {
	ID             int       `json:"id" gorm:"column:id;type:bigserial;primary_key;index:users_ukey"`
	UID            uuid.UUID `json:"uid" gorm:"column:uid;type:uuid;not null;index:users_ukey"`
	Nickname       string    `json:"nickname" gorm:"column:nickname;type:varchar;size:15"`
	PasswordDigest string    `json:"passworddigest" gorm:"column:password_digest;type:varchar;size:60;not null"`
	Email          string    `json:"email" gorm:"column:email;type:varchar;size:255;not null"`
	Avatar         string    `json:"avatar" gorm:"column:avatar;type:varchar;size:255;not null"`
	SSOConfig      string    `json:"SSOConfig" gorm:"column:sso_config;type:jsonb"`        // for single sign-on data
	Customization  string    `json:"customization" gorm:"column:customization;type:jsonb"` // for user itself customization config, including: Language, IsSubscribed
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp"`
}

type UserForExport struct {
	ID           int       `json:"id"`
	UID          uuid.UUID `json:"uid"`
	TeamMemberID int       `json:"teamMemberID"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	Avatar       string    `json:"avatar"`
	Language     string    `json:"language"`
	IsSubscribed bool      `json:"isSubscribed"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUser() *User {
	return &User{}
}

func (ufe *UserForExport) SetTeamMemberID(teamMemberID int) {
	ufe.TeamMemberID = teamMemberID
}

func (u *User) Export() *UserForExport {
	customization := u.ExportUserCustomization()
	ret := &UserForExport{}
	ret.ID = u.ID
	ret.UID = u.UID
	if u.Nickname == PENDING_USER_NICKNAME { // reset pending user name for frontend display
		ret.Nickname = ""
	} else {
		ret.Nickname = u.Nickname
	}
	ret.Email = u.Email
	ret.Avatar = u.Avatar
	ret.Language = customization.Language
	ret.IsSubscribed = customization.IsSubscribed
	ret.CreatedAt = u.CreatedAt
	ret.UpdatedAt = u.UpdatedAt
	return ret
}

func (u *User) ConstructByJSON(userJSON []byte) error {
	if err := json.Unmarshal(userJSON, u); err != nil {
		return err
	}
	return nil
}

func (u *User) InitUID() {
	u.UID = uuid.New()
}

func (u *User) InitCreatedAt() {
	u.CreatedAt = time.Now().UTC()
}

func (u *User) InitUpdatedAt() {
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) SetID(id int) {
	u.ID = id
}

func (u *User) SetNickname(nickname string) {
	u.Nickname = nickname
	u.InitUpdatedAt()
}

func (u *User) SetAvatar(avatar string) {
	u.Avatar = avatar
	u.InitUpdatedAt()
}

func (u *User) SetPasswordByByte(password []byte) {
	u.PasswordDigest = string(password)
	u.InitUpdatedAt()
}

func (u *User) SetLanguage(language string) {
	userCustomization := u.ExportUserCustomization()
	userCustomization.SetLanguage(language)
	u.SetUserCustomization(userCustomization)
	u.InitUpdatedAt()
}

func (u *User) SetIsTutorialViewed(isTutorialViewed bool) {
	userCustomization := u.ExportUserCustomization()
	userCustomization.SetIsTutorialViewed(isTutorialViewed)
	u.SetUserCustomization(userCustomization)
	u.InitUpdatedAt()
}

func (u *User) SetUserCustomization(userCustomization *UserCustomization) {
	u.Customization, _ = userCustomization.Export()
}

func (u *User) UpdateByUpdateUserAvatarRequest(req *UpdateAvatarRequest) {
	u.Avatar = req.Avatar
	u.InitUpdatedAt()
}

func (u *User) GetUIDInString() string {
	return u.UID.String()
}

func (u *User) ExportLanguage() string {
	userCustomization := u.ExportUserCustomization()
	return userCustomization.Language
}

func (u *User) ExportID() int {
	return u.ID
}

func (u *User) ExportEmail() string {
	return u.Email
}

func (u *User) ExportUserCustomization() *UserCustomization {
	userCustomization := NewUserCustomization()
	json.Unmarshal([]byte(u.Customization), &userCustomization)
	return userCustomization
}

type UserCustomization struct {
	Language         string
	IsSubscribed     bool
	IsTutorialViewed bool
}

func NewUserCustomization() *UserCustomization {
	return &UserCustomization{
		Language:         CUSTOMIZATION_LANGUAGE_EN_US,
		IsSubscribed:     false,
		IsTutorialViewed: false,
	}
}

func (c *UserCustomization) Export() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *UserCustomization) SetLanguage(language string) {
	c.Language = language
}

func (c *UserCustomization) SetIsSubscribed(isSubscribed bool) {
	c.IsSubscribed = isSubscribed
}

func (c *UserCustomization) SetIsTutorialViewed(isTutorialViewed bool) {
	c.IsTutorialViewed = isTutorialViewed
}

func BuildLookUpTableForUserExport(users []*User) map[int]*UserForExport {
	usersNum := len(users)
	lt := make(map[int]*UserForExport, usersNum)
	for _, user := range users {
		lt[user.ID] = user.Export()
	}
	return lt
}

type UserSSOConfig struct {
	Github string
}

func NewUserSSOConfig() *UserSSOConfig {
	return &UserSSOConfig{}
}

func (u *UserSSOConfig) Export() string {
	r, _ := json.Marshal(u)
	return string(r)
}
