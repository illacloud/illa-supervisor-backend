package model

import (
	"github.com/google/uuid"
)

/**
 * team drive folder design
 *
 * - /{bucket-name}
 *   - /team-{team_uid}
 *     - /system  (for team system object storage)
 *       - /icon
 *     - /team (for team drive data)
 */

const TEAM_FOLDER_PREFIX = "team-"
const TEAM_SYSTEM_FOLDER = "/system"
const TEAM_ICON_FOLDER = "/icon"
const TEAM_SPACE_FOLDER = "/team"

type TeamDrive struct {
	UID              uuid.UUID  `json:"uid"`
	Drive            S3Instance `json:"-"`
	TeamSystemFolder string     `json:"teamsystemfolder"`
	TeamSpaceFolder  string     `json:"teamspacefolder"`
}

func NewTeamDrive(drive *Drive) *TeamDrive {
	return &TeamDrive{
		Drive: drive.TeamDriveS3Instance,
	}
}

func (d *TeamDrive) SetTeam(team *Team) {
	d.UID = team.GetUID()
	d.TeamSystemFolder = TEAM_FOLDER_PREFIX + team.GetUIDInString() + TEAM_SYSTEM_FOLDER
	d.TeamSpaceFolder = TEAM_FOLDER_PREFIX + team.GetUIDInString() + TEAM_SPACE_FOLDER
}

func (d *TeamDrive) GetIconUploadPreSignedURL(fileName string) (string, error) {
	path := d.TeamSystemFolder + TEAM_ICON_FOLDER + "/" + fileName
	return d.Drive.GetPreSignedPutURL(path)
}
