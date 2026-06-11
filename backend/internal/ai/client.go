// Package ai powers the Brand Studio and Legal AI features. The service is
// conditionally enabled: it is only constructed when ANTHROPIC_API_KEY is set,
// matching the platform-provider pattern used elsewhere in this codebase.
//
// Models are routed per task and configured via env vars (AI_MODEL_*).
// "claude-*" models go through the Anthropic SDK; any other model id goes
// through an OpenAI-compatible endpoint (AI_COMPAT_BASE_URL + AI_COMPAT_API_KEY,
// e.g. DeepSeek direct or a US host like Fireworks/Together). If a non-Claude
// model is configured without compat credentials, the task falls back to its
// Claude default so a missing key never breaks the feature.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

// Per-task Claude defaults: Opus for the quality-visible logo design, Sonnet
// for liability-adjacent legal drafting, Haiku for short-form text work.
const (
	defaultModelLogos  = "claude-opus-4-8"
	defaultModelCopy   = "claude-haiku-4-5"
	defaultModelRefine = "claude-haiku-4-5"
	defaultModelLegal  = "claude-sonnet-4-6"
)

type Config struct {
	ModelLogos    string
	ModelCopy     string
	ModelRefine   string
	ModelLegal    string
	CompatBaseURL string // OpenAI-compatible endpoint, e.g. https://api.deepseek.com/v1
	CompatAPIKey  string
}

type Service struct {
	client anthropic.Client
	cfg    Config
	http   *http.Client
}

func New(cfg Config) *Service {
	if cfg.ModelLogos == "" {
		cfg.ModelLogos = defaultModelLogos
	}
	if cfg.ModelCopy == "" {
		cfg.ModelCopy = defaultModelCopy
	}
	if cfg.ModelRefine == "" {
		cfg.ModelRefine = defaultModelRefine
	}
	if cfg.ModelLegal == "" {
		cfg.ModelLegal = defaultModelLegal
	}

	compatReady := cfg.CompatBaseURL != "" && cfg.CompatAPIKey != ""
	fallback := func(model, def, task string) string {
		if isClaude(model) || compatReady {
			return model
		}
		log.Printf("AI: %s model %q needs AI_COMPAT_BASE_URL/AI_COMPAT_API_KEY; falling back to %s", task, model, def)
		return def
	}
	cfg.ModelLogos = fallback(cfg.ModelLogos, defaultModelLogos, "logos")
	cfg.ModelCopy = fallback(cfg.ModelCopy, defaultModelCopy, "copy")
	cfg.ModelRefine = fallback(cfg.ModelRefine, defaultModelRefine, "refine")
	cfg.ModelLegal = fallback(cfg.ModelLegal, defaultModelLegal, "legal")

	log.Printf("AI models — logos: %s, copy: %s, refine: %s, legal: %s",
		cfg.ModelLogos, cfg.ModelCopy, cfg.ModelRefine, cfg.ModelLegal)

	return &Service{
		// Anthropic client reads ANTHROPIC_API_KEY from the environment.
		client: anthropic.NewClient(),
		cfg:    cfg,
		http:   &http.Client{Timeout: 3 * time.Minute},
	}
}

func isClaude(model string) bool {
	return strings.HasPrefix(model, "claude")
}

func (s *Service) complete(ctx context.Context, model, system, prompt string, maxTokens int64) (string, error) {
	if isClaude(model) {
		return s.completeAnthropic(ctx, model, system, prompt, maxTokens)
	}
	return s.completeCompat(ctx, model, system, prompt, maxTokens)
}

func (s *Service) completeAnthropic(ctx context.Context, model, system, prompt string, maxTokens int64) (string, error) {
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: maxTokens,
		System:    []anthropic.TextBlockParam{{Text: system}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	}
	// Adaptive thinking is supported on Opus 4.6+ / Sonnet 4.6 but not Haiku 4.5.
	if !strings.Contains(model, "haiku") {
		params.Thinking = anthropic.ThinkingConfigParamUnion{
			OfAdaptive: &anthropic.ThinkingConfigAdaptiveParam{},
		}
	}

	msg, err := s.client.Messages.New(ctx, params)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	for _, block := range msg.Content {
		if block.Type == "text" {
			out.WriteString(block.Text)
		}
	}
	return out.String(), nil
}

// completeCompat calls an OpenAI-compatible /chat/completions endpoint
// (DeepSeek, Fireworks, Together, ...).
func (s *Service) completeCompat(ctx context.Context, model, system, prompt string, maxTokens int64) (string, error) {
	body, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": prompt},
		},
		"max_tokens": maxTokens,
	})
	if err != nil {
		return "", err
	}

	url := strings.TrimSuffix(s.cfg.CompatBaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.CompatAPIKey)

	resp, err := s.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("compat API %s returned %d: %.300s", model, resp.StatusCode, raw)
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("compat API response parse: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("compat API returned no choices")
	}
	return parsed.Choices[0].Message.Content, nil
}

