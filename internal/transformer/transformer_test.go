package transformer

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/sharadregoti/devops/internal/transformer/testdata"
	"github.com/sharadregoti/devops/model"
)

func TestGetResourceInTableFormat(t *testing.T) {

	data := map[string]interface{}{}
	_ = json.Unmarshal([]byte(testdata.PodJSON), &data)

	type args struct {
		t         *model.ResourceTransfomer
		resources []interface{}
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "pass",
			args: args{
				t: &model.ResourceTransfomer{
					Operations: []model.Operations{
						{
							Name: "namespace",
							JSONPaths: []model.JSONPaths{
								{
									Path: "metadata.namespace",
								},
							},
							OutputFormat: "",
						},
						{
							Name: "name",
							JSONPaths: []model.JSONPaths{
								{
									Path: "metadata.name",
								},
							},
							OutputFormat: "",
						},
						{
							Name: "ready",
							JSONPaths: []model.JSONPaths{
								{
									Path: "status.containerStatuses.#(ready==true)#|#",
								},
								{
									Path: "status.containerStatuses.#",
								},
							},
							OutputFormat: "%v/%v",
						},
						{
							Name: "restarts",
							JSONPaths: []model.JSONPaths{
								{
									Path: "status.containerStatuses.0.restartCount",
								},
							},
							OutputFormat: "",
						},
						{
							Name: "status",
							JSONPaths: []model.JSONPaths{
								{
									Path: "status.phase",
								},
							},
							OutputFormat: "",
						},
						{
							Name: "ip",
							JSONPaths: []model.JSONPaths{
								{
									Path: "status.podIP",
								},
							},
							OutputFormat: "",
						},
						{
							Name: "node",
							JSONPaths: []model.JSONPaths{
								{
									Path: "spec.nodeName",
								},
							},
							OutputFormat: "",
						},
						{
							Name: "age",
							JSONPaths: []model.JSONPaths{
								{
									Path: "status.startTime|@age",
								},
							},
							OutputFormat: "",
						},
					},
				},
				resources: []interface{}{
					data,
				},
			},
			want: [][]string{
				{
					"NAMESPACE", "NAME", "READY", "RESTARTS", "STATUS", "IP", "NODE", "AGE",
				},
				{
					"default", "httpbin-74fb669cc6-8g9vf", "2/2", "33", "Running", "10.1.17.12", "docker-desktop", "145d",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetResourceInTableFormat(tt.args.t, tt.args.resources)
			if arr := deep.Equal(got, tt.want); len(arr) > 0 {
				t.Errorf("GetResourceInTableFormat() = %v", arr)
			}
		})
	}
}
