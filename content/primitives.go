package content

import (
	"bytes"
	"embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/lovyou-ai/site/views"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

//go:embed reference/primitives.md
var primitivesRaw []byte

//go:embed reference/agent-primitives.md
var agentPrimitivesRaw []byte

//go:embed reference/product-layers.md
var productLayersRaw []byte

//go:embed reference/layers/*.md
var layerDocsFS embed.FS

//go:embed reference/grammars/*.md
var grammarsFS embed.FS

var (
	layerHeading = regexp.MustCompile(`^## Layer (\d+): (.+)`)
	groupHeading = regexp.MustCompile(`^### Group \d+ — (.+)`)
	primHeading  = regexp.MustCompile(`^#### (.+)`)
	gapLine      = regexp.MustCompile(`^\*\*Gap from Layer \d+:\*\* (.+)`)
	transLine    = regexp.MustCompile(`^\*\*Transition:\*\* (.+)`)
	tableRow     = regexp.MustCompile(`^\| \*\*(.+?)\*\* \| (.+) \|$`)
	primMD       = goldmark.New(goldmark.WithExtensions(extension.Table))
)

// LoadLayers parses primitives.md + product-layers.md + per-layer docs
// into Layer structs with their primitives.
func LoadLayers() []views.Layer {
	// 1. Parse primitives from primitives.md (Layers 1-13).
	layers := parsePrimitivesFile()

	// 2. Parse Layer 0 from 00-foundation.md (full specs).
	layer0 := parseFoundationFile()
	if layer0 != nil {
		layers = append([]views.Layer{*layer0}, layers...)
	}

	// 3. Enrich layer descriptions from product-layers.md.
	productDescs := parseProductLayers()
	for i := range layers {
		if desc, ok := productDescs[layers[i].Number]; ok {
			layers[i].Description = desc
		}
	}

	// 4. Enrich layer descriptions from per-layer derivation docs.
	for i := range layers {
		if derivation := parseLayerDerivation(layers[i].Number); derivation != "" {
			if layers[i].Description != "" {
				layers[i].Description += "\n" + derivation
			} else {
				layers[i].Description = derivation
			}
		}
	}

	return layers
}

// parsePrimitivesFile extracts Layers 1-13 with full primitive specs.
func parsePrimitivesFile() []views.Layer {
	lines := strings.Split(string(primitivesRaw), "\n")

	var layers []views.Layer
	var curLayer *views.Layer
	var curGroup string
	var curPrim *views.Primitive
	var bodyLines []string
	inPrimitive := false

	flushPrimitive := func() {
		if curPrim != nil && curLayer != nil {
			// Collect any body text (notes after spec table) and render as HTML.
			if len(bodyLines) > 0 {
				md := strings.TrimSpace(strings.Join(bodyLines, "\n"))
				var buf bytes.Buffer
				primMD.Convert([]byte(md), &buf)
				curPrim.Notes = buf.String()
			}
			curLayer.Primitives = append(curLayer.Primitives, *curPrim)
			curPrim = nil
		}
		bodyLines = nil
		inPrimitive = false
	}

	flushLayer := func() {
		flushPrimitive()
		if curLayer != nil {
			layers = append(layers, *curLayer)
			curLayer = nil
		}
	}

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		// Layer heading (skip Layer 0).
		if m := layerHeading.FindStringSubmatch(line); m != nil {
			flushLayer()
			n, _ := strconv.Atoi(m[1])
			if n == 0 {
				continue
			}
			curLayer = &views.Layer{Number: n, Name: m[2]}
			curGroup = ""
			continue
		}

		if m := groupHeading.FindStringSubmatch(line); m != nil {
			flushPrimitive()
			curGroup = m[1]
			continue
		}

		if m := primHeading.FindStringSubmatch(line); m != nil {
			flushPrimitive()
			inPrimitive = true
			name := m[1]
			curPrim = &views.Primitive{
				Name:  name,
				Slug:  slugify(name),
				Group: curGroup,
			}
			if curLayer != nil {
				curPrim.Layer = curLayer.Number
				curPrim.LayerName = curLayer.Name
			}
			continue
		}

		if curLayer != nil && !inPrimitive {
			if m := gapLine.FindStringSubmatch(line); m != nil {
				curLayer.Gap = m[1]
			}
			if m := transLine.FindStringSubmatch(line); m != nil {
				curLayer.Transition = m[1]
			}
		}

		if curPrim != nil {
			if m := tableRow.FindStringSubmatch(line); m != nil {
				key := m[1]
				val := strings.TrimSpace(m[2])
				switch key {
				case "Subscribes to":
					curPrim.SubscribesTo = val
				case "Emits":
					curPrim.Emits = val
				case "Depends on":
					curPrim.DependsOn = val
				case "State":
					curPrim.State = val
				case "Intelligent", "Mechanical", "Both":
					curPrim.Intelligent = key + ": " + val
				}
				continue
			}
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || trimmed == "| | |" || trimmed == "|---|---|" {
				continue
			}
			if !strings.HasPrefix(trimmed, "##") && !strings.HasPrefix(trimmed, "---") && !strings.HasPrefix(trimmed, "**Full specification") {
				if curPrim.Description == "" {
					curPrim.Description = trimmed
				} else {
					bodyLines = append(bodyLines, trimmed)
				}
			}
		}
	}
	flushLayer()
	return layers
}

