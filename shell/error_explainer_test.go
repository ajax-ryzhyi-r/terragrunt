package shell_test

import (
	"bytes"
	"testing"

	"github.com/gruntwork-io/terragrunt/internal/errors"
	"github.com/gruntwork-io/terragrunt/shell"
	"github.com/gruntwork-io/terragrunt/util"
	"github.com/stretchr/testify/assert"
)

func TestExplainError(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		errorOutput string
		explanation string
	}{
		{
			errorOutput: "Error refreshing state: AccessDenied: Access Denied",
			explanation: "Check your credentials and permissions",
		},
		{
			errorOutput: "Error: Initialization required",
			explanation: "You need to run terragrunt (run --all) init to initialize working directory",
		},
		{
			errorOutput: "Unit source has changed",
			explanation: "You need to run terragrunt (run --all) init install all required modules",
		},
		{
			errorOutput: "Error: Failed to get existing workspaces: Unable to list objects in S3 bucket \"mybucket\": operation error S3: ListObjectsV2, https response error StatusCode: 301, RequestID: GH67DSB7KB8H578N, HostID: vofohiXBwNhR8Im+Dj7RpUPCPnOq9IDfn1rsUHHCzN9HgVMFfuIH5epndgLQvDeJPz2DrlUh0tA=, requested bucket from \"us-east-1\", actual location \"eu-west-1\"\n",
			explanation: "You don't have access to the S3 bucket where the state is stored. Check your credentials and permissions.",
		},
		{
			errorOutput: "exec: \"tofu\": executable file not found in $PATH",
			explanation: "The executables 'terraform' and 'tofu' are missing from your $PATH. Please add at least one of these to your $PATH.",
		},
	}

	for _, tt := range testCases {

		t.Run(tt.errorOutput, func(t *testing.T) {
			t.Parallel()

			output := util.CmdOutput{}
			output.Stderr = *bytes.NewBufferString(tt.errorOutput)

			errs := new(errors.MultiError)
			errs = errs.Append(util.ProcessExecutionError{
				Err:    errors.New(""),
				Output: output,
			})
			explanation := shell.ExplainError(errs)
			assert.Contains(t, explanation, tt.explanation)
		})
	}
}
