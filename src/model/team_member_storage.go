package model

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TeamMemberStorage struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewTeamMemberStorage(db *gorm.DB, logger *zap.SugaredLogger) *TeamMemberStorage {
	return &TeamMemberStorage{
		logger: logger,
		db:     db,
	}
}

func (d *TeamMemberStorage) Create(u *TeamMember) (int, error) {
	if err := d.db.Create(u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (d *TeamMemberStorage) RetrieveByID(id int) (*TeamMember, error) {
	u := &TeamMember{}
	if err := d.db.First(u, id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *TeamMemberStorage) RetrieveByUID(uid string) (*TeamMember, error) {
	u := &TeamMember{}
	if err := d.db.Where("uid = ?", uid).First(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *TeamMemberStorage) RetrieveByTeamID(team_id int) ([]*TeamMember, error) {
	var teamMembers []*TeamMember
	if err := d.db.Where("team_id = ?", team_id).Find(&teamMembers).Error; err != nil {
		return nil, err
	}
	return teamMembers, nil
}

func (d *TeamMemberStorage) RetrieveByTeamIDAndID(team_id int, id int) (*TeamMember, error) {
	var teamMember *TeamMember
	if err := d.db.Where("team_id = ? AND id = ?", team_id, id).First(&teamMember).Error; err != nil {
		return nil, err
	}
	return teamMember, nil
}

func (d *TeamMemberStorage) RetrieveByUserID(user_id int) ([]*TeamMember, error) {
	var teamMembers []*TeamMember
	if err := d.db.Where("user_id = ?", user_id).Find(&teamMembers).Error; err != nil {
		return nil, err
	}
	return teamMembers, nil
}

func (d *TeamMemberStorage) RetrieveByTeamIDAndUserID(teamID int, userID int) (*TeamMember, error) {
	var teamMember *TeamMember
	if err := d.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&teamMember).Error; err != nil {
		return nil, err
	}
	return teamMember, nil
}

func (d *TeamMemberStorage) RetrieveTeamMemberByTeamIDAndID(teamID int, id int) (*TeamMember, error) {
	var teamMember *TeamMember
	if err := d.db.Where("team_id = ? AND id = ?", teamID, id).First(&teamMember).Error; err != nil {
		return nil, err
	}
	return teamMember, nil
}

func (d *TeamMemberStorage) DoesTeamIncludedTargetUser(teamID int, userID int) (bool, error) {
	var teamMember *TeamMember
	var count int64
	if err := d.db.Model(&teamMember).Where("team_id = ? AND user_id = ?", teamID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (d *TeamMemberStorage) IsNowUserIsTeamOwner(userID int) (bool, error) {
	var teamMember *TeamMember
	var count int64
	if err := d.db.Model(&teamMember).Where("user_id = ? AND user_role = ?", userID, USER_ROLE_OWNER).Count(&count).Error; err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (d *TeamMemberStorage) Update(u *TeamMember) error {
	if err := d.db.Model(&TeamMember{}).Where("id = ?", u.ID).UpdateColumns(u).Error; err != nil {
		return err
	}
	return nil
}

func (d *TeamMemberStorage) DeleteByID(id int) error {
	if err := d.db.Delete(&TeamMember{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (d *TeamMemberStorage) DeleteByUID(uid string) error {
	if err := d.db.Delete(&TeamMember{}).Where("uid = ?", uid).Error; err != nil {
		return err
	}
	return nil
}

func (d *TeamMemberStorage) DeleteByTeamID(teamID int) error {
	if err := d.db.Where("team_id = ?", teamID).Delete(&TeamMember{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *TeamMemberStorage) DeleteByUserID(userID int) error {
	if err := d.db.Where("user_id = ?", userID).Delete(&TeamMember{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *TeamMemberStorage) DeleteByTeamIDAndUserID(teamID int, userID int) error {
	if err := d.db.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&TeamMember{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *TeamMemberStorage) DeleteByIDAndTeamID(id int, teamID int) error {
	if err := d.db.Where("id = ? AND team_id = ?", id, teamID).Delete(&TeamMember{}).Error; err != nil {
		return err
	}
	return nil
}
