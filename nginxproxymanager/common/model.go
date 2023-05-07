package common

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type IModel[T any] interface {
	Load(ctx context.Context, resource *T) diag.Diagnostics
}
