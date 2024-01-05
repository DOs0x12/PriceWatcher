package configer

import (
	"PriceWatcher/internal/entities/config"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigDto struct {
	BotKey   string           `yaml:"botKey"`
	Services []ServiceConfDto `yaml:"services"`
}

type ServiceConfDto struct {
	PriceType    string            `yaml:"price_type"`
	SendingHours []int             `yaml:"sending_hours"`
	Items        map[string]string `yaml:"items"`
	Marketplace  string            `yaml:"marketplace"`
	Email        EmailDto          `yaml:"email"`
}

type EmailDto struct {
	From     string `yaml:"from"`
	Pass     string `yaml:"password"`
	To       string `yaml:"to"`
	SmtpHost string `yaml:"smtp_host"`
	SmtpPort int    `yaml:"smtp_port"`
}

type Configer struct {
	path string
}

func (c Configer) GetConfig() (config.Config, error) {
	confFile, err := os.ReadFile(c.path)
	if err != nil {
		return config.Config{}, err
	}

	return unmarshalConf(confFile)
}

func NewConfiger(path string) Configer {
	return Configer{path: path}
}

func (c Configer) AddItemToWatch(address, name, priceType string) error {
	str, err := os.ReadFile(c.path)
	if err != nil {
		return fmt.Errorf("cannot read the file %v: %w", c.path, err)
	}

	var config yaml.Node
	yaml.Unmarshal(str, &config)

	itemsNode := getItems(&config, priceType)
	if itemsNode == nil {
		return fmt.Errorf("cannot find the items field in the price_type %v: %w", priceType, err)
	}

	nameNode := yaml.Node{Kind: yaml.ScalarNode, Style: yaml.SingleQuotedStyle, Tag: "!!str", Value: name}
	addressNode := yaml.Node{Kind: yaml.ScalarNode, Style: yaml.SingleQuotedStyle, Tag: "!!str", Value: address}

	itemsNode.Content = append(itemsNode.Content, &nameNode)
	itemsNode.Content = append(itemsNode.Content, &addressNode)

	file, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("cannot create a file: %v", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	enc := yaml.NewEncoder(file)

	return enc.Encode(config.Content[0])
}

func getItems(node *yaml.Node, priceType string) *yaml.Node {
	itemsParam := "items"
	marketplaceParam := "marketplace"

	for i, cont := range node.Content {
		if cont.Value == marketplaceParam && cont.Kind == yaml.ScalarNode && node.Content[i+1].Value == priceType {
			for j, ptCont := range node.Content {
				if ptCont.Value == itemsParam {
					return node.Content[j+1]
				}
			}

			return nil
		}

		if itemsNode := getItems(cont, priceType); itemsNode != nil {
			return itemsNode
		}
	}

	return nil
}

func (c Configer) RemoveItemFromWatching(address, name, priceType string) error {
	return nil
}
