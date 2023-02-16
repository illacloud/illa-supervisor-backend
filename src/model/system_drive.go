package model

/**
 * system drive folder design
 *
 * - /{bucket-name}
 *   - /user-{user_uid} (for user object storage)
 *     - /avatar (for user avatar)
 */

const USER_FOLDER_PREFIX = "user-"
const USER_AVATAR_FOLVER = "/avatar"

type SystemDrive struct {
	Drive      S3Instance `json:"-"`
	UserFolder string     `json:"userFolder"`
}

func NewSystemDrive(drive *Drive) *SystemDrive {
	return &SystemDrive{
		Drive: drive.SystemDriveS3Instance,
	}
}

func (d *SystemDrive) SetUser(user *User) {
	d.UserFolder = USER_FOLDER_PREFIX + user.GetUIDInString()
}

func (d *SystemDrive) GetUserAvatarUploadPreSignedURL(fileName string) (string, error) {
	path := d.UserFolder + USER_AVATAR_FOLVER + "/" + fileName
	return d.Drive.GetPreSignedPutURL(path)
}
