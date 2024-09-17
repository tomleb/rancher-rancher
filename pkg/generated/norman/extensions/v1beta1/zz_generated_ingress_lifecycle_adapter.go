package v1beta1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ingressLifecycleConverter struct {
	lifecycle IngressLifecycle
}

func (w *ingressLifecycleConverter) CreateContext(_ context.Context, obj *v1beta1.Ingress) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *ingressLifecycleConverter) RemoveContext(_ context.Context, obj *v1beta1.Ingress) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *ingressLifecycleConverter) UpdatedContext(_ context.Context, obj *v1beta1.Ingress) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type IngressLifecycle interface {
	Create(obj *v1beta1.Ingress) (runtime.Object, error)
	Remove(obj *v1beta1.Ingress) (runtime.Object, error)
	Updated(obj *v1beta1.Ingress) (runtime.Object, error)
}

type IngressLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1beta1.Ingress) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1beta1.Ingress) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1beta1.Ingress) (runtime.Object, error)
}

type ingressLifecycleAdapter struct {
	lifecycle IngressLifecycleContext
}

func (w *ingressLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *ingressLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *ingressLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *ingressLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1beta1.Ingress))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *ingressLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *ingressLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1beta1.Ingress))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *ingressLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *ingressLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1beta1.Ingress))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewIngressLifecycleAdapter(name string, clusterScoped bool, client IngressInterface, l IngressLifecycle) IngressHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(IngressGroupVersionResource)
	}
	adapter := &ingressLifecycleAdapter{lifecycle: &ingressLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1beta1.Ingress) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewIngressLifecycleAdapterContext(name string, clusterScoped bool, client IngressInterface, l IngressLifecycleContext) IngressHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(IngressGroupVersionResource)
	}
	adapter := &ingressLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1beta1.Ingress) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