// parseFoundationFile parses 00-foundation.md for Layer 0 primitives with full specs.
func parseFoundationFile() *views.Layer {
	raw, err := layerDocsFS.ReadFile("reference/layers/00-foundation.md")
	if err != nil {
		return nil
	}

	layer := &views.Layer{
		Number:     0,
		Name:       "Foundation",
		Gap:        "None — this is the base layer.",
		Transition: "Nothing → Something",
	}

	lines := strings.Split(string(raw), "\n")
	var curGroup string
	var curPrim *views.Primitive
	var bodyLines []string

	// Foundation uses ### for primitives (not ####).
	foundationPrimHeading := regexp.MustCompile(`^### (.+)`)
	foundationGroupHeading := regexp.MustCompile(`^## Group \d+ — (.+)`)

	flushPrim := func() {
		if curPrim != nil {
			if len(bodyLines) > 0 {
				md := strings.TrimSpace(strings.Join(bodyLines, "\n"))
				var buf bytes.Buffer
				primMD.Convert([]byte(md), &buf)
				curPrim.Notes = buf.String()
			}
			layer.Primitives = append(layer.Primitives, *curPrim)
			curPrim = nil
		}
		bodyLines = nil
	}

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if m := foundationGroupHeading.FindStringSubmatch(line); m != nil {
			flushPrim()
			curGroup = m[1]
			continue
		}

		if m := foundationPrimHeading.FindStringSubmatch(line); m != nil {
			flushPrim()
			name := m[1]
			curPrim = &views.Primitive{
				Name:      name,
				Slug:      slugify(name),
				Layer:     0,
				LayerName: "Foundation",
				Group:     curGroup,
			}
			continue
		}

		if curPrim != nil {
			if m := tableRow.FindStringSubmatch(line); m != nil {
				key := m[1]
				val := strings.TrimSpace(m[2])
				switch key {
				case "Subscribes to":
					curPrim.SubscribesTo = val
				case "Emits":
					curPrim.Emits = val
				case "Depends on":
					curPrim.DependsOn = val
				case "State":
					curPrim.State = val
				case "Intelligent", "Mechanical", "Both":
					curPrim.Intelligent = key + ": " + val
				}
				continue
			}
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || trimmed == "| | |" || trimmed == "|---|---|" {
				continue
			}
			if !strings.HasPrefix(trimmed, "##") && !strings.HasPrefix(trimmed, "---") {
				if curPrim.Description == "" {
					curPrim.Description = trimmed
				} else {
					bodyLines = append(bodyLines, trimmed)
				}
			}
		}
	}
	flushPrim()

	if len(layer.Primitives) == 0 {
		return nil
	}
	return layer
}

