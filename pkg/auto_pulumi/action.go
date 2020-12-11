package auto_pulumi

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optup"
	"github.com/zchase/stevie/pkg/utils"
)

// outputActionSteps prints out the steps of an action to the terminal.
func outputActionSteps(steps []auto.PreviewStep) {
	utils.Print(utils.TextColor("\nPreview Resuls:\n", color.FgGreen))

	// Create the output table.
	table := utils.CreateTerminalTable()
	table.AddRow("  Name\tProvider\tDescriptor\tResource\tAction\n")

	for _, step := range steps {
		textColor := color.FgBlue
		switch step.Op {
		case "create":
			textColor = color.FgGreen
		case "update":
			textColor = color.FgYellow
		case "delete":
			textColor = color.FgRed
		}

		// Get the action and URN of the resource.
		action := utils.TextColor(strings.ToUpper(step.Op), textColor)
		urnParts := strings.Split(string(step.URN), "::")

		// Get the resource info
		resourceInfo := strings.Split(urnParts[2], ":")
		provider := resourceInfo[0]
		resourceType := resourceInfo[1]
		resourceIdentifier := resourceInfo[2]

		// Get the resource name
		resourceName := urnParts[3]

		table.AddRow("  %s\t%s\t%s\t%s\t%s\n", resourceName, provider, resourceType, resourceIdentifier, action)
	}

	table.Print()
	utils.Print("")
}

type PulumiAction struct {
	CreateDeployment   func(environment string) (auto.Stack, error)
	Environment        string
	Stack              auto.Stack
	TemporaryDirectory *utils.TemporaryDirectory
}

// SetUp sets up the action to run.
func (a *PulumiAction) SetUp(ctx context.Context, configPath string) error {
	// Create a spinner to show the user what is happening.
	checkEnvSpinner := utils.TerminalSpinner{
		SpinnerText:   "Checking environment",
		CompletedText: "✅ Environment check complete.",
		FailureText:   "❌ Environment check failed.",
	}
	checkEnvSpinner.Create()

	// Read config in.
	projectConfig, err := utils.ReadConfigFile(configPath, a.Environment)
	if err != nil {
		checkEnvSpinner.Fail()
		return err
	}

	// Create a tmp directory for compiliing the TypeScript lambdas.
	tmp := &utils.TemporaryDirectory{
		Name: tmpDirName,
	}
	err = tmp.Create()
	if err != nil {
		checkEnvSpinner.Fail()
		return fmt.Errorf("Error creating tmp directory: %v", err)
	}
	a.TemporaryDirectory = tmp

	// Stop the environment spinner and start the API deployment function
	// spinner.
	checkEnvSpinner.Stop()
	createAPISpinner := utils.TerminalSpinner{
		SpinnerText:   "Generating Pulumi Program",
		CompletedText: "✅ Pulumi Program generated.",
		FailureText:   "❌ Pulumi Program failed to generate.",
	}
	createAPISpinner.Create()

	stack, err := a.CreateDeployment(projectConfig.Environment)
	if err != nil {
		createAPISpinner.Fail()
		return fmt.Errorf("Error creating API Deployment: %v", err)
	}

	// Set the stack.
	a.Stack = stack

	// Stop the API Deployment spinner and start a spinner for setting
	// up the preview.
	createAPISpinner.Stop()
	settingUpPreviewSpinner := utils.TerminalSpinner{
		SpinnerText:   "Setting up program execution environment.",
		CompletedText: "✅ Program execution environment set up successfully.",
		FailureText:   "❌ Failed to set up program execution environment.",
	}
	settingUpPreviewSpinner.Create()

	// Install plugins and set the config values.
	workspace := stack.Workspace()

	err = workspace.InstallPlugin(ctx, "aws", "v3.2.1")
	if err != nil {
		settingUpPreviewSpinner.Fail()
		return fmt.Errorf("Failed to install program plugins: %v", err)
	}

	// Set stack configuration specifying the AWS region to deploy
	stack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "us-west-2"})

	// Stop the setting up spinner and create a spinner for the preview.
	settingUpPreviewSpinner.Stop()

	return nil
}

// Preview runs a preview of the infrastructure changes.
func (a *PulumiAction) Preview(ctx context.Context) error {
	// Create a spinner.
	actionSpinner := utils.TerminalSpinner{
		SpinnerText:   "Running infrastructure preview",
		CompletedText: "✅ Preview completed successfully.",
		FailureText:   "❌ Preview failed.",
	}
	actionSpinner.Create()

	// Run the action.
	result, err := a.Stack.Preview(ctx)
	if err != nil {
		actionSpinner.Fail()
		a.TemporaryDirectory.Clean()
		return fmt.Errorf("Error running stack preview: %v", err)
	}

	// Stop the spinner, clean the temp directory, and output the results.
	actionSpinner.Stop()
	outputActionSteps(result.Steps)
	a.TemporaryDirectory.Clean()
	return nil
}

// Update runs an update of the infrastructure.
func (a *PulumiAction) Update(ctx context.Context) error {
	// Create a spinner.
	actionSpinner := utils.TerminalSpinner{
		SpinnerText:   "Running infrastructure update",
		CompletedText: "✅ Update completed successfully.",
		FailureText:   "❌ Update failed.",
	}
	actionSpinner.Create()
	actionSpinner.SetOutput(os.Stdout)
	outputLogger := utils.TerminalSpinnerLogger{}

	// Stream the results to the terminal.
	outputStreamer := optup.ProgressStreams(&outputLogger)

	// Run the action.
	result, err := a.Stack.Up(ctx, outputStreamer)
	if err != nil {
		actionSpinner.Fail()
		a.TemporaryDirectory.Clean()
		return fmt.Errorf("Error running stack update: %v", err)
	}

	// Stop the spinner and clean the temp directory.
	actionSpinner.Stop()
	a.TemporaryDirectory.Clean()
	utils.ClearLine()
	utils.Print("")

	// Print out the outputs
	utils.Print(utils.TextColor("Outputs", color.FgGreen))
	for key, value := range result.Outputs {
		utils.Printf("    - %s: %v\n", key, value.Value)
	}

	utils.Print("")
	return nil
}

// Destroy runs a destroy of the infrastructure.
func (a *PulumiAction) Destroy(ctx context.Context) error {
	// Create a spinner.
	actionSpinner := utils.TerminalSpinner{
		SpinnerText:   "Running infrastructure destroy",
		CompletedText: "✅ Destroy completed successfully.",
		FailureText:   "❌ Destroy failed.",
	}
	actionSpinner.Create()
	outputLogger := utils.TerminalSpinnerLogger{}

	// Stream the results to the terminal.
	outputStreamer := optdestroy.ProgressStreams(&outputLogger)

	// Run the action.
	_, err := a.Stack.Destroy(ctx, outputStreamer)
	if err != nil {
		actionSpinner.Fail()
		a.TemporaryDirectory.Clean()
		return fmt.Errorf("Error running stack destroy: %v", err)
	}

	// Stop the spinner and clean the temp directory.
	actionSpinner.Stop()
	a.TemporaryDirectory.Clean()
	return nil
}
