package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	cmd := &cobra.Command{
		Use:   "git:newTag",
		Short: "",
		Run:   runGitnewTag,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	// cmd.Flags().StringP("param", "p", "value", "--param=x | -p x")
	appendCommand(cmd)
}

func runGitnewTag(_ *cobra.Command, args []string) {
	// param, _ := cmd.Flags().GetString("param")

	latestTag, err := getLatestTag()
	if err != nil {
		fmt.Println("Error getting latest tag:", err)
		return
	}
	fmt.Println("Latest Tag:\t", latestTag)

	nextTag, err := generateNextTag(latestTag)
	if err != nil {
		fmt.Println("Error generating next tag:", err)
		return
	}
	fmt.Println("Next Tag:\t", nextTag)
	fmt.Println(fmt.Sprintf("Git command:\t git tag %v && git push --tags ", nextTag))
}

func getLatestTag() (string, error) {
	cmd := exec.Command("git", "tag", "--sort=-v:refname")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Split the output by newlines and take the first one (latest tag)
	tags := strings.Split(out.String(), "\n")
	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}
	latestTag := strings.TrimSpace(tags[0])

	return latestTag, nil
}

func generateNextTag(latestTag string) (string, error) {
	// Assuming the tag format is "vX.Y.Z"
	re := regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`)
	matches := re.FindStringSubmatch(latestTag)
	if len(matches) != 4 {
		return "", fmt.Errorf("invalid tag format: %s", latestTag)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", err
	}
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return "", err
	}

	// Increment the patch number
	patch++

	// Format the new tag
	nextTag := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	return nextTag, nil
}
