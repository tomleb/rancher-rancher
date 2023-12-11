package v3

import (
	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type EncryptedTokenLifecycle interface {
	Create(obj *v3.EncryptedToken) (runtime.Object, error)
	Remove(obj *v3.EncryptedToken) (runtime.Object, error)
	Updated(obj *v3.EncryptedToken) (runtime.Object, error)
}

type encryptedTokenLifecycleAdapter struct {
	lifecycle EncryptedTokenLifecycle
}

func (w *encryptedTokenLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *encryptedTokenLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *encryptedTokenLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v3.EncryptedToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *encryptedTokenLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v3.EncryptedToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *encryptedTokenLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v3.EncryptedToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewEncryptedTokenLifecycleAdapter(name string, clusterScoped bool, client EncryptedTokenInterface, l EncryptedTokenLifecycle) EncryptedTokenHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(EncryptedTokenGroupVersionResource)
	}
	adapter := &encryptedTokenLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.EncryptedToken) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
