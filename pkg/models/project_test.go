// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-2021 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public Licensee as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public Licensee for more details.
//
// You should have received a copy of the GNU Affero General Public Licensee
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"reflect"
	"testing"

	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/files"
	"code.vikunja.io/api/pkg/user"
	"github.com/stretchr/testify/assert"
)

func TestProject_CreateOrUpdate(t *testing.T) {
	usr := &user.User{
		ID:       1,
		Username: "user1",
		Email:    "user1@example.com",
	}

	t.Run("create", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			db.LoadAndAssertFixtures(t)
			s := db.NewSession()
			project := Project{
				Title:       "test",
				Description: "Lorem Ipsum",
				NamespaceID: 1,
			}
			err := project.Create(s, usr)
			assert.NoError(t, err)
			err = s.Commit()
			assert.NoError(t, err)
			db.AssertExists(t, "projects", map[string]interface{}{
				"id":           project.ID,
				"title":        project.Title,
				"description":  project.Description,
				"namespace_id": project.NamespaceID,
			}, false)
		})
		t.Run("nonexistant namespace", func(t *testing.T) {
			db.LoadAndAssertFixtures(t)
			s := db.NewSession()
			project := Project{
				Title:       "test",
				Description: "Lorem Ipsum",
				NamespaceID: 999999,
			}
			err := project.Create(s, usr)
			assert.Error(t, err)
			assert.True(t, IsErrNamespaceDoesNotExist(err))
			_ = s.Close()
		})
		t.Run("nonexistant owner", func(t *testing.T) {
			db.LoadAndAssertFixtures(t)
			s := db.NewSession()
			usr := &user.User{ID: 9482385}
			project := Project{
				Title:       "test",
				Description: "Lorem Ipsum",
				NamespaceID: 1,
			}
			err := project.Create(s, usr)
			assert.Error(t, err)
			assert.True(t, user.IsErrUserDoesNotExist(err))
			_ = s.Close()
		})
		t.Run("existing identifier", func(t *testing.T) {
			db.LoadAndAssertFixtures(t)
			s := db.NewSession()
			project := Project{
				Title:       "test",
				Description: "Lorem Ipsum",
				Identifier:  "test1",
				NamespaceID: 1,
			}
			err := project.Create(s, usr)
			assert.Error(t, err)
			assert.True(t, IsErrProjectIdentifierIsNotUnique(err))
			_ = s.Close()
		})
		t.Run("non ascii characters", func(t *testing.T) {
			db.LoadAndAssertFixtures(t)
			s := db.NewSession()
			project := Project{
				Title:       "приффки фсем",
				Description: "Lorem Ipsum",
				NamespaceID: 1,
			}
			err := project.Create(s, usr)
			assert.NoError(t, err)
			err = s.Commit()
			assert.NoError(t, err)
			db.AssertExists(t, "projects", map[string]interface{}{
				"id":           project.ID,
				"title":        project.Title,
				"description":  project.Description,
				"namespace_id": project.NamespaceID,
			}, false)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			db.LoadAndAssertFixtures(t)
			s := db.NewSession()
			project := Project{
				ID:          1,
				Title:       "test",
				Description: "Lorem Ipsum",
				NamespaceID: 1,
			}
			project.Description = "Lorem Ipsum dolor sit amet."
			err := project.Update(s, usr)
			assert.NoError(t, err)
			err = s.Commit()
			assert.NoError(t, err)
			db.AssertExists(t, "projects", map[string]interface{}{
				"id":           project.ID,
				"title":        project.Title,
				"description":  project.Description,
				"namespace_id": project.NamespaceID,
			}, false)
		})
		t.Run("nonexistant", func(t *testing.T) {
			db.LoadAndAssertFixtures(t)
			s := db.NewSession()
			project := Project{
				ID:          99999999,
				Title:       "test",
				NamespaceID: 1,
			}
			err := project.Update(s, usr)
			assert.Error(t, err)
			assert.True(t, IsErrProjectDoesNotExist(err))
			_ = s.Close()

		})
		t.Run("existing identifier", func(t *testing.T) {
			db.LoadAndAssertFixtures(t)
			s := db.NewSession()
			project := Project{
				Title:       "test",
				Description: "Lorem Ipsum",
				Identifier:  "test1",
				NamespaceID: 1,
			}
			err := project.Create(s, usr)
			assert.Error(t, err)
			assert.True(t, IsErrProjectIdentifierIsNotUnique(err))
			_ = s.Close()
		})
		t.Run("change namespace", func(t *testing.T) {
			t.Run("own", func(t *testing.T) {
				usr := &user.User{
					ID:       6,
					Username: "user6",
					Email:    "user6@example.com",
				}

				db.LoadAndAssertFixtures(t)
				s := db.NewSession()
				project := Project{
					ID:          6,
					Title:       "Test6",
					Description: "Lorem Ipsum",
					NamespaceID: 7, // from 6
				}
				can, err := project.CanUpdate(s, usr)
				assert.NoError(t, err)
				assert.True(t, can)
				err = project.Update(s, usr)
				assert.NoError(t, err)
				err = s.Commit()
				assert.NoError(t, err)
				db.AssertExists(t, "projects", map[string]interface{}{
					"id":           project.ID,
					"title":        project.Title,
					"description":  project.Description,
					"namespace_id": project.NamespaceID,
				}, false)
			})
			// FIXME: The check for whether the namespace is archived is missing in namespace.CanWrite
			// t.Run("archived own", func(t *testing.T) {
			// 	db.LoadAndAssertFixtures(t)
			// 	s := db.NewSession()
			// 	project := Project{
			// 		ID:          1,
			// 		Title:       "Test1",
			// 		Description: "Lorem Ipsum",
			// 		NamespaceID: 16, // from 1
			// 	}
			// 	can, err := project.CanUpdate(s, usr)
			// 	assert.NoError(t, err)
			// 	assert.False(t, can) // namespace is archived and thus not writeable
			// 	_ = s.Close()
			// })
			t.Run("others", func(t *testing.T) {
				db.LoadAndAssertFixtures(t)
				s := db.NewSession()
				project := Project{
					ID:          1,
					Title:       "Test1",
					Description: "Lorem Ipsum",
					NamespaceID: 2, // from 1
				}
				can, _ := project.CanUpdate(s, usr)
				assert.False(t, can) // namespace is not writeable by us
				_ = s.Close()
			})
			t.Run("pseudo namespace", func(t *testing.T) {
				usr := &user.User{
					ID:       6,
					Username: "user6",
					Email:    "user6@example.com",
				}

				db.LoadAndAssertFixtures(t)
				s := db.NewSession()
				project := Project{
					ID:          6,
					Title:       "Test6",
					Description: "Lorem Ipsum",
					NamespaceID: -1,
				}
				err := project.Update(s, usr)
				assert.Error(t, err)
				assert.True(t, IsErrProjectCannotBelongToAPseudoNamespace(err))
			})
		})
	})
}

