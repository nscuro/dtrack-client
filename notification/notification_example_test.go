package notification_test

import (
	"fmt"
	"os"

	"github.com/nscuro/dtrack-client/notification"
)

// This example demonstrates how to parse and process notifications.
func Example_parse() {
	file, err := os.Open("./testdata/new-vulnerability.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	n, err := notification.Parse(file)
	if err != nil {
		panic(err)
	}

	switch subject := n.Subject.(type) {
	case *notification.NewVulnerabilitySubject:
		fmt.Printf("new vulnerability identified: %s\n", subject.Vulnerability.VulnID)
		for _, project := range subject.AffectedProjects {
			fmt.Printf("=> Project: %s %s\n", project.Name, project.Version)
			fmt.Printf("   Component: %s %s\n", subject.Component.Name, subject.Component.Version)
		}
	}

	// Output:
	// new vulnerability identified: CVE-2012-5784
	// => Project: Acme Example 1.0.0
	//    Component: axis 1.4
}
