package command

import (
	"strings"

	"github.com/jrasell/levant/levant/structs"
	"github.com/jrasell/levant/logging"
	"github.com/jrasell/levant/scale"
)

// ScaleOutCommand is the command implementation that allows users to scale a
// Nomad job out.
type ScaleOutCommand struct {
	Meta
}

// Help provides the help information for the scale-out command.
func (c *ScaleOutCommand) Help() string {
	helpText := `
Usage: levant scale-out [options] <job-id>

  Scale a Nomad job and optional task group out.

General Options:

  -address=<http_address>
    The Nomad HTTP API address including port which Levant will use to make
    calls.
  
  -log-level=<level>
    Specify the verbosity level of Levant's logs. Valid values include DEBUG,
    INFO, and WARN, in decreasing order of verbosity. The default is INFO.
  
  -log-format=<format>
    Specify the format of Levant's logs. Valid values are HUMAN or JSON. The
    default is HUMAN.
	
Scale Out Options:

  -count=<num>
    The count by which the job and task groups should be scaled out by. Only
    one of count or percent can be passed.

  -percent=<num>
    A percentage value by which the job and task groups should be scaled out
    by. Counts will be rounded up, to ensure required capacity is met. Only 
    one of count or percent can be passed.

  -task-group=<name>
    The name of the task group you wish to target for scaling. Is this is not
    speicified all task groups within the job will be scaled.
`
	return strings.TrimSpace(helpText)
}

// Synopsis is provides a brief summary of the scale-out command.
func (c *ScaleOutCommand) Synopsis() string {
	return "Scale out a Nomad job"
}

// Run triggers a run of the Levant scale-out functions.
func (c *ScaleOutCommand) Run(args []string) int {

	var err error
	var logL, logF string

	config := &structs.ScalingConfig{}
	config.Direction = structs.ScalingDirectionOut

	flags := c.Meta.FlagSet("scale-out", FlagSetVars)
	flags.Usage = func() { c.UI.Output(c.Help()) }

	flags.StringVar(&config.Addr, "address", "", "")
	flags.StringVar(&logL, "log-level", "INFO", "")
	flags.StringVar(&logF, "log-format", "HUMAN", "")
	flags.IntVar(&config.Count, "count", 0, "")
	flags.IntVar(&config.Percent, "percent", 0, "")
	flags.StringVar(&config.TaskGroup, "task-group", "", "")

	if err = flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()

	if len(args) != 1 {
		c.UI.Error("This command takes one argument: <job-name>")
		return 1
	}

	config.JobID = args[0]

	if config.Count == 0 && config.Percent == 0 || config.Count > 0 && config.Percent > 0 {
		c.UI.Error("You must set either -count or -percent flag to scale-out")
		return 1
	}

	if config.Count > 0 {
		config.DirectionType = structs.ScalingDirectionTypeCount
	}

	if config.Percent > 0 {
		config.DirectionType = structs.ScalingDirectionTypePercent
	}

	if err = logging.SetupLogger(logL, logF); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	success := scale.TriggerScalingEvent(config)
	if !success {
		return 1
	}

	return 0
}
