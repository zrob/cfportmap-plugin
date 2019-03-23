package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	. "github.com/zrob/cfportmap-plugin/util"
	. "github.com/zrob/cfportmap-plugin/models"
	"errors"
	"strconv"
)

type CFPortMapPlugin struct{}

func main() {
	plugin.Start(new(CFPortMapPlugin))
}

func (c *CFPortMapPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "map-route-port" {
		if len(args) != 5 {
			fmt.Println(c.GetMetadata().Commands[0].UsageDetails.Usage)
		} else {
			c.mapRoute(cliConnection, args)
		}
	}
}

func (c *CFPortMapPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "cfportmap",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "map-route-port",
				HelpText: "Map a route to a non-standard app port",
				UsageDetails: plugin.Usage{
					Usage: "cf map-route-port APP DOMAIN HOST PORT",
				},
			},
		},
	}
}

func (c *CFPortMapPlugin) mapRoute(cliConnection plugin.CliConnection, args []string) {
	app := args[1]
	domain := args[2]
	host := args[3]
	portString := args[4]

	port, err := strconv.Atoi(portString)
	FreakOut(err)

	fmt.Printf("Mapping route %s.%s to app %s on port %v...\r\n\r\n", host, domain, app, port)

	appGuid, err := getAppGuid(cliConnection, app)
	FreakOut(err)

	domainGuid, err := getDomainGuid(cliConnection, domain)
	FreakOut(err)

	route, err := createRoute(cliConnection, domainGuid, host)
	FreakOut(err)

	err = enableAppPort(cliConnection, appGuid, port)
	FreakOut(err)

	err = createMapping(cliConnection, route, appGuid, port)
	FreakOut(err)

	fmt.Println("OK")
}

func enableAppPort(cliConnection plugin.CliConnection, appGuid string, port int) (err error) {
	_, err = cliConnection.CliCommandWithoutTerminalOutput("curl", fmt.Sprintf("v2/apps/%s", appGuid),
		"-X", "PUT", "-d", fmt.Sprintf(`{"ports": [%v]}`, port))
	return
}

func createMapping(cliConnection plugin.CliConnection, route Route, appGuid string, port int) (err error) {
	mappingRequest := RouteMapping{
		AppGuid:   appGuid,
		RouteGuid: route.Metadata.Guid,
		AppPort:   port,
	}
	mappingBody, _ := json.Marshal(mappingRequest)

	_, err = cliConnection.CliCommandWithoutTerminalOutput("curl", "v2/route_mappings", "-X", "POST", "-d", string(mappingBody))

	return
}

func createRoute(cliConnection plugin.CliConnection, domainGuid string, host string) (route Route, err error) {
	mySpace, _ := cliConnection.GetCurrentSpace()

	output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", fmt.Sprintf("v2/routes?q=domain_guid:%s&host:%s", domainGuid, host))
	if (err != nil) {
		return
	}
	response := stringifyCurlResponse(output)

	routes := RouteList{}
	err = json.Unmarshal([]byte(response), &routes)
	if (err != nil) {
		return
	}

	if (len(routes.Resources) > 0) {
		route = routes.Resources[0]
		return
	}

	routeRequest := Route{
		DomainGuid: domainGuid,
		SpaceGuid:  mySpace.Guid,
		Host:       host,
	}
	routeBody, _ := json.Marshal(routeRequest)

	output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", "v2/routes", "-X", "POST", "-d", string(routeBody))
	if (err != nil) {
		return
	}
	response = stringifyCurlResponse(output)
	err = json.Unmarshal([]byte(response), &route)

	return
}

func getAppGuid(cliConnection plugin.CliConnection, app string) (appGuid string, err error) {
	mySpace, _ := cliConnection.GetCurrentSpace()

	output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", fmt.Sprintf("v2/apps?q=name:%s&q=space_guid:%s", app, mySpace.Guid))
	if err != nil {
		return
	}
	response := stringifyCurlResponse(output)

	apps := AppModelList{}
	err = json.Unmarshal([]byte(response), &apps)
	if (err != nil) {
		return
	}

	if len(apps.Resources) == 0 {
		err = errors.New(fmt.Sprintf("App %s not found", app))
		return
	}

	appGuid = apps.Resources[0].Metadata.Guid
	return
}

func getDomainGuid(cliConnection plugin.CliConnection, domain string) (domainGuid string, err error) {
	output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", fmt.Sprintf("v2/domains?q=name:%s", domain))
	if err != nil {
		return
	}
	response := stringifyCurlResponse(output)

	domains := DomainList{}
	err = json.Unmarshal([]byte(response), &domains)
	if (err != nil) {
		return
	}

	if len(domains.Resources) == 0 {
		err = errors.New(fmt.Sprintf("Domain %s not found", domain))
		return
	}

	domainGuid = domains.Resources[0].Metadata.Guid
	return
}

func stringifyCurlResponse(output []string) string {
	var responseString string
	for _, part := range output {
		responseString += part
	}
	return responseString
}
