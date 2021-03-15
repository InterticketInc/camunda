package deploy

type ListOptions struct {
	// id Filter by deployment id.
	Id string `url:"id,omitempty"`
	// Name Filter by the deployment name. Exact match.
	Name string `url:"name,omitempty"`
	// NameLike Filter by the deployment name that the parameter is a substring of. The
	//parameter can include the wildcard % to express like-strategy such as: starts with (%name), ends with (name%) or contains (%name%).
	NameLike string `url:"nameLike,omitempty"`
	// source Filter by the deployment source.
	Source string `url:"source,omitempty"`
	// WithoutSource Filter by the deployment source whereby source is equal to null.
	WithoutSource string `url:"withoutSource,omitempty"`
	// TenantIdIn Filter by a comma-separated list of tenant ids. A deployment must have
	// one of the given tenant ids.
	TenantIdIn string `url:"tenantIdIn,omitempty"`
	// WithoutTenantID Only include deployments which belong to no tenant. Value may only
	// be true, as false is the default behavior.
	WithoutTenantID bool `url:"withoutTenantId,omitempty"`
	// IncludeDeploymentsWithoutTenantID Include deployments which belong to no tenant. Can
	// be used in combination with tenantIdIn. Value may only be true, as false is the default behavior.
	IncludeDeploymentsWithoutTenantID bool `url:"includeDeploymentsWithoutTenantId,omitempty"`
	// After Restricts to all deployments after the given date. By default*, the date must
	// have the format yyyy-MM-dd'T'HH:mm:ss.SSSZ, e.g., 2013-01-23T14:42:45.000+0200.
	After string `url:"after,omitempty"`
	// Before Restricts to all deployments before the given date. By default*, the date must
	// have the format yyyy-MM-dd'T'HH:mm:ss.SSSZ, e.g., 2013-01-23T14:42:45.000+0200.
	Before string `url:"before,omitempty"`
	// SortBy Sort the results lexicographically by a given criterion. Valid values are id,
	// name, deploymentTime and tenantId. Must be used in conjunction with the sortOrder parameter.
	SortBy string `url:"sortBy,omitempty"`
	// SortOrder Sort the results in a given order. Values may be asc for ascending order or
	// desc for descending order. Must be used in conjunction with the sortBy parameter.
	SortOrder string `url:"sortOrder,omitempty"`
	// FirstResult Pagination of results. Specifies the index of the first result to return.
	FirstResult string `url:"firstResult,omitempty"`
	// MaxResults Pagination of results. Specifies the maximum number of results to return.
	// Will return less results if there are no more results left.
	MaxResults int `url:"maxResults,omitempty"`
}
