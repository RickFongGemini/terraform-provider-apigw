package apigw
  
import (
    "github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var osMutexKV = mutexkv.NewMutexKV()

type PConfig struct {
    Config
}

// Provider returns a schema.Provider for APIGW.
func Provider() terraform.ResourceProvider {
    provider := &schema.Provider{
        Schema: map[string]*schema.Schema{
            "apikey": {
                Type:		schema.TypeString,
                Required:	true,
                DefaultFunc:	schema.EnvDefaultFunc("APIGW_APIKEY", ""),
                Description:	descriptions["apikey"],
            },
            "apigw_url": {       
                Type:		schema.TypeString,
                Required:	true,
                DefaultFunc:	schema.EnvDefaultFunc("APIGW_URL", ""),
                Description:	descriptions["apigw_url"],
            },
        },

        DataSourcesMap: map[string]*schema.Resource{
            "apigw_network":			dataSourceNetwork(),
            "apigw_project":			dataSourceProject(),
            "apigw_solution":			dataSourceSolution(),
            "apigw_firewall":			dataSourceFirewall(),
            "apigw_firewall_rule":		dataSourceFirewallRule(),
            "apigw_vcs":			dataSourceVCS(),
            "apigw_volume":			dataSourceVolume(),
            "apigw_volume_snapshot":		dataSourceVolumeSnapshot(),
            "apigw_ike_policy":			dataSourceIKEPolicy(),
            "apigw_ipsec_policy":		dataSourceIPSecPolicy(),
            "apigw_vpn":			dataSourceVPN(),
            "apigw_container":			dataSourceContainer(),
            "apigw_s3_key":			dataSourceS3Key(),
            "apigw_security_group":		dataSourceSecurityGroup(),
            "apigw_loadbalancer":		dataSourceLoadBalancer(),
            "apigw_auto_scaling_policy":	dataSourceAutoScalingPolicy(),
            "apigw_extra_property":		dataSourceExtraProperty(),
            "apigw_waf":			dataSourceWAF(),
        },

        ResourcesMap: map[string]*schema.Resource{
            "apigw_auto_scaling_policy":	resourceAutoScalingPolicy(),
            "apigw_auto_scaling_relation":	resourceAutoScalingRelation(),
            "apigw_container":			resourceContainer(),
            "apigw_firewall":			resourceFirewall(),
            "apigw_firewall_rule":		resourceFirewallRule(),
            "apigw_ike_policy":			resourceIKEPolicy(),
            "apigw_ipsec_policy":		resourceIPSecPolicy(),
            "apigw_loadbalancer":		resourceLoadBalancer(),
            "apigw_network":			resourceNetwork(),
            "apigw_vcs":			resourceVCS(),
            "apigw_vcs_image":			resourceVCSImage(),
            "apigw_volume":			resourceVolume(),
            "apigw_volume_attachment":		resourceVolumeAttachment(),
            "apigw_volume_snapshot":		resourceVolumeSnapshot(),
            "apigw_vpn":			resourceVPN(),
            "apigw_vpn_connection":		resourceVPNConnection(),
            "apigw_s3_key":			resourceS3Key(),
            "apigw_security_group_rule":	resourceSecurityGroupRule(),
            "apigw_waf":			resourceWAF(),
        },
    }

    provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
        terraformVersion := provider.TerraformVersion
        if terraformVersion == "" {
            // Terraform 0.12 introduced this field to the protocol
            // We can therefore assume that if it's missing it's 0.10 or 0.11
            terraformVersion = "0.11+compatible"
        }
        return configureProvider(d, terraformVersion)
    }
    return provider
}

var descriptions map[string]string

func init() {
    descriptions = map[string]string{
        "apikey": "APIKey to login with.",
        "apigw_url": "APIGW endpoint to request to.",
    }
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
    config := PConfig{
        Config{
            APIGW_APIKEY:	d.Get("apikey").(string),
            APIGW_URL:		d.Get("apigw_url").(string),
        },
    }

    if err := config.LoadAndValidate(); err != nil {
        return nil, err
    }

    return &config, nil
}
