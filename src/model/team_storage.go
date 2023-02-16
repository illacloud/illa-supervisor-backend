package model

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TeamStorage struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewTeamStorage(db *gorm.DB, logger *zap.SugaredLogger) *TeamStorage {
	return &TeamStorage{
		logger: logger,
		db:     db,
	}
}

func (d *TeamStorage) Create(u *Team) (int, error) {
	if err := d.db.Create(u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (d *TeamStorage) RetrieveByID(id int) (*Team, error) {
	u := &Team{}
	if err := d.db.First(u, id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *TeamStorage) RetrieveByIDs(ids []int) ([]*Team, error) {
	teams := []*Team{}
	if err := d.db.Where("(id) IN ?", ids).Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (d *TeamStorage) RetrieveByUID(uid string) (*Team, error) {
	u := &Team{}
	if err := d.db.Where("uid = ?", uid).First(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *TeamStorage) RetrieveByIdentifier(identifier string) (*Team, error) {
	u := &Team{}
	if err := d.db.Where("identifier = ?", identifier).First(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *TeamStorage) IsIdentifierExists(identifier string) bool {
	var count int64
	d.db.Where("identifier = ?", identifier).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

func (d *TeamStorage) UpdateByID(u *Team) error {
	if err := d.db.Model(&Team{}).Where("id = ?", u.ID).UpdateColumns(u).Error; err != nil {
		return err
	}
	return nil
}

func (d *TeamStorage) DeleteByID(id int) error {
	if err := d.db.Delete(&Team{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (d *TeamStorage) DeleteByUID(uid string) error {
	if err := d.db.Delete(&Team{}).Where("uid = ?", uid).Error; err != nil {
		return err
	}
	return nil
}
