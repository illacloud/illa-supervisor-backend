package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

const TEAM_DEFAULT_ID = 0

const TEAM_FIELD_NAME = "name"
const TEAM_FIELD_IDENTIFIER = "identifier"
const TEAM_FIELD_ICON = "icon"

const TEAM_P_FIELD_ALLOW_EDITOR_INVITE = "allowEditorInvite"
const TEAM_P_FIELD_ALLOW_VIEWER_INVITE = "allowViewerInvite"
const TEAM_P_FIELD_ALLOW_EDITOR_MANAGE_TEAM_MEMBER = "allowEditorManageTeamMember"
const TEAM_P_FIELD_ALLOW_VIEWER_MANAGE_TEAM_MEMBER = "allowViewerManageTeamMember"
const TEAM_P_FIELD_INVITE_LINK_ENABLED = "inviteLinkEnabled"
const TEAM_P_FIELD_BLOCK_REGISTER = "blockRegister"

type Team struct {
	ID         int       `json:"id" gorm:"column:id;type:bigserial;primary_key;index:teams_ukey"`
	UID        uuid.UUID `json:"uid" gorm:"column:uid;type:uuid;not null;index:teams_ukey"`
	Name       string    `json:"name" gorm:"column:name;type:varchar;size:255;not null"`
	Identifier string    `json:"identifier" gorm:"column:identifier;type:varchar;size:255;not null"`
	Icon       string    `json:"icon" gorm:"column:icon;type:varchar;size:255;not null"`
	Permission string    `json:"permission" gorm:"column:permission;type:jsonb"` // for team permission config
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp"`
}

