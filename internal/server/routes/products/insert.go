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

package products

import (
	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/rileyr/middleware"
)

var createProductsHandler = middlewares.InsertEndpoint(
	"products",
	func() interface{} { return &db.Products{} },
	[]middleware.Middleware{
		middlewares.SetUserID,
	},
	[]middleware.Middleware{},
)

var createSuppliersHandler = middlewares.InsertEndpoint(
	"suppliers",
	func() interface{} { return &db.Suppliers{} },
	[]middleware.Middleware{
		middlewares.SetUserID,
	},
	[]middleware.Middleware{},
)

var createProductSuppliersHandler = middlewares.InsertEndpoint(
	"productsuppliers",
	func() interface{} { return &db.ProductSuppliers{} },
	[]middleware.Middleware{
		middlewares.SetUserID,
	},
	[]middleware.Middleware{},
)
