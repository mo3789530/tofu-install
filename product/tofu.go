package product

import (
	"context"
	"encoding/json"
	"os/exec"
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/mo3789530/hc-install/internal/build"
)

var Tofu = Product{
	Name: "opentofu",
	BinaryName: func() string {
		if runtime.GOOS == "windows" {
			return "opentofu.exe"
		}
		return "opentofu"
	},
	GetVersion: func(ctx context.Context, path string) (*version.Version, error) {
		v, err := tofuJsonVersion(ctx, path)
		if err == nil {
			return v, nil
		}
		return nil, err

	},
	BuildInstructions: &BuildInstructions{
		GitRepoURL:    "https://github.com/opentofu/opentofu.git",
		PreCloneCheck: &build.GoIsInstalled{},
		Build:         &build.GoBuild{DetectVendoring: true},
	},
}

type tofuJsonVersionOutput struct {
	// open tofu is stall use `terraform_version`
	Version *version.Version `json:"terraform_version"`
}

func tofuJsonVersion(ctx context.Context, path string) (*version.Version, error) {
	cmd := exec.CommandContext(ctx, path, "version", "-json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var vOut tofuJsonVersionOutput
	err = json.Unmarshal(out, &vOut)
	if err != nil {
		return nil, err
	}

	return vOut.Version, nil

}
