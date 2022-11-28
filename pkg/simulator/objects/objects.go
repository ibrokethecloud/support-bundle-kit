package objects

import (
	"bytes"
	"fmt"
	wranglerunstructured "github.com/rancher/wrangler/pkg/unstructured"
	"github.com/rancher/wrangler/pkg/yaml"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"path/filepath"
	"strings"
)

// GenerateClusterScopedRuntimeObjects will parse the yaml directory
// in the bundle directory and apply the cluster and namespaced objects
func GenerateClusterScopedRuntimeObjects(path string) (crd []runtime.Object, clusterObjs []runtime.Object, err error) {

	var crdList, noncrdList []string
	dir, err := filepath.Abs(filepath.Join(path, "yamls", "cluster"))
	if err != nil {
		return crd, clusterObjs, fmt.Errorf("error generating absolute path %v", err)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if strings.Contains(absPath, "apiextensions.k8s.io") {
				crdList = append(crdList, absPath)
			} else {
				noncrdList = append(noncrdList, absPath)
			}

		}
		return nil
	})

	if err != nil {
		return crd, clusterObjs, fmt.Errorf("error during dir walk %v", err)
	}

	// generate objects //
	for _, v := range crdList {
		obj, err := GenerateObjects(v)
		if err != nil {
			return crd, clusterObjs, err
		}
		crd = append(crd, obj...)
	}

	for _, v := range noncrdList {
		obj, err := GenerateObjects(v)
		if err != nil {
			return crd, clusterObjs, err
		}
		clusterObjs = append(clusterObjs, obj...)
	}

	return crd, clusterObjs, nil
}

// GenerateNamespacedRuntimeObjects will return a map[string][]runtime.Object.
// the map key is the namespace and the list of objects associated with this namespaced.
// Two maps to split workloads into pods and nonpod types as pods may have dependency on other objects like service accounts.
func GenerateNamespacedRuntimeObjects(path string) (nonpods []runtime.Object, pods []runtime.Object, err error) {

	var podList, nonPodList, eventsList []string
	dir, err := filepath.Abs(filepath.Join(path, "yamls", "namespaced"))
	if err != nil {
		return nonpods, pods, fmt.Errorf("error generating absolute path %v", err)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if strings.Contains(absPath, "pods.yaml") && !strings.Contains(absPath, "metrics.k8s.io") {
				podList = append(podList, absPath)
			} else if strings.Contains(absPath, "events") {
				eventsList = append(eventsList, absPath)
			} else {
				nonPodList = append(nonPodList, absPath)
			}

		}
		return nil
	})

	if err != nil {
		return nonpods, pods, fmt.Errorf("error during dir walk %v", err)
	}

	// append events to pods to ensure they are created once all pods have been setup
	podList = append(podList, eventsList...)

	// walk each list to get the runtime objects and populate the result
	// generate objects //
	for _, v := range podList {
		obj, err := GenerateObjects(v)
		if err != nil {
			return nonpods, pods, err
		}
		pods = append(pods, obj...)
	}

	for _, v := range nonPodList {
		obj, err := GenerateObjects(v)
		if err != nil {
			return nonpods, pods, err
		}
		nonpods = append(nonpods, obj...)
	}

	return nonpods, pods, err
}

func GenerateObjects(file string) (obj []runtime.Object, err error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return obj, err
	}

	obj, err = yaml.ToObjects(bytes.NewReader(content))
	return obj, err
}

func GenerateUnstructuredObjects(file string) (objs []*unstructured.Unstructured, err error) {
	runObjs, err := GenerateObjects(file)
	if err != nil {
		return objs, err
	}

	for _, runObj := range runObjs {
		obj, err := wranglerunstructured.ToUnstructured(runObj)
		if err != nil {
			return objs, err
		}
		objs = append(objs, obj)
	}

	return objs, err
}

// GenerateUnstructuredObjectsFromString is a helper used by tests to generated objects from embedded yamls available in variables
func GenerateUnstructuredObjectsFromString(contents string) (obj []*unstructured.Unstructured, err error) {
	tmpFile, err := ioutil.TempFile("/tmp", "obj_from_var")
	if err != nil {
		return nil, fmt.Errorf("error creating tmpfile: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write([]byte(contents))
	if err != nil {
		return nil, fmt.Errorf("error writing contents to tmpFile: %v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing file: %v", err)
	}
	objs, err := GenerateObjects(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error during object generation: %v", err)
	}
	for _, v := range objs {
		unstructObj, err := wranglerunstructured.ToUnstructured(v)
		if err != nil {
			return nil, fmt.Errorf("error converting object to unstructured obj: %v", err)
		}
		obj = append(obj, unstructObj)
	}

	return obj, nil
}
