package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"encoding/base64"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisior-backend/src/utils/idconvertor"
)

// invite record category
const CATEGORY_INVITE_BY_LINK = 1
const CATEGORY_INVITE_BY_EMAIL = 2

// invite record status
const INVITE_RECORD_STATUS_OK = 1

const INVITE_URI_TEMPLATE = "?inviteToken=%s"
const INVITE_EMAIL_TEMPLATE = "&teamIdentifier=%s&email=%s"
const INVITE_SHARE_APP_TEMPLATE = "&teamIdentifier=%s&appID=%s"
const INVITE_EMAIL_AND_SHARE_APP_TEMPLATE = "&teamIdentifier=%s&email=%s&appID=%s"

type Invite struct {
	ID             int       `json:"id" gorm:"column:id;type:bigserial;primary_key;index:invite_ukey"`
	UID            uuid.UUID `json:"uid" gorm:"column:uid;type:uuid;not null;index:invite_uid"`
	Category       int       `json:"category" gorm:"column:category;type:smallint;not null"`
	TeamID         int       `json:"teamID" gorm:"column:team_id;type:bigserial;index:users_ukey"`
	TeamMemberID   int       `json:"teamMemberID" gorm:"column:team_member_id;type:bigserial"`
	TeamIdentifier string    `json:"teamIdentifier" sql:"-" gorm:"-"`
	AppID          int       `json:"appID" sql:"-" gorm:"-"`
	Email          string    `json:"email" gorm:"column:email;type:varchar;size:255;;index:invite_email"`
	EmailStatus    bool      `json:"emailStatus" gorm:"column:email_status;type:bool"`
	UserRole       int       `json:"userRole" gorm:"column:user_role;type:smallint;index:invite_user_role"`
	Status         int       `json:"status" gorm:"column:status;type:smallint;index:users_ukey"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp"`
}

