package pkmnquizbot

import (
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ChatConfig struct {
	Id              int64 `dynamo:"id"`
	ShadowMessageId int   `dynamo:"shadowMessageId"`
	CurrentPokemon  int   `dynamo:"currentPokemon"`
}

//StoredAnswers holds the current Pokemon for any given chat
var DynamoClient *dynamodb.DynamoDB

const tableName string = "ChatConfig"

func getGameState(chat int64) ChatConfig {
	id := strconv.FormatInt(chat, 10)
	log.Printf("Getting chat state for chat %s", id)
	result, err := DynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Print(err)
		return ChatConfig{}
	}
	config := ChatConfig{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &config)
	if err != nil {
		log.Print(err)
		return ChatConfig{}
	}
	log.Printf("Chat state: %+v", config)
	return config
}

func saveGameState(chat int64, message int, pokemonId int) {
	id := strconv.FormatInt(chat, 10)
	log.Printf("Saving state for chat %s", id)
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":currentPokemon": {
				N: aws.String(strconv.Itoa(pokemonId)),
			},
			":shadowMessageId": {
				N: aws.String(strconv.Itoa(message)),
			},
		},
		UpdateExpression: aws.String("set currentPokemon = :currentPokemon, shadowMessageId = :shadowMessageId"),
		TableName:        aws.String(tableName),
	}

	_, err := DynamoClient.UpdateItem(input)
	if err != nil {
		log.Print(err)
	}
}
