package db

import (
	"context"
	"testing"

	"github.com/gooferOrm/goofer/test"
	"github.com/gooferOrm/gooferest/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "FindFirst",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.FindFirst(
				User.Email.Equals("john@example.com"),
			).Exec(ctx)
			massert.Equal(t, ErrNotFound, err)
		},
	}, {
		name: "FindUnique",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.FindUnique(
				User.Email.Equals("john@example.com"),
			).Exec(ctx)
			massert.Equal(t, ErrNotFound, err)
		},
	}, {
		name: "Update",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.FindUnique(
				User.Email.Equals("john@example.com"),
			).Update(
				User.Email.Set("asdf"),
			).Exec(ctx)
			massert.Equal(t, ErrNotFound, err)
		},
	}, {
		name: "Delete",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.FindUnique(
				User.Email.Equals("john@example.com"),
			).Delete().Exec(ctx)
			massert.Equal(t, ErrNotFound, err)
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
