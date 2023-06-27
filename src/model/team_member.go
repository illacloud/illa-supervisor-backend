package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// User Role ID in Team
// @note: this will extend as role system later.
const (
	USER_ROLE_ANONYMOUS = -1
	USER_ROLE_OWNER     = 1
	USER_ROLE_ADMIN     = 2
	USER_ROLE_EDITOR    = 3
	USER_ROLE_VIEWER    = 4
)

const TEAM_MEMBER_STATUS_OK = 1
const TEAM_MEMBER_STATUS_PENDING = 2

type TeamMember struct {
	ID         int       `json:"id" gorm:"column:id;type:bigserial;primary_key;index:team_members_ukey"`
	TeamID     int       `json:"team_id" gorm:"column:team_id;type:bigserial;index:team_members_team_and_user_id"`
	UserID     int       `json:"user_id" gorm:"column:user_id;type:bigserial;index:team_members_team_and_user_id"`
	UserRole   int       `json:"user_role" gorm:"column:user_role;type:smallint"`
	Permission string    `json:"permission" gorm:"column:permission;type:jsonb"` // for user permission config
	Status     int       `json:"status" gorm:"column:status;type:smallint"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp"`
}

type TeamMemberWithUserInfoForExport struct {
	ID           int                   `json:"userID"`
	UID          uuid.UUID             `json:"uid"`
	TeamMemberID int                   `json:"teamMemberID"`
	Nickname     string                `json:"nickname"`
	Email        string                `json:"email"`
	Avatar       string                `json:"avatar"`
	Language     string                `json:"language"`
	IsSubscribed bool                  `json:"isSubscribed"`
	UserRole     int                   `json:"userRole"`
	Permission   *TeamMemberPermission `json:"permission"` // for user permission config
	UserStatus   int                   `json:"userStatus"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

type TeamMemberForExport struct {
	ID         int                   `json:"id"`
	TeamID     int                   `json:"teamID"`
	UserID     int                   `json:"userID"`
	UserRole   int                   `json:"userRole"`
	Permission *TeamMemberPermission `json:"permission"` // for user permission config
	Status     int                   `json:"status"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

func NewTeamMember() *TeamMember {
	return &TeamMember{}
}

func (u *TeamMember) ConstructByJSON(TeamMemberJSON []byte) error {
	if err := json.Unmarshal(TeamMemberJSON, u); err != nil {
		return err
	}
	return nil
}

func (u *TeamMember) SetID(id int) {
	u.ID = id
}

func (u *TeamMember) SetUserID(userID int) {
	u.UserID = userID
}

func (u *TeamMember) SetUserRole(userRole int) {
	u.UserRole = userRole
}

func (u *TeamMember) InitCreatedAt() {
	u.CreatedAt = time.Now().UTC()
}

func (u *TeamMember) InitUpdatedAt() {
	u.UpdatedAt = time.Now().UTC()
}

func (u *TeamMember) ActiveUser() {
	u.Status = TEAM_MEMBER_STATUS_OK
}

func (u *TeamMember) ExportUserRole() int {
	return u.UserRole
}

func (u *TeamMember) ExportUserID() int {
	return u.UserID
}

func (u *TeamMember) ExportWithUserInfo(userForExport *UserForExport) *TeamMemberWithUserInfoForExport {
	return &TeamMemberWithUserInfoForExport{
		ID:           userForExport.ID,
		UID:          userForExport.UID,
		TeamMemberID: userForExport.TeamMemberID,
		Nickname:     userForExport.Nickname,
		Email:        userForExport.Email,
		Avatar:       userForExport.Avatar,
		Language:     userForExport.Language,
		IsSubscribed: userForExport.IsSubscribed,
		UserRole:     u.UserRole,
		Permission:   u.ExportPermission(),
		UserStatus:   u.Status,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (u *TeamMember) ExportWithPendingUserInfo(userForExport *UserForExport) *TeamMemberWithUserInfoForExport {
	return &TeamMemberWithUserInfoForExport{
		ID:           userForExport.ID,
		UID:          userForExport.UID,
		TeamMemberID: userForExport.TeamMemberID,
		Nickname:     userForExport.Nickname,
		Email:        userForExport.Email,
		Avatar:       userForExport.Avatar,
		Language:     userForExport.Language,
		IsSubscribed: userForExport.IsSubscribed,
		UserRole:     u.UserRole,
		Permission:   u.ExportPermission(),
		UserStatus:   u.Status,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (u *TeamMember) Export() *TeamMemberForExport {
	return &TeamMemberForExport{
		ID:         u.ID,
		TeamID:     u.TeamID,
		UserID:     u.UserID,
		UserRole:   u.UserRole,
		Permission: u.ExportPermission(),
		Status:     u.Status,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (u *TeamMember) ExportID() int {
	return u.ID
}

func (u *TeamMember) UpdateByUpdateTeamMemberRoleRequest(req *UpdateTeamMemberRoleRequest) {
	u.UserRole = req.UserRole
	u.InitUpdatedAt()
}

func (u *TeamMember) UpdateTeamMemberRole(userRole int) {
	u.UserRole = userRole
	u.InitUpdatedAt()
}

func (u *TeamMember) IsOwner() bool {
	if u.UserRole == USER_ROLE_OWNER {
		return true
	}
	return false
}

func (u *TeamMember) IsAdmin() bool {
	if u.UserRole == USER_ROLE_ADMIN {
		return true
	}
	return false
}

func (u *TeamMember) IsEditor() bool {
	if u.UserRole == USER_ROLE_EDITOR {
		return true
	}
	return false
}

func (u *TeamMember) IsViewer() bool {
	if u.UserRole == USER_ROLE_VIEWER {
		return true
	}
	return false
}

func (u *TeamMember) ExportPermission() *TeamMemberPermission {
	tmp := NewTeamMemberPermission()
	json.Unmarshal([]byte(u.Permission), tmp)
	return tmp
}

func (u *TeamMember) IsStatusPending() bool {
	if u.Status == TEAM_MEMBER_STATUS_PENDING {
		return true
	}
	return false
}

func (u *TeamMember) IsStatusOK() bool {
	if u.Status == TEAM_MEMBER_STATUS_OK {
		return true
	}
	return false
}

func NewTeamMemberByCreateTeamMemberRequest(req *CreateTeamMemberRequest) *TeamMember {
	TeamMember := NewTeamMember()
	TeamMember.TeamID = req.TeamID
	TeamMember.UserID = req.UserID
	TeamMember.UserRole = req.UserRole
	TeamMember.Permission = req.Permission
	TeamMember.InitCreatedAt()
	TeamMember.InitUpdatedAt()
	return TeamMember
}

func NewTeamMemberByUpdateTeamMemberRequest(req *UpdateTeamMemberRequest) *TeamMember {
	TeamMember := NewTeamMember()
	TeamMember.ID = req.ID
	TeamMember.TeamID = req.TeamID
	TeamMember.UserID = req.UserID
	TeamMember.UserRole = req.UserRole
	TeamMember.Permission = req.Permission
	TeamMember.InitUpdatedAt()
	return TeamMember
}

func NewTeamMemberByInviteAndUserID(invite *Invite, userID int) *TeamMember {
	tmp := NewTeamMemberPermission()
	TeamMember := NewTeamMember()
	TeamMember.TeamID = invite.ExportTeamID()
	TeamMember.UserID = userID
	TeamMember.UserRole = invite.ExportUserRole()
	TeamMember.Permission = tmp.ExportForTeam()
	TeamMember.Status = TEAM_MEMBER_STATUS_OK
	TeamMember.InitCreatedAt()
	TeamMember.InitUpdatedAt()
	return TeamMember
}

func NewPendingTeamMemberByInvite(invite *Invite) *TeamMember {
	tmp := NewTeamMemberPermission()
	TeamMember := NewTeamMember()
	TeamMember.TeamID = invite.ExportTeamID()
	TeamMember.UserID = PENDING_USER_ID
	TeamMember.UserRole = invite.ExportUserRole()
	TeamMember.Permission = tmp.ExportForTeam()
	TeamMember.Status = TEAM_MEMBER_STATUS_PENDING
	TeamMember.InitCreatedAt()
	TeamMember.InitUpdatedAt()
	return TeamMember
}

func NewTeamMemberByTeamIDAndUserIDAndUserRole(teamID, userID, userRole int) *TeamMember {
	tmp := NewTeamMemberPermission()
	TeamMember := NewTeamMember()
	TeamMember.TeamID = teamID
	TeamMember.UserID = userID
	TeamMember.UserRole = userRole
	TeamMember.Permission = tmp.ExportForTeam()
	TeamMember.Status = TEAM_MEMBER_STATUS_OK
	TeamMember.InitCreatedAt()
	TeamMember.InitUpdatedAt()
	return TeamMember
}

func NewEditorTeamMemberByUserID(userID int) *TeamMember {
	tmp := NewTeamMemberPermission()
	TeamMember := NewTeamMember()
	TeamMember.TeamID = TEAM_DEFAULT_ID
	TeamMember.UserID = userID
	TeamMember.UserRole = USER_ROLE_EDITOR
	TeamMember.Permission = tmp.ExportForTeam()
	TeamMember.Status = TEAM_MEMBER_STATUS_OK
	TeamMember.InitCreatedAt()
	TeamMember.InitUpdatedAt()
	return TeamMember
}

func PickUpTeamIDsInTeamMembers(teamMembers []*TeamMember) []int {
	idlen := len(teamMembers)
	ids := make([]int, idlen, idlen)
	for serial, teamMember := range teamMembers {
		ids[serial] = teamMember.TeamID
	}
	return ids
}

func PickUpTeamMemberIDsInTeamMembers(teamMembers []*TeamMember) []int {
	idlen := len(teamMembers)
	ids := make([]int, idlen, idlen)
	for serial, teamMember := range teamMembers {
		ids[serial] = teamMember.ID
	}
	return ids
}

func PickUpUserIDsInUserMembers(teamMembers []*TeamMember) []int {
	idlen := len(teamMembers)
	ids := make([]int, idlen, idlen)
	for serial, teamMember := range teamMembers {
		if teamMember.UserID == PENDING_USER_ID { // skip pending use, they where not in the database.
			continue
		}
		ids[serial] = teamMember.UserID
	}
	return ids
}

type TeamMemberPermission struct {
	Config int `json:"config"`
}

func NewTeamMemberPermission() *TeamMemberPermission {
	return &TeamMemberPermission{}
}
func (tmp *TeamMemberPermission) ExportForTeam() string {
	r, _ := json.Marshal(tmp)
	return string(r)
}

func BuildTeamIDLookUpTableForTeamMemberExport(teamMembers []*TeamMember) map[int]*TeamMemberForExport {
	teamMembersNum := len(teamMembers)
	lt := make(map[int]*TeamMemberForExport, teamMembersNum)
	for _, teamMember := range teamMembers {
		lt[teamMember.TeamID] = teamMember.Export()
	}
	return lt
}
