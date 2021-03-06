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
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
)

// Init -
func Init(router *httprouter.Router) {
	anon := middlewares.AnonStack()
	auth := middlewares.AuthStack()

	router.POST("/product", auth.Wrap(createProductsHandler))
	router.POST("/supplier", auth.Wrap(createSuppliersHandler))
	router.POST("/productsupplier", auth.Wrap(createProductSuppliersHandler))

	router.GET("/products/search", anon.Wrap(searchProducts))
}
