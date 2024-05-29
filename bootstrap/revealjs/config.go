package revealjs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
)

type Config struct {
	Slides   []string               `json:"slides"`
	Title    string                 `json:"title"`
	Theme    string                 `json:"theme"`
	RevealJS map[string]interface{} `json:"revealjs"`
}

func LoadConfigFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var c Config
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) RevealJSConfig() map[string]string {
	m := map[string]string{}

	// config from env
	p := regexp.MustCompile(`REVEALJS_(.*?)=(.*)`)
	for _, env := range os.Environ() {
		if match := p.FindStringSubmatch(env); len(match) > 0 {
			key := match[1]
			value := match[2]
			m[key] = c.valueToString(key, value)
		}
	}
	// config from file
	for k, v := range c.RevealJS {
		m[k] = c.valueToString(k, v)
	}
	return m
}

func (c *Config) valueToString(k string, v interface{}) string {
	switch k {
	case "controlsLayout", "controlsBackArrows", "transition", "transitionSpeed", "backgroundTransition", "parallaxBackgroundImage", "parallaxBackgroundSize", "display", "showSlideNumber", "parallaxBackgroundPosition", "parallaxBackgroundRepeat", "autoAnimateMatcher", "navigationMode", "preloadIframes", "autoAnimateEasing":
		if v == nil {
			return "null"
		}
		return fmt.Sprintf(`'%v'`, v)
	case "autoPlayMedia", "autoSlideMethod", "defaultTiming", "parallaxBackgroundHorizontal", "parallaxBackgroundVertical", "keyboardCondition":
		if v == nil {
			return "null"
		}
		return fmt.Sprint(v)
	case "autoAnimateStyles":
		b, _ := json.Marshal(v)
		return string(b)
	case "plugins":
		b, _ := json.Marshal(v)
		return strings.ReplaceAll(string(b), `"`, ``)
	default:
		return fmt.Sprint(v)
	}
}
