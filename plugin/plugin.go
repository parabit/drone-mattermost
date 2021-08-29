package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/drone-plugins/drone-plugin-lib/drone"
	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/flowchartsman/handlebars/v3"
	mattermost "github.com/mattermost/mattermost-server/v5/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Plugin implements drone.Plugin to provide the plugin implementation.
type Plugin struct {
	URL      string
	Token    string
	Team     string
	Channel  string
	Template string
	Replace  string
	replacer func(string) string
}

// New creates the drone mattermost plugin.
func New() *Plugin {
	return new(Plugin)
}

// Run is the cli run entry.
func (p *Plugin) Run(ctx *cli.Context) error {
	urfave.LoggingFromContext(ctx)
	if err := p.Execute(ctx); err != nil {
		if e, ok := err.(errors.ExitCoder); ok {
			return e
		}
		return errors.ExitMessagef("execution failed: %w", err)
	}
	return nil
}

// Execute executes the plugin.
func (p *Plugin) Execute(ctx *cli.Context) error {
	// validate
	if p.URL == "" || p.Token == "" {
		return ErrMissingURLOrToken
	}
	if p.Team == "" || p.Channel == "" {
		return ErrMissingTeamOrChannel
	}
	if err := p.BuildReplacer(); err != nil {
		return err
	}
	// execute
	return p.CreatePost(urfave.PipelineFromContext(ctx), urfave.NetworkFromContext(ctx))
}

// BuildReplacer builds a replacer.
func (p *Plugin) BuildReplacer() error {
	if p.Replace == "" {
		return nil
	}
	var replace []struct {
		Regexp  string         `json:"regexp,omitempty"`
		Replace string         `json:"replace,omitempty"`
		RE      *regexp.Regexp `json:"-"`
	}
	// unmarshal json
	d := json.NewDecoder(strings.NewReader(p.Replace))
	d.DisallowUnknownFields()
	if err := d.Decode(&replace); err != nil {
		return fmt.Errorf("unable to unmarshal replace from json: %v", err)
	}
	// compile
	for i := 0; i < len(replace); i++ {
		var err error
		if replace[i].RE, err = regexp.Compile(replace[i].Regexp); err != nil {
			return fmt.Errorf("unable to compile regexp %d (%q): %v", i, replace[i].Regexp, err)
		}
	}
	// set replacer
	p.replacer = func(s string) string {
		for i := 0; i < len(replace); i++ {
			s = replace[i].RE.ReplaceAllString(s, replace[i].Replace)
		}
		return s
	}
	return nil
}

// CreatePost creates the post.
func (p *Plugin) CreatePost(pipeline drone.Pipeline, network drone.Network) error {
	// replace
	if p.replacer != nil {
		pipeline.Commit.Message.Title = p.replacer(pipeline.Commit.Message.Title)
		pipeline.Commit.Message.Body = p.replacer(pipeline.Commit.Message.Body)
	}
	// template
	template :=
		"# Push `{{repo.owner}}/{{repo.name}}:{{truncate commit 7}}`\n" +
			"Pipeline for [branch `{{commit.branch}}` by `{{commit.author}}`]({{build.link}}): **{{build.status}}**!\n" +
			"> {{commit.message.title}}{{#if commit.message.body}}\n" +
			">\n" +
			"{{{regexReplace \"(?m)^\" commit.message.body \"> \"}}}{{/if}}"
	if p.Template != "" {
		template = p.Template
	}
	// render
	message, err := handlebars.Render(template, pipeline)
	if err != nil {
		return fmt.Errorf("could not render message template: %w", err)
	}
	// create client
	cl := mattermost.NewAPIv4Client(strings.TrimSpace(p.URL))
	cl.SetToken(strings.TrimSpace(p.Token))
	// retrieve team
	teamName, channelName := strings.TrimSpace(p.Team), strings.TrimSpace(p.Channel)
	logrus.WithFields(logrus.Fields{
		"team":    teamName,
		"channel": channelName,
		"message": strings.TrimSpace(message),
	}).Info("sending message")
	team, res := cl.GetTeamByName(teamName, "")
	switch {
	case res.Error != nil:
		return res.Error
	case res.StatusCode != http.StatusOK:
		return fmt.Errorf("could not retrieve team: status code (%d) != 200", res.StatusCode)
	case team.Id == "":
		return fmt.Errorf("could not determine team id: team id is blank: status code == %d", res.StatusCode)
	}
	// retrieve channel
	channel, res := cl.GetChannelByName(channelName, team.Id, "")
	switch {
	case res.Error != nil:
		return res.Error
	case res.StatusCode != http.StatusOK:
		return fmt.Errorf("could not retrieve channel: status code (%d) != 200", res.StatusCode)
	case channel.Id == "":
		return fmt.Errorf("could not determine channel id: channel id is blank: status code == %d", res.StatusCode)
	}
	// create post
	_, res = cl.CreatePost(&mattermost.Post{
		ChannelId: channel.Id,
		Message:   message,
	})
	switch {
	case res.Error != nil:
		return res.Error
	case res.StatusCode != http.StatusCreated:
		return fmt.Errorf("could not create post: status code (%d) != 201", res.StatusCode)
	}
	return nil
}

// Flags returns the configuration flags for the plugin.
func (p *Plugin) Flags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "mattermost.url",
			Usage:       "mattermost url",
			EnvVars:     []string{"MATTERMOST_URL", "PLUGIN_URL"},
			Destination: &p.URL,
		},
		&cli.StringFlag{
			Name:        "mattermost.token",
			Usage:       "mattermost token",
			EnvVars:     []string{"MATTERMOST_TOKEN", "PLUGIN_TOKEN"},
			Destination: &p.Token,
		},
		&cli.StringFlag{
			Name:        "mattermost.team",
			Usage:       "mattermost team",
			EnvVars:     []string{"MATTERMOST_TEAM", "PLUGIN_TEAM"},
			Destination: &p.Team,
		},
		&cli.StringFlag{
			Name:        "mattermost.channel",
			Usage:       "mattermost channel",
			EnvVars:     []string{"MATTERMOST_CHANNEL", "PLUGIN_CHANNEL"},
			Destination: &p.Channel,
		},
		&cli.StringFlag{
			Name:        "mattermost.template",
			Usage:       "mattermost template",
			EnvVars:     []string{"MATTERMOST_TEMPLATE", "PLUGIN_TEMPLATE"},
			Destination: &p.Template,
		},
		&cli.StringFlag{
			Name:        "mattermost.replace",
			Usage:       "mattermost replace",
			EnvVars:     []string{"MATTERMOST_REPLACE", "PLUGIN_REPLACE"},
			Destination: &p.Replace,
		},
	}
	return append(flags, urfave.Flags()...)
}

// Error is a plugin error.
type Error string

// Error satisfies the error interface.
func (err Error) Error() string {
	return string(err)
}

// Error values.
const (
	// ErrMissingURLOrToken is the missing url or token error.
	ErrMissingURLOrToken Error = "missing url or token"
	// ErrMissingTeamOrChannel is the missing team or channel error.
	ErrMissingTeamOrChannel Error = "missing team or channel"
)
