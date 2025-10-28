package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	commands       []string
	commandScripts []string
	force          bool
)

var rootCmd = &cobra.Command{
	Use:     "rena [flags] [file]...",
	Short:   "A simple utility to rename multiple files",
	Args:    cobra.ArbitraryArgs,
	Version: "v0.2.0",
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 && cmd.Flags().NArg() == 0 {
			cmd.Help()
			os.Exit(0)
		}

		scripts := []string{}
		for _, script := range commandScripts {
			result, err := readFile(script)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			scripts = append(scripts, result...)
		}

		// commands from flag executed first
		commands = append(scripts, commands...)

	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		parsed, err := ParseCommands(commands)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		old := []file{}
		new := []file{}
		conflictsMap := make(map[string]string)
		for _, v := range args {
			newFile := CreateFile(v)
			oldFile := newFile

			// do transformation
			for _, c := range parsed {
				if c != nil {
					err := c.Execute(&newFile)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			}

			newFilePath := newFile.getFullPath()
			oldFilePath := oldFile.getFullPath()

			if !IsSafeName(newFile) {
				fmt.Printf("Skipped because the end result contains a forbidden character: '%v' <- '%v'\n", newFilePath, oldFilePath)
				continue
			}

			if newFilePath == oldFilePath || newFilePath == "" || oldFilePath == "" {
				continue
			}

			if value, conflict := conflictsMap[newFilePath]; conflict {
				fmt.Println("Conflicts have been detected and must be resolved manually:")
				fmt.Printf("- '%v' -> '%v'\n", value, newFilePath)
				fmt.Printf("- '%v' -> '%v'\n", oldFilePath, newFilePath)
				os.Exit(1)
			}

			old = append(old, oldFile)
			new = append(new, newFile)
			conflictsMap[newFilePath] = oldFilePath
		}

		if len(old) == 0 || len(new) == 0 {
			return
		}

		if !force {
			for i, value := range old {
				newFilePath := new[i].getFullPath()
				oldFilePath := value.getFullPath()
				if _, err := os.Stat(newFilePath); !os.IsNotExist(err) {
					fmt.Printf("overwrite '%v' <- '%v'\n", newFilePath, oldFilePath)
				} else {
					fmt.Printf("rename '%v' -> '%v'\n", oldFilePath, newFilePath)
				}
			}

			fmt.Print("Are you sure? (y/n): ")

			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			input = strings.ToLower(input)

			if !(input == "" || input == "y" || input == "yes") {
				return
			}
		}

		for i, value := range old {
			oldFilePath := value.getFullPath()
			newFilePath := new[i].getFullPath()

			if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
				fmt.Printf("%v: no such file or directory\n", oldFilePath)
				continue
			}

			var message string
			if _, err := os.Stat(newFilePath); !os.IsNotExist(err) {
				message = fmt.Sprintf("overwritten '%v' <- '%v'", newFilePath, oldFilePath)
			} else {
				message = fmt.Sprintf("renamed '%v' -> '%v'", oldFilePath, newFilePath)
			}

			err := RenameFile(value, new[i])
			if err != nil {
				fmt.Printf("%v ('%v' -> '%v')", err, oldFilePath, newFilePath)
				continue
			}

			fmt.Println(message)
		}
	},
}

func init() {
	rootCmd.Flags().StringSliceVarP(&commands, "command", "c", nil, "specify the rename command(s) to execute")
	rootCmd.Flags().StringSliceVarP(&commandScripts, "command-script", "s", nil, "load commands from a script file")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "execute without confirmation prompts")
}

func readFile(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []string{}, err
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
