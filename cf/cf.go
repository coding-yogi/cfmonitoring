package cf

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/coding-yogi/cfmonitoring/config"
	"github.com/coding-yogi/cfmonitoring/log"
)

type IClient interface {
	GetClient() error
	GetApp(appName string) (cfclient.App, error)
}

type Client struct {
	orgName   string
	spaceName string
	apps      []string
	client    *cfclient.Client
}

// GetClient ...
func (c *Client) GetClient() error {

	settings, err := config.GetConfig()
	if err != nil {
		return err
	}

	config := &cfclient.Config{
		ApiAddress: settings.Cf.API,
		Username:   settings.Cf.Username,
		Password:   settings.Cf.Password,
	}

	log.Info("Connecting to CF")
	if c.client, err = cfclient.NewClient(config); err != nil {
		return err
	}

	c.orgName = settings.Cf.Org
	c.spaceName = settings.Cf.Space
	c.apps = settings.Cf.Apps

	return nil
}

// GetApps ...
func (c *Client) GetApps() []cfclient.App {

	apps := make([]cfclient.App, len(c.apps))

	log.Info("Getting org: " + c.orgName)
	org, err := c.client.GetOrgByName(c.orgName)
	if err != nil {
		return apps
	}

	log.Info("Getting space: " + c.spaceName)
	space, err := c.client.GetSpaceByName(c.spaceName, org.Guid)
	if err != nil {
		return apps
	}

	for i := 0; i < len(c.apps); i++ {
		apps[i], _ = c.client.AppByName(c.apps[i], space.Guid, org.Guid)
	}

	return apps
}

// GetAppStats
func (c *Client) GetAppStats(app cfclient.App) (map[string]cfclient.AppStats, error) {

	log.Debug("Get App Stats by for App: " + app.Name)
	appStats, err := c.client.GetAppStats(app.Guid)
	if err != nil {
		return map[string]cfclient.AppStats{}, err
	}

	return appStats, nil
}
