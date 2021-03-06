package task

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestGitOutput(t *testing.T) {
	t.Skip("Functional testing harness. Needs to be completed.")

	opts := testutil.MockMetadataAPIOptions{
		SpecResponse: map[string]interface{}{
			"value": map[string]interface{}{
				"git": map[string]interface{}{
					"ssh_key":     "<ssh_key>",
					"known_hosts": "<known_hosts>",
					"name":        uuid.New().String(),
					"repository":  "<repository>",
				},
			},
		},
	}

	//nolint:gocritic
	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		// opts := taskutil.DefaultPlanOptions{
		// 	Client:  ts.Client(),
		// 	SpecURL: fmt.Sprintf("%s/spec", ts.URL),
		// }

		// task := NewTaskInterface(opts)

		// err := task.CloneRepository("master", "output")
		// require.Nil(t, err, "err is not nil")
	}, opts)
}

func TestGitSSHKeyBackwardCompatibility(t *testing.T) {
	tests := []struct {
		Spec           map[string]interface{}
		ExpectedSSHKey string
		ExpectedError  bool
	}{
		{
			Spec: map[string]interface{}{
				"ssh_key": "VEVTVCBLRVk=",
			},
			ExpectedSSHKey: "TEST KEY",
		},
		{
			Spec: map[string]interface{}{},
		},
		{
			Spec: map[string]interface{}{
				"ssh_key": "invalid",
			},
			ExpectedSSHKey: "invalid",
		},
		{
			Spec: map[string]interface{}{
				"connection": map[string]interface{}{
					"sshKey": "TEST KEY",
				},
			},
			ExpectedSSHKey: "TEST KEY",
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			opts := testutil.MockMetadataAPIOptions{
				SpecResponse: map[string]interface{}{
					"value": test.Spec,
				},
			}

			testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
				opts := taskutil.DefaultPlanOptions{
					Client:  ts.Client(),
					SpecURL: fmt.Sprintf("%s/spec", ts.URL),
				}

				var spec model.GitDetails
				require.NoError(t, taskutil.PopulateSpecFromDefaultPlan(&spec, opts))

				key, found, err := spec.ConfiguredSSHKey()
				if test.ExpectedError {
					require.Error(t, err)
				} else if test.ExpectedSSHKey != "" {
					require.True(t, found)
					require.Equal(t, test.ExpectedSSHKey, key)
				} else {
					require.False(t, found)
				}
			}, opts)
		})
	}
}

func TestGitSSHConfiguration(t *testing.T) {
	cases := []struct {
		gitURL      string
		host        string
		shouldMatch bool
	}{
		{gitURL: "git@github.com:example/example-repo", host: "github.com", shouldMatch: true},
		{gitURL: "git@github.com:example/example-repo.git", host: "github.com", shouldMatch: true},
		{gitURL: "git@github.com/example/example-repo.git", shouldMatch: false},
		{gitURL: "foobar@github.com:example/example-repo", host: "github.com", shouldMatch: true},
		{gitURL: "foobar@github.com:example/example-repo.git", host: "github.com", shouldMatch: true},
		{gitURL: "foobar@github.com/example/example-repo.git", shouldMatch: false},
		{gitURL: "git@github.com:example", shouldMatch: false},
		{gitURL: "git@github.com/example", shouldMatch: false},
		{gitURL: "git@example.com:example-org/example-repo", host: "example.com", shouldMatch: true},
		{gitURL: "", shouldMatch: false},
		{gitURL: "https://example.com", shouldMatch: false},
		{gitURL: "https://example.com/example", shouldMatch: false},
		{gitURL: "https://example.com:example", shouldMatch: false},
		{gitURL: "https://example.com:example/example-repo.git", shouldMatch: false},
		{gitURL: "https://example.com/example-org/example-repo", shouldMatch: false},
		{gitURL: "https://example.com/example-org/example-repo.git", shouldMatch: false},
	}

	for _, c := range cases {
		t.Run(c.gitURL, func(t *testing.T) {
			gd := &model.GitDetails{
				Repository: c.gitURL,
				SSHKey:     "<ssh_key>",
				KnownHosts: "<known_hosts>",
			}

			ssh, err := configuredSSH(gd)

			require.NoError(t, err)

			if !c.shouldMatch {
				require.Empty(t, ssh)
			} else {
				require.Equal(t, c.host, ssh.Host)
				require.Equal(t, gd.SSHKey, ssh.SSHKey)
				require.Equal(t, gd.KnownHosts, ssh.KnownHosts)
			}
		})
	}
}
