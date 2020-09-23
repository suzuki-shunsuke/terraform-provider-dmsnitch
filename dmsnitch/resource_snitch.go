package dmsnitch

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

type Snitch struct {
	Token    string   `json:"token,omitempty"`
	URL      string   `json:"check_in_url,omitempty"`
	Name     string   `json:"name,omitempty"`
	Status   string   `json:"status,omitempty"`
	Interval string   `json:"interval,omitempty"`
	Type     string   `json:"alert_type,omitempty"`
	Notes    string   `json:"notes,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func resourceSnitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceSnitchCreate,
		Update: resourceSnitchUpdate,
		Read:   resourceSnitchRead,
		Delete: resourceSnitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "basic",
			},

			"interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "daily",
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func newSnitchFromResource(d *schema.ResourceData) Snitch {
	tags := make([]string, 0, len(d.Get("tags").(*schema.Set).List()))

	for _, item := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, item.(string))
	}

	return Snitch{
		Name:     d.Get("name").(string),
		Token:    d.Get("token").(string),
		URL:      d.Get("url").(string),
		Status:   d.Get("status").(string),
		Interval: d.Get("interval").(string),
		Type:     d.Get("type").(string),
		Notes:    d.Get("notes").(string),
		Tags:     tags,
	}
}

func resourceSnitchCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)
	snitch := newSnitchFromResource(d)

	ctx := context.Background()
	respSnitch, _, err := client.Post(ctx, snitch) //nolint:bodyclose
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] ID received: %s", respSnitch.Token)
	d.SetId(respSnitch.Token)
	if err := d.Set("url", respSnitch.URL); err != nil {
		return err
	}
	if err := d.Set("token", respSnitch.Token); err != nil {
		return err
	}
	return nil
}

func resourceSnitchRead(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)
	ctx := context.Background()
	snitch, _, err := client.Get(ctx, d.Id()) //nolint:bodyclose
	if err != nil {
		return err
	}

	if err := d.Set("name", snitch.Name); err != nil {
		return err
	}
	if err := d.Set("token", snitch.Token); err != nil {
		return err
	}
	if err := d.Set("url", snitch.URL); err != nil {
		return err
	}
	if err := d.Set("status", snitch.Status); err != nil {
		return err
	}
	if err := d.Set("interval", snitch.Interval); err != nil {
		return err
	}
	if err := d.Set("type", snitch.Type); err != nil {
		return err
	}
	if err := d.Set("notes", snitch.Notes); err != nil {
		return err
	}
	if err := d.Set("tags", snitch.Tags); err != nil {
		return err
	}

	return nil
}

func resourceSnitchUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)
	snitch := newSnitchFromResource(d)

	ctx := context.Background()
	if _, err := client.Patch(ctx, d.Id(), snitch); err != nil { //nolint:bodyclose
		return err
	}

	return nil
}

func resourceSnitchDelete(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	client := m.(Client)
	_, err := client.Delete(ctx, d.Id()) //nolint:bodyclose
	return err
}
