package cloudscale

import "context"

type operationPathKey struct{}

// WithOperationPath attaches the path template for the current operation to ctx.
// The template should use :id placeholders for dynamic segments, e.g.
// "v1/servers/:id" or "v1/servers/:id/start".
func WithOperationPath(ctx context.Context, template string) context.Context {
	return context.WithValue(ctx, operationPathKey{}, template)
}

// OperationPath returns the path template previously attached to ctx via
// WithOperationPath. If no template is present, it returns the empty string.
func OperationPath(ctx context.Context) string {
	template, _ := ctx.Value(operationPathKey{}).(string)
	return template
}
