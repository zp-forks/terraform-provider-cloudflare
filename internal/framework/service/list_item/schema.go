package list_item

import (
	"context"
	"fmt"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *ListItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: heredoc.Doc(`
			Provides individual list items (IPs, Redirects, ASNs, Hostnames) to be used in Edge Rules Engine
			across all zones within the same account.
	`),

		Attributes: map[string]schema.Attribute{
			consts.AccountIDSchemaKey: schema.StringAttribute{
				MarkdownDescription: consts.AccountIDSchemaDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"list_id": schema.StringAttribute{
				MarkdownDescription: "The list identifier to target for the resource.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The list item identifier.",
				Computed:            true,
			},
			"ip": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("IP address to include in the list. %s", utils.RenderMustProviderOnlyOneOfDocumentationValuesStringSlice([]string{"ip", "asn", "redirect", "hostname"})),
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("redirect"),
						path.MatchRelative().AtParent().AtName("asn"),
						path.MatchRelative().AtParent().AtName("hostname"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"asn": schema.Int64Attribute{
				MarkdownDescription: fmt.Sprintf("Autonomous system number to include in the list. %s", utils.RenderMustProviderOnlyOneOfDocumentationValuesStringSlice([]string{"ip", "asn", "redirect", "hostname"})),
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("redirect"),
						path.MatchRelative().AtParent().AtName("ip"),
						path.MatchRelative().AtParent().AtName("hostname"),
					),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "An optional comment for the item.",
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"hostname": schema.ListNestedBlock{
				MarkdownDescription: fmt.Sprintf("Hostname to store in the list. %s", utils.RenderMustProviderOnlyOneOfDocumentationValuesStringSlice([]string{"ip", "asn", "redirect", "hostname"})),
				Validators: []validator.List{
					listvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("redirect"),
						path.MatchRelative().AtParent().AtName("asn"),
						path.MatchRelative().AtParent().AtName("ip"),
					),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"url_hostname": schema.StringAttribute{
							MarkdownDescription: "The FQDN to match on.",
							Required:            true,
						},
					},
				},
			},
			"redirect": schema.ListNestedBlock{
				MarkdownDescription: fmt.Sprintf("Redirect configuration to store in the list. %s", utils.RenderMustProviderOnlyOneOfDocumentationValuesStringSlice([]string{"ip", "asn", "redirect", "hostname"})),
				Validators: []validator.List{
					listvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("asn"),
						path.MatchRelative().AtParent().AtName("hostname"),
						path.MatchRelative().AtParent().AtName("ip"),
					),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"source_url": schema.StringAttribute{
							MarkdownDescription: "The source url of the redirect.",
							Required:            true,
						},
						"target_url": schema.StringAttribute{
							MarkdownDescription: "The target url of the redirect.",
							Required:            true,
						},
						"include_subdomains": schema.BoolAttribute{
							MarkdownDescription: "Whether the redirect also matches subdomains of the source url.",
							Optional:            true,
						},
						"subpath_matching": schema.BoolAttribute{
							MarkdownDescription: "Whether the redirect also matches subpaths of the source url.",
							Optional:            true,
						},
						"status_code": schema.Int64Attribute{
							MarkdownDescription: "The status code to be used when redirecting a request.",
							Optional:            true,
						},
						"preserve_path_suffix": schema.BoolAttribute{
							MarkdownDescription: "Whether the redirect target url should keep the query string of the request's url.",
							Optional:            true,
						},
						"preserve_query_string": schema.BoolAttribute{
							MarkdownDescription: "Whether the redirect target url should keep the query string of the request's url.",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}
