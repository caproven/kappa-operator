/*
Copyright 2021 caproven.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	caproveninfov1 "github.com/caproven/kappa-operator/api/v1"
)

// KappaReconciler reconciles a Kappa object
type KappaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=caproven.info,resources=kappas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=caproven.info,resources=kappas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=caproven.info,resources=kappas/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps;secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *KappaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	kappaLog := log.FromContext(ctx)

	kappaLog.Info("Reconciling Kappa")

	// Fetch the Kappa instance
	kappa := &caproveninfov1.Kappa{}
	err := r.Get(ctx, req.NamespacedName, kappa)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Return if the object is being deleted
	if kappa.GetDeletionTimestamp() != nil {
		return ctrl.Result{}, nil
	}

	if err := r.reconcileStatus(ctx, kappa); err != nil {
		return ctrl.Result{}, err
	}

	if err = r.reconcileConfigMap(ctx, kappa); err != nil {
		return ctrl.Result{}, err
	}

	if err = r.reconcileSecret(ctx, kappa); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *KappaReconciler) reconcileStatus(ctx context.Context, kappa *caproveninfov1.Kappa) error {
	kappa.Status.HasCucumber = true
	return r.Status().Update(ctx, kappa)
}

func (r *KappaReconciler) reconcileConfigMap(ctx context.Context, kappa *caproveninfov1.Kappa) error {
	kappaLog := log.FromContext(ctx).WithName("reconcileConfigMap")

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.Join([]string{kappa.Name, "config"}, "-"),
			Namespace: kappa.Namespace,
		},
		Data: map[string]string{
			"hello": SayHello(kappa.Name),
		},
	}
	ctrl.SetControllerReference(kappa, cm, r.Scheme)
	// cm.SetFinalizers([]string{"caproven.info/kappa-operator"})

	err := r.Get(ctx, client.ObjectKey{Name: cm.Name, Namespace: cm.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			// create
			err = r.Create(ctx, cm)
			if err != nil {
				return err
			}
			kappaLog.Info("Created ConfigMap")
			return nil
		}
		return err
	}
	// update
	err = r.Update(ctx, cm)
	if err != nil {
		return err
	}
	kappaLog.Info("Updated ConfigMap")
	return nil
}

func (r *KappaReconciler) reconcileSecret(ctx context.Context, kappa *caproveninfov1.Kappa) error {
	kappaLog := log.FromContext(ctx).WithName("reconcileSecret")

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.Join([]string{kappa.Name, "secret"}, "-"),
			Namespace: kappa.Namespace,
		},
		Data: map[string][]byte{
			"hello": []byte(SayHello(kappa.Name)),
		},
	}
	ctrl.SetControllerReference(kappa, secret, r.Scheme)

	if err := r.Create(ctx, secret); err != nil {
		if errors.IsAlreadyExists(err) {
			return nil
		}
		return err
	} else {
		kappaLog.Info("Created Secret")
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KappaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&caproveninfov1.Kappa{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func SayHello(name string) string {
	if name == "" {
		return "Hello!"
	}
	return fmt.Sprintf("Hello, %s", name)
}

func IsOwnedBy(obj, owner client.Object) bool {
	ownerRefs := obj.GetOwnerReferences()
	for _, ownerRef := range ownerRefs {
		if ownerRef.UID == owner.GetUID() {
			return true
		}
	}
	return false
}
