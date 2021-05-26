package rule

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/entql"
	"github.com/roffe/gin-ent/ent/privacy"
	"github.com/roffe/gin-ent/ent/todo"
	"github.com/roffe/gin-ent/ent/user"
	"github.com/roffe/gin-ent/pkg/viewer"
)

func FilterOnlyOwnTodos() privacy.QueryMutationRule {
	type TodoFilter interface {
		Where(p entql.P)
	}

	return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
		view := viewer.FromContext(ctx)
		if view == nil {
			return privacy.Denyf("missing viewer in context")
		}

		tf, ok := f.(TodoFilter)
		if !ok {
			return privacy.Denyf("unexpected filter type %T", f)
		}

		tf.Where(
			entql.HasEdgeWith(todo.EdgeOwner,
				sqlgraph.WrapFunc(func(s *sql.Selector) {
					s.Where(sql.EQ(s.C(user.FieldID), view.GetID()))
				}),
			),
		)

		return privacy.Skip
	})
}
