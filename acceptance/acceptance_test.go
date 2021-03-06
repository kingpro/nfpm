package acceptance

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goreleaser/nfpm"
	// shut up
	_ "github.com/goreleaser/nfpm/deb"
	_ "github.com/goreleaser/nfpm/rpm"
)

var formats = []string{"deb", "rpm"}

func TestSimple(t *testing.T) {
	for _, format := range formats {
		t.Run("amd64", func(t *testing.T) {
			t.Parallel()
			accept(t, fmt.Sprintf("simple_%s", format), "simple.yaml", format, fmt.Sprintf("%s.dockerfile", format))
		})
		t.Run("i386", func(t *testing.T) {
			t.Parallel()
			accept(t, fmt.Sprintf("simple_%s_386", format), "simple.386.yaml", format, fmt.Sprintf("%s.386.dockerfile", format))
		})
	}
}

func TestComplex(t *testing.T) {
	for _, format := range formats {
		t.Run("amd64", func(t *testing.T) {
			t.Parallel()
			accept(t, fmt.Sprintf("complex_%s", format), "complex.yaml", format, fmt.Sprintf("%s.complex.dockerfile", format))
		})
		t.Run("i386", func(t *testing.T) {
			t.Parallel()
			accept(t, fmt.Sprintf("complex_%s_386", format), "complex.386.yaml", format, fmt.Sprintf("%s.386.complex.dockerfile", format))
		})
	}
}

func TestComplexOverridesDeb(t *testing.T) {
	for _, format := range formats {
		t.Run("amd64", func(t *testing.T) {
			t.Parallel()
			accept(t, fmt.Sprintf("overrides_%s", format), "overrides.yaml", format, fmt.Sprintf("%s.overrides.dockerfile", format))
		})
	}
}

func TestMinDeb(t *testing.T) {
	for _, format := range formats {
		t.Run("amd64", func(t *testing.T) {
			t.Parallel()
			accept(t, fmt.Sprintf("min_%s", format), "min.yaml", format, fmt.Sprintf("%s.min.dockerfile", format))
		})
	}
}

func accept(t *testing.T, name, conf, format, dockerfile string) {
	var configFile = filepath.Join("./testdata", conf)
	tmp, err := filepath.Abs("./testdata/tmp")
	require.NoError(t, err)
	var packageName = name + "." + format
	var target = filepath.Join(tmp, packageName)

	require.NoError(t, os.MkdirAll(tmp, 0700))

	config, err := nfpm.ParseFile(configFile)
	require.NoError(t, err)

	info, err := config.Get(format)
	require.NoError(t, err)
	require.NoError(t, nfpm.Validate(info))

	pkg, err := nfpm.Get(format)
	require.NoError(t, err)

	f, err := os.Create(target)
	require.NoError(t, err)
	require.NoError(t, pkg.Package(nfpm.WithDefaults(info), f))
	bts, _ := exec.Command("pwd").CombinedOutput()
	t.Log(string(bts))
	cmd := exec.Command(
		"docker", "build", "--rm", "--force-rm",
		"-f", dockerfile,
		"--build-arg", "package="+filepath.Join("tmp", packageName),
		".",
	)
	cmd.Dir = "./testdata"
	t.Log("will exec:", cmd.Args)
	bts, err = cmd.CombinedOutput()
	t.Log("output:", string(bts))
	require.NoError(t, err)
}