func TestProject_Delete(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		project := Project{
			ID: 1,
		}
		err := project.Delete(s, &user.User{ID: 1})
		assert.NoError(t, err)
		err = s.Commit()
		assert.NoError(t, err)
		db.AssertMissing(t, "projects", map[string]interface{}{
			"id": 1,
		})
	})
	t.Run("with background", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		files.InitTestFileFixtures(t)
		s := db.NewSession()
		project := Project{
			ID: 25,
		}
		err := project.Delete(s, &user.User{ID: 6})
		assert.NoError(t, err)
		err = s.Commit()
		assert.NoError(t, err)
		db.AssertMissing(t, "projects", map[string]interface{}{
			"id": 25,
		})
		db.AssertMissing(t, "files", map[string]interface{}{
			"id": 1,
		})
	})
}

func TestProject_DeleteBackgroundFileIfExists(t *testing.T) {
	t.Run("project with background", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		files.InitTestFileFixtures(t)
		s := db.NewSession()
		file := &files.File{ID: 1}
		project := Project{
			ID:               1,
			BackgroundFileID: file.ID,
		}
		err := SetProjectBackground(s, project.ID, file, "")
		assert.NoError(t, err)
		err = project.DeleteBackgroundFileIfExists()
		assert.NoError(t, err)
	})
	t.Run("project with invalid background", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		files.InitTestFileFixtures(t)
		s := db.NewSession()
		file := &files.File{ID: 9999}
		project := Project{
			ID:               1,
			BackgroundFileID: file.ID,
		}
		err := SetProjectBackground(s, project.ID, file, "")
		assert.NoError(t, err)
		err = project.DeleteBackgroundFileIfExists()
		assert.Error(t, err)
		assert.True(t, files.IsErrFileDoesNotExist(err))
	})
	t.Run("project without background", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		files.InitTestFileFixtures(t)
		project := Project{ID: 1}
		err := project.DeleteBackgroundFileIfExists()
		assert.NoError(t, err)
	})
}

func TestProject_ReadAll(t *testing.T) {
	t.Run("all in namespace", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		// Get all projects for our namespace
		projects, err := GetProjectsByNamespaceID(s, 1, &user.User{})
		assert.NoError(t, err)
		assert.Equal(t, len(projects), 2)
		_ = s.Close()
	})
	t.Run("all projects for user", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		u := &user.User{ID: 1}
		project := Project{}
		projects3, _, _, err := project.ReadAll(s, u, "", 1, 50)

		assert.NoError(t, err)
		assert.Equal(t, reflect.TypeOf(projects3).Kind(), reflect.Slice)
		ls := projects3.([]*Project)
		assert.Equal(t, 16, len(ls))
		assert.Equal(t, int64(3), ls[0].ID) // Project 3 has a position of 1 and should be sorted first
		assert.Equal(t, int64(1), ls[1].ID)
		assert.Equal(t, int64(4), ls[2].ID)
		_ = s.Close()
	})
	t.Run("projects for nonexistant user", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		usr := &user.User{ID: 999999}
		project := Project{}
		_, _, _, err := project.ReadAll(s, usr, "", 1, 50)
		assert.Error(t, err)
		assert.True(t, user.IsErrUserDoesNotExist(err))
		_ = s.Close()
	})
	t.Run("search", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		u := &user.User{ID: 1}
		project := Project{}
		projects3, _, _, err := project.ReadAll(s, u, "TEST10", 1, 50)

		assert.NoError(t, err)
		ls := projects3.([]*Project)
		assert.Equal(t, 1, len(ls))
		assert.Equal(t, int64(10), ls[0].ID)
		_ = s.Close()
	})
}

func TestProject_ReadOne(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		u := &user.User{ID: 1}
		l := &Project{ID: 1}
		can, _, err := l.CanRead(s, u)
		assert.NoError(t, err)
		assert.True(t, can)
		err = l.ReadOne(s, u)
		assert.NoError(t, err)
		assert.Equal(t, "Test1", l.Title)
	})
	t.Run("with subscription", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		u := &user.User{ID: 6}
		l := &Project{ID: 12}
		can, _, err := l.CanRead(s, u)
		assert.NoError(t, err)
		assert.True(t, can)
		err = l.ReadOne(s, u)
		assert.NoError(t, err)
		assert.NotNil(t, l.Subscription)
	})
}