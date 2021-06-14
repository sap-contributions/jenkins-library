// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/spf13/cobra"
)

type npmExecuteScriptsOptions struct {
	Install                    bool     `json:"install,omitempty"`
	RunScripts                 []string `json:"runScripts,omitempty"`
	DefaultNpmRegistry         string   `json:"defaultNpmRegistry,omitempty"`
	VirtualFrameBuffer         bool     `json:"virtualFrameBuffer,omitempty"`
	ScriptOptions              []string `json:"scriptOptions,omitempty"`
	BuildDescriptorExcludeList []string `json:"buildDescriptorExcludeList,omitempty"`
	BuildDescriptorList        []string `json:"buildDescriptorList,omitempty"`
	CreateBOM                  bool     `json:"createBOM,omitempty"`
	Publish                    bool     `json:"publish,omitempty"`
	RepositoryURL              string   `json:"repositoryUrl,omitempty"`
	RepositoryPassword         string   `json:"repositoryPassword,omitempty"`
	RepositoryUsername         string   `json:"repositoryUsername,omitempty"`
}

// NpmExecuteScriptsCommand Execute npm run scripts on all npm packages in a project
func NpmExecuteScriptsCommand() *cobra.Command {
	const STEP_NAME = "npmExecuteScripts"

	metadata := npmExecuteScriptsMetadata()
	var stepConfig npmExecuteScriptsOptions
	var startTime time.Time
	var logCollector *log.CollectorHook

	var createNpmExecuteScriptsCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Execute npm run scripts on all npm packages in a project",
		Long:  `Execute npm run scripts in all package json files, if they implement the scripts.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.RepositoryPassword)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			telemetryData := telemetry.CustomData{}
			telemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				telemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				telemetryData.ErrorCategory = log.GetErrorCategory().String()
				telemetry.Send(&telemetryData)
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunk.Send(&telemetryData, logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetry.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunk.Initialize(GeneralConfig.CorrelationID,
					GeneralConfig.HookConfig.SplunkConfig.Dsn,
					GeneralConfig.HookConfig.SplunkConfig.Token,
					GeneralConfig.HookConfig.SplunkConfig.Index,
					GeneralConfig.HookConfig.SplunkConfig.SendLogs)
			}
			npmExecuteScripts(stepConfig, &telemetryData)
			telemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addNpmExecuteScriptsFlags(createNpmExecuteScriptsCmd, &stepConfig)
	return createNpmExecuteScriptsCmd
}

func addNpmExecuteScriptsFlags(cmd *cobra.Command, stepConfig *npmExecuteScriptsOptions) {
	cmd.Flags().BoolVar(&stepConfig.Install, "install", true, "Run npm install or similar commands depending on the project structure.")
	cmd.Flags().StringSliceVar(&stepConfig.RunScripts, "runScripts", []string{}, "List of additional run scripts to execute from package.json.")
	cmd.Flags().StringVar(&stepConfig.DefaultNpmRegistry, "defaultNpmRegistry", os.Getenv("PIPER_defaultNpmRegistry"), "URL of the npm registry to use. Defaults to https://registry.npmjs.org/")
	cmd.Flags().BoolVar(&stepConfig.VirtualFrameBuffer, "virtualFrameBuffer", false, "(Linux only) Start a virtual frame buffer in the background. This allows you to run a web browser without the need for an X server. Note that xvfb needs to be installed in the execution environment.")
	cmd.Flags().StringSliceVar(&stepConfig.ScriptOptions, "scriptOptions", []string{}, "Options are passed to all runScripts calls separated by a '--'. './piper npmExecuteScripts --runScripts ci-e2e --scriptOptions '--tag1' will correspond to 'npm run ci-e2e -- --tag1'")
	cmd.Flags().StringSliceVar(&stepConfig.BuildDescriptorExcludeList, "buildDescriptorExcludeList", []string{`deployment/**`}, "List of build descriptors and therefore modules to exclude from execution of the npm scripts. The elements can either be a path to the build descriptor or a pattern.")
	cmd.Flags().StringSliceVar(&stepConfig.BuildDescriptorList, "buildDescriptorList", []string{}, "List of build descriptors and therefore modules for execution of the npm scripts. The elements have to be paths to the build descriptors. **If set, buildDescriptorExcludeList will be ignored.**")
	cmd.Flags().BoolVar(&stepConfig.CreateBOM, "createBOM", false, "Create a BOM xml using CycloneDX.")
	cmd.Flags().BoolVar(&stepConfig.Publish, "publish", false, "Configures npm to publish the artifact to a repository.")
	cmd.Flags().StringVar(&stepConfig.RepositoryURL, "repositoryUrl", os.Getenv("PIPER_repositoryUrl"), "")
	cmd.Flags().StringVar(&stepConfig.RepositoryPassword, "repositoryPassword", os.Getenv("PIPER_repositoryPassword"), "Password for the alternative deployment repository to which the project artifacts should be deployed ( other than those specified in <distributionManagement> ). This password will be updated in settings.xml . When no settings.xml is provided a new one is created corresponding with <servers> tag")
	cmd.Flags().StringVar(&stepConfig.RepositoryUsername, "repositoryUsername", os.Getenv("PIPER_repositoryUsername"), "User for the alternative deployment repository to which the project artifacts should be deployed ( other than those specified in <distributionManagement> ). This user will be updated in settings.xml . When no settings.xml is provided a new one is created corresponding with <servers> tag")

}

// retrieve step metadata
func npmExecuteScriptsMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "npmExecuteScripts",
			Aliases:     []config.Alias{{Name: "executeNpm", Deprecated: false}},
			Description: "Execute npm run scripts on all npm packages in a project",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Resources: []config.StepResources{
					{Name: "source", Type: "stash"},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "install",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "runScripts",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "defaultNpmRegistry",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "GENERAL", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "npm/defaultNpmRegistry"}},
					},
					{
						Name:        "virtualFrameBuffer",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "scriptOptions",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "buildDescriptorExcludeList",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "buildDescriptorList",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "createBOM",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "STEPS", "STAGES", "PARAMETERS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "publish",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"STEPS", "STAGES", "PARAMETERS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name: "repositoryUrl",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/repositoryUrl",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name: "repositoryPassword",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/repositoryPassword",
							},

							{
								Name:  "",
								Paths: []string{"$(vaultPath)/repository-password", "$(vaultBasePath)/$(vaultPipelineName)/repository-password", "$(vaultBasePath)/GROUP-SECRETS/repository-password"},
								Type:  "vaultSecretFile",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name: "repositoryUsername",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/repositoryUsername",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
				},
			},
			Containers: []config.Container{
				{Name: "node", Image: "node:12-buster-slim"},
			},
		},
	}
	return theMetaData
}
