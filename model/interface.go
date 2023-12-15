package model

import "context"

type Interface interface {
	Complete(context context.Context, prompt string) string
}