// parseProductLayers extracts per-layer product descriptions from product-layers.md.
// Returns map[layerNumber] → rendered HTML description.
func parseProductLayers() map[int]string {
	lines := strings.Split(string(productLayersRaw), "\n")
	result := make(map[int]string)

	plHeading := regexp.MustCompile(`^### Layer (\d+):`)
	var curNum int = -1
	var curLines []string

	flush := func() {
		if curNum >= 0 && len(curLines) > 0 {
			md := strings.Join(curLines, "\n")
			var buf bytes.Buffer
			primMD.Convert([]byte(md), &buf)
			result[curNum] = buf.String()
		}
		curLines = nil
	}

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if m := plHeading.FindStringSubmatch(line); m != nil {
			flush()
			curNum, _ = strconv.Atoi(m[1])
			continue
		}
		if curNum >= 0 {
			// Stop at next ### or --- separator.
			if strings.HasPrefix(line, "### ") || strings.HasPrefix(line, "## ") {
				flush()
				curNum = -1
				continue
			}
			if strings.TrimSpace(line) == "---" {
				flush()
				curNum = -1
				continue
			}
			curLines = append(curLines, line)
		}
	}
	flush()
	return result
}

// parseLayerDerivation reads the per-layer derivation doc and returns rendered HTML.
func parseLayerDerivation(layerNum int) string {
	filenames := []string{
		"00-foundation.md", "01-agency.md", "02-exchange.md", "03-society.md",
		"04-legal.md", "05-technology.md", "06-information.md", "07-ethics.md",
		"08-identity.md", "09-relationship.md", "10-community.md", "11-culture.md",
		"12-emergence.md", "13-existence.md",
	}
	if layerNum < 0 || layerNum >= len(filenames) {
		return ""
	}

	raw, err := layerDocsFS.ReadFile("reference/layers/" + filenames[layerNum])
	if err != nil {
		return ""
	}

	// Extract the derivation section (everything up to the first primitive spec).
	lines := strings.Split(string(raw), "\n")
	var derivLines []string
	pastTitle := false
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		// Skip the # title.
		if strings.HasPrefix(line, "# ") && !pastTitle {
			pastTitle = true
			continue
		}
		if !pastTitle {
			continue
		}
		// Stop at the first primitive spec (### heading that's not a derivation section).
		// Derivation sections use ## headings; primitive specs use ### headings.
		// For foundation (Layer 0), stop at "## Group 0".
		// For other layers, stop at "### Group 0".
		if layerNum == 0 && strings.HasPrefix(line, "## Group") {
			break
		}
		if layerNum > 0 && strings.HasPrefix(line, "## Primitives") {
			break
		}
		// Also stop at dimensional analysis table continuation (we want the narrative).
		derivLines = append(derivLines, line)
	}

	if len(derivLines) == 0 {
		return ""
	}

	md := strings.Join(derivLines, "\n")
	var buf bytes.Buffer
	primMD.Convert([]byte(md), &buf)
	return buf.String()
}

// LoadAgentPrimitives parses agent-primitives.md into primitives.
func LoadAgentPrimitives() []views.Primitive {
	lines := strings.Split(string(agentPrimitivesRaw), "\n")

	var prims []views.Primitive
	var category string

	catHeading := regexp.MustCompile(`^### (Structural|Operational|Relational|Modal) \((\d+)\)`)
	agentRow := regexp.MustCompile(`^\| \*\*(.+?)\*\* \| (.+?) \| (.+?) \| (.+?) \| (.+?) \| (.+?) \| (.+?) \|$`)

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if m := catHeading.FindStringSubmatch(line); m != nil {
			category = m[1]
			continue
		}

		if m := agentRow.FindStringSubmatch(line); m != nil {
			prims = append(prims, views.Primitive{
				Name:         m[1],
				Slug:         "agent-" + slugify(m[1]),
				Layer:        -1,
				LayerName:    "Agent",
				Group:        category,
				Description:  m[2],
				SubscribesTo: m[3],
				Emits:        m[4],
				DependsOn:    m[5],
				State:        m[6],
				Intelligent:  m[7],
			})
		}
	}

	return prims
}

func slugify(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	var b strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

// LoadGrammars reads composition grammar markdown files.
func LoadGrammars() ([]views.RefPage, error) {
	entries, err := grammarsFS.ReadDir("reference/grammars")
	if err != nil {
		return nil, fmt.Errorf("read grammars dir: %w", err)
	}

	var pages []views.RefPage
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		raw, err := grammarsFS.ReadFile("reference/grammars/" + e.Name())
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", e.Name(), err)
		}

		page, err := parseGrammarPage(e.Name(), raw)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		pages = append(pages, page)
	}

	return pages, nil
}
