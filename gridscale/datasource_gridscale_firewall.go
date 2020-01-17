package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscaleFirewall() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleFirewallRead,

		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Computed:    true,
			},
			"rules_v4_in": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: getFirewallRuleCommonSchema(),
				},
			},
			"rules_v4_out": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: getFirewallRuleCommonSchema(),
				},
			},
			"rules_v6_in": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: getFirewallRuleCommonSchema(),
				},
			},
			"rules_v6_out": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: getFirewallRuleCommonSchema(),
				},
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status indicates the status of the object",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "The date and time the object was initially created.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "The date and time of the last object change.",
				Computed:    true,
			},
			"private": {
				Type:        schema.TypeBool,
				Description: "The object is private, the value will be true. Otherwise the value will be false.",
				Computed:    true,
			},
			"network": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the Firewall.",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGridscaleFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	fw, err := client.GetFirewall(emptyCtx, id)

	if err == nil {
		props := fw.Properties
		d.SetId(props.ObjectUUID)

		d.Set("name", props.Name)
		d.Set("location_name", props.LocationName)
		d.Set("status", props.Status)
		d.Set("private", props.Private)
		d.Set("create_time", props.CreateTime)
		d.Set("change_time", props.ChangeTime)
		d.Set("description", props.Description)

		//Get network relating to this firewall
		networks := make([]interface{}, 0)
		for _, value := range props.Relations.Networks {
			rule := map[string]interface{}{
				"network_uuid": value.NetworkUUID,
				"object_uuid":  value.ObjectUUID,
				"network_name": value.NetworkName,
				"object_name":  value.ObjectName,
				"create_time":  value.CreateTime.String(),
			}
			networks = append(networks, rule)
		}
		if err = d.Set("network", networks); err != nil {
			return fmt.Errorf("Error setting network: %v", err)
		}

		//Get rules_v4_in
		rulesV4In := convFirewallRuleSliceToInterfaceSlice(props.Rules.RulesV4In)
		if err = d.Set("rules_v4_in", rulesV4In); err != nil {
			return fmt.Errorf("Error setting rules_v4_in: %v", err)
		}

		//Get rules_v4_out
		rulesV4Out := convFirewallRuleSliceToInterfaceSlice(props.Rules.RulesV4Out)
		if err = d.Set("rules_v4_out", rulesV4Out); err != nil {
			return fmt.Errorf("Error setting rules_v4_out: %v", err)
		}

		//Get rules_v6_in
		rulesV6In := convFirewallRuleSliceToInterfaceSlice(props.Rules.RulesV6In)
		if err = d.Set("rules_v6_in", rulesV6In); err != nil {
			return fmt.Errorf("Error setting rules_v6_in: %v", err)
		}

		//Get rules_v6_out
		rulesV6Out := convFirewallRuleSliceToInterfaceSlice(props.Rules.RulesV6Out)
		if err = d.Set("rules_v6_out", rulesV6Out); err != nil {
			return fmt.Errorf("Error setting rules_v6_out: %v", err)
		}

		if err = d.Set("labels", props.Labels); err != nil {
			return fmt.Errorf("Error setting labels: %v", err)
		}
	}

	return err
}