/*
Copyright 2023.

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
	"errors"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kapps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields" // Required for Watching
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types" // Required for Watching
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	// "sigs.k8s.io/controller-runtime/pkg/builder" // Required for Watching
	"sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/handler"   // Required for Watching
	// "sigs.k8s.io/controller-runtime/pkg/predicate" // Required for Watching
	"sigs.k8s.io/controller-runtime/pkg/reconcile" // Required for Watching
	// "sigs.k8s.io/controller-runtime/pkg/source"    // Required for Watching

	gethtestv1 "github.com/michaelgugino/gethtest/operator/api/v1"
)

// RacecourseReconciler reconciles a Racecourse object
type RacecourseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

/*
Determine the path of the field in the ConfigDeployment CRD that we wish to use as the "object reference".
This will be used in both the indexing and watching.
*/
const (
	racecourseField = ".spec.deploymentName"
)

//+kubebuilder:rbac:groups=gethtest.michaelgugino.com,resources=racecourses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gethtest.michaelgugino.com,resources=racecourses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gethtest.michaelgugino.com,resources=racecourses/finalizers,verbs=update

//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *RacecourseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Error(errors.New("testing"), "reconciled")
	// TODO(user): your logic here
	rc := &gethtestv1.Racecourse{}
	if err := r.Get(ctx, req.NamespacedName, rc); err != nil {
		log.Error(err, "unable to fetch Racecourse")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// name of our custom finalizer
	myFinalizerName := "racecourse.gethtest.michaelgugino.com/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if rc.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(rc, myFinalizerName) {
			controllerutil.AddFinalizer(rc, myFinalizerName)
			if err := r.Update(ctx, rc); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(rc, myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteApp(rc); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(rc, myFinalizerName)
			if err := r.Update(ctx, rc); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	if rc.Spec.DeploymentName == "" {
		rc.Spec.DeploymentName = rc.ObjectMeta.Name
		if err := r.Update(ctx, rc); err != nil {
			return ctrl.Result{}, err
		}
	}

	if err := r.createApp(rc); err != nil {
		return ctrl.Result{}, err
	}

	// Add some status logic here
	// Check deployment rollout status here

	return ctrl.Result{}, nil
}

func (r *RacecourseReconciler) createApp(rc *gethtestv1.Racecourse) error {
	replicas := int32(1)
	dep := &kapps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rc.Spec.DeploymentName,
			Namespace: rc.ObjectMeta.Namespace,
		},
		Spec: kapps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "racecourse-" + rc.Spec.DeploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"k8s-app": "racecourse-" + rc.Spec.DeploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "racecourse-app",
							Image: "quay.io/michaelgugino/gethtest:racecourse",
							//' Command:   []string{"npm"},
							// Args:      args,
							Ports: []corev1.ContainerPort{
								{
									Name:          "raceapp",
									ContainerPort: 3000,
								},
							},
							/*
								ReadinessProbe: &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/healthz",
											Port: intstr.Parse("app"),
										},
									},
								},
								LivenessProbe: &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/readyz",
											Port: intstr.Parse("app"),
										},
									},
								},
							*/
						},
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(rc, dep, r.Scheme); err != nil {
		return err
	}

	depobjkey := client.ObjectKey{
		Name:      rc.Spec.DeploymentName,
		Namespace: rc.ObjectMeta.Namespace,
	}
	if err := createIfNotPresent(r.Client, depobjkey, dep); err != nil {
		return err
	}

	if !dep.ObjectMeta.DeletionTimestamp.IsZero() {
		return errors.New("Deployment marked as deleted, will create another.")
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rc.Spec.DeploymentName,
			Namespace: rc.ObjectMeta.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "http",
					Port:       int32(3000),
					TargetPort: intstr.FromString("raceapp"),
				},
			},
			Selector: map[string]string{
				"k8s-app": "racecourse-" + rc.Spec.DeploymentName,
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	if err := controllerutil.SetControllerReference(rc, svc, r.Scheme); err != nil {
		return err
	}

	svcobjkey := client.ObjectKey{
		Name:      rc.Spec.DeploymentName,
		Namespace: rc.ObjectMeta.Namespace,
	}
	if err := createIfNotPresent(r.Client, svcobjkey, svc); err != nil {
		return err
	}

	if !svc.ObjectMeta.DeletionTimestamp.IsZero() {
		return errors.New("Service marked as deleted, will create another.")
	}

	return nil
}

func (r *RacecourseReconciler) deleteApp(_ *gethtestv1.Racecourse) error {

	return nil
}

func (r *RacecourseReconciler) findObjectsForRacecourse(dependency client.Object) []reconcile.Request {
	Racecourses := &gethtestv1.RacecourseList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(racecourseField, dependency.GetName()),
		Namespace:     dependency.GetNamespace(),
	}
	err := r.List(context.TODO(), Racecourses, listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(Racecourses.Items))
	for i, item := range Racecourses.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      item.GetName(),
				Namespace: item.GetNamespace(),
			},
		}
	}
	return requests
}

// SetupWithManager sets up the controller with the Manager.
func (r *RacecourseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gethtestv1.Racecourse{}).
		Owns(&kapps.Deployment{}).
		/* Watches(
			&source.Kind{Type: &kapps.Deployment{}},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsForRacecourse),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		). */
		Owns(&corev1.Service{}).
		/* Watches(
			&source.Kind{Type: &corev1.Service{}},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsForRacecourse),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		). */
		Complete(r)
}

func createIfNotPresent(kclient client.Client, objKey client.ObjectKey, obj client.Object) error {
	getErr := kclient.Get(context.Background(), objKey, obj)
	if getErr != nil {
		if apierrors.IsNotFound(getErr) {
			if createErr := kclient.Create(context.TODO(), obj); createErr != nil {
				return createErr
			}
		} else {
			return getErr
		}
	}
	return nil
}
