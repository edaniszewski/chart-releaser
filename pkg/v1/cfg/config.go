package v1

import (
	"errors"
	"fmt"
	"strings"

	"github.com/edaniszewski/chart-releaser/pkg/errs"
	"github.com/edaniszewski/chart-releaser/pkg/strategies"
	"sigs.k8s.io/yaml"
)

// ConfigVersion defines the version for the v1 configuration scheme.
// This is always expected to be "v1" for v1 configs.
const ConfigVersion = "v1"

// Errors for v1 configuration parsing and validation.
var (
	ErrNoChart         = errors.New("required option 'chart' missing from config")
	ErrNoChartName     = errors.New("required option 'chart.name' missing from config")
	ErrNoChartRepo     = errors.New("required option 'chart.repo' missing from config")
	ErrUnsupportedRepo = errors.New("unsupported repo specified in 'chart.repo'. currently supported repos include: github.com")
)

// Config contains the configuration options for chart-releaser's
// v1 configuration scheme.
type Config struct {
	Version string          `yaml:"version,omitempty"`
	Chart   *ChartConfig    `yaml:"chart,omitempty"`
	Publish *PublishConfig  `yaml:"publish,omitempty"`
	Commit  *CommitConfig   `yaml:"commit,omitempty"`
	Release *ReleaseConfig  `yaml:"release,omitempty"`
	Extras  []*ExtrasConfig `yaml:"extras,omitempty"`
}

