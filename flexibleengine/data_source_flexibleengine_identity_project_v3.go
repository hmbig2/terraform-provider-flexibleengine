package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIdentityProjectV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityProjectV3Read,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"is_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// dataSourceIdentityProjectV3Read performs the project lookup.
func dataSourceIdentityProjectV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	listOpts := projects.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Name:     d.Get("name").(string),
		ParentID: d.Get("parent_id").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var project projects.Project
	allPages, err := projects.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query projects: %s", err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve projects: %s", err)
	}

	if len(allProjects) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allProjects) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allProjects)
		return fmt.Errorf("Your query returned more than one result")
	}
	project = allProjects[0]

	log.Printf("[DEBUG] Single project found: %s", project.ID)
	return dataSourceIdentityProjectV3Attributes(d, &project)
}

// dataSourceIdentityProjectV3Attributes populates the fields of an Project resource.
func dataSourceIdentityProjectV3Attributes(d *schema.ResourceData, project *projects.Project) error {
	log.Printf("[DEBUG] flexibleengine_identity_project_v3 details: %#v", project)

	d.SetId(project.ID)
	d.Set("is_domain", project.IsDomain)
	d.Set("description", project.Description)
	d.Set("domain_id", project.DomainID)
	d.Set("enabled", project.Enabled)
	d.Set("name", project.Name)
	d.Set("parent_id", project.ParentID)

	return nil
}
