package model

import (
	"github.com/google/uuid"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InviteStorage struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewInviteStorage(db *gorm.DB, logger *zap.SugaredLogger) *InviteStorage {
	return &InviteStorage{
		logger: logger,
		db:     db,
	}
}

func (d *InviteStorage) Create(u *Invite) (int, error) {
	if err := d.db.Create(u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (d *InviteStorage) RetrieveByID(id int) (*Invite, error) {
	u := &Invite{}
	if err := d.db.First(u, id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *InviteStorage) RetrieveByUID(uid uuid.UUID) (*Invite, error) {
	u := &Invite{}
	if err := d.db.Where("uid = ?", uid).First(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *InviteStorage) RetrieveByIDs(ids []int) ([]*Invite, error) {
	invites := []*Invite{}
	if err := d.db.Where("(team_member_id) IN ?", ids).Find(&invites).Error; err != nil {
		return nil, err
	}
	return invites, nil
}

func (d *InviteStorage) RetrieveInviteByEmailByIDs(ids []int) ([]*Invite, error) {
	invites := []*Invite{}
	if err := d.db.Where("category = ? AND (team_member_id) IN ?", CATEGORY_INVITE_BY_EMAIL, ids).Find(&invites).Error; err != nil {
		return nil, err
	}
	return invites, nil
}

func (d *InviteStorage) RetrieveAvaliableInviteByTeamIDAndTeamMemberID(teamID int, teamMemberID int) (*Invite, error) {
	invite := &Invite{}
	if err := d.db.Where("team_id = ? AND team_member_id = ? AND status = ?", teamID, teamMemberID, INVITE_RECORD_STATUS_OK).First(&invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

func (d *InviteStorage) RetrieveAvaliableInviteLinkByTeamIDAndUserRole(teamID int, userRole int) (*Invite, error) {
	invite := &Invite{}
	if err := d.db.Where("team_id = ? AND user_role = ? AND category = ? AND status = ?", teamID, userRole, CATEGORY_INVITE_BY_LINK, INVITE_RECORD_STATUS_OK).First(&invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

func (d *InviteStorage) RetrieveInviteByTeamIDAndEmail(teamID int, email string) (*Invite, error) {
	invite := &Invite{}
	if err := d.db.Where("team_id = ? AND email = ? AND category = ?", teamID, email, CATEGORY_INVITE_BY_EMAIL).First(&invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

func (d *InviteStorage) RetrieveAvaliableInviteByTeamIDAndEmail(teamID int, email string) (*Invite, error) {
	invite := &Invite{}
	if err := d.db.Where("team_id = ? AND email = ? AND category = ? AND status = ?", teamID, email, CATEGORY_INVITE_BY_EMAIL, INVITE_RECORD_STATUS_OK).First(&invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

// Note: this method is not avaliable invite record inluded.
func (d *InviteStorage) RetrieveInviteByEmailAndTeamIDAndUserRole(teamID int, email string, userRole int) (*Invite, error) {
	invite := &Invite{}
	if err := d.db.Where("team_id = ? AND email = ? AND user_role = ? AND category = ?", teamID, email, userRole, CATEGORY_INVITE_BY_EMAIL).First(&invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

func (d *InviteStorage) Update(u *Invite) error {
	if err := d.db.Model(&Invite{}).Where("id = ?", u.ID).UpdateColumns(u).Error; err != nil {
		return err
	}
	return nil
}

func (d *InviteStorage) DeleteByID(id int) error {
	if err := d.db.Delete(&Invite{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (d *InviteStorage) DeleteByUID(uid uuid.UUID) error {
	if err := d.db.Where("uid = ?", uid).Delete(&Invite{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *InviteStorage) DeleteByTeamID(teamID int) error {
	if err := d.db.Where("team_id = ?", teamID).Delete(&Invite{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *InviteStorage) DeleteByTeamIDAndTeamMemberID(teamID int, teamMemberID int) error {
	if err := d.db.Where("team_id = ? AND team_member_id = ?", teamID, teamMemberID).Delete(&Invite{}).Error; err != nil {
		return err
	}
	return nil
}
