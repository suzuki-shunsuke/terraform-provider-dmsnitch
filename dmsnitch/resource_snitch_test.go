package dmsnitch_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/plukevdh/terraform-provider-dmsnitch/dmsnitch"
	"github.com/suzuki-shunsuke/flute/v2/flute"
)

func setEnv() error {
	envs := map[string]string{
		"DMS_TOKEN": "dms-token",
	}
	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

func testHeader() http.Header {
	return http.Header{
		"Content-Type": []string{"application/json"},
	}
}

func setHTTPClient(t *testing.T, httpClient *http.Client, routes ...flute.Route) {
	transport := flute.Transport{
		T: t,
		Services: []flute.Service{
			{
				Endpoint: "https://api.deadmanssnitch.com",
				Routes:   routes,
			},
		},
	}

	httpClient.Transport = transport
}

func TestAccSnitch(t *testing.T) { //nolint:funlen
	if err := setEnv(); err != nil {
		t.Fatal(err)
	}

	remoteBody := ""

	getRoute := flute.Route{
		Name: "get a snitch",
		Matcher: flute.Matcher{
			Method: "GET",
		},
		Tester: flute.Tester{
			Path:         "/v1/snitches/xxx",
			PartOfHeader: testHeader(),
		},
		Response: flute.Response{
			Response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(remoteBody)),
				}, nil
			},
		},
	}

	postRoute := flute.Route{
		Name: "create a snitch",
		Matcher: flute.Matcher{
			Method: "POST",
		},
		Tester: flute.Tester{
			Path:         "/v1/snitches",
			PartOfHeader: testHeader(),
			BodyJSONString: `{
  "name":"My Important Service",
	"notes": "Description or other notes about this snitch.",
  "interval": "daily",
	"alert_type": "basic",
	"tags": ["one", "two"]
}`,
			Test: func(t *testing.T, req *http.Request, svc flute.Service, route flute.Route) {
				remoteBody = `{
  "token": "xxx",
  "href": "/v1/snitches/xxx",
  "name":"My Important Service",
	"notes": "Description or other notes about this snitch.",
  "interval": "daily",
	"alert_type": "basic",
	"tags": ["one", "two"],
  "status": "pending",
  "checked_in_at": null,
  "alert_type": "basic",
  "alert_email": [],
  "check_in_url": "https://nosnch.in/xxx",
  "created_at": "2014-04-02T15:54:54.784Z"
}`
			},
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 201,
			},
			BodyString: `{
  "token": "xxx",
  "href": "/v1/snitches/xxx",
  "name":"My Important Service",
	"notes": "Description or other notes about this snitch.",
  "interval": "daily",
	"alert_type": "basic",
	"tags": ["one", "two"],
  "status": "pending",
  "checked_in_at": null,
  "alert_type": "basic",
  "alert_email": [],
  "check_in_url": "https://nosnch.in/xxx",
  "created_at": "2014-04-02T15:54:54.784Z"
}`,
		},
	}

	httpClient := &http.Client{}

	createStep := resource.TestStep{
		ResourceName: "dmsnitch_snitch.test",
		PreConfig: func() {
			setHTTPClient(t, httpClient, getRoute, postRoute)
		},
		Config: `
resource "dmsnitch_snitch" "test" {
	name = "My Important Service"
  notes = "Description or other notes about this snitch."
  
  interval = "daily" 
  type = "basic"
  tags = ["one", "two"]
}
`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "name", "My Important Service"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "notes", "Description or other notes about this snitch."),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "interval", "daily"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "type", "basic"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "token", "xxx"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "url", "https://nosnch.in/xxx"),
		),
	}

	updateRoute := flute.Route{
		Name: "update a snitch",
		Matcher: flute.Matcher{
			Method: "PATCH",
		},
		Tester: flute.Tester{
			Path:         "/v1/snitches/xxx",
			PartOfHeader: testHeader(),
			BodyJSONString: `{
  "name":"My Important Service",
	"notes": "updated notes",
  "interval": "daily",
	"alert_type": "basic",
  "check_in_url": "https://nosnch.in/xxx",
  "token": "xxx",
	"tags": ["one"]
}`,
			Test: func(t *testing.T, req *http.Request, svc flute.Service, route flute.Route) {
				remoteBody = `{
  "token": "xxx",
  "href": "/v1/snitches/xxx",
  "name":"My Important Service",
	"notes": "updated notes",
  "interval": "daily",
	"alert_type": "basic",
	"tags": ["one"],
  "status": "pending",
  "checked_in_at": null,
  "alert_type": "basic",
  "alert_email": [],
  "check_in_url": "https://nosnch.in/xxx",
  "created_at": "2014-04-02T15:54:54.784Z"
}`
			},
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 204,
			},
		},
	}

	deleteRoute := flute.Route{
		Name: "delete a snitch",
		Matcher: flute.Matcher{
			Method: "DELETE",
		},
		Tester: flute.Tester{
			Path:         "/v1/snitches/xxx",
			PartOfHeader: testHeader(),
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 204,
			},
		},
	}

	updateStep := resource.TestStep{
		ResourceName: "dmsnitch_snitch.test",
		PreConfig: func() {
			setHTTPClient(t, httpClient, getRoute, updateRoute, deleteRoute)
		},
		Config: `
resource "dmsnitch_snitch" "test" {
	name = "My Important Service"
  notes = "updated notes"
  
  interval = "daily" 
  type = "basic"
  tags = ["one"]
}
`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "name", "My Important Service"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "notes", "updated notes"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "interval", "daily"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "type", "basic"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "token", "xxx"),
			resource.TestCheckResourceAttr("dmsnitch_snitch.test", "url", "https://nosnch.in/xxx"),
		),
	}

	provider := dmsnitch.Provider().(*schema.Provider)
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		client := dmsnitch.NewClient(d.Get("api_key").(string))
		client.Client.HTTPClient = httpClient
		return client, nil
	}

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"dmsnitch": provider,
		},
		Steps: []resource.TestStep{
			createStep,
			updateStep,
		},
	})
}
