package test

import (
	"testing"
)

func TestCloseIncidentTest(t *testing.T) {
	// _, _, rep := internal.New()
	// channelId := strconv.Itoa(rand.Intn(100000))
	// newinc := model.Incident{ChannelId: channelId, DescriptionResolved: "test"}

	// _, err := rep.InsertIncident(&newinc)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// now := time.Now()
	// channelID := "channelID test"
	// responsibility := "responsibility test"
	// team := "team test"
	// rootCause := "rootCause test"
	// feature := "feature test"
	// startDate := time.Now()
	// impact := int64(2)
	// severly := int64(3)
	// err = rep.CloseIncident(channelID, responsibility, rootCause, feature, team, impact, severly, startDate)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// incDb, err := rep.GetIncident(channelId)

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if incDb.CustomerImpact != impact ||
	// 	incDb.SeverityLevel != severly ||
	// 	incDb.StartTimestamp.Equal(now) {
	// 	t.Fatal("Not updated states in database")
	// }
}

func TestCanceIncidentTest(t *testing.T) {
	// _, _, rep := internal.New()
	// channelId := strconv.Itoa(rand.Intn(100000))
	// newinc := model.Incident{ChannelName: channelId, DescriptionCancelled: "test"}

	// _, err := rep.InsertIncident(&newinc)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// err = rep.CancelIncident(channelId, "cancelamento teste")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// incDb, err := rep.GetIncident(channelId)

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if incDb.Status != model.StatusCancel || incDb.DescriptionCancelled != "cancelamento teste" {
	// 	t.Fatal("Not updated states in database")
	// }
}

func TestGetIncidentTest(t *testing.T) {
	// _, _, rep := internal.New()
	// incidents, err := rep.ListActiveIncidents()

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if len(incidents) == 0 {
	// 	t.Fatal("Not incidents for execute integrated tests")
	// }

	// inc := incidents[6]
	// incResult, err := rep.GetIncident(inc.ChannelName)
	// if incResult.ChannelId != inc.ChannelId {
	// 	t.Fatal("incident not equals")
	// }
}
