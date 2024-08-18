/*
https://docs.aws.amazon.com/sdk-for-go/api/service/route53/#pkg-overview
*/
package route53

import (
	cmdArgs "cmd-args"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

var f = fmt.Sprintf
var logger *blablo.Logger

func extractDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return domain
}

func getHostedZoneId(svc *route53.Route53, domain string) (string, error) {
	logger.Debug("Requesting hosted zone list")

	input1 := &route53.ListHostedZonesByNameInput{}
	zones, err := svc.ListHostedZonesByName(input1)
	if err != nil {
		logger.Error(f("Error receiving public ip: %s", c.WithRed(err.Error())))
		return "", err
	}
	logger.Debug(f("Found (%s) zones", c.WithGreen(fmt.Sprint(len(zones.HostedZones)))))

	var mainDomain = extractDomain(domain)
	logger.Debug(f("Looking for hosted zone for main domain: %s", c.WithCyan(mainDomain)))

	for _, zone := range zones.HostedZones {
		if strings.Contains(*zone.Name, mainDomain) {
			return strings.Split(*zone.Id, "/")[2], nil
		}
	}

	return "", nil
}

func UpdateDomainIp(domain, ip, dryRun string) bool {
	logger = blablo.NewLogger("route53", cmdArgs.LogLevel, false)

	logger.Debug("Creating AWS CDK session")
	mySession := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
	}))
	svc := route53.New(mySession)

	awsHostedZoneId, err := getHostedZoneId(svc, domain)
	if err != nil || awsHostedZoneId == "" {
		return false
	}
	logger.Debug(f("Hosted zone ID for (%s): %s", c.WithCyan(domain), c.WithGreen(awsHostedZoneId)))

	logger.Info(f(
		"Request to update 'A' rec for (%s) in (%s) zone",
		c.WithCyan(domain), c.WithCyan(awsHostedZoneId),
	))
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: &domain,
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: &ip,
							},
						},
						TTL:  aws.Int64(60),
						Type: aws.String("A"),
					},
				},
			},
			Comment: aws.String("nrp-cli public IP update"),
		},
		HostedZoneId: &awsHostedZoneId,
	}

	if dryRun == "yes" {
		logger.Info(f("Dry run mode. Skipping making request to update A record"))
		return true
	}

	result, err := svc.ChangeResourceRecordSets(input)
	if err != nil {
		logger.Error(f("Error updating A record: %s", c.WithRed(err.Error())))
	}

	logger.Debug(f("Result: %s%#v%s", c.Yellow, result, c.Reset))

	return true
}
