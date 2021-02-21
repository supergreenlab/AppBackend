/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
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

package db

import (
	"time"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v3"
)

// User -
type User struct {
	ID       uuid.NullUUID `db:"id,omitempty" json:"id"`
	Nickname string        `db:"nickname" json:"nickname"`
	Password string        `db:"password,omitempty" json:"password"`

	Pic   null.String `db:"pic,omitempty" json:"pic,omitempty"`
	Liked bool        `db:"liked,omitempty" json:"liked,omitempty"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

func (u *User) GetID() uuid.NullUUID {
	return u.ID
}
