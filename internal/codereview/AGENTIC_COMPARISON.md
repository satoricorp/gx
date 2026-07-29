# Agentic Review Comparison

Agentic review is intentionally opt-in. Do not make it the default until repeated comparisons show better findings without unacceptable latency or false-positive cost.

## Manual Benchmark

Run both modes against the same pull request and commit range:

```sh
TOTALITY_REVIEW_AGENTIC=0 tl review --verbose > /tmp/tl-review-single-shot.md
TOTALITY_REVIEW_AGENTIC=1 tl review --verbose > /tmp/tl-review-agentic.md
```

Record:

- Findings kept by both modes.
- Findings found only by single-shot mode.
- Findings found only by agentic mode.
- False positives in each mode.
- Whether anchors map to useful changed lines.
- Wall-clock latency for each run.

Agentic mode should remain disabled by default unless it consistently improves finding quality enough to justify the added tool-loop latency.
