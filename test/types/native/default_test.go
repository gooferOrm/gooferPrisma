w
package db

import (
	"context"
	"testing"

	"github.com/gooferOrm/goofer/test"
	"github.com/gooferOrm/gooferest/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestNative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "native create",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			a := "a"
			b := "b"
			created, err := client.User.CreateOne(
				User.A.Set(a),
				User.B.Set(b),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID: "123",
					A:  a,
					B:  &b,
				},
			}

			massert.Equal(t, expected, created)

			actual, err := client.User.FindUnique(User.ID.Equals(created.ID)).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "native find by native field",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			a := "a"
			b := "b"
			created, err := client.User.CreateOne(
				User.A.Set(a),
				User.B.Set(b),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID: "123",
					A:  a,
					B:  &b,
				},
			}

			massert.Equal(t, expected, created)

			actual, err := client.User.FindFirst(
				User.A.Equals(a),
				User.B.Equals(b),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
