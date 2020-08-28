package pv

import (
	"k8s-pv-provisioner/cmd/provisioner/config"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

var predicates = []func(*v1.PersistentVolume) bool{
	IsProperProvisionedByAnnotation,
	HasProperReclaimPolicy,
}

func IsReleasedPV(pv *v1.PersistentVolume) bool {
	if pv.Status.Phase == v1.VolumeReleased {
		return true
	}

	klog.V(2).Infof("PersistentVolume: %v is not released yet", pv.Name)
	return false

}

func IsProperProvisionedByAnnotation(pv *v1.PersistentVolume) bool {
	storageClassName := pv.Spec.StorageClassName
	currentStorageClass := appConfig.StorageClasses[storageClassName]

	value, ok := pv.Annotations[config.AnnotationProvisionedBy]
	if ok && value == currentStorageClass.Provisioner {
		return true
	}

	klog.V(2).Infof("PersistentVolume: %v does not have annotation specifying onto the provisioner: %v", pv.Name, currentStorageClass.Provisioner)
	return false
}

func HasProperReclaimPolicy(pv *v1.PersistentVolume) bool {
	if pv.Spec.PersistentVolumeReclaimPolicy == v1.PersistentVolumeReclaimDelete {
		return true
	}

	klog.V(2).Infof("PersistentVolume: %v does not have proper reclaimPolicy", pv.Name)
	return false

}

func HasProperClassName(pv *v1.PersistentVolume) bool {
	storageClassName := pv.Spec.StorageClassName
	if _, present := appConfig.StorageClasses[storageClassName]; present {
		return true
	}
	return false
}

func allChecksPassed(predicates []func(*v1.PersistentVolume) bool, input *v1.PersistentVolume) bool {
	var result = true
	for _, predicate := range predicates {
		result = result && predicate(input)
	}

	return result
}