type InviteForExport struct {
	ID             int       `json:"id"`
	UID            uuid.UUID `json:"uid"`
	Category       int       `json:"category"`
	TeamID         int       `json:"teamID"`
	TeamMemberID   int       `json:"teamMemberID"`
	TeamIdentifier string    `json:"teamIdentifier"`
	AppID          int       `json:"appID"`
	Email          string    `json:"email"`
	EmailStatus    bool      `json:"emailStatus"`
	UserRole       int       `json:"userRole"`
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func NewInvite() *Invite {
	return &Invite{}
}

func (u *Invite) ConstructByJSON(InviteJSON []byte) error {
	if err := json.Unmarshal(InviteJSON, u); err != nil {
		return err
	}
	return nil
}

func (u *Invite) InitUID() {
	u.UID = uuid.New()
}

func (u *Invite) InitCreatedAt() {
	u.CreatedAt = time.Now().UTC()
}

func (u *Invite) InitUpdatedAt() {
	u.UpdatedAt = time.Now().UTC()
}

func (u *Invite) SetEmailStatusFailed() {
	u.EmailStatus = false
}

func (u *Invite) SetEmailStatusSuccess() {
	u.EmailStatus = true
}

func (u *Invite) SetUserRole(role int) {
	u.UserRole = role
}

func (u *Invite) SetTeamMemberID(teamMemberID int) {
	u.TeamMemberID = teamMemberID
}

func (u *Invite) SetTeamIdentifier(teamIdentifier string) {
	u.TeamIdentifier = teamIdentifier
}

func (u *Invite) SetAppID(appID int) {
	fmt.Printf("appID: %v, u.AppID: %v\n", appID, u.AppID)
	u.AppID = appID
}

func (u *Invite) ExportID() int {
	return u.ID
}

func (u *Invite) Export() *InviteForExport {
	return &InviteForExport{
		ID:             u.ID,
		UID:            u.UID,
		Category:       u.Category,
		TeamID:         u.TeamID,
		TeamMemberID:   u.TeamMemberID,
		TeamIdentifier: u.TeamIdentifier,
		AppID:          u.AppID,
		Email:          u.Email,
		EmailStatus:    u.EmailStatus,
		UserRole:       u.UserRole,
		Status:         u.Status,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func (u *Invite) ExportUID() uuid.UUID {
	return u.UID
}

func (u *Invite) ExportInviteLink() string {
	template := ""
	template = Config.GetServeHTTPAddress() + INVITE_URI_TEMPLATE
	return fmt.Sprintf(template, base64.StdEncoding.EncodeToString([]byte(u.UID.String())))
}

func (u *Invite) ExportShareAppLink() string {
	return u.ExportInviteLink() + fmt.Sprintf(INVITE_SHARE_APP_TEMPLATE, u.TeamIdentifier, idconvertor.ConvertIntToString(u.AppID))
}

func (u *Invite) ExportInviteLinkWithEmailAndTeam(team *Team) string {
	return u.ExportInviteLink() + fmt.Sprintf(INVITE_EMAIL_TEMPLATE, team.GetIdentifier(), u.Email)
}

func (u *Invite) ExportShareAppLinkWithEmailANDTeamIdentifier(team *Team) string {
	return u.ExportInviteLink() + fmt.Sprintf(INVITE_EMAIL_AND_SHARE_APP_TEMPLATE, team.GetIdentifier(), u.Email, idconvertor.ConvertIntToString(u.AppID))
}

func (u *Invite) ExportTeamID() int {
	return u.TeamID
}

func (u *Invite) ExportTeamMemberID() int {
	return u.TeamMemberID
}

func (u *Invite) ExportAppIDInString() string {
	return strconv.Itoa(u.AppID)
}

func (u *Invite) ExportEmail() string {
	return u.Email
}

func (u *Invite) ExportUserRole() int {
	return u.UserRole
}

func (u *Invite) ExportEmailStatus() bool {
	return u.EmailStatus
}

func (u *Invite) ImportInviteLink(link string) error {
	uuidString, err := base64.StdEncoding.DecodeString(link)
	if err != nil {
		return err
	}
	var errInParse error
	u.UID, errInParse = uuid.Parse(string(uuidString))
	if errInParse != nil {
		return errInParse
	}
	return nil
}

func (u *Invite) IsAvaliable() bool {
	if u.Status == INVITE_RECORD_STATUS_OK {
		return true
	}
	return false
}

func (u *Invite) IsEmailInviteLink() bool {
	if u.Category == CATEGORY_INVITE_BY_EMAIL {
		return true
	}
	return false
}

func (u *Invite) IsInviteLink() bool {
	if u.Category == CATEGORY_INVITE_BY_LINK {
		return true
	}
	return false
}

func (u *Invite) IsShareAppInvite() bool {
	if u.AppID > 0 {
		return true
	}
	return false
}

func NewInviteLinkByTeamAndUserRole(team *Team, userRole int) *Invite {
	Invite := NewInvite()
	Invite.Category = CATEGORY_INVITE_BY_LINK
	Invite.TeamID = team.ID
	Invite.TeamIdentifier = team.Identifier
	Invite.AppID = 0
	Invite.Email = "" // invite by link, the email was not setted
	Invite.UserRole = userRole
	Invite.Status = INVITE_RECORD_STATUS_OK
	Invite.InitUID()
	Invite.InitCreatedAt()
	Invite.InitUpdatedAt()
	return Invite
}

func NewInviteEmailLinkByTeamAndRequest(team *Team, req *InviteMemberByEmailRequest) *Invite {
	Invite := NewInvite()
	Invite.Category = CATEGORY_INVITE_BY_EMAIL
	Invite.TeamID = team.ID
	Invite.TeamIdentifier = team.Identifier
	Invite.AppID = 0
	Invite.Email = req.ExportEmail()
	Invite.UserRole = req.ExportUserRole()
	Invite.Status = INVITE_RECORD_STATUS_OK
	Invite.InitUID()
	Invite.InitCreatedAt()
	Invite.InitUpdatedAt()
	return Invite
}

func NewInviteEmailLinkByTeamIDAndTeamMemberIDAndRequest(teamID int, teamMemberID int, req *InviteMemberByEmailRequest) *Invite {
	Invite := NewInvite()
	Invite.Category = CATEGORY_INVITE_BY_EMAIL
	Invite.TeamID = teamID
	Invite.AppID = 0
	Invite.TeamMemberID = teamMemberID
	Invite.Email = req.ExportEmail()
	Invite.UserRole = req.ExportUserRole()
	Invite.Status = INVITE_RECORD_STATUS_OK
	Invite.InitUID()
	Invite.InitCreatedAt()
	Invite.InitUpdatedAt()
	return Invite
}

func NewInviteByCreateInviteRequest(req *CreateInviteRequest) *Invite {
	Invite := NewInvite()
	Invite.TeamID = req.TeamID
	Invite.AppID = 0
	Invite.Email = req.Email
	Invite.UserRole = req.UserRole
	Invite.Status = INVITE_RECORD_STATUS_OK
	Invite.InitUID()
	Invite.InitCreatedAt()
	Invite.InitUpdatedAt()
	return Invite
}

func NewInviteByUpdateInviteRequest(req *UpdateInviteRequest) *Invite {
	Invite := NewInvite()
	Invite.ID = req.ID
	Invite.UID = req.UID
	Invite.TeamID = req.TeamID
	Invite.AppID = 0
	Invite.Email = req.Email
	Invite.UserRole = req.UserRole
	Invite.Status = req.Status
	Invite.InitUpdatedAt()
	return Invite
}

type EmailInviteMessage struct {
	UserName      string `json:"userName"`
	TeamName      string `json:"teamName"`
	TeamIcon      string `json:"teamIcon"`
	Email         string `json:"email"`
	JoinLink      string `json:"joinLink"`
	Language      string `json:"language"`
	ValidateToken string `json:"validateToken"` // token for query authorize, base64.Encoding(md5(...param + ROTOR_TOKEN))
}

func NewEmailInviteMessage(invite *Invite, team *Team, user *User) *EmailInviteMessage {
	return &EmailInviteMessage{
		UserName: user.Nickname,
		TeamName: team.Name,
		TeamIcon: team.Icon,
		Email:    invite.Email,
		JoinLink: invite.ExportInviteLinkWithEmailAndTeam(team),
		Language: user.ExportLanguage(),
	}
}

func (m *EmailInviteMessage) Export() map[string]string {
	payload := map[string]string{
		"userName":      m.UserName,
		"teamName":      m.TeamName,
		"teamIcon":      m.TeamIcon,
		"email":         m.Email,
		"joinLink":      m.JoinLink,
		"language":      m.Language,
		"validateToken": m.ValidateToken,
	}
	return payload
}

func (m *EmailInviteMessage) SetValidateToken(token string) {
	m.ValidateToken = token
}

type EmailShareAppMessage struct {
	UserName      string `json:"userName"`
	TeamName      string `json:"teamName"`
	TeamIcon      string `json:"teamIcon"`
	Email         string `json:"email"`
	AppLink       string `json:"appLink"`
	Language      string `json:"language"`
	ValidateToken string `json:"validateToken"` // token for query authorize, base64.Encoding(md5(...param + ROTOR_TOKEN))
}

func NewEmailShareAppMessage(invite *Invite, team *Team, user *User) *EmailShareAppMessage {
	return &EmailShareAppMessage{
		UserName: user.Nickname,
		TeamName: team.Name,
		TeamIcon: team.Icon,
		Email:    invite.Email,
		AppLink:  invite.ExportShareAppLinkWithEmailANDTeamIdentifier(team),
		Language: user.ExportLanguage(),
	}
}

func (m *EmailShareAppMessage) Export() map[string]string {
	payload := map[string]string{
		"userName":      m.UserName,
		"teamName":      m.TeamName,
		"teamIcon":      m.TeamIcon,
		"email":         m.Email,
		"appLink":       m.AppLink,
		"language":      m.Language,
		"validateToken": m.ValidateToken,
	}
	return payload
}

func (m *EmailShareAppMessage) SetValidateToken(token string) {
	m.ValidateToken = token
}

func BuildLookUpTableForInvitesExport(invites []*Invite) map[int]*InviteForExport {
	invitesNum := len(invites)
	lt := make(map[int]*InviteForExport, invitesNum)
	for _, invite := range invites {
		lt[invite.TeamMemberID] = invite.Export()
	}
	return lt
}
