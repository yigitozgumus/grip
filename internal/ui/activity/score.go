package activity

import (
	"fmt"
	"math"
	"time"

	"github.com/yigitozgumus/grip/internal/config"
)

// CalculateScore computes an activity score for a repository.
// Score factors:
//   - Recency (50%): how recently was the repo updated
//   - Branch activity (30%): implied by local branch count
//   - Remote integration (20%): ratio of remote to local branches
//
// Score range: 0.0 - 5.0
func CalculateScore(repo *config.Repository) float64 {
	if repo == nil {
		return 0.0
	}

	var score float64

	// Factor 1: Recency (0-2.5 points)
	recencyScore := calculateRecencyScore(repo.LastUpdated)
	score += recencyScore * 2.5

	// Factor 2: Branch activity (0-1.5 points)
	branchScore := calculateBranchScore(len(repo.LocalBranches))
	score += branchScore * 1.5

	// Factor 3: Remote integration (0-1.0 points)
	remoteScore := calculateRemoteScore(len(repo.RemoteBranches), len(repo.LocalBranches))
	score += remoteScore * 1.0

	return math.Min(5.0, math.Max(0.0, score))
}

// calculateRecencyScore returns 0-1 based on how recent LastUpdated is.
// Uses exponential decay with ~2 week half-life.
func calculateRecencyScore(lastUpdated time.Time) float64 {
	if lastUpdated.IsZero() {
		return 0.0
	}

	elapsed := time.Since(lastUpdated)
	elapsedDays := elapsed.Hours() / 24.0
	decayConstant := 14.0 // Half-life roughly at 2 weeks

	return math.Exp(-elapsedDays / decayConstant)
}

// calculateBranchScore returns 0-1 based on branch count.
// Uses logarithmic scale: 2 branches = 0.5, 5 branches = 0.8, 10+ = 1.0
func calculateBranchScore(branchCount int) float64 {
	if branchCount <= 1 {
		return 0.2
	}

	return math.Min(1.0, math.Log2(float64(branchCount))/math.Log2(10))
}

// calculateRemoteScore returns 0-1 based on remote branch presence.
func calculateRemoteScore(remoteCount, localCount int) float64 {
	if remoteCount == 0 {
		return 0.0
	}

	if localCount == 0 {
		localCount = 1
	}

	ratio := float64(remoteCount) / float64(localCount)
	return math.Min(1.0, ratio*0.5)
}

// FormatScore returns a formatted score string (e.g., "4.2")
func FormatScore(score float64) string {
	return fmt.Sprintf("%.1f", score)
}
