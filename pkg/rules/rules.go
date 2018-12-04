package rules

import (
	"io"
	"os"

	"strings"

	"encoding/json"

	"github.com/nuclio/logger"
	"github.com/v3io/go-errors"
	"github.com/v3io/logfwd/pkg/record"
	"gopkg.in/yaml.v2"
)

type RuleConfig struct {
	log logger.Logger

	rules *Rules

	sinks              map[string]*HTTPSink
	adminSink          *HTTPSink
	defaultSink        *HTTPSink
	deadLetters        bool
	errorOnDeadLetters bool
}

func NewRuleConfig(log logger.Logger, ruleConfigPath string) (*RuleConfig, error) {
	ruleConfig := &RuleConfig{log: log.GetChild("rules")}
	ruleConfig.log.DebugWith("Reading config file", "ruleConfigPath", ruleConfigPath)
	rules, err := fromRuleConfigPath(ruleConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read rules from ruleConfig file")
	}
	ruleConfig.log.InfoWith("Read config file",
		"output-rules", len(rules.Output),
		"admin-nss", len(rules.Admin.Namespaces),
		"allow-dead-letters", rules.Admin.DeadLetters,
		"error-on-dead-letters", rules.Admin.ErrorOnDeadLetters)
	ruleConfig.rules = rules

	if err := ruleConfig.init(); err != nil {
		return nil, errors.Wrap(err, "Unable to initialize ruleConfig")
	}
	return ruleConfig, nil
}

func (c *RuleConfig) Send(record *record.LogRecord) error {
	ns := record.Kubernetes.Namespace
	c.log.DebugWith("Looking for output sink", "ns", ns)
	httpOutput, exists := c.sinks[strings.ToLower(ns)]
	if !exists {
		if c.deadLetters {
			c.log.DebugWith("Using default sink", "ns", ns)
			return c.sendToSink(c.adminSink, record)
		}
		if c.errorOnDeadLetters {
			return errors.Errorf("Unable to find output for %s namespace", ns)
		}
		c.log.DebugWith("Ignoring missing sink", "ns", ns)
		return nil
	}
	return c.sendToSink(httpOutput, record)
}

func (c *RuleConfig) init() error {
	c.sinks = make(map[string]*HTTPSink)
	c.log.Debug("Starting Admin Sink")
	c.adminSink = &HTTPSink{
		output: &c.rules.Admin.Output.HTTP,
		log:    c.log.GetChild("sink.admin")}
	c.adminSink.Start()
	if c.rules.Admin.DeadLetters {
		c.log.Debug("Using Default Sink (dead letters)")
		c.deadLetters = true
		c.defaultSink = c.adminSink
	}
	if c.rules.Admin.ErrorOnDeadLetters {
		c.errorOnDeadLetters = true
	}
	for _, ns := range c.rules.Admin.Namespaces {
		c.log.DebugWith("Admin namespace added", "ns", ns)
		c.sinks[strings.ToLower(ns)] = c.adminSink
	}

	for _, output := range c.rules.Output {
		c.log.DebugWith("Output namespace sink started", "ns", output.Namespace)
		sink := &HTTPSink{
			output: &output.HTTP,
			log:    c.log.GetChild("sink." + output.Namespace)}
		sink.Start()
		c.sinks[strings.ToLower(output.Namespace)] = sink
	}

	return nil
}

func (c *RuleConfig) sendToSink(sink *HTTPSink, record *record.LogRecord) error {
	dataAsJSON, err := json.Marshal(record)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal partial data")
	}
	sink.Send(dataAsJSON)
	return nil
}

func fromRuleConfigPath(ruleConfigPath string) (*Rules, error) {
	openFile, err := os.Open(ruleConfigPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to open rules file %s", ruleConfigPath)
	}
	return fromReader(openFile)
}

func fromReader(reader io.Reader) (*Rules, error) {
	decoder := yaml.NewDecoder(reader)
	result := Rules{}
	if err := decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "Unable to read rules file")
	}
	return &result, nil
}
