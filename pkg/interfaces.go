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

package appbackend

import "github.com/gofrs/uuid"

// Object -
type Object interface {
	GetID() uuid.NullUUID
}

// Objects -
type Objects interface {
	Each(func(Object))
}

// UserObject -
type UserObject interface {
	Object
	SetUserID(uuid.UUID)
	GetUserID() uuid.UUID
}

type S3Path struct {
	Path   *string
	Bucket string
}

type S3FileHolder interface {
	SetURLs(paths []string)
	GetURLs() (paths []S3Path)
}

type S3FileHolders interface {
	AsFeedMediasArray() []S3FileHolder
}
