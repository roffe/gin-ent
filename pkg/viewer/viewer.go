package viewer

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Role for viewer actions.
type Role int

// List of roles.
const (
	_ Role = 1 << iota
	Admin
	View
)

// Viewer describes the query/mutation viewer-context.
type Viewer interface {
	Admin() bool // If viewer is admin.
	GetID() int  // Get the viewer user id
}

// UserViewer describes a user-viewer.
type UserViewer struct {
	Role Role // Attached roles.
	ID   int  // Attached user id
}

func (v UserViewer) Admin() bool {
	return v.Role&Admin != 0
}

func (v UserViewer) GetID() int {
	return v.ID
}

type CtxKey struct{}

// FromContext returns the Viewer stored in a context.
func FromContext(ctx context.Context) Viewer {
	v, _ := ctx.Value(CtxKey{}).(Viewer)
	return v
}

// NewContext returns a copy of parent context with the given Viewer attached with it.
func NewContext(parent context.Context, v Viewer) context.Context {
	return context.WithValue(parent, CtxKey{}, v)
}

func FromGinContext(c *gin.Context) context.Context {
	v, _ := c.Get("viewer")
	return context.WithValue(c.Request.Context(), CtxKey{}, v)
}
