package auto_pulumi

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/fatih/color"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optup"
	"github.com/zchase/stevie/pkg/utils"
)

type PulumiActionResourceOutput struct {
	Parts []string
	Order int
}

type PulumiActionOutputLogger struct {
	Headers              []string
	Rows                 map[string]PulumiActionResourceOutput
	NewRowsAdded         int
	TableWriteInProgress bool
	FirstPrintDone       bool
}

func (p *PulumiActionOutputLogger) Write(msg []byte) (int, error) {
	p.AddRow(string(msg))

	if len(p.Rows) > 0 && p.TableWriteInProgress == false {
		// Clear the line above so the spinner doesn't log.
		fmt.Printf("\033[2K")
		fmt.Println()
		fmt.Printf("\033[1A")
		p.WriteTable()
	}
	return len(msg), nil
}

func (p *PulumiActionOutputLogger) WriteTable() {
	p.TableWriteInProgress = true
	// Create the table.
	table := utils.CreateTerminalTable()

	rowLength := len(p.Rows)
	if p.NewRowsAdded > 0 {
		oldRows := rowLength - p.NewRowsAdded
		tm.MoveCursorUp(oldRows + 2)
		p.NewRowsAdded = 0
	} else {
		tm.MoveCursorUp(rowLength + 2)
		p.FirstPrintDone = true
	}

	// Add the table headers.
	table.AddRow("\033[2K\r    %s\n", strings.Join(p.Headers, "\t"))

	// Sort the table rows so we can keep consistent ordering. Otherwise it can be hard
	// to track what is actually happenning.
	var sortedRows []PulumiActionResourceOutput
	for _, row := range p.Rows {
		sortedRows = append(sortedRows, row)
	}
	sort.Slice(sortedRows, func(i, j int) bool {
		return sortedRows[i].Order < sortedRows[j].Order
	})

	// Add the table rows.
	for _, row := range sortedRows {
		table.AddRow("\033[2K\r    %s\n", strings.Join(row.Parts, "\t"))
	}

	tm.Printf("%s\n", table.Table)
	tm.Flush()
	p.TableWriteInProgress = false
}

func (p *PulumiActionOutputLogger) AddRow(msg string) {
	msgParts := strings.Split(msg, " ")

	if len(msgParts) == 7 {

		urn := msgParts[3]
		name := msgParts[4]
		status := msgParts[5]

		var colorfulStatus string
		switch status {
		case "creating", "deleting", "updating":
			colorfulStatus = utils.TextColor(status, color.FgYellow)
			break
		case "created", "updated":
			colorfulStatus = utils.TextColor(status, color.FgGreen)
			break
		case "deleted":
			colorfulStatus = utils.TextColor(status, color.FgRed)
		default:
			colorfulStatus = status
		}

		// The Pulumi Stack resource never reports back as created from the streaming
		// ourput so lets ignore that row from our table.
		if urn == "pulumi:pulumi:Stack" {
			return
		}

		hash := utils.HashStringMD5(fmt.Sprintf("%s%s", name, urn))
		if val, ok := p.Rows[hash]; ok {
			p.Rows[hash] = PulumiActionResourceOutput{
				Order: val.Order,
				Parts: []string{name, urn, colorfulStatus},
			}

		} else {
			p.NewRowsAdded += 1
			p.Rows[hash] = PulumiActionResourceOutput{
				Order: len(p.Rows) + 1,
				Parts: []string{name, urn, colorfulStatus},
			}
		}
	}
}

