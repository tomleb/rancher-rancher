package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type cloudCredentialLifecycleConverter struct {
	lifecycle CloudCredentialLifecycle
}

func (w *cloudCredentialLifecycleConverter) CreateContext(_ context.Context, obj *v3.CloudCredential) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *cloudCredentialLifecycleConverter) RemoveContext(_ context.Context, obj *v3.CloudCredential) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *cloudCredentialLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.CloudCredential) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type CloudCredentialLifecycle interface {
	Create(obj *v3.CloudCredential) (runtime.Object, error)
	Remove(obj *v3.CloudCredential) (runtime.Object, error)
	Updated(obj *v3.CloudCredential) (runtime.Object, error)
}

type CloudCredentialLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.CloudCredential) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.CloudCredential) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.CloudCredential) (runtime.Object, error)
}

type cloudCredentialLifecycleAdapter struct {
	lifecycle CloudCredentialLifecycleContext
}

func (w *cloudCredentialLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *cloudCredentialLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *cloudCredentialLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *cloudCredentialLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.CloudCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *cloudCredentialLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *cloudCredentialLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.CloudCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *cloudCredentialLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *cloudCredentialLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.CloudCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewCloudCredentialLifecycleAdapter(name string, clusterScoped bool, client CloudCredentialInterface, l CloudCredentialLifecycle) CloudCredentialHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(CloudCredentialGroupVersionResource)
	}
	adapter := &cloudCredentialLifecycleAdapter{lifecycle: &cloudCredentialLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.CloudCredential) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewCloudCredentialLifecycleAdapterContext(name string, clusterScoped bool, client CloudCredentialInterface, l CloudCredentialLifecycleContext) CloudCredentialHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(CloudCredentialGroupVersionResource)
	}
	adapter := &cloudCredentialLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.CloudCredential) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
