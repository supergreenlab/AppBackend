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

package tools

import (
	"fmt"
	"reflect"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"

	"github.com/gofrs/uuid"
	"upper.io/db.v3/lib/sqlbuilder"
)

// CheckUserID - checks a given field value against a userID
func CheckUserID(sess sqlbuilder.Database, uid uuid.UUID, o db.UserObject, collection, field string, optional bool, factory func() db.UserObject) error {
	var id uuid.UUID
	idFieldValue := reflect.ValueOf(o).Elem().FieldByName(field).Interface()
	if v, ok := idFieldValue.(uuid.UUID); ok == true {
		id = v
	} else if v, ok := idFieldValue.(uuid.NullUUID); ok == true {
		if !v.Valid && !optional {
			return fmt.Errorf("Missing value for field %s", field)
		} else if !v.Valid && optional {
			return nil
		}
		id = v.UUID
	}

	parent := factory()
	err := sess.Collection(collection).Find("id", id).One(parent)
	if err != nil {
		return err
	}

	uidParent := parent.GetUserID()

	if uid != uidParent {
		return fmt.Errorf("Parent is owned by another user")
	}
	return nil
}
