package e2e_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"flag"
	"testing"
)

// Test suite required arguments
var (
	KubectlPath           string
	ClusterctlPath        string
	DumpPath              string
	WorkingDir            string
	TenantKubeConfig      string
	InfraClusterNamespace string
	InfraKubeConfig       string
)

// Initialize test required arguments
func init() {
	flag.StringVar(&KubectlPath, "kubectl-path", "", "Path to the kubectl binary")
	flag.StringVar(&ClusterctlPath, "clusterctl-path", "", "Path to the clusterctl binary")
	flag.StringVar(&DumpPath, "dump-path", "", "Path to the kubevirt artifacts dump cmd binary")
	flag.StringVar(&WorkingDir, "working-dir", "", "Path used for e2e test files")
	flag.StringVar(&TenantKubeConfig, "tenant-kubeconfig", "", "Path to tenant kubeconfig")
	flag.StringVar(&InfraKubeConfig, "infra-kubeconfig", "", "Path to infra kubeconfig")
	flag.StringVar(&InfraClusterNamespace, "infra-cluster-namespace", "kv-guest-cluster", "Namespace of the guest cluster in the infra cluster")
}

func TestE2E(t *testing.T) {
	if len(TenantKubeConfig) == 0 {
		// Make sure that valid arguments have been passed for this test suite run.
		if KubectlPath == "" {
			t.Fatal("kubectl-path or tenant-kubeconfig required")
		} else if _, err := os.Stat(KubectlPath); os.IsNotExist(err) {
			t.Fatalf("invalid kubectl-path path: %s doesn't exist", KubectlPath)
		}
		if ClusterctlPath == "" {
			t.Fatal("clusterctl-path required")
		} else if _, err := os.Stat(ClusterctlPath); os.IsNotExist(err) {
			t.Fatalf("invalid clusterctl-path path: %s doesn't exist", ClusterctlPath)
		}
		if WorkingDir == "" {
			t.Fatal("working-dir required")
		} else if _, err := os.Stat(WorkingDir); os.IsNotExist(err) {
			t.Fatalf("invalid working-dir path: %s doesn't exist", WorkingDir)
		}
	}
	if DumpPath != "" {
		if _, err := os.Stat(DumpPath); os.IsNotExist(err) {
			t.Fatalf("invalid dump-path: %s doesn't exist", DumpPath)
		}
	}

	if len(InfraKubeConfig) == 0 {
		t.Fatal("infra kubeconfig required")
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
	// parse test suite arguments
	flag.Parse()
})

var _ = JustAfterEach(func() {
	if CurrentSpecReport().Failed() && DumpPath != "" {
		dump(os.Getenv("KUBECONFIG"), "")
	}
})

func dump(kubeconfig, artifactsSuffix string) {
	cmd := exec.Command(DumpPath, "--kubeconfig", kubeconfig)

	failureLocation := CurrentSpecReport().Failure.Location
	artifactsPath := filepath.Join(os.Getenv("ARTIFACTS"), fmt.Sprintf("%s:%d", filepath.Base(failureLocation.FileName), failureLocation.LineNumber), artifactsSuffix)
	cmd.Env = append(cmd.Env, fmt.Sprintf("ARTIFACTS=%s", artifactsPath))

	By(fmt.Sprintf("dumping k8s artifacts to %s", artifactsPath))
	output, err := cmd.CombinedOutput()
	Expect(err).ToNot(HaveOccurred(), string(output))
}
