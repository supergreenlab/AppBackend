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

// Object -
type Object interface {
	GetID() uuid.UUID
}

// Objects -
type Objects interface {
	Each(func(Object))
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

// GetID -
func (o Box) GetID() uuid.UUID {
	return o.ID.UUID
}

// Boxes -
type Boxes []Box

// Each -
func (os Boxes) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
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

// GetID -
func (o Plant) GetID() uuid.UUID {
	return o.ID.UUID
}

// Plants -
type Plants []Plant

// Each -
func (os Plants) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
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

// GetID -
func (o Timelapse) GetID() uuid.UUID {
	return o.ID.UUID
}

// Timelapses -
type Timelapses []Timelapse

// Each -
func (os Timelapses) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
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

// GetID -
func (o Device) GetID() uuid.UUID {
	return o.ID.UUID
}

// Devices -
type Devices []Device

// Each -
func (os Devices) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
}

// Feed -
type Feed struct {
	ID     uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID uuid.UUID     `db:"userid" json:"userID"`
	Name   string        `db:"name" json:"name"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o Feed) GetID() uuid.UUID {
	return o.ID.UUID
}

// Feeds -
type Feeds []Feed

// Each -
func (os Feeds) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
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

// GetID -
func (o FeedEntry) GetID() uuid.UUID {
	return o.ID.UUID
}

// FeedEntries -
type FeedEntries []FeedEntry

// Each -
func (os FeedEntries) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
}

// FeedMedia -
type FeedMedia struct {
	ID            uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID        uuid.UUID     `db:"userid" json:"userID"`
	FeedEntryID   uuid.UUID     `db:"feedentryid" json:"feedEntryID"`
	FilePath      string        `db:"filepath" json:"filePath"`
	ThumbnailPath string        `db:"thumbnailpath" json:"thumbnailPath"`

	Params string `db:"params" json:"params"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o FeedMedia) GetID() uuid.UUID {
	return o.ID.UUID
}

// FeedMedias -
type FeedMedias []FeedMedia

// Each -
func (os FeedMedias) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
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

// UserEndObject -
type UserEndObject interface {
	SetUserEndID(uuid.UUID)
	SetObjectID(uuid.UUID)
	SetDirty(bool)
	SetSent(bool)
}

// UserEndBox -
type UserEndBox struct {
	UserEndID uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	BoxID     uuid.UUID `db:"boxid" json:"boxID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// SetUserEndID -
func (ueo *UserEndBox) SetUserEndID(id uuid.UUID) {
	ueo.UserEndID = id
}

// SetObjectID -
func (ueo *UserEndBox) SetObjectID(id uuid.UUID) {
	ueo.BoxID = id
}

// SetDirty -
func (ueo *UserEndBox) SetDirty(dirty bool) {
	ueo.Dirty = dirty
}

// SetSent -
func (ueo *UserEndBox) SetSent(sent bool) {
	ueo.Sent = sent
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

// SetUserEndID -
func (ueo *UserEndPlant) SetUserEndID(id uuid.UUID) {
	ueo.UserEndID = id
}

// SetObjectID -
func (ueo *UserEndPlant) SetObjectID(id uuid.UUID) {
	ueo.PlantID = id
}

// SetDirty -
func (ueo *UserEndPlant) SetDirty(dirty bool) {
	ueo.Dirty = dirty
}

// SetSent -
func (ueo *UserEndPlant) SetSent(sent bool) {
	ueo.Sent = sent
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

// SetUserEndID -
func (ueo *UserEndTimelapse) SetUserEndID(id uuid.UUID) {
	ueo.UserEndID = id
}

// SetObjectID -
func (ueo *UserEndTimelapse) SetObjectID(id uuid.UUID) {
	ueo.TimelapseID = id
}

// SetDirty -
func (ueo *UserEndTimelapse) SetDirty(dirty bool) {
	ueo.Dirty = dirty
}

// SetSent -
func (ueo *UserEndTimelapse) SetSent(sent bool) {
	ueo.Sent = sent
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

// SetUserEndID -
func (ueo *UserEndDevice) SetUserEndID(id uuid.UUID) {
	ueo.UserEndID = id
}

// SetObjectID -
func (ueo *UserEndDevice) SetObjectID(id uuid.UUID) {
	ueo.DeviceID = id
}

// SetDirty -
func (ueo *UserEndDevice) SetDirty(dirty bool) {
	ueo.Dirty = dirty
}

// SetSent -
func (ueo *UserEndDevice) SetSent(sent bool) {
	ueo.Sent = sent
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

// SetUserEndID -
func (ueo *UserEndFeed) SetUserEndID(id uuid.UUID) {
	ueo.UserEndID = id
}

// SetObjectID -
func (ueo *UserEndFeed) SetObjectID(id uuid.UUID) {
	ueo.FeedID = id
}

// SetDirty -
func (ueo *UserEndFeed) SetDirty(dirty bool) {
	ueo.Dirty = dirty
}

// SetSent -
func (ueo *UserEndFeed) SetSent(sent bool) {
	ueo.Sent = sent
}

// UserEndFeedEntry -
type UserEndFeedEntry struct {
	UserEndID   uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	FeedEntryID uuid.UUID `db:"feedentryid" json:"feedEntryID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// SetUserEndID -
func (ueo *UserEndFeedEntry) SetUserEndID(id uuid.UUID) {
	ueo.UserEndID = id
}

// SetObjectID -
func (ueo *UserEndFeedEntry) SetObjectID(id uuid.UUID) {
	ueo.FeedEntryID = id
}

// SetDirty -
func (ueo *UserEndFeedEntry) SetDirty(dirty bool) {
	ueo.Dirty = dirty
}

// SetSent -
func (ueo *UserEndFeedEntry) SetSent(sent bool) {
	ueo.Sent = sent
}

// UserEndFeedMedia -
type UserEndFeedMedia struct {
	UserEndID   uuid.UUID `db:"userendid,omitempty" json:"userEndID"`
	FeedMediaID uuid.UUID `db:"feedmediaid" json:"feedMediaID"`

	Sent  bool `db:"sent" json:"sent"`
	Dirty bool `db:"dirty" json:"dirty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// SetUserEndID -
func (ueo *UserEndFeedMedia) SetUserEndID(id uuid.UUID) {
	ueo.UserEndID = id
}

// SetObjectID -
func (ueo *UserEndFeedMedia) SetObjectID(id uuid.UUID) {
	ueo.FeedMediaID = id
}

// SetDirty -
func (ueo *UserEndFeedMedia) SetDirty(dirty bool) {
	ueo.Dirty = dirty
}

// SetSent -
func (ueo *UserEndFeedMedia) SetSent(sent bool) {
	ueo.Sent = sent
}
