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
	PlantID      string `db:"boxid" json:"boxID"`
	ControllerID string `db:"controllerid" json:"controllerID"`
	Rotate       string `db:"rotate" json:"rotate"`
	Name         string `db:"name" json:"name"`
	Strain       string `db:"strain" json:"strain"`
	UploadName   string `db:"uploadName" json:"uploadName"`
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
