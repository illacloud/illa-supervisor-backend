package model

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Storage struct {
	UserStorage           *UserStorage
	TeamStorage           *TeamStorage
	TeamMemberStorage     *TeamMemberStorage
	InviteStorage         *InviteStorage
	InvitationCodeStorage *InvitationCodeStorage
}

func NewStorage(postgresDriver *gorm.DB, logger *zap.SugaredLogger) *Storage {
	userStorage := NewUserStorage(postgresDriver, logger)
	teamStorage := NewTeamStorage(postgresDriver, logger)
	teamMemberStorage := NewTeamMemberStorage(postgresDriver, logger)
	inviteStorage := NewInviteStorage(postgresDriver, logger)
	invitationCodeStorage := NewInvitationCodeStorage(postgresDriver, logger)
	return &Storage{
		UserStorage:           userStorage,
		TeamStorage:           teamStorage,
		TeamMemberStorage:     teamMemberStorage,
		InviteStorage:         inviteStorage,
		InvitationCodeStorage: invitationCodeStorage,
	}
}
