/*
 * Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package appbackend

import (
	"time"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v3"
)

// Box -
type Box struct {
	ID        uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID    uuid.UUID     `db:"userid" json:"userID"`
	DeviceID  uuid.NullUUID `db:"deviceid" json:"deviceID"`
	DeviceBox *uint         `db:"devicebox,omitempty" json:"deviceBox,omitempty"`
	FeedID    uuid.NullUUID `db:"feedid" json:"feedID"`
	Name      string        `db:"name" json:"name"`

	Settings string `db:"settings" json:"settings"`

	Deleted bool `db:"deleted" json:"deleted"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o Box) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *Box) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o Box) GetUserID() uuid.UUID {
	return o.UserID
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
	ID            uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID        uuid.UUID     `db:"userid" json:"userID"`
	BoxID         uuid.UUID     `db:"boxid" json:"boxID"`
	FeedID        uuid.UUID     `db:"feedid" json:"feedID"`
	Name          string        `db:"name" json:"name"`
	Single        bool          `db:"single" json:"single"` // TODO remove this field
	Public        bool          `db:"is_public" json:"public"`
	AlertsEnabled bool          `db:"alerts_enabled" json:"alertsEnabled"`

	Settings string `db:"settings" json:"settings"`

	Deleted  bool `db:"deleted" json:"deleted"`
	Archived bool `db:"archived" json:"archived"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o Plant) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *Plant) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o Plant) GetUserID() uuid.UUID {
	return o.UserID
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
	ID      uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID  uuid.UUID     `db:"userid" json:"userID"`
	PlantID uuid.UUID     `db:"plantid" json:"plantID"`

	Name     string `db:"name" json:"name"`
	Type     string `db:"ttype" json:"type"`
	Settings string `db:"settings" json:"settings"`

	Deleted bool `db:"deleted" json:"deleted"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o Timelapse) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *Timelapse) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o Timelapse) GetUserID() uuid.UUID {
	return o.UserID
}

// Timelapses -
type Timelapses []Timelapse

// Each -
func (os Timelapses) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
}

// TimelapseFrame -
type TimelapseFrame struct {
	ID          uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID      uuid.UUID     `db:"userid" json:"userID"`
	TimelapseID uuid.UUID     `db:"timelapseid" json:"timelapseID"`

	FilePath string `db:"filepath" json:"filePath"`
	Meta     string `db:"meta" json:"meta"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

type TimelapseFrameMeta struct {
	MetricsMeta
}

func (r *TimelapseFrame) SetURLs(paths []string) {
	r.FilePath = paths[0]
}

func (r TimelapseFrame) GetURLs() []S3Path {
	return []S3Path{
		S3Path{
			Path:   &r.FilePath,
			Bucket: "timelapses",
		},
	}
}

// GetID -
func (o TimelapseFrame) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *TimelapseFrame) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o TimelapseFrame) GetUserID() uuid.UUID {
	return o.UserID
}

// TimelapseFrames -
type TimelapseFrames []TimelapseFrame

// Each -
func (os TimelapseFrames) Each(fn func(Object)) {
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

	Deleted bool `db:"deleted" json:"deleted"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o Device) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *Device) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o Device) GetUserID() uuid.UUID {
	return o.UserID
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
	ID         uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID     uuid.UUID     `db:"userid" json:"userID"`
	Name       string        `db:"name" json:"name"`
	IsNewsFeed bool          `db:"isnewsfeed" json:"isNewsFeed"`

	Deleted bool `db:"deleted" json:"deleted"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o Feed) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *Feed) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o Feed) GetUserID() uuid.UUID {
	return o.UserID
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
	Date   time.Time     `db:"createdat" json:"date"`
	Type   string        `db:"etype" json:"type"`

	Params string      `db:"params" json:"params"`
	Meta   null.String `db:"meta,omitempty" json:"meta,omitempty"`

	Deleted bool `db:"deleted" json:"deleted"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

type FeedEntryMeta struct {
	MetricsMeta
}

// GetID -
func (o FeedEntry) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *FeedEntry) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o FeedEntry) GetUserID() uuid.UUID {
	return o.UserID
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

	Deleted bool `db:"deleted" json:"deleted"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o FeedMedia) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *FeedMedia) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o FeedMedia) GetUserID() uuid.UUID {
	return o.UserID
}

// FeedMedias -
type FeedMedias []FeedMedia

// Each -
func (os FeedMedias) Each(fn func(Object)) {
	for _, o := range os {
		fn(o)
	}
}
