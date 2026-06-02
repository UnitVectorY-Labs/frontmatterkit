package frontmatter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

const schemaResource = "frontmatter.schema.json"

var ErrInvalidJSONSchema = errors.New("invalid JSON schema")

// ValidateJSONSchema validates the document's front matter against a JSON Schema.
// Missing front matter is validated as an empty object.
func ValidateJSONSchema(doc *Document, schemaReader io.Reader) error {
	decoder := json.NewDecoder(schemaReader)
	var schemaDoc any
	if err := decoder.Decode(&schemaDoc); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidJSONSchema, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("multiple JSON values")
		}
		return fmt.Errorf("%w: %s", ErrInvalidJSONSchema, err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaResource, schemaDoc); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidJSONSchema, err)
	}

	schema, err := compiler.Compile(schemaResource)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidJSONSchema, err)
	}

	value, err := jsonValue(doc)
	if err != nil {
		return err
	}
	if err := schema.Validate(value); err != nil {
		return fmt.Errorf("front matter does not match schema: %w", err)
	}

	return nil
}

func jsonValue(doc *Document) (any, error) {
	if doc == nil || !doc.HasFrontMatter || doc.Node == nil {
		return map[string]any{}, nil
	}

	node := doc.Node
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return map[string]any{}, nil
		}
		node = node.Content[0]
	}

	return yamlNodeToJSONValue(node)
}

func yamlNodeToJSONValue(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.MappingNode:
		value := make(map[string]any, len(node.Content)/2)
		for i := 0; i < len(node.Content)-1; i += 2 {
			key := node.Content[i].Value
			child, err := yamlNodeToJSONValue(node.Content[i+1])
			if err != nil {
				return nil, err
			}
			value[key] = child
		}
		return value, nil
	case yaml.SequenceNode:
		value := make([]any, 0, len(node.Content))
		for _, childNode := range node.Content {
			child, err := yamlNodeToJSONValue(childNode)
			if err != nil {
				return nil, err
			}
			value = append(value, child)
		}
		return value, nil
	case yaml.ScalarNode:
		return yamlScalarToJSONValue(node)
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return map[string]any{}, nil
		}
		return yamlNodeToJSONValue(node.Content[0])
	case yaml.AliasNode:
		return nil, fmt.Errorf("aliases are not allowed")
	default:
		return nil, fmt.Errorf("unsupported YAML node kind %d", node.Kind)
	}
}

func yamlScalarToJSONValue(node *yaml.Node) (any, error) {
	switch node.Tag {
	case "!!null":
		return nil, nil
	case "!!bool":
		return strconv.ParseBool(node.Value)
	case "!!int":
		return strconv.ParseInt(strings.ReplaceAll(node.Value, "_", ""), 0, 64)
	case "!!float":
		return strconv.ParseFloat(strings.ReplaceAll(node.Value, "_", ""), 64)
	case "!!str", "!!timestamp", "!!binary":
		return node.Value, nil
	default:
		var value any
		if err := node.Decode(&value); err != nil {
			return nil, fmt.Errorf("convert YAML value to JSON-compatible value: %w", err)
		}
		return value, nil
	}
}
