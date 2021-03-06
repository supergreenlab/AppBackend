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

package users

import (
	cmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
)

// Init -
func Init(router *httprouter.Router) {
	anon := cmiddlewares.AnonStack()
	auth := cmiddlewares.AuthStack()

	router.POST("/login", anon.Wrap(loginHandler()))
	router.POST("/user", anon.Wrap(createUserHandler))

	router.PUT("/user", auth.Wrap(updateUserHandler))
	router.GET("/users/me", auth.Wrap(meHandler)) // TODO remove this one:/
	router.GET("/user/me", auth.Wrap(meHandler))

	router.POST("/profilePicUploadURL", auth.Wrap(profilePicUploadURLHandler))
}
