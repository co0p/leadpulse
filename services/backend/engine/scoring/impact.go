package scoring

// ComputeImpactWeighted computes the impact-weighted score for a signal.
// Takes the four impact ratings (IG, IP, IT, IO) and returns a 0-100 score.
//
// PRD 5.3: ImpactWeighted = 0.20*IG + 0.35*IP + 0.25*IT + 0.20*IO (0–5)
//
//	ImpactWeighted100 = ImpactWeighted * 20 (0–100)
//
// Impact layer weights (default for all signals):
// - Growth (IG) 20%
// - Project (IP) 35%
// - Team (IT) 25%
// - Organization (IO) 20%
func ComputeImpactWeighted(ig, ip, it, io int) float64 {
	// Convert to float and apply weights
	weighted := 0.20*float64(ig) + 0.35*float64(ip) + 0.25*float64(it) + 0.20*float64(io)
	// Convert 0-5 scale to 0-100 scale
	return weighted * 20.0
}
