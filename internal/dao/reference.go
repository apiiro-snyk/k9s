package dao

import (
	"context"
	"errors"

	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	_ Accessor = (*Reference)(nil)
)

type Reference struct {
	NonResource
}

func (r *Reference) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	gvr, ok := ctx.Value(internal.KeyGVR).(string)
	if !ok {
		return nil, errors.New("No context GVR found")
	}
	switch gvr {
	case "v1/serviceaccounts":
		return r.ScanSA(ctx)
	default:
		return r.Scan(ctx)
	}
}

func (c *Reference) Get(ctx context.Context, path string) (runtime.Object, error) {
	panic("NYI")
}

func (r *Reference) Scan(ctx context.Context) ([]runtime.Object, error) {
	refs, err := ScanForRefs(ctx, r.Factory)
	if err != nil {
		return nil, err
	}

	fqn, ok := ctx.Value(internal.KeyPath).(string)
	if !ok {
		return nil, errors.New("expecting context Path")
	}
	ns, _ := client.Namespaced(fqn)
	oo := make([]runtime.Object, 0, len(refs))
	for _, ref := range refs {
		_, n := client.Namespaced(ref.FQN)
		oo = append(oo, render.ReferenceRes{
			Namespace: ns,
			Name:      n,
			GVR:       ref.GVR,
		})
	}

	return oo, nil
}

func (r *Reference) ScanSA(ctx context.Context) ([]runtime.Object, error) {
	refs, err := ScanForSARefs(ctx, r.Factory)
	if err != nil {
		return nil, err
	}

	fqn, ok := ctx.Value(internal.KeyPath).(string)
	if !ok {
		return nil, errors.New("expecting context Path")
	}
	ns, _ := client.Namespaced(fqn)
	oo := make([]runtime.Object, 0, len(refs))
	for _, ref := range refs {
		_, n := client.Namespaced(ref.FQN)
		oo = append(oo, render.ReferenceRes{
			Namespace: ns,
			Name:      n,
			GVR:       ref.GVR,
		})
	}

	return oo, nil
}