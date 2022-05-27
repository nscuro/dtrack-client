package notification_test

import (
	"fmt"
	"os"

	"github.com/nscuro/dtrack-client/notification"
)

func ExampleParse() {
	file, err := os.Open("./testdata/new-vulnerability.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	n, err := notification.Parse(file)
	if err != nil {
		panic(err)
	}

	subject, ok := n.Subject.(*notification.NewVulnerabilitySubject)
	if !ok {
		panic("unexpected subject type")
	}

	fmt.Println(subject.Vulnerability.VulnID)

	// Output:
	// CVE-2012-5784
}
