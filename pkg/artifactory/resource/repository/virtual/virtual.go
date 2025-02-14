package virtual

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jfrog/terraform-provider-artifactory/v6/pkg/artifactory/resource/repository"
	"github.com/jfrog/terraform-provider-shared/util"
	"github.com/jfrog/terraform-provider-shared/validator"
)

type VirtualRepositoryBaseParams struct {
	Key                                           string   `hcl:"key" json:"key,omitempty"`
	ProjectKey                                    string   `json:"projectKey"`
	ProjectEnvironments                           []string `json:"environments"`
	Rclass                                        string   `json:"rclass"`
	PackageType                                   string   `hcl:"package_type" json:"packageType,omitempty"`
	Description                                   string   `hcl:"description" json:"description,omitempty"`
	Notes                                         string   `hcl:"notes" json:"notes,omitempty"`
	IncludesPattern                               string   `hcl:"includes_pattern" json:"includesPattern,omitempty"`
	ExcludesPattern                               string   `hcl:"excludes_pattern" json:"excludesPattern,omitempty"`
	RepoLayoutRef                                 string   `hcl:"repo_layout_ref" json:"repoLayoutRef,omitempty"`
	Repositories                                  []string `hcl:"repositories" json:"repositories,omitempty"`
	ArtifactoryRequestsCanRetrieveRemoteArtifacts bool     `hcl:"artifactory_requests_can_retrieve_remote_artifacts" json:"artifactoryRequestsCanRetrieveRemoteArtifacts,omitempty"`
	DefaultDeploymentRepo                         string   `hcl:"default_deployment_repo" json:"defaultDeploymentRepo,omitempty"`
}

type VirtualRepositoryBaseParamsWithRetrievalCachePeriodSecs struct {
	VirtualRepositoryBaseParams
	VirtualRetrievalCachePeriodSecs int `hcl:"retrieval_cache_period_seconds" json:"virtualRetrievalCachePeriodSecs"`
}

func (bp VirtualRepositoryBaseParams) Id() string {
	return bp.Key
}

var VirtualRepoTypesLikeGeneric = []string{
	"docker",
	"gems",
	"generic",
	"gitlfs",
	"composer",
	"p2",
	"pub",
	"puppet",
	"pypi",
}

var VirtualRepoTypesLikeGenericWithRetrievalCachePeriodSecs = []string{
	"chef",
	"conan",
	"conda",
	"cran",
	"npm",
}

var BaseVirtualRepoSchema = map[string]*schema.Schema{
	"key": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The Repository Key. A mandatory identifier for the repository and must be unique. It cannot begin with a number or contain spaces or special characters. For local repositories, we recommend using a '-local' suffix (e.g. 'libs-release-local').",
	},
	"project_key": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validator.ProjectKey,
		Description:      "Project key for assigning this repository to. Must be 3 - 10 lowercase alphanumeric characters. When assigning repository to a project, repository key must be prefixed with project key, separated by a dash.",
	},
	"project_environments": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		MaxItems:    2,
		Set:         schema.HashString,
		Optional:    true,
		Description: `Project environment for assigning this repository to. Allow values: "DEV" or "PROD"`,
	},
	"package_type": {
		Type:        schema.TypeString,
		Required:    false,
		Computed:    true,
		ForceNew:    true,
		Description: "The Package Type. This must be specified when the repository is created, and once set, cannot be changed.",
	},
	"description": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A free text field that describes the content and purpose of the repository.\nIf you choose to insert a link into this field, clicking the link will prompt the user to confirm that they might be redirected to a new domain.",
	},
	"notes": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A free text field to add additional notes about the repository. These are only visible to the administrator.",
	},
	"includes_pattern": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "**/*",
		Description: "List of artifact patterns to include when evaluating artifact requests in the form of x/y/**/z/*. " +
			"When used, only artifacts matching one of the include patterns are served. By default, all artifacts are included (**/*).",
	},
	"excludes_pattern": {
		Type:     schema.TypeString,
		Optional: true,
		Description: "List of artifact patterns to exclude when evaluating artifact requests, in the form of x/y/**/z/*." +
			"By default no artifacts are excluded.",
	},
	"repo_layout_ref": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: repository.ValidateRepoLayoutRefSchemaOverride,
		Description:      "Sets the layout that the repository should use for storing and identifying modules. A recommended layout that corresponds to the package type defined is suggested, and index packages uploaded and calculate metadata accordingly.",
	},
	"repositories": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "The effective list of actual repositories included in this virtual repository.",
	},

	"artifactory_requests_can_retrieve_remote_artifacts": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether the virtual repository should search through remote repositories when trying to resolve an artifact requested by another Artifactory instance.",
	},
	"default_deployment_repo": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Default repository to deploy artifacts.",
	},
	"retrieval_cache_period_seconds": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      7200,
		Description:  "This value refers to the number of seconds to cache metadata files before checking for newer versions on aggregated repositories. A value of 0 indicates no caching.",
		ValidateFunc: validation.IntAtLeast(0),
	},
}

func UnpackBaseVirtRepo(s *schema.ResourceData, packageType string) VirtualRepositoryBaseParams {
	d := &util.ResourceData{s}

	return VirtualRepositoryBaseParams{
		Key:                 d.GetString("key", false),
		Rclass:              "virtual",
		ProjectKey:          d.GetString("project_key", false),
		ProjectEnvironments: d.GetSet("project_environments"),
		PackageType:         packageType, // must be set independently
		IncludesPattern:     d.GetString("includes_pattern", false),
		ExcludesPattern:     d.GetString("excludes_pattern", false),
		RepoLayoutRef:       d.GetString("repo_layout_ref", false),
		ArtifactoryRequestsCanRetrieveRemoteArtifacts: d.GetBool("artifactory_requests_can_retrieve_remote_artifacts", false),
		Repositories:          d.GetList("repositories"),
		Description:           d.GetString("description", false),
		Notes:                 d.GetString("notes", false),
		DefaultDeploymentRepo: repository.HandleResetWithNonExistantValue(d, "default_deployment_repo"),
	}
}

func UnpackBaseVirtRepoWithRetrievalCachePeriodSecs(s *schema.ResourceData, packageType string) VirtualRepositoryBaseParamsWithRetrievalCachePeriodSecs {
	d := &util.ResourceData{s}

	return VirtualRepositoryBaseParamsWithRetrievalCachePeriodSecs{
		VirtualRepositoryBaseParams:     UnpackBaseVirtRepo(s, packageType),
		VirtualRetrievalCachePeriodSecs: d.GetInt("retrieval_cache_period_seconds", false),
	}
}
