package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "triggr",
	Short: "reactive cli manager/tool??",
	Long:  `treactive cli manager/tool??".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, config_err := cmd.PersistentFlags().GetString("config")

		if config_err != nil {
			return config_err
		}

		if !cmd.Flags().Changed("config") {
			config = getDefaultConfigPath()
		}

		cc, cc_err := readConfigFile(config)

		if cc_err != nil {
			return cc_err
		}

		fmt.Print(cc)

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "specify the config file")
}

func getDefaultConfigPath() string {
	xdg_config_home := "XDG_CONFIG_HOME"
	home := "HOME"

	if _, ok := os.LookupEnv(xdg_config_home); ok {
		return os.ExpandEnv(filepath.Join(fmt.Sprintf("$%v", xdg_config_home), "hypr-go", "triggr.json"))
	}
	return os.ExpandEnv(filepath.Join(fmt.Sprintf("$%v", home), ".config", "hypr-go", "triggr.json"))

}

func readConfigFile(config_path string) (*Config, error) {

	if _, err := os.Stat(config_path); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist and if it does not exists just return the defaults

		return getDefaultConfigVals(), nil
	}

	file, err := os.Open(config_path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var config Config

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

type Config struct {
	Test string `json:"test"`
}

func getDefaultConfigVals() *Config {

	return &Config{
		Test: "",
	}
}
