package engine

import "context"

type Server struct {
	Ctx    context.Context
	cancel context.CancelFunc
}
