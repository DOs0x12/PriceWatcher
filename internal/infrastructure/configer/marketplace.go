package configer

import (
	"PriceWatcher/internal/entities/config"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

func (c Configer) GetMarketplaceConfig(name string) (config.ServiceConf, error) {
	fullConfig, err := c.GetConfig()
	if err != nil {
		return config.ServiceConf{}, err
	}

	for _, serviceConfig := range fullConfig.Services {
		if serviceConfig.Marketplace == name {
			return serviceConfig, nil
		}
	}

	return config.ServiceConf{}, fmt.Errorf("do not find a service config with the name: %v", name)
}

func (c Configer) AddItemToWatch(name, address, priceType string) error {
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

func (c Configer) RemoveItemFromWatching(name, priceType string) error {
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

	if err := removeItem(name, itemsNode); err != nil {
		return err
	}

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

func removeItem(name string, itemsNode *yaml.Node) error {
	for i, item := range itemsNode.Content {
		if item.Value == name {
			itemsNode.Content = slices.Delete(itemsNode.Content, i, i+2)

			return nil
		}
	}

	return fmt.Errorf("cannot find an item with the name: %v", name)
}
