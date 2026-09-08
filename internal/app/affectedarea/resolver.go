package affectedarea

import "context"

const (
	DecisionMatched   = "matched"
	DecisionNew       = "new"
	DecisionUncertain = "uncertain"

	MatchReasonProximity     = "proximity"
	MatchReasonWithinPolygon = "within_polygon"
	MatchReasonOverlap       = "overlap"
	MatchReasonUnmatchable   = "unmatchable"
	MatchReasonNone          = "none"
)

type ResolvedArea struct {
	Area        *AffectedArea
	Confidence  float64
	Decision    string
	MatchReason string
}

// AffectedAreaResolver decides whether a new report belongs to an existing
// affected area or should create a new one. The spatial matching is a
// placeholder for now — it always creates a new area.
type AffectedAreaResolver struct{}

func NewAffectedAreaResolver() *AffectedAreaResolver {
	return &AffectedAreaResolver{}
}

func (r *AffectedAreaResolver) Resolve(ctx context.Context, disasterID int32, loc LocationPayload) (*ResolvedArea, error) {
	// TODO: spatial matching against existing affected areas.
	//   gps/place_search → point-in-polygon, then proximity to nearest centroid.
	//   map             → polygon overlap ratio.
	//   text            → never matchable.
	//
	// Return DecisionMatched (with the matched area + confidence) when the
	// report clearly belongs to an existing area.
	matchReason := MatchReasonNone
	if loc.Source == SourceText {
		matchReason = MatchReasonUnmatchable
	}

	return &ResolvedArea{
		Confidence:  0,
		Decision:    DecisionNew,
		MatchReason: matchReason,
	}, nil
}
