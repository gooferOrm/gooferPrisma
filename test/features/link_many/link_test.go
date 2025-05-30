package linkmany

import (
	"context"
	"testing"

	"github.com/gooferOrm/goofer/test"
	"github.com/gooferOrm/gooferest/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func str(v string) *string {
	return &v
}

func TestLinkMany(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "link many in create",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "unrelated",
					email: "unrelated",
					username: "unrelated",
					name: "unrelated",
					posts: {
						create: [{
							id: "a",
							title: "common",
							content: "a",
						}, {
							id: "b",
							title: "common",
							content: "b",
						}, {
							id: "c",
							title: "common",
							content: "c",
						}, {
							id: "d",
							title: "common",
							content: "d",
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			created, err := client.User.CreateOne(
				User.Email.Set("new"),
				User.Username.Set("new"),
				User.ID.Set("new"),
				User.Posts.Link(
					Post.ID.Equals("a"),
					Post.ID.Equals("b"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			newUser := InnerUser{
				ID:       "new",
				Email:    "new",
				Username: "new",
			}

			massert.Equal(t, created, &UserModel{InnerUser: newUser})

			actual, err := client.User.FindUnique(
				User.ID.Equals("new"),
			).With(
				User.Posts.Fetch().Take(2),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			authorID := "new"
			expected := &UserModel{
				InnerUser: newUser,
				RelationsUser: RelationsUser{
					Posts: []PostModel{{
						InnerPost: InnerPost{
							ID:       "a",
							Title:    "common",
							Content:  str("a"),
							AuthorID: &authorID,
						},
					}, {
						InnerPost: InnerPost{
							ID:       "b",
							Title:    "common",
							Content:  str("b"),
							AuthorID: &authorID,
						},
					}},
				},
			}

			massert.Equal(t, expected, actual)
			massert.Equal(t, 2, len(actual.Posts()))
		},
	}, {
		name: "link many in update",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "new",
					email: "new",
					username: "new",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOnePost(data: {
					id: "a",
					title: "common",
					content: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOnePost(data: {
					id: "b",
					title: "common",
					content: "b",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOnePost(data: {
					id: "c",
					title: "common",
					content: "c",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			created, err := client.User.FindUnique(
				User.ID.Equals("new"),
			).Update(
				User.Posts.Link(
					Post.ID.Equals("a"),
					Post.ID.Equals("b"),
					Post.ID.Equals("c"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			newUser := InnerUser{
				ID:       "new",
				Email:    "new",
				Username: "new",
			}

			massert.Equal(t, created, &UserModel{InnerUser: newUser})

			actual, err := client.User.FindUnique(
				User.ID.Equals("new"),
			).With(
				User.Posts.Fetch().Take(2),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			authorID := "new"
			expected := &UserModel{
				InnerUser: newUser,
				RelationsUser: RelationsUser{
					Posts: []PostModel{{
						InnerPost: InnerPost{
							ID:       "a",
							Title:    "common",
							Content:  str("a"),
							AuthorID: &authorID,
						},
					}, {
						InnerPost: InnerPost{
							ID:       "b",
							Title:    "common",
							Content:  str("b"),
							AuthorID: &authorID,
						},
					}},
				},
			}

			massert.Equal(t, expected, actual)
			massert.Equal(t, 2, len(actual.Posts()))
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, test.Databases, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
