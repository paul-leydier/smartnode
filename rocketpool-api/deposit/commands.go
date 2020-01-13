package deposit

import (
    "github.com/urfave/cli"

    cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)


// Register deposit commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
    app.Commands = append(app.Commands, cli.Command{
        Name:      name,
        Aliases:   aliases,
        Usage:     "Manage node deposits",
        Subcommands: []cli.Command{

            // Get the current deposit RPL requirement
            cli.Command{
                Name:      "required",
                Aliases:   []string{"q"},
                Usage:     "Get the current RPL requirement information",
                UsageText: "rocketpool deposit required",
                Action: func(c *cli.Context) error {

                    // Validate arguments
                    if err := cliutils.ValidateAPIArgs(c, 0, nil); err != nil {
                        return err
                    }

                    // Run command
                    return getRplRequired(c)

                },
            },

            // Get the current deposit status
            cli.Command{
                Name:      "status",
                Aliases:   []string{"s"},
                Usage:     "Get the current deposit status information",
                UsageText: "rocketpool deposit status",
                Action: func(c *cli.Context) error {

                    // Validate arguments
                    if err := cliutils.ValidateAPIArgs(c, 0, nil); err != nil {
                        return err
                    }

                    // Run command
                    return getDepositStatus(c)

                },
            },

            // Reserve a deposit
            cli.Command{
                Name:      "reserve",
                Aliases:   []string{"r"},
                Usage:     "Reserve a node deposit",
                UsageText: "rocketpool deposit reserve durationID",
                Action: func(c *cli.Context) error {

                    // Arguments
                    var durationId string

                    // Validate arguments
                    if err := cliutils.ValidateAPIArgs(c, 1, func(messages *[]string) {

                        // Get duration ID
                        durationId = c.Args().Get(0)

                    }); err != nil {
                        return err
                    }

                    // Run command
                    return reserveDeposit(c, durationId)

                },
            },

            // Complete a deposit
            cli.Command{
                Name:      "complete",
                Aliases:   []string{"c"},
                Usage:     "Complete a reserved node deposit",
                UsageText: "rocketpool deposit complete",
                Action: func(c *cli.Context) error {

                    // Validate arguments
                    if err := cliutils.ValidateAPIArgs(c, 0, nil); err != nil {
                        return err
                    }

                    // Run command
                    return completeDeposit(c)

                },
            },

            // Cancel a deposit
            cli.Command{
                Name:      "cancel",
                Aliases:   []string{"a"},
                Usage:     "Cancel a reserved node deposit",
                UsageText: "rocketpool deposit cancel",
                Action: func(c *cli.Context) error {

                    // Validate arguments
                    if err := cliutils.ValidateAPIArgs(c, 0, nil); err != nil {
                        return err
                    }

                    // Run command
                    return cancelDeposit(c)

                },
            },

        },
    })
}

