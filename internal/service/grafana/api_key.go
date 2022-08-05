package grafana

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/managedgrafana"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiKeyCreate,
		Read:   resourceApiKeyRead,
		Delete: resourceApiKeyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"seconds_to_live": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GrafanaConn

	input := &managedgrafana.CreateWorkspaceApiKeyInput{
		KeyName:       aws.String(d.Get("key_name").(string)),
		KeyRole:       aws.String(d.Get("key_role").(string)),
		SecondsToLive: aws.Int64(int64(d.Get("seconds_to_live").(int))),
		WorkspaceId:   aws.String(d.Get("workspace_id").(string)),
	}

	log.Printf("[DEBUG] Creating Grafana API Key: %s", input)
	output, err := conn.CreateWorkspaceApiKey(input)

	if err != nil {
		return fmt.Errorf("error creating Grafana API Key: %w", err)
	}

	d.Set("key", output.Key)

	//if _, err := waitApiKeyCreated(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
	//return fmt.Errorf("error waiting for Grafana Api Key (%s) create: %w", d.Id(), err)
	//}

	return resourceApiKeyRead(d, meta)
}

func resourceApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceApiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GrafanaConn

	log.Printf("[DEBUG] Deleting Grafana Api Key: %s", d.Id())
	_, err := conn.DeleteWorkspaceApiKey(&managedgrafana.DeleteWorkspaceApiKeyInput{
		KeyName:     aws.String(d.Get("key_name").(string)),
		WorkspaceId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, managedgrafana.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Grafana Api Key (%s): %w", d.Id(), err)
	}

	if _, err := waitWorkspaceUpdated(conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for Grafana Api Key (%s) delete: %w", d.Id(), err)
	}

	return nil
}