// LogoConcept is one generated logo idea.
type LogoConcept struct {
	Name      string `json:"name"`
	Rationale string `json:"rationale"`
	SVG       string `json:"svg"`
}

const logoSystem = `You are a senior brand designer for Creatrid, a creator identity platform.
You design logo concepts as clean, self-contained SVG markup.

Rules for every SVG you produce:
- A single <svg> element with viewBox="0 0 120 120", no width/height attributes.
- Vector shapes only: path, circle, rect, ellipse, polygon, line, text, defs, linearGradient, radialGradient, stop, g.
- Absolutely no <script>, <foreignObject>, <image>, event handlers (onload etc.), or external references (href/xlink:href/url() to external resources).
- Distinctive, professional marks — avoid generic clip-art. Limit each mark to 2-4 colors.

Respond with ONLY a JSON array (no markdown fences, no prose) of exactly 4 objects:
[{"name": "...", "rationale": "one sentence", "svg": "<svg ...>...</svg>"}]`

// GenerateLogos returns logo concepts for a creator's brand prompt.
func (s *Service) GenerateLogos(ctx context.Context, brief string) ([]LogoConcept, error) {
	raw, err := s.complete(ctx, s.cfg.ModelLogos, logoSystem, brief, 16000)
	if err != nil {
		return nil, err
	}
	raw = stripFences(raw)
	var concepts []LogoConcept
	if err := json.Unmarshal([]byte(raw), &concepts); err != nil {
		return nil, fmt.Errorf("parse logo concepts: %w", err)
	}
	out := concepts[:0]
	for _, c := range concepts {
		c.SVG = sanitizeSVG(c.SVG)
		if c.SVG != "" {
			out = append(out, c)
		}
	}
	return out, nil
}

const copySystem = `You are a brand copywriter for creators on Creatrid (a verified creator identity platform).
Write sharp, specific marketing copy in the creator's voice. No clichés, no filler.
Format your answer in Markdown with clear section headings.`

// GenerateCopy produces marketing copy (bios, taglines, post copy) from a brief.
func (s *Service) GenerateCopy(ctx context.Context, brief string) (string, error) {
	return s.complete(ctx, s.cfg.ModelCopy, copySystem, brief, 4000)
}

const refineSystem = `You are an editor helping a creator refine their own written content
(descriptions, captions, bios, scripts). Improve clarity, energy, and flow while preserving
the creator's voice and intent. Return only the refined text, optionally followed by a short
"Notes" section in Markdown explaining the key changes.`

// RefineText polishes user-provided creative text per an instruction.
func (s *Service) RefineText(ctx context.Context, text, instruction string) (string, error) {
	prompt := fmt.Sprintf("Instruction: %s\n\n---\n\n%s", instruction, text)
	return s.complete(ctx, s.cfg.ModelRefine, refineSystem, prompt, 8000)
}

const legalSystem = `You are a legal information assistant for independent creators, built into the
Creatrid platform. You help with copyright, IP, licensing, DMCA takedowns, NDAs, and
collaboration agreements — in plain language a non-lawyer can act on.

Constraints:
- You provide general legal information and document drafts, NOT legal advice.
- Always note jurisdiction matters and recommend a licensed attorney for binding decisions.
- When drafting documents, produce complete, usable drafts with [BRACKETED PLACEHOLDERS] for
  party-specific details.
- Where relevant, mention that Creatrid's blockchain content anchors can serve as timestamped
  evidence of authorship.

Format answers in Markdown. End every response with: "_This is general information, not legal
advice. Consult a licensed attorney for your specific situation._"`

// LegalAssist answers a creator's legal question or drafts a document.
func (s *Service) LegalAssist(ctx context.Context, question string) (string, error) {
	return s.complete(ctx, s.cfg.ModelLegal, legalSystem, question, 8000)
}

// stripFences removes a ```json ... ``` wrapper if the model added one.
func stripFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		if i := strings.Index(s, "\n"); i >= 0 {
			s = s[i+1:]
		}
		s = strings.TrimSuffix(strings.TrimSpace(s), "```")
	}
	return strings.TrimSpace(s)
}

var (
	svgDangerous = regexp.MustCompile(`(?is)<\s*(script|foreignObject|image|iframe|embed|object|use)\b|on[a-z]+\s*=|javascript:|xlink:href|\bhref\s*=`)
	svgShape     = regexp.MustCompile(`(?is)^<svg\b[^>]*>.*</svg>$`)
)

// sanitizeSVG rejects any SVG containing scriptable or external-reference
// content. Model-generated markup is rendered inline in the frontend, so this
// is a hard allow/deny gate rather than a rewriter: suspicious input returns "".
func sanitizeSVG(svg string) string {
	svg = strings.TrimSpace(svg)
	if !svgShape.MatchString(svg) || svgDangerous.MatchString(svg) {
		return ""
	}
	return svg
}
