package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
	controllers "github.com/flux-framework/flux-operator/controllers/flux"
	"github.com/flux-framework/flux-operator/pkg/flux"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// go run cmd/gen/gen.go -f examples/tests/lammps/minicluster.yaml

var (
	l         = log.New(os.Stderr, "🟧️fluxoperator: ", log.Ldate|log.Ltime|log.Lshortfile)
	separator = "----"
	filename  string
	includes  string
	scheme    = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(api.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(batchv1.AddToScheme(scheme))
}

type IncludePreference struct {

	// By default, we generate all
	GenConfigs bool
	GenService bool
	GenJob     bool
}

// determineIncludes determines what the user wants to print
func determineIncludes(includes string) IncludePreference {

	// By default, we print all!
	prefs := IncludePreference{true, true, true}
	if includes == "" {
		return prefs
	}
	// Config Maps
	if !strings.Contains(includes, "c") {
		prefs.GenConfigs = false
	}
	// Service
	if !strings.Contains(includes, "s") {
		prefs.GenService = false
	}
	// Job
	if !strings.Contains(includes, "j") {
		prefs.GenJob = false
	}
	return prefs
}

func main() {
	flag.StringVar(&filename, "f", "", "YAML filename to read it")
	flag.StringVar(&includes, "i", "", "Custom list of includes (csjv) for cm, svc, job, volume")
	flag.Parse()

	if filename == "" {
		l.Fatalf("Please provide a yaml file with -f")
	}

	// Unless the user has asked
	decode := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode
	stream, err := os.ReadFile(filename)
	if err != nil {
		l.Fatalf("Issue reading file %s", err.Error())
	}

	obj, _, err := decode(stream, nil, nil)
	if err != nil {
		l.Fatalf("Issue decoding yaml stream %s", err.Error())
	}

	// Determine what the person wants to generate
	prefs := determineIncludes(includes)

	switch cluster := obj.(type) {
	case *api.MiniCluster:

		// Generate the MiniCluster assets (config maps and indexed job)
		if prefs.GenConfigs {
			cms, err := generateMiniClusterConfigs(cluster)
			if err != nil {
				l.Fatalf("Issue generating MiniCluster configs %s", err.Error())
			}
			for _, cm := range cms {
				printYaml(cm)
			}
		}
		// We no longer generate volumes - up to the users
		// Service
		if prefs.GenService {

			svc := generateMiniClusterService(cluster)
			printYaml(svc)
		}
		// Finally the indexed job
		if prefs.GenJob {
			job, err := controllers.NewMiniClusterJob(cluster)
			if err != nil {
				l.Fatalf("Issue generating MiniCluster job %s", err.Error())
			}
			if err != nil {
				l.Fatalf("Issue generating MiniCluster job YAML %s", err.Error())
			}
			printYaml(job)
		}
	default:
		l.Fatalf("Type %s is not a MiniCluster", cluster)
	}
}

func printYaml(obj runtime.Object) {

	out, err := yaml.Marshal(obj)
	if err != nil {
		l.Fatalf("Issue generating MiniCluster job YAML %s", err.Error())
	}
	fmt.Println(separator)
	fmt.Println(string(out))
}

// Generate the service for the MiniCluster
func generateMiniClusterService(cluster *api.MiniCluster) *corev1.Service {

	// Create headless service for the MiniCluster OR single service for the broker
	selector := map[string]string{"job-name": cluster.Name}

	// If we are adding a minimal service to the index 0 pod only (not supported here)
	if cluster.Spec.Flux.MinimalService {
		l.Fatalf("We do not currently support minimal service for this command.")
	}

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Spec.Network.HeadlessName,
			Namespace: cluster.Namespace,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  selector,
		},
	}
}

// generateMinicluster creates complete Minicluster yaml (Jobset and config maps)
func generateMiniClusterConfigs(cluster *api.MiniCluster) ([]*corev1.ConfigMap, error) {

	// We will return a set of config maps and a jobset
	cms := []*corev1.ConfigMap{}

	// Generate entrypoints first
	data, err := flux.GenerateEntrypoints(cluster)
	if err != nil {
		return cms, err
	}
	cm := createConfigMap(cluster, cluster.EntrypointConfigMapName(), data)
	cms = append(cms, cm)
	return cms, nil
}

// createConfigMap generates a config map with some kind of data
func createConfigMap(
	cluster *api.MiniCluster,
	configName string,
	data map[string]string,
) *corev1.ConfigMap {

	// Create the config map with respective data!
	// Likely we shouldn't hard code this :)
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: cluster.Namespace,
		},
		Data: data,
	}
	return cm
}
