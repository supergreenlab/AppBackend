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
)

// Products -
type Products struct {
	ID     uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID uuid.UUID     `db:"userid" json:"userID"`

	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`

	FilePath      string `db:"filepath" json:"filePath"`
	ThumbnailPath string `db:"thumbnailpath" json:"thumbnailPath"`

	Categories string `db:"categories" json:"categories"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o Products) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *Products) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o Products) GetUserID() uuid.UUID {
	return o.UserID
}

// Suppliers -
type Suppliers struct {
	ID     uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID uuid.UUID     `db:"userid" json:"userID"`

	Name        string `db:"name" json:"name"`
	URL         string `db:"url" json:"url"`
	Description string `db:"description" json:"description"`
	Locals      string `db:"locals" json:"locals"`

	FilePath      string `db:"filepath" json:"filePath"`
	ThumbnailPath string `db:"thumbnailpath" json:"thumbnailPath"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o Suppliers) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *Suppliers) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o Suppliers) GetUserID() uuid.UUID {
	return o.UserID
}

// ProductsSuppliers -
type ProductsSuppliers struct {
	ID         uuid.NullUUID `db:"id,omitempty" json:"id"`
	UserID     uuid.UUID     `db:"userid" json:"userID"`
	ProductID  uuid.UUID     `db:"productid" json:"productID"`
	SupplierID uuid.UUID     `db:"supplierid" json:"supplierID"`

	URL   string  `db:"url" json:"url"`
	Price float64 `db:"price" json:"price"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

// GetID -
func (o ProductsSuppliers) GetID() uuid.NullUUID {
	return o.ID
}

// SetUserID -
func (o *ProductsSuppliers) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

// GetUserID -
func (o ProductsSuppliers) GetUserID() uuid.UUID {
	return o.UserID
}
