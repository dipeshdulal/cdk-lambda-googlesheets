package main

import (
	"io"
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/joho/godotenv"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
)

type SqsTestStackProps struct {
	awscdk.StackProps
}

func NewSqsTestStack(scope constructs.Construct, id string, props *SqsTestStackProps) awscdk.Stack {
	godotenv.Load(".env")

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	credsFile, err := os.OpenFile("credentials.json", os.O_RDONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(credsFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(b))

	env := map[string]*string{
		"REFRESH_TOKEN": jsii.String(os.Getenv("REFRESH_TOKEN")),
		"CREDENTIALS":   jsii.String(string(b)),
	}

	lambdaApp := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("go-function"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:       jsii.String("lambda-app/cmd"),
		Environment: &env,
	})

	topic := awssns.NewTopic(stack, jsii.String("test-topic"), &awssns.TopicProps{
		DisplayName: jsii.String("Test Topic"),
	})

	subscription := awssnssubscriptions.NewLambdaSubscription(lambdaApp, &awssnssubscriptions.LambdaSubscriptionProps{})
	topic.AddSubscription(subscription)

	rule := awsevents.NewRule(stack, jsii.String("Rule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Rate(awscdk.Duration_Millis(jsii.Number(60 * 1000))), // 1 minute
	})
	rule.AddTarget(awseventstargets.NewSnsTopic(topic, &awseventstargets.SnsTopicProps{}))

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewSqsTestStack(app, "SqsTestStack", &SqsTestStackProps{
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
