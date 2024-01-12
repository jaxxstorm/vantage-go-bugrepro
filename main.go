package main

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	modelsv2 "github.com/vantage-sh/vantage-go/vantagev2/models"
	segmentsv2 "github.com/vantage-sh/vantage-go/vantagev2/vantage/segments"
	"go.uber.org/zap"
)

var (
	name  = kingpin.Flag("name", "Segment Name.").Required().Short('n').String()
	token = kingpin.Flag("token", "Vantage token.").OverrideDefaultFromEnvar("VANTAGE_API_TOKEN").Required().Short('t').String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	logger, _ := zap.NewProduction()

	outputLog := logger.Sugar()

	outputLog.Infow("Creating vantage segment", "name", *name)

	params := segmentsv2.NewCreateSegmentParams()
	body := &modelsv2.PostSegments{
		Title:            name,
		Priority:         "100",
		TrackUnallocated: false, // create a segment with value false
	}

	r, err := NewClient("https://api.vantage.sh", *token)
	if err != nil {
		fmt.Println(err)
		return
	}

	params.WithSegments(body)

	outputLog.Infow("Sending POST request", "name", name, "title", body.Title, "priority", body.Priority, "track_unallocated", body.TrackUnallocated)

	out, err := r.V2.Segments.CreateSegment(params, r.Auth)
	if err != nil {
		fmt.Println(err)
		return
	}

	segmentToken := out.Payload.Token

	outputLog.Infow("Post request response", "title", out.Payload.Title, "priority", out.Payload.Priority, "track_unallocated", out.Payload.TrackUnallocated, "token", segmentToken)

	outputLog.Infow("Changing track unallocated to true", "name", name, "token", segmentToken)

	updateParams := segmentsv2.NewUpdateSegmentParams()
	updateParams.SetSegmentToken(segmentToken)


	model := &modelsv2.PutSegments{
		Title:              *name,
		TrackUnallocated:   true,
	}
	updateParams.WithSegments(model)

	outputLog.Infow("Sending PUT request", "name", name, "title", updateParams.Segments.Title, "track_unallocated", updateParams.Segments.TrackUnallocated)

	updateOut, err := r.V2.Segments.UpdateSegment(updateParams, r.Auth)
	if err != nil {
		fmt.Println(err)
		return
	}

	outputLog.Infow("Result of setting track unallocated to true", "title", updateOut.Payload.Title, "priority", updateOut.Payload.Priority, "track_unallocated", updateOut.Payload.TrackUnallocated, "token", updateOut.Payload.Token)
	outputLog.Infow("Changing track unallocated back to false", "name", name, "token", segmentToken)

	model = &modelsv2.PutSegments{
		Title:              *name,
		TrackUnallocated:   false,
	}
	updateParams.WithSegments(model)

	outputLog.Infow("Sending PUT request", "name", name, "title", updateParams.Segments.Title, "track_unallocated", updateParams.Segments.TrackUnallocated)

	finalOut, err := r.V2.Segments.UpdateSegment(updateParams, r.Auth)
	if err != nil {
		fmt.Println(err)
		return
	}

	outputLog.Infow("Final result", "title", finalOut.Payload.Title, "priority", finalOut.Payload.Priority, "track_unallocated", finalOut.Payload.TrackUnallocated, "token", finalOut.Payload.Token)

	if finalOut.Payload.TrackUnallocated {
		panic("track unallocated should be false")
	}



}
