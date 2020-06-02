// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-2020 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package todoist

import (
	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/files"
	"code.vikunja.io/api/pkg/models"
	"code.vikunja.io/api/pkg/timeutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/d4l3k/messagediff.v1"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

func TestConvertTodoistToVikunja(t *testing.T) {

	config.InitConfig()

	time1, err := time.Parse(time.RFC3339Nano, "2014-09-26T08:25:05Z")
	assert.NoError(t, err)
	time3, err := time.Parse(time.RFC3339Nano, "2014-10-21T08:25:05Z")
	assert.NoError(t, err)
	dueTime, err := time.Parse(time.RFC3339Nano, "2020-05-31T00:00:00Z")
	assert.NoError(t, err)
	nilTime, err := time.Parse(time.RFC3339Nano, "1970-01-01T00:00:00Z")
	assert.NoError(t, err)
	exampleFile, err := ioutil.ReadFile(config.ServiceRootpath.GetString() + "/pkg/modules/migration/wunderlist/testimage.jpg")
	assert.NoError(t, err)

	makeTestItem := func(id, projectId int, hasDueDate, hasLabels, done bool) *item {
		item := &item{
			ID:            id,
			UserID:        1855589,
			ProjectID:     projectId,
			Content:       "Task" + strconv.Itoa(id),
			Priority:      1,
			ParentID:      0,
			ChildOrder:    1,
			DateAdded:     time1,
			DateCompleted: nilTime,
		}

		if done {
			item.Checked = 1
			item.DateCompleted = time3
		}

		if hasLabels {
			item.Labels = []int{
				80000,
				80001,
				80002,
				80003,
			}
		}

		if hasDueDate {
			item.Due = &dueDate{
				Date:        "2020-05-31",
				Timezone:    nil,
				IsRecurring: false,
			}
		}

		return item
	}

	testSync := &sync{
		Projects: []*project{
			{
				ID:         396936926,
				Name:       "Project1",
				Color:      30,
				ChildOrder: 1,
				Collapsed:  0,
				Shared:     false,
				IsDeleted:  0,
				IsArchived: 0,
				IsFavorite: 0,
			},
			{
				ID:         396936927,
				Name:       "Project2",
				Color:      37,
				ChildOrder: 1,
				Collapsed:  0,
				Shared:     false,
				IsDeleted:  0,
				IsArchived: 0,
				IsFavorite: 0,
			},
			{
				ID:         396936928,
				Name:       "Project3 - Archived",
				Color:      37,
				ChildOrder: 1,
				Collapsed:  0,
				Shared:     false,
				IsDeleted:  0,
				IsArchived: 1,
				IsFavorite: 0,
			},
		},
		Items: []*item{
			makeTestItem(400000000, 396936926, false, false, false),
			makeTestItem(400000001, 396936926, false, false, false),
			makeTestItem(400000002, 396936926, false, false, false),
			makeTestItem(400000003, 396936926, true, true, true),
			makeTestItem(400000004, 396936926, false, true, false),
			makeTestItem(400000005, 396936926, true, false, true),
			makeTestItem(400000006, 396936926, true, false, true),
			{
				ID:         400000110,
				UserID:     1855589,
				ProjectID:  396936926,
				Content:    "Task with parent",
				Priority:   2,
				ParentID:   400000006,
				ChildOrder: 1,
				Checked:    0,
				DateAdded:  time1,
			},
			makeTestItem(400000106, 396936926, true, true, true),
			makeTestItem(400000107, 396936926, false, false, true),
			makeTestItem(400000108, 396936926, false, false, true),
			makeTestItem(400000109, 396936926, false, false, true),

			makeTestItem(400000007, 396936927, true, false, false),
			makeTestItem(400000008, 396936927, true, false, false),
			makeTestItem(400000009, 396936927, false, false, false),
			makeTestItem(400000010, 396936927, false, false, true),
			makeTestItem(400000101, 396936927, false, false, false),
			makeTestItem(400000102, 396936927, true, true, false),
			makeTestItem(400000103, 396936927, false, true, false),
			makeTestItem(400000104, 396936927, false, true, false),
			makeTestItem(400000105, 396936927, true, true, false),

			makeTestItem(400000111, 396936928, false, false, true),
		},
		Labels: []*label{
			{
				ID:    80000,
				Name:  "Label1",
				Color: 30,
			},
			{
				ID:    80001,
				Name:  "Label2",
				Color: 31,
			},
			{
				ID:    80002,
				Name:  "Label3",
				Color: 32,
			},
			{
				ID:    80003,
				Name:  "Label4",
				Color: 33,
			},
		},
		Notes: []*note{
			{
				ID:        101476,
				PostedUID: 1855589,
				ItemID:    400000000,
				Content:   "Lorem Ipsum dolor sit amet",
				Posted:    time1,
			},
			{
				ID:        101477,
				PostedUID: 1855589,
				ItemID:    400000001,
				Content:   "Lorem Ipsum dolor sit amet",
				Posted:    time1,
			},
			{
				ID:        101478,
				PostedUID: 1855589,
				ItemID:    400000003,
				Content:   "Lorem Ipsum dolor sit amet",
				Posted:    time1,
			},
			{
				ID:        101479,
				PostedUID: 1855589,
				ItemID:    400000010,
				Content:   "Lorem Ipsum dolor sit amet",
				Posted:    time1,
			},
			{
				ID:        101480,
				PostedUID: 1855589,
				ItemID:    400000101,
				Content:   "Lorem Ipsum dolor sit amet",
				FileAttachment: &fileAttachment{
					FileName:    "file.md",
					FileType:    "text/plain",
					FileSize:    12345,
					FileURL:     "https://vikunja.io/testimage.jpg", // Using an image which we are hosting, so it'll still be up
					UploadState: "completed",
				},
				Posted: time1,
			},
		},
		ProjectNotes: []*projectNote{
			{
				ID:        102000,
				Content:   "Lorem Ipsum dolor sit amet",
				ProjectID: 396936926,
				Posted:    time3,
				PostedUID: 1855589,
			},
			{
				ID:        102001,
				Content:   "Lorem Ipsum dolor sit amet 2",
				ProjectID: 396936926,
				Posted:    time3,
				PostedUID: 1855589,
			},
			{
				ID:        102002,
				Content:   "Lorem Ipsum dolor sit amet 3",
				ProjectID: 396936926,
				Posted:    time3,
				PostedUID: 1855589,
			},
			{
				ID:        102003,
				Content:   "Lorem Ipsum dolor sit amet 4",
				ProjectID: 396936927,
				Posted:    time3,
				PostedUID: 1855589,
			},
			{
				ID:        102004,
				Content:   "Lorem Ipsum dolor sit amet 5",
				ProjectID: 396936927,
				Posted:    time3,
				PostedUID: 1855589,
			},
		},
		Reminders: []*reminder{
			{
				ID:     103000,
				ItemID: 400000000,
				Due: &dueDate{
					Date:        "2020-06-15",
					IsRecurring: false,
				},
				MmOffset: 180,
			},
			{
				ID:     103001,
				ItemID: 400000000,
				Due: &dueDate{
					Date:        "2020-06-16",
					IsRecurring: false,
				},
			},
			{
				ID:     103002,
				ItemID: 400000002,
				Due: &dueDate{
					Date:        "2020-07-15",
					IsRecurring: true,
				},
			},
			{
				ID:     103003,
				ItemID: 400000003,
				Due: &dueDate{
					Date:        "2020-06-15",
					IsRecurring: false,
				},
			},
			{
				ID:     103004,
				ItemID: 400000005,
				Due: &dueDate{
					Date:        "2020-06-15",
					IsRecurring: false,
				},
			},
			{
				ID:     103006,
				ItemID: 400000009,
				Due: &dueDate{
					Date:        "2020-06-15",
					IsRecurring: false,
				},
			},
		},
	}

	vikunjaLabels := []*models.Label{
		{
			Title:    "Label1",
			HexColor: todoistColors[30],
		},
		{
			Title:    "Label2",
			HexColor: todoistColors[31],
		},
		{
			Title:    "Label3",
			HexColor: todoistColors[32],
		},
		{
			Title:    "Label4",
			HexColor: todoistColors[33],
		},
	}

	expectedHierachie := []*models.NamespaceWithLists{
		{
			Namespace: models.Namespace{
				Title: "Migrated from todoist",
			},
			Lists: []*models.List{
				{
					Title:       "Project1",
					Description: "Lorem Ipsum dolor sit amet\nLorem Ipsum dolor sit amet 2\nLorem Ipsum dolor sit amet 3",
					HexColor:    todoistColors[30],
					Tasks: []*models.Task{
						{
							Title:       "Task400000000",
							Description: "Lorem Ipsum dolor sit amet",
							Done:        false,
							Created:     timeutil.FromTime(time1),
							Reminders: []timeutil.TimeStamp{
								timeutil.FromTime(time.Date(2020, time.June, 15, 0, 0, 0, 0, time.UTC)),
								timeutil.FromTime(time.Date(2020, time.June, 16, 0, 0, 0, 0, time.UTC)),
							},
						},
						{
							Title:       "Task400000001",
							Description: "Lorem Ipsum dolor sit amet",
							Done:        false,
							Created:     timeutil.FromTime(time1),
						},
						{
							Title:   "Task400000002",
							Done:    false,
							Created: timeutil.FromTime(time1),
							Reminders: []timeutil.TimeStamp{
								timeutil.FromTime(time.Date(2020, time.July, 15, 0, 0, 0, 0, time.UTC)),
							},
						},
						{
							Title:       "Task400000003",
							Description: "Lorem Ipsum dolor sit amet",
							Done:        true,
							DueDate:     timeutil.FromTime(dueTime),
							Created:     timeutil.FromTime(time1),
							DoneAt:      timeutil.FromTime(time3),
							Labels:      vikunjaLabels,
							Reminders: []timeutil.TimeStamp{
								timeutil.FromTime(time.Date(2020, time.June, 15, 0, 0, 0, 0, time.UTC)),
							},
						},
						{
							Title:   "Task400000004",
							Done:    false,
							Created: timeutil.FromTime(time1),
							Labels:  vikunjaLabels,
						},
						{
							Title:   "Task400000005",
							Done:    true,
							DueDate: timeutil.FromTime(dueTime),
							Created: timeutil.FromTime(time1),
							DoneAt:  timeutil.FromTime(time3),
							Reminders: []timeutil.TimeStamp{
								timeutil.FromTime(time.Date(2020, time.June, 15, 0, 0, 0, 0, time.UTC)),
							},
						},
						{
							Title:   "Task400000006",
							Done:    true,
							DueDate: timeutil.FromTime(dueTime),
							Created: timeutil.FromTime(time1),
							DoneAt:  timeutil.FromTime(time3),
							RelatedTasks: map[models.RelationKind][]*models.Task{
								models.RelationKindSubtask: {
									{
										Title:    "Task with parent",
										Done:     false,
										Priority: 2,
										Created:  timeutil.FromTime(time1),
										DoneAt:   timeutil.FromTime(nilTime),
									},
								},
							},
						},
						{
							Title:    "Task with parent",
							Done:     false,
							Priority: 2,
							Created:  timeutil.FromTime(time1),
							DoneAt:   timeutil.FromTime(nilTime),
						},
						{
							Title:   "Task400000106",
							Done:    true,
							DueDate: timeutil.FromTime(dueTime),
							Created: timeutil.FromTime(time1),
							DoneAt:  timeutil.FromTime(time3),
							Labels:  vikunjaLabels,
						},
						{
							Title:   "Task400000107",
							Done:    true,
							Created: timeutil.FromTime(time1),
							DoneAt:  timeutil.FromTime(time3),
						},
						{
							Title:   "Task400000108",
							Done:    true,
							Created: timeutil.FromTime(time1),
							DoneAt:  timeutil.FromTime(time3),
						},
						{
							Title:   "Task400000109",
							Done:    true,
							Created: timeutil.FromTime(time1),
							DoneAt:  timeutil.FromTime(time3),
						},
					},
				},
				{
					Title:       "Project2",
					Description: "Lorem Ipsum dolor sit amet 4\nLorem Ipsum dolor sit amet 5",
					HexColor:    todoistColors[37],
					Tasks: []*models.Task{
						{
							Title:   "Task400000007",
							Done:    false,
							DueDate: timeutil.FromTime(dueTime),
							Created: timeutil.FromTime(time1),
						},
						{
							Title:   "Task400000008",
							Done:    false,
							DueDate: timeutil.FromTime(dueTime),
							Created: timeutil.FromTime(time1),
						},
						{
							Title:   "Task400000009",
							Done:    false,
							Created: timeutil.FromTime(time1),
							Reminders: []timeutil.TimeStamp{
								timeutil.FromTime(time.Date(2020, time.June, 15, 0, 0, 0, 0, time.UTC)),
							},
						},
						{
							Title:       "Task400000010",
							Description: "Lorem Ipsum dolor sit amet",
							Done:        true,
							Created:     timeutil.FromTime(time1),
							DoneAt:      timeutil.FromTime(time3),
						},
						{
							Title:       "Task400000101",
							Description: "Lorem Ipsum dolor sit amet",
							Done:        false,
							Created:     timeutil.FromTime(time1),
							Attachments: []*models.TaskAttachment{
								{
									File: &files.File{
										Name:        "file.md",
										Mime:        "text/plain",
										Size:        12345,
										Created:     time1,
										CreatedUnix: timeutil.FromTime(time1),
										FileContent: exampleFile,
									},
									Created: timeutil.FromTime(time1),
								},
							},
						},
						{
							Title:   "Task400000102",
							Done:    false,
							DueDate: timeutil.FromTime(dueTime),
							Created: timeutil.FromTime(time1),
							Labels:  vikunjaLabels,
						},
						{
							Title:   "Task400000103",
							Done:    false,
							Created: timeutil.FromTime(time1),
							Labels:  vikunjaLabels,
						},
						{
							Title:   "Task400000104",
							Done:    false,
							Created: timeutil.FromTime(time1),
							Labels:  vikunjaLabels,
						},
						{
							Title:   "Task400000105",
							Done:    false,
							DueDate: timeutil.FromTime(dueTime),
							Created: timeutil.FromTime(time1),
							Labels:  vikunjaLabels,
						},
					},
				},
				{
					Title:      "Project3 - Archived",
					HexColor:   todoistColors[37],
					IsArchived: true,
					Tasks: []*models.Task{
						{
							Title:   "Task400000111",
							Done:    true,
							Created: timeutil.FromTime(time1),
							DoneAt:  timeutil.FromTime(time3),
						},
					},
				},
			},
		},
	}

	hierachie, err := convertTodoistToVikunja(testSync)
	assert.NoError(t, err)
	assert.NotNil(t, hierachie)
	if diff, equal := messagediff.PrettyDiff(hierachie, expectedHierachie); !equal {
		t.Errorf("converted todoist data = %v, want %v, diff: %v", hierachie, expectedHierachie, diff)
	}
}