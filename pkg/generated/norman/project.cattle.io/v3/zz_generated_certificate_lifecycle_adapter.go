package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type certificateLifecycleConverter struct {
	lifecycle CertificateLifecycle
}

func (w *certificateLifecycleConverter) CreateContext(_ context.Context, obj *v3.Certificate) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *certificateLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Certificate) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *certificateLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Certificate) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type CertificateLifecycle interface {
	Create(obj *v3.Certificate) (runtime.Object, error)
	Remove(obj *v3.Certificate) (runtime.Object, error)
	Updated(obj *v3.Certificate) (runtime.Object, error)
}

type CertificateLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Certificate) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Certificate) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Certificate) (runtime.Object, error)
}

type certificateLifecycleAdapter struct {
	lifecycle CertificateLifecycleContext
}

func (w *certificateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *certificateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *certificateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *certificateLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Certificate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *certificateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *certificateLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Certificate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *certificateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *certificateLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Certificate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewCertificateLifecycleAdapter(name string, clusterScoped bool, client CertificateInterface, l CertificateLifecycle) CertificateHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(CertificateGroupVersionResource)
	}
	adapter := &certificateLifecycleAdapter{lifecycle: &certificateLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Certificate) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewCertificateLifecycleAdapterContext(name string, clusterScoped bool, client CertificateInterface, l CertificateLifecycleContext) CertificateHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(CertificateGroupVersionResource)
	}
	adapter := &certificateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Certificate) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
