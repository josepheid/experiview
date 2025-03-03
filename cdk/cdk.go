package main

import (
	"net/http"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"

	golambda "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	experimentsTable := awsdynamodb.NewTableV2(stack, jsii.String("experimentsTable"), &awsdynamodb.TablePropsV2{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	// Allow querying all experiments by date
	experimentsTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexPropsV2{
		IndexName: jsii.String("allExperimentsIndex"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("type"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	createExperiment := golambda.NewGoFunction(stack, jsii.String("createExperiment"), &golambda.GoFunctionProps{
		Entry:         jsii.String("../backend/api/handlers/createexperiment/post"),
		Description:   jsii.String("lambda responsible for creating experiment records"),
		InitialPolicy: &[]awsiam.PolicyStatement{},
		Environment: &map[string]*string{
			"EXPERIMENTS_TABLE_NAME": experimentsTable.TableName(),
		},
	})

	experimentsTable.GrantFullAccess(createExperiment)

	notFound := golambda.NewGoFunction(stack, jsii.String("notFound"), &golambda.GoFunctionProps{
		Description: jsii.String("Returns a not found response."),
		Entry:       jsii.String("../backend/api/handlers/notfound"),
		MemorySize:  jsii.Number(128),
	})

	apiResourceOpts := &awsapigateway.ResourceOptions{}
	apiLambdaOpts := &awsapigateway.LambdaIntegrationOptions{}
	api := awsapigateway.NewLambdaRestApi(stack, jsii.String("experiview-api"), &awsapigateway.LambdaRestApiProps{
		CloudWatchRole: jsii.Bool(false),
		Handler:        notFound,
		Proxy:          jsii.Bool(false),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowHeaders: jsii.Strings("Content-Type", "Authorization", "x-api-key"),
		},

		ApiKeySourceType: awsapigateway.ApiKeySourceType_HEADER,
	})

	apiKey := awsapigateway.NewApiKey(stack, jsii.String("experiview-apikey"), &awsapigateway.ApiKeyProps{})

	usagePlan := awsapigateway.NewUsagePlan(stack, jsii.String("experiview-usagePlan"), &awsapigateway.UsagePlanProps{
		Name:      jsii.String("ExperiviewUsagePlan"),
		ApiStages: &[]*awsapigateway.UsagePlanPerApiStage{{Api: api, Stage: api.DeploymentStage()}},
	})

	usagePlan.AddApiKey(apiKey, &awsapigateway.AddApiKeyOptions{})

	experiview := api.Root().AddResource(jsii.String("experiview"), apiResourceOpts)
	experiments := experiview.AddResource(jsii.String("experiments"), apiResourceOpts)
	createExperimentPostIntegration := awsapigateway.NewLambdaIntegration(createExperiment, apiLambdaOpts)
	experiments.AddMethod(jsii.String(http.MethodPost), createExperimentPostIntegration, &awsapigateway.MethodOptions{ApiKeyRequired: jsii.Bool(true)})
	// The code that defines your stack goes here

	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("CdkQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "ExperiviewStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
