package feeds

import (
	"time"

	"github.com/gofrs/uuid"
)

// User -
type User struct {
	ID       uuid.NullUUID `db:"id,omitempty" json:"id"`
	Nickname string        `db:"nickname" json:"nickname"`
	Password string        `db:"password,omitempty" json:"password"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// Box -
type Box struct {
	ID        uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID    uuid.UUID     `db:"userid" json:"userID"`
	DeviceID  uuid.NullUUID `db:"deviceid,omitempty" json:"deviceID,omitempty"`
	DeviceBox *uint         `db:"devicebox,omitempty" json:"deviceBox,omitempty"`
	Name      string        `db:"name" json:"name"`

	Settings string `db:"settings" json:"settings"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// Plant -
type Plant struct {
	ID     uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID uuid.UUID     `db:"userid" json:"userID"`
	BoxID  uuid.UUID     `db:"boxid" json:"boxID"`
	FeedID uuid.UUID     `db:"feedid" json:"feedID"`
	Name   string        `db:"name" json:"name"`

	Settings string `db:"settings" json:"settings"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// Timelapse -
type Timelapse struct {
	ID           uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID       uuid.UUID     `db:"userid" json:"userID"`
	PlantID      uuid.UUID     `db:"plantid" json:"plantID"`
	ControllerID string        `db:"controllerid" json:"controllerID"`
	Rotate       string        `db:"rotate" json:"rotate"`
	Name         string        `db:"name" json:"name"`
	Strain       string        `db:"strain" json:"strain"`
	DropboxToken string        `db:"dropboxtoken" json:"dropboxToken"`
	UploadName   string        `db:"uploadname" json:"uploadName"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// Device -
type Device struct {
	ID         uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID     uuid.UUID     `db:"userid" json:"userID"`
	Identifier string        `db:"identifier" json:"identifier"`
	Name       string        `db:"name" json:"name"`
	IP         string        `db:"ip" json:"ip"`
	Mdns       string        `db:"mdns" json:"mdns"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// Feed -
type Feed struct {
	ID     uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID uuid.UUID     `db:"userid" json:"userID"`
	Name   string        `db:"name" json:"name"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// FeedEntry -
type FeedEntry struct {
	ID     uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID uuid.UUID     `db:"userid" json:"userID"`
	FeedID uuid.UUID     `db:"feedid" json:"feedID"`
	Date   string        `db:"createdat" json:"date"`
	Type   string        `db:"etype" json:"type"`

	Params string `db:"params" json:"params"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// FeedMedia -
type FeedMedia struct {
	ID          uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID      uuid.UUID     `db:"userid" json:"userID"`
	FeedEntryID uuid.UUID     `db:"feedentryid" json:"feedEntryID"`
	FileRef     string        `db:"fileref" json:"fileRef"`

	Params string `db:"params" json:"params"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// PlantSharing -
type PlantSharing struct {
	UserID   uuid.NullUUID `db:"userid" json:"userID"`
	PlantID  uuid.UUID     `db:"plantid" json:"plantID"`
	ToUserID uuid.UUID     `db:"touserid" json:"toUserID"`

	Params string `db:"params" json:"params"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// UserEnd -
type UserEnd struct {
	ID     uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID uuid.UUID     `db:"userid" json:"userID"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// UserEndPlant -
type UserEndPlant struct {
	UserEndID uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	PlantID   uuid.UUID `db:"plantid" json:"plantID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// UserEndTimelapse -
type UserEndTimelapse struct {
	UserEndID   uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	TimelapseID uuid.UUID `db:"timelapseid" json:"timelapseID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// UserEndDevice -
type UserEndDevice struct {
	UserEndID uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	DeviceID  uuid.UUID `db:"deviceid" json:"deviceID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// UserEndFeed -
type UserEndFeed struct {
	UserEndID uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	FeedID    uuid.UUID `db:"feedid" json:"feedID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// UserEndFeedEntries -
type UserEndFeedEntries struct {
	UserEndID   uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	FeedEntryID uuid.UUID `db:"feedentryid" json:"feedEntryID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// UserEndFeedMedias -
type UserEndFeedMedias struct {
	UserEndID   uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	FeedMediaID uuid.UUID `db:"feedmediaid" json:"feedMediaID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}