// LoadFromBytes attempts to load raw bytes into a Config struct.
func LoadFromBytes(b []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// GetVersion returns the version for the configuration format.
func (c *Config) GetVersion() string {
	return c.Version
}

// Validate the Config is correct.
func (c *Config) Validate() error {
	collector := errs.NewCollector()

	if c.Version != ConfigVersion {
		collector.Add(fmt.Errorf("error: using v1 config parser for non v1 config"))
	}

	if c.Chart == nil {
		collector.Add(ErrNoChart)
	} else {
		if err := c.Chart.validate(); err != nil {
			collector.Add(err)
		}
	}

	if c.Publish == nil {
		c.Publish = &PublishConfig{}
	}
	if err := c.Publish.validate(); err != nil {
		collector.Add(err)
	}

	if c.Commit == nil {
		c.Commit = &CommitConfig{}
	}
	if err := c.Commit.validate(); err != nil {
		collector.Add(err)
	}

	if c.Release == nil {
		c.Release = &ReleaseConfig{}
	}
	if err := c.Release.validate(); err != nil {
		collector.Add(err)
	}

	if c.Extras == nil {
		c.Extras = []*ExtrasConfig{}
	}
	for _, extra := range c.Extras {
		if err := extra.validate(); err != nil {
			collector.Add(err)
		}
	}

	if collector.HasErrors() {
		return collector
	}
	return nil
}

// ChartConfig contains the options for the v1 configuration's "chart"
// section. These options provide definitions for where chart-releaser
// can locate the Helm Chart for the configured project.
type ChartConfig struct {
	Name string `yaml:"name,omitempty"`
	Repo string `yaml:"repo,omitempty"`
	Path string `yaml:"path,omitempty"`
}

// validate the ChartConfig is correct.
func (c *ChartConfig) validate() error {
	collector := errs.NewCollector()

	if c.Name == "" {
		collector.Add(ErrNoChartName)
	}
	if c.Repo == "" {
		collector.Add(ErrNoChartRepo)
	} else {
		// TODO: I don't know that this is the right check to have here...
		if !strings.HasPrefix(c.Repo, "github.com") {
			collector.Add(ErrUnsupportedRepo)
		}
	}

	if collector.HasErrors() {
		return collector
	}
	return nil
}

// PublishConfig contains the options for the v1 configuration's "publish"
// section. These options provide definitions for how chart-releaser should
// behave when publishing changes to the chart repo for a new application version.
type PublishConfig struct {
	Commit *PublishCommitConfig `yaml:"commit,omitempty"`
	PR     *PublishPRConfig     `yaml:"pr,omitempty"`
}

// validate the PublishConfig is correct.
func (c *PublishConfig) validate() error {
	collector := errs.NewCollector()

	if c.Commit != nil && c.PR != nil {
		collector.Add(fmt.Errorf("invalid publish config: cannot define both 'commit' and 'pr' blocks"))
	} else if c.Commit == nil && c.PR == nil {
		// Default behavior is to use PR when no config is specified
		c.PR = &PublishPRConfig{}
	}

	if c.Commit != nil {
		if err := c.Commit.validate(); err != nil {
			collector.Add(err)
		}
	}

	if c.PR != nil {
		if err := c.PR.validate(); err != nil {
			collector.Add(err)
		}
	}

	if collector.HasErrors() {
		return collector
	}
	return nil
}

// CommitConfig contains the options for the v1 configuration's "commit"
// section. These options provide definitions for how chart-releaser should
// format commits as well as metadata about the committer.
type CommitConfig struct {
	Author    *CommitAuthorConfig   `yaml:"author,omitempty"`
	Templates *CommitTemplateConfig `yaml:"templates,omitempty"`
}

// validate the CommitConfig is correct.
func (c *CommitConfig) validate() error {
	collector := errs.NewCollector()

	if c.Author == nil {
		c.Author = &CommitAuthorConfig{}
	}
	if err := c.Author.validate(); err != nil {
		collector.Add(err)
	}

	if c.Templates == nil {
		c.Templates = &CommitTemplateConfig{}
	}
	if err := c.Templates.validate(); err != nil {
		collector.Add(err)
	}

	if collector.HasErrors() {
		return collector
	}
	return nil
}

// CommitTemplateConfig defines the templates for commit messages for different
// operations.
type CommitTemplateConfig struct {
	Update string `yaml:"update,omitempty"`
	Extras string `yaml:"extras,omitempty"`
}

// validate the CommitTemplateConfig is correct.
func (c *CommitTemplateConfig) validate() error {
	// TODO (etd): could check if templates are valid upfront
	return nil
}

// ReleaseConfig contains the options for the v1 configuration's "release"
// section. These options provide definitions for how chart-releaser should
// operate when a new release of the application is cut.
type ReleaseConfig struct {
	Matches  []string `yaml:"matches,omitempty"`
	Ignores  []string `yaml:"ignores,omitempty"`
	Strategy string   `yaml:"strategy,omitempty"`
}

// validate the ReleaseConfig is correct.
func (c *ReleaseConfig) validate() error {
	collector := errs.NewCollector()

	// Only validate the strategy if one is set. If not set, this will use
	// the "default" strategy once the update pipeline is run.
	if c.Strategy != "" {
		if _, err := strategies.UpdateStrategyFromString(c.Strategy); err != nil {
			collector.Add(fmt.Errorf("invalid release strategy '%v', should be one of: %v", c.Strategy, strategies.ListUpdateStrategies()))
		}
	}

	if collector.HasErrors() {
		return collector
	}
	return nil
}

// PublishCommitConfig contains additional options for how chart-releaser
// should behave when using the "commit" update strategy.
type PublishCommitConfig struct {
	Branch string `yaml:"branch,omitempty"`
	Base   string `yaml:"base,omitempty"`
}

// validate the UpdateCommitConfig is correct.
func (c *PublishCommitConfig) validate() error {
	return nil
}

// PublishPRConfig contains additional options for how chart-releaser should
// behave when using the "pull request" update strategy.
type PublishPRConfig struct {
	BranchTemplate string `yaml:"branch_template,omitempty"`
	Base           string `yaml:"base,omitempty"`
	TitleTemplate  string `yaml:"title_template,omitempty"`
	BodyTemplate   string `yaml:"body_template,omitempty"`
}

// validate the PublishPRConfig is correct.
func (c *PublishPRConfig) validate() error {
	// TODO (etd): could check if templates are valid upfront
	return nil
}

// CommitAuthorConfig provides the commit metadata for who the author of
// the commits made by chart-releaser will be.
type CommitAuthorConfig struct {
	Name  string `yaml:"name,omitempty"`
	Email string `yaml:"email,omitempty"`
}

// validate the CommitAuthorConfig is correct.
func (c *CommitAuthorConfig) validate() error {
	collector := errs.NewCollector()

	if c.Name != "" && c.Email == "" {
		collector.Add(fmt.Errorf("commit author specifies name, but no email"))
	}

	if c.Name == "" && c.Email != "" {
		collector.Add(fmt.Errorf("commit author specified email, but no name"))
	}

	if collector.HasErrors() {
		return collector
	}
	return nil
}

// ExtrasConfig is used to specify additional files within the configured
// repository to update via a regular expression (regex).
type ExtrasConfig struct {
	Path    string           `yaml:"path,omitempty"`
	Updates []*SearchReplace `yaml:"updates,omitempty"`
}

// validate the ExtrasConfig is correct.
func (c *ExtrasConfig) validate() error {
	collector := errs.NewCollector()

	if c.Path != "" && len(c.Updates) == 0 {
		collector.Add(fmt.Errorf("extras config specifies path but no options for search/replace updates"))
	} else if c.Path == "" && len(c.Updates) != 0 {
		collector.Add(fmt.Errorf("extras config specifies search/replace options, but no file path"))
	}

	for _, u := range c.Updates {
		if err := u.validate(); err != nil {
			// TODO (etd): If there are multiple updates with error, there is no
			//   way of saying which one is in error, as there is no context collected.
			//   perhaps consider a way to pass along context when adding an error so
			//   it can be augmented in the error message? e.g. a map[string]interface{}
			collector.Add(err)
		}
	}

	if collector.HasErrors() {
		return collector
	}
	return nil
}

// SearchReplace defines a regex to search for and a value to replace the found
// match(es) to the regex.
type SearchReplace struct {
	Search  string `yaml:"search,omitempty"`
	Replace string `yaml:"replace,omitempty"`
	Limit   int    `yaml:"limit,omitempty"`
}

// validate the SearchReplace is correct.
func (c *SearchReplace) validate() error {
	collector := errs.NewCollector()

	if c.Search == "" {
		collector.Add(fmt.Errorf("search and replace config has no search string set"))
	}
	if c.Replace == "" {
		collector.Add(fmt.Errorf("search and replace config has no replace string set"))
	}
	if c.Limit < 0 {
		collector.Add(fmt.Errorf("search and replace config must have a non-negative limit"))
	}

	if collector.HasErrors() {
		return collector
	}
	return nil
}
