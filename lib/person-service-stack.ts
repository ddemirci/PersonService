import * as cdk from 'aws-cdk-lib/core';
import { Construct } from 'constructs';
// import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as events from 'aws-cdk-lib/aws-events';
import * as lambdaEventSources from 'aws-cdk-lib/aws-lambda-event-sources';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as sqs from 'aws-cdk-lib/aws-sqs';


export class PersonServiceStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

   const table = new dynamodb.Table(this, 'PersonsTable', {
      tableName: `persons-${this.node.tryGetContext('env') || 'dev'}`,
      partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      stream: dynamodb.StreamViewType.NEW_IMAGE,
    });

    const getLambda = new lambda.Function(this, 'GetPersonsLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      handler: 'main',
      code: lambda.Code.fromAsset('service'),
    });

    const postLambda = new lambda.Function(this, 'PostPersonLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      handler: 'main',
      code: lambda.Code.fromAsset('service'),
    });

    const api = new apigateway.RestApi(this, 'PersonApi', {
      restApiName: 'Person Service',
    });

    const eventBus = new events.EventBus(this, 'PersonEventBus', {
      eventBusName: `person-bus-${this.node.tryGetContext('env') || 'dev'}`,
    });

    const dlq = new sqs.Queue(this, 'StreamDLQ', {
      queueName: `person-dlq-${this.node.tryGetContext('env') || 'dev'}`,
      retentionPeriod: cdk.Duration.days(14),
    });

    const streamLambda = new lambda.Function(this, 'PersonStreamLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset('stream-service'),
      deadLetterQueue: dlq
    });

    const logGroup = new logs.LogGroup(this, 'EventLogs');

    // Pass bus name to Lambda
    streamLambda.addEnvironment('EVENT_BUS_NAME', eventBus.eventBusName);

    streamLambda.addEventSource(new lambdaEventSources.DynamoEventSource(table, {
      startingPosition: lambda.StartingPosition.LATEST,
      batchSize: 1,
      retryAttempts: 3,
    }));

    new events.Rule(this, 'LogPersonEvents', {
      eventBus: eventBus,
      eventPattern: {
        detailType: ['PersonCreated'],
      },
      targets: [new targets.CloudWatchLogGroup(logGroup)],
    });

    // Give Lambda access to DynamoDB
    table.grantReadData(getLambda);
    table.grantWriteData(postLambda);

    // Pass table names to Lambda
    getLambda.addEnvironment('TABLE_NAME', table.tableName);
    postLambda.addEnvironment('TABLE_NAME', table.tableName);

    // Create /persons resource
    const persons = api.root.addResource('persons');
    
    // Connect GET → Lambda
    persons.addMethod('GET', new apigateway.LambdaIntegration(getLambda));
    
    // Connect POST → Lambda
    persons.addMethod('POST', new apigateway.LambdaIntegration(postLambda));

    // Allow the stream Lambda to send events to EventBridge
    eventBus.grantPutEventsTo(streamLambda);
  }
}
