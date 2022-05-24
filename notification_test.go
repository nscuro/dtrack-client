package dtrack

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNotification(t *testing.T) {
	t.Run("BomConsumed", func(t *testing.T) {
		notification, err := ParseNotification(strings.NewReader(`
		{
		  "notification": {
			"level": "INFORMATIONAL",
			"scope": "PORTFOLIO",
			"group": "BOM_CONSUMED",
			"timestamp": "2019-08-23T21:57:57.418",
			"title": "Bill-of-Materials Consumed",
			"content": "A CycloneDX BOM was consumed and will be processed",
			"subject": {
			  "project": {
				"uuid": "6fb1820f-5280-4577-ac51-40124aabe307",
				"name": "Acme Example",
				"version": "1.0.0"
			  },
			  "bom": {
				"content": "<base64 encoded bom>",
				"format": "CycloneDX",
				"specVersion": "1.1"
			  }
			}
		  }
		}
		`))
		require.NoError(t, err)

		require.IsType(t, &BOMSubject{}, notification.Subject)
		subject := notification.Subject.(*BOMSubject)

		require.Equal(t, "6fb1820f-5280-4577-ac51-40124aabe307", subject.Project.UUID.String())
		require.Equal(t, "<base64 encoded bom>", subject.BOM.Content)
	})

	t.Run("BomProcessed", func(t *testing.T) {
		notification, err := ParseNotification(strings.NewReader(`
		{
		  "notification": {
			"level": "INFORMATIONAL",
			"scope": "PORTFOLIO",
			"group": "BOM_PROCESSED",
			"timestamp": "2019-08-23T21:57:57.418",
			"title": "Bill-of-Materials Consumed",
			"content": "A CycloneDX BOM was consumed and will be processed",
			"subject": {
			  "project": {
				"uuid": "6fb1820f-5280-4577-ac51-40124aabe307",
				"name": "Acme Example",
				"version": "1.0.0"
			  },
			  "bom": {
				"content": "<base64 encoded bom>",
				"format": "CycloneDX",
				"specVersion": "1.1"
			  }
			}
		  }
		}
		`))
		require.NoError(t, err)

		require.IsType(t, &BOMSubject{}, notification.Subject)
		subject := notification.Subject.(*BOMSubject)

		require.Equal(t, "6fb1820f-5280-4577-ac51-40124aabe307", subject.Project.UUID.String())
		require.Equal(t, "<base64 encoded bom>", subject.BOM.Content)
	})

	t.Run("NewVulnerableDependency", func(t *testing.T) {
		notification, err := ParseNotification(strings.NewReader(`
		{
		  "notification": {
			"level": "INFORMATIONAL",
			"scope": "PORTFOLIO",
			"group": "NEW_VULNERABLE_DEPENDENCY",
			"timestamp": "2018-08-27T23:23:00.145",
			"title": "Vulnerable Dependency Introduced",
			"content": "A dependency was introduced that contains 1 known vulnerability",
			"subject": {
			  "project": {
				"uuid": "6fb1820f-5280-4577-ac51-40124aabe307",
				"name": "Acme Example",
				"version": "1.0.0"
			  },
			  "component": {
				"uuid": "4d5cd8df-cff7-4212-a038-91ae4ab79396",
				"group": "apache",
				"name": "axis",
				"version": "1.4",
				"md5": "03dcfdd88502505cc5a805a128bfdd8d",
				"sha1": "94a9ce681a42d0352b3ad22659f67835e560d107",
				"sha256": "05aebb421d0615875b4bf03497e041fe861bf0556c3045d8dda47e29241ffdd3",
				"purl": "pkg:maven/apache/axis@1.4"
			  },
			  "vulnerabilities": [
				{
				  "uuid": "941a93f5-e06b-4304-84de-4d788eeb4969",
				  "vulnId": "CVE-2012-5784",
				  "source": "NVD",
				  "description": "Apache Axis 1.4 and earlier, as used in PayPal Payments Pro, PayPal Mass Pay, PayPal Transactional Information SOAP, the Java Message Service implementation in Apache ActiveMQ, and other products, does not verify that the server hostname matches a domain name in the subject's Common Name (CN) or subjectAltName field of the X.509 certificate, which allows man-in-the-middle attackers to spoof SSL servers via an arbitrary valid certificate.",
				  "cvssv2": 5.8,
				  "severity": "MEDIUM",
				  "cwes": [
					{
				  		"cweId": 20,
				  		"name": "Improper Input Validation"
					},
					{
						"cweId": 66,
						"name": "Foobar"
					}
				  ]
				},
				{
				  "uuid": "ca318ca7-616f-4af0-9c6b-15b8e208c586",
				  "vulnId": "CVE-2014-3596",
				  "source": "NVD",
				  "description": "The getCN function in Apache Axis 1.4 and earlier does not properly verify that the server hostname matches a domain name in the subject's Common Name (CN) or subjectAltName field of the X.509 certificate, which allows man-in-the-middle attackers to spoof SSL servers via a certificate with a subject that specifies a common name in a field that is not the CN field.  NOTE: this issue exists because of an incomplete fix for CVE-2012-5784.\n\n<a href=\"http://cwe.mitre.org/data/definitions/297.html\" target=\"_blank\">CWE-297: Improper Validation of Certificate with Host Mismatch</a>",
				  "cvssv2": 5.8,
				  "severity": "MEDIUM"
				}
			  ]
			}
		  }
		}
		`))
		require.NoError(t, err)

		require.IsType(t, &NewVulnerableDependencySubject{}, notification.Subject)
		subject := notification.Subject.(*NewVulnerableDependencySubject)

		require.Equal(t, "6fb1820f-5280-4577-ac51-40124aabe307", subject.Project.UUID.String())
		require.Equal(t, "4d5cd8df-cff7-4212-a038-91ae4ab79396", subject.Component.UUID.String())
		require.Len(t, subject.Vulnerabilities, 2)
		require.Equal(t, "941a93f5-e06b-4304-84de-4d788eeb4969", subject.Vulnerabilities[0].UUID.String())
		require.Equal(t, "ca318ca7-616f-4af0-9c6b-15b8e208c586", subject.Vulnerabilities[1].UUID.String())
	})

	t.Run("NewVulnerability", func(t *testing.T) {
		notification, err := ParseNotification(strings.NewReader(`
		{
		  "notification": {
			"level": "INFORMATIONAL",
			"scope": "PORTFOLIO",
			"group": "NEW_VULNERABILITY",
			"timestamp": "2018-08-27T23:26:22.961",
			"title": "New Vulnerability Identified",
			"content": "Apache Axis 1.4 and earlier, as used in PayPal Payments Pro, PayPal Mass Pay, PayPal Transactional Information SOAP, the Java Message Service implementation in Apache ActiveMQ, and other products, does not verify that the server hostname matches a domain name in the subject's Common Name (CN) or subjectAltName field of the X.509 certificate, which allows man-in-the-middle attackers to spoof SSL servers via an arbitrary valid certificate.",
			"subject": {
			  "component": {
				"uuid": "4d5cd8df-cff7-4212-a038-91ae4ab79396",
				"group": "apache",
				"name": "axis",
				"version": "1.4",
				"md5": "03dcfdd88502505cc5a805a128bfdd8d",
				"sha1": "94a9ce681a42d0352b3ad22659f67835e560d107",
				"sha256": "05aebb421d0615875b4bf03497e041fe861bf0556c3045d8dda47e29241ffdd3",
				"purl": "pkg:maven/apache/axis@1.4"
			  },
			  "vulnerability": {
				"uuid": "941a93f5-e06b-4304-84de-4d788eeb4969",
				"vulnId": "CVE-2012-5784",
				"source": "NVD",
				"description": "Apache Axis 1.4 and earlier, as used in PayPal Payments Pro, PayPal Mass Pay, PayPal Transactional Information SOAP, the Java Message Service implementation in Apache ActiveMQ, and other products, does not verify that the server hostname matches a domain name in the subject's Common Name (CN) or subjectAltName field of the X.509 certificate, which allows man-in-the-middle attackers to spoof SSL servers via an arbitrary valid certificate.",
				"cvssv2": 5.8,
				"severity": "MEDIUM",
				"cwes": [
					{
				  		"cweId": 20,
				  		"name": "Improper Input Validation"
					},
					{
						"cweId": 66,
						"name": "Foobar"
					}
				]
			  },
			  "affectedProjects": [
				{
				  "uuid": "6fb1820f-5280-4577-ac51-40124aabe307",
				  "name": "Acme Example",
				  "version": "1.0.0"
				}
			  ]
			}
		  }
		}
		`))
		require.NoError(t, err)

		require.IsType(t, &NewVulnerabilitySubject{}, notification.Subject)
		subject := notification.Subject.(*NewVulnerabilitySubject)

		require.Equal(t, "4d5cd8df-cff7-4212-a038-91ae4ab79396", subject.Component.UUID.String())
		require.Equal(t, "941a93f5-e06b-4304-84de-4d788eeb4969", subject.Vulnerability.UUID.String())
		require.Len(t, subject.AffectedProjects, 1)
		require.Equal(t, "6fb1820f-5280-4577-ac51-40124aabe307", subject.AffectedProjects[0].UUID.String())
	})

	t.Run("PolicyViolation", func(t *testing.T) {
		notification, err := ParseNotification(strings.NewReader(`
		{
		  "notification": {
			"level": "INFORMATIONAL",
			"scope": "PORTFOLIO",
			"group": "POLICY_VIOLATION",
			"timestamp": "2022-05-12T23:07:59.611303",
			"title": "Policy Violation",
			"content": "A operational policy violation occurred",
			"subject": {
			  "project": {
				"uuid": "7a36e5c0-9f09-42dd-b401-360da56c2abe",
				"name": "Acme Example",
				"version": "1.0.0"
			  },
			  "component": {
				"uuid": "4e04c695-9acd-46fc-9bf6-ed23d7eb551e",
				"group": "apache",
				"name": "axis",
				"version": "1.4"
			  },
			  "policyViolation": {
				"uuid": "c82fcb50-029a-4636-a657-96242b20680e",
				"type": "OPERATIONAL",
				"timestamp": "2022-05-12T20:34:46Z",
				"policyCondition": {
				  "uuid": "8e5c0a5b-71fb-45c5-afac-6c6a99742cbe",
				  "subject": "COORDINATES",
				  "operator": "MATCHES",
				  "value": "{\"group\":\"apache\",\"name\":\"axis\",\"version\":\"*\"}",
				  "policy": {
					"uuid": "6d4c7398-689a-4ec7-b5c5-9abb6b5393e9",
					"name": "Banned Components",
					"violationState": "FAIL"
				  }
				}
			  }
			}
		  }
		}
		`))
		require.NoError(t, err)

		require.IsType(t, &PolicyViolationSubject{}, notification.Subject)
		subject := notification.Subject.(*PolicyViolationSubject)

		require.Equal(t, "4e04c695-9acd-46fc-9bf6-ed23d7eb551e", subject.Component.UUID.String())
		require.Equal(t, "7a36e5c0-9f09-42dd-b401-360da56c2abe", subject.Project.UUID.String())
		require.Equal(t, "c82fcb50-029a-4636-a657-96242b20680e", subject.PolicyViolation.UUID.String())
		require.Equal(t, "8e5c0a5b-71fb-45c5-afac-6c6a99742cbe", subject.PolicyViolation.PolicyCondition.UUID.String())
		require.Equal(t, "6d4c7398-689a-4ec7-b5c5-9abb6b5393e9", subject.PolicyViolation.PolicyCondition.Policy.UUID.String())
	})
}
