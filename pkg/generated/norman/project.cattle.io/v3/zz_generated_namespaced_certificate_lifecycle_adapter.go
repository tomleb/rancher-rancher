package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type namespacedCertificateLifecycleConverter struct {
	lifecycle NamespacedCertificateLifecycle
}

func (w *namespacedCertificateLifecycleConverter) CreateContext(_ context.Context, obj *v3.NamespacedCertificate) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *namespacedCertificateLifecycleConverter) RemoveContext(_ context.Context, obj *v3.NamespacedCertificate) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *namespacedCertificateLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.NamespacedCertificate) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NamespacedCertificateLifecycle interface {
	Create(obj *v3.NamespacedCertificate) (runtime.Object, error)
	Remove(obj *v3.NamespacedCertificate) (runtime.Object, error)
	Updated(obj *v3.NamespacedCertificate) (runtime.Object, error)
}

type NamespacedCertificateLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.NamespacedCertificate) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.NamespacedCertificate) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.NamespacedCertificate) (runtime.Object, error)
}

type namespacedCertificateLifecycleAdapter struct {
	lifecycle NamespacedCertificateLifecycleContext
}

func (w *namespacedCertificateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *namespacedCertificateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *namespacedCertificateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *namespacedCertificateLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.NamespacedCertificate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *namespacedCertificateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *namespacedCertificateLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.NamespacedCertificate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *namespacedCertificateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *namespacedCertificateLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.NamespacedCertificate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNamespacedCertificateLifecycleAdapter(name string, clusterScoped bool, client NamespacedCertificateInterface, l NamespacedCertificateLifecycle) NamespacedCertificateHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NamespacedCertificateGroupVersionResource)
	}
	adapter := &namespacedCertificateLifecycleAdapter{lifecycle: &namespacedCertificateLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.NamespacedCertificate) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNamespacedCertificateLifecycleAdapterContext(name string, clusterScoped bool, client NamespacedCertificateInterface, l NamespacedCertificateLifecycleContext) NamespacedCertificateHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NamespacedCertificateGroupVersionResource)
	}
	adapter := &namespacedCertificateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.NamespacedCertificate) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
