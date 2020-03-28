package feeds

// User -
type User struct {
	ID       string `db:"id,omitempty" json:"id"`
	Nickname string `db:"nickname" json:"nickname"`
	Password string `db:"password,omitempty" json:"password"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}

// Plant -
type Plant struct {
	ID        string `db:"id,omitempty" json:"id"`
	UserID    string `db:"userid,omitempty" json:"userID"`
	FeedID    string `db:"feedid" json:"feedID"`
	DeviceID  string `db:"deviceid,omitempty" json:"deviceID"`
	DeviceBox int    `db:"deviceBox,omitempty" json:"deviceBox"`
	Name      string `db:"name" json:"name"`
	Settings  string `db:"settings" json:"settings"`
}

// Timelapse -
type Timelapse struct {
	ID           string `db:"id,omitempty" json:"id"`
	UserID       string `db:"userid,omitempty" json:"userID"`
	PlantID      string `db:"boxid" json:"boxID"`
	ControllerID string `db:"controllerid" json:"controllerID"`
	Rotate       string `db:"rotate" json:"rotate"`
	Name         string `db:"name" json:"name"`
	Strain       string `db:"strain" json:"strain"`
	DropboxToken string `db:"dropboxtoken" json:"dropboxToken"`
	UploadName   string `db:"uploadname" json:"uploadName"`
}

// Device -
type Device struct {
	ID         string `db:"id,omitempty" json:"id"`
	Identifier string `db:"identifier" json:"identifier"`
	Name       string `db:"name" json:"name"`
	IP         string `db:"ip" json:"ip"`
	Mdns       string `db:"mdns" json:"mdns"`
}

// Feed -
type Feed struct {
	ID     string `db:"id,omitempty" json:"id"`
	UserID string `db:"userid,omitempty" json:"userID"`
	Name   string `db:"name" json:"name"`
}

// FeedEntry -
type FeedEntry struct {
	ID        string `db:"id,omitempty" json:"id"`
	FeedID    string `db:"feedid" json:"feedID"`
	CreatedAt string `db:"createdAt" json:"createdAt"`
}

// PlantSharing -
type PlantSharing struct {
	UserID  string `db:"userid,omitempty" json:"userID"`
	PlantID string `db:"plantid" json:"plantID"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}

// UserEnd -
type UserEnd struct {
	ID     string `db:"id,omitempty" json:"id"`
	UserID string `db:"userid" json:"userID"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}

// UserEndPlant -
type UserEndPlant struct {
	UserEndID string `db:"userendid,omitempty" json:"userEndID"`
	PlantID   string `db:"plantid" json:"plantID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}

// UserEndTimelapse -
type UserEndTimelapse struct {
	UserEndID   string `db:"userendid,omitempty" json:"userEndID"`
	TimelapseID string `db:"timelapseid" json:"timelapseID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}

// UserEndDevice -
type UserEndDevice struct {
	UserEndID string `db:"userendid,omitempty" json:"userEndID"`
	DeviceID  string `db:"deviceid" json:"deviceID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}

// UserEndFeed -
type UserEndFeed struct {
	UserEndID string `db:"userendid,omitempty" json:"userEndID"`
	FeedID    string `db:"feedid" json:"feedID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}

// UserEndFeedEntries -
type UserEndFeedEntries struct {
	UserEndID   string `db:"userendid,omitempty" json:"userEndID"`
	FeedEntryID string `db:"feedentryid" json:"feedEntryID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}

// UserEndFeedMedias -
type UserEndFeedMedias struct {
	UserEndID   string `db:"userendid,omitempty" json:"userEndID"`
	FeedMediaID string `db:"feedmediaid" json:"feedMediaID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt string `db:"cat,omitempty" json:"cat"`
	UpdatedAt string `db:"uat,omitempty" json:"uat"`
}