// outputActionSteps prints out the steps of an action to the terminal.
func outputActionSteps(steps []auto.PreviewStep) {
	utils.Print(utils.TextColor("\nPreview Results:\n", color.FgGreen))

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

// Fail handles the failure operations like cleaning the tmp dir.
func (a *PulumiAction) Fail(err error, message string) error {
	a.TemporaryDirectory.Clean()
	return utils.NewErrorMessage(message, err)
}

// SetUp sets up the action to run.
func (a *PulumiAction) SetUp(ctx context.Context, configPath string) error {
	// Create a spinner to show the user what is happening.
	checkEnvSpinner := utils.CreateNewTerminalSpinner(
		"Checking environment",
		"Environment check complete.",
		"Environment check failed.",
	)

	// Read config in.
	projectConfig, err := utils.ReadConfigFile(configPath, a.Environment)
	if err != nil {
		checkEnvSpinner.Fail()
		return err
	}

	// Create a tmp directory for compiliing the TypeScript lambdas.
	tmp := &utils.TemporaryDirectory{Name: tmpDirName}
	err = tmp.Create()
	if err != nil {
		return checkEnvSpinner.FailWithMessage("Error creating tmp directory", err)
	}
	a.TemporaryDirectory = tmp

	// Stop the environment spinner and start the API deployment function
	// spinner.
	checkEnvSpinner.Stop()
	createAPISpinner := utils.CreateNewTerminalSpinner(
		"Generating Pulumi Program",
		"Pulumi Program generated.",
		"Pulumi Program failed to generate.",
	)

	stack, err := a.CreateDeployment(projectConfig.Environment)
	if err != nil {
		return createAPISpinner.FailWithMessage("Error creating API Deployment", err)
	}

	// Set the stack.
	a.Stack = stack

	// Stop the API Deployment spinner and start a spinner for setting
	// up the preview.
	createAPISpinner.Stop()
	settingUpPreviewSpinner := utils.CreateNewTerminalSpinner(
		"Setting up program execution environment",
		"Program execution environment set up successfully.",
		"Failed to set up program execution environment.",
	)

	// Install plugins and set the config values.
	workspace := stack.Workspace()

	err = workspace.InstallPlugin(ctx, "aws", "v3.2.1")
	if err != nil {
		return settingUpPreviewSpinner.FailWithMessage("Failed to install program plugins", err)
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
	actionSpinner := utils.CreateNewTerminalSpinner(
		"Running infrastructure preview",
		"Preview completed successfully.",
		"Preview failed.",
	)

	// Run the action.
	result, err := a.Stack.Preview(ctx)
	if err != nil {
		actionSpinner.Fail()
		return a.Fail(err, "Error running stack preview")
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
	actionSpinner := utils.CreateNewTerminalSpinner(
		"Running infrastructure update",
		"Update completed successfully.",
		"Update failed.",
	)
	actionSpinner.SetOutput(os.Stdout)
	outputLogger := &PulumiActionOutputLogger{
		Headers: []string{"Name", "URN", "Status"},
		Rows:    make(map[string]PulumiActionResourceOutput),
	}

	// The terminal table constantly replaces itself so we need to create
	// two buffer lines to handle the behavior.
	utils.Print("")
	utils.Print("")
	utils.Print("")

	// Stream the results to the terminal.
	outputStreamer := optup.ProgressStreams(outputLogger)

	// Run the action.
	result, err := a.Stack.Up(ctx, outputStreamer)
	if err != nil {
		actionSpinner.Fail()
		return a.Fail(err, "Error running stack update")
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
	actionSpinner := utils.CreateNewTerminalSpinner(
		"Running infrastructure destroy",
		"Destroy completed successfully.",
		"Destroy failed.",
	)
	outputLogger := PulumiActionOutputLogger{
		Headers: []string{"Name", "URN", "Status"},
		Rows:    make(map[string]PulumiActionResourceOutput),
	}

	// The terminal table constantly replaces itself so we need to create
	// two buffer lines to handle the behavior.
	utils.Print("")
	utils.Print("")
	utils.Print("")

	// Stream the results to the terminal.
	outputStreamer := optdestroy.ProgressStreams(&outputLogger)

	// Run the action.
	_, err := a.Stack.Destroy(ctx, outputStreamer)
	if err != nil {
		actionSpinner.Fail()
		return a.Fail(err, "Error running stack destroy")
	}

	// Stop the spinner and clean the temp directory.
	actionSpinner.Stop()
	a.TemporaryDirectory.Clean()
	return nil
}