type TeamForExport struct {
	ID         int             `json:"id"`
	UID        uuid.UUID       `json:"uid"`
	Name       string          `json:"name"`
	Identifier string          `json:"identifier"`
	Icon       string          `json:"icon"`
	Permission *TeamPermission `json:"permission"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type TeamsForExport []*TeamForExport

func NewTeamsForExport(teams []*Team) (r TeamsForExport) {
	for _, team := range teams {
		r = append(r, team.Export())
	}
	return r
}

func (t *Team) Export() *TeamForExport {
	return &TeamForExport{
		ID:         t.ID,
		UID:        t.UID,
		Name:       t.Name,
		Identifier: t.Identifier,
		Icon:       t.Icon,
		Permission: t.ExportTeamPermission(),
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func NewTeam() *Team {
	return &Team{}
}

func (u *Team) ConstructByJSON(TeamJSON []byte) error {
	if err := json.Unmarshal(TeamJSON, u); err != nil {
		return err
	}
	return nil
}

func (u *Team) InitUID() {
	u.UID = uuid.New()
}

func (u *Team) InitCreatedAt() {
	u.CreatedAt = time.Now().UTC()
}

func (u *Team) InitUpdatedAt() {
	u.UpdatedAt = time.Now().UTC()
}

func (u *Team) GetUID() uuid.UUID {
	return u.UID
}

func (u *Team) GetUIDInString() string {
	return u.UID.String()
}

func (u *Team) GetIdentifier() string {
	return u.Identifier
}

func (u *Team) UpdateByUpdateTeamConfigRawRequest(rawReq map[string]interface{}) error {
	var assertPass bool
	for key, value := range rawReq {
		switch key {
		case TEAM_FIELD_NAME:
			u.Name, assertPass = value.(string)
			if !assertPass {
				return errors.New("update team config failed due to assert failed.")
			}
		case TEAM_FIELD_IDENTIFIER:
			u.Identifier, assertPass = value.(string)
			if !assertPass {
				return errors.New("update team config failed due to assert failed.")
			}
		case TEAM_FIELD_ICON:
			u.Icon, assertPass = value.(string)
			if !assertPass {
				return errors.New("update team config failed due to assert failed.")
			}
		default:
		}
	}
	u.InitUpdatedAt()
	return nil
}

func (u *Team) UpdateByUpdateTeamPermissionRawRequest(rawReq map[string]interface{}) error {
	tp := u.ExportTeamPermission()
	var assertPass bool
	for key, value := range rawReq {
		switch key {
		case TEAM_P_FIELD_ALLOW_EDITOR_INVITE:
			tp.AllowEditorInvite, assertPass = value.(bool)
			if !assertPass {
				return errors.New("update team permission failed due to assert failed.")
			}
		case TEAM_P_FIELD_ALLOW_VIEWER_INVITE:
			tp.AllowViewerInvite, assertPass = value.(bool)
			if !assertPass {
				return errors.New("update team permission failed due to assert failed.")
			}
		case TEAM_P_FIELD_ALLOW_EDITOR_MANAGE_TEAM_MEMBER:
			tp.AllowEditorManageTeamMember, assertPass = value.(bool)
			if !assertPass {
				return errors.New("update team permission failed due to assert failed.")
			}
		case TEAM_P_FIELD_ALLOW_VIEWER_MANAGE_TEAM_MEMBER:
			tp.AllowViewerManageTeamMember, assertPass = value.(bool)
			if !assertPass {
				return errors.New("update team permission failed due to assert failed.")
			}
		case TEAM_P_FIELD_INVITE_LINK_ENABLED:
			tp.InviteLinkEnabled, assertPass = value.(bool)
			if !assertPass {
				return errors.New("update team permission failed due to assert failed.")
			}
		case TEAM_P_FIELD_BLOCK_REGISTER:
			tp.BlockRegister, assertPass = value.(bool)
			if !assertPass {
				return errors.New("update team permission failed due to assert failed.")
			}
		default:
		}
	}
	u.Permission = tp.ExportForTeam()
	u.InitUpdatedAt()
	return nil
}

func (u *Team) SetTeamPermission(tp *TeamPermission) {
	u.Permission = tp.ExportForTeam()
	u.InitUpdatedAt()
}

func (u *Team) ConfigInviteLinkByRequest(req *ConfigInviteLinkRequest) {
	tp := u.ExportTeamPermission()
	tp.InviteLinkEnabled = req.InviteLinkEnabled
	u.Permission = tp.ExportForTeam()
	u.InitUpdatedAt()
}

func (u *Team) ExportTeamPermission() *TeamPermission {
	tp := &TeamPermission{}
	json.Unmarshal([]byte(u.Permission), tp)
	return tp
}

func (u *Team) ExportID() int {
	return u.ID
}

func (u *Team) DoesEditorOrViewerCanInviteMember() bool {
	tp := u.ExportTeamPermission()
	if tp.AllowEditorInvite && tp.AllowViewerInvite {
		return true
	}
	return false
}

type TeamPermission struct {
	AllowEditorInvite           bool `json:"allowEditorInvite"`
	AllowViewerInvite           bool `json:"allowViewerInvite"`
	AllowEditorManageTeamMember bool `json:"allowEditorManageTeamMember"`
	AllowViewerManageTeamMember bool `json:"allowViewerManageTeamMember"`
	InviteLinkEnabled           bool `json:"inviteLinkEnabled"`
	BlockRegister               bool `json:"blockRegister"`
}

func NewTeamPermission() *TeamPermission {
	return &TeamPermission{ // switch opend in team created
		AllowEditorInvite:           true,
		AllowViewerInvite:           true,
		AllowEditorManageTeamMember: true,
		AllowViewerManageTeamMember: true,
		InviteLinkEnabled:           true,
		BlockRegister:               false,
	}
}

func (tp *TeamPermission) ExportForTeam() string {
	r, _ := json.Marshal(tp)
	return string(r)
}

func (tp *TeamPermission) ImportFromTeam(team *Team) {
	ttp := team.ExportTeamPermission()
	tp.AllowEditorInvite = ttp.AllowEditorInvite
	tp.AllowViewerInvite = ttp.AllowViewerInvite
	tp.AllowEditorManageTeamMember = ttp.AllowEditorManageTeamMember
	tp.AllowViewerManageTeamMember = ttp.AllowViewerManageTeamMember
	tp.InviteLinkEnabled = ttp.InviteLinkEnabled
	tp.BlockRegister = ttp.BlockRegister
}

func (tp *TeamPermission) EnableInviteLink() {
	tp.InviteLinkEnabled = true
}

func (tp *TeamPermission) DisableInviteLink() {
	tp.InviteLinkEnabled = true
}

func (tp *TeamPermission) DoesInviteLinkEnabled() bool {
	return tp.InviteLinkEnabled
}

func (tp *TeamPermission) DoesEditorCanManageTeamMember() bool {
	return tp.AllowEditorManageTeamMember
}

func (tp *TeamPermission) DoesViewerCanManageTeamMember() bool {
	return tp.AllowViewerManageTeamMember
}

func (tp *TeamPermission) DoesBlockRegister() bool {
	return tp.BlockRegister
}
