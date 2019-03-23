package models

type MetadataModel struct {
	Guid string `json:"guid,omitempty"`
}

type AppModel struct {
	Metadata MetadataModel
}
type AppModelList struct {
	Resources []AppModel `json:"resources"`
}

type RouteMapping struct {
	AppGuid   string `json:"app_guid"`
	RouteGuid string `json:"route_guid"`
	AppPort   int   `json:"app_port"`
}

type Route struct {
	Metadata MetadataModel
	DomainGuid string `json:"domain_guid"`
	SpaceGuid  string `json:"space_guid"`
	Host       string `json:"host"`
}
type RouteList struct {
	Resources []Route `json:"resources"`

}

type Domain struct {
	Metadata MetadataModel
}
type DomainList struct {
	Resources []Domain `json:"resources"`
}