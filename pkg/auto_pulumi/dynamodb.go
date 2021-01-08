package auto_pulumi

import (
	"fmt"
	"path"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

// convertSchemaTypeToDynamoType converts a type from the schema to a Dynamo key type.
func convertSchemaTypeToDynamoType(typeValue string) (string, error) {
	switch typeValue {
	case "string":
		return "S", nil
	case "number":
		return "N", nil
	default:
		return "", utils.NewErrorMessage("Unknown or supported type provided", nil)
	}
}

type DynamoDBTable struct {
	Name        pulumi.StringOutput
	DisplayName string
}

// BuildDynamoDBTables builds the tables during deployment.
func BuildDynamoDBTables(ctx *pulumi.Context, modelDirectoryPath, environment string) ([]DynamoDBTable, error) {
	modelDirContents, err := utils.ReadDirectoryContents(modelDirectoryPath)
	if err != nil {
		return nil, utils.NewErrorMessage("Error reading model directory contents", err)
	}

	var modelNames []DynamoDBTable
	for _, fileName := range modelDirContents {
		rawModelName := strings.Split(fileName, ".")[0]
		modelName := fmt.Sprintf("%s_%s", environment, rawModelName)
		filePath := path.Join(modelDirectoryPath, fileName)
		modelSchema, err := utils.GenerateModelSchemaFromFile(filePath)
		if err != nil {
			return nil, utils.NewErrorMessage("Error generating model schema", err)
		}

		var hashKey DynamoDBKey
		var rangeKey DynamoDBKey
		for name, item := range modelSchema.Definitions[utils.DashCaseToSentenceCase(rawModelName)].Properties {
			if item.HashKey == true {
				hashKey.Name = name
				hashKey.Type, err = convertSchemaTypeToDynamoType(item.Type)
				if err != nil {
					return nil, utils.NewErrorMessage("Error converting hashKey type", err)
				}
			}

			if item.RangeKey == true {
				rangeKey.Name = name
				rangeKey.Type, err = convertSchemaTypeToDynamoType(item.Type)
				if err != nil {
					return nil, utils.NewErrorMessage("Error converting rangeKey type", err)
				}
			}
		}

		createdTable, err := CreateDynamoDBTable(ctx, modelName, hashKey, rangeKey)
		if err != nil {
			return nil, utils.NewErrorMessage("Error creating DynamoDB table", err)
		}
		modelNames = append(modelNames, DynamoDBTable{
			DisplayName: modelName,
			Name:        createdTable.Name,
		})
	}

	return modelNames, nil
}

type DynamoDBKey struct {
	Name string
	Type string
}

// CreateDynamoDBTable provisions/updates a DynamoDB table.
func CreateDynamoDBTable(ctx *pulumi.Context, name string, hashKey, rangeKey DynamoDBKey) (*dynamodb.Table, error) {
	tableArgs := &dynamodb.TableArgs{
		HashKey:       pulumi.String(hashKey.Name),
		WriteCapacity: pulumi.Int(1),
		ReadCapacity:  pulumi.Int(1),
		BillingMode:   pulumi.String("PAY_PER_REQUEST"),
	}

	tableAttributes := dynamodb.TableAttributeArray{
		&dynamodb.TableAttributeArgs{
			Name: pulumi.String(hashKey.Name),
			Type: pulumi.String(hashKey.Type),
		},
	}

	if rangeKey.Name != "" {
		tableArgs.RangeKey = pulumi.String(rangeKey.Name)
		tableAttributes = append(tableAttributes, &dynamodb.TableAttributeArgs{
			Name: pulumi.String(rangeKey.Name),
			Type: pulumi.String(rangeKey.Type),
		})
	}

	tableArgs.Attributes = tableAttributes

	table, err := dynamodb.NewTable(ctx, name, tableArgs)
	if err != nil {
		return nil, utils.NewErrorMessage("Error creating DynamoDB table", err)
	}

	return table, nil
}
