package gomplate

import (
	"io"
	"os"
	"text/template"
	"time"

	"github.com/hairyhenderson/gomplate/data"
)

// Config - values necessary for rendering templates with gomplate.
// Mainly for use by the CLI
type Config struct {
	Input       string
	InputFiles  []string
	InputDir    string
	ExcludeGlob []string
	OutputFiles []string
	OutputDir   string

	DataSources       []string
	DataSourceHeaders []string

	LDelim string
	RDelim string
}

// gomplate -
type gomplate struct {
	funcMap    template.FuncMap
	leftDelim  string
	rightDelim string
}

// runTemplate -
func (g *gomplate) runTemplate(t *tplate) error {
	context := &context{}
	tmpl, err := t.toGoTemplate(g)
	if err != nil {
		return err
	}

	switch t.target.(type) {
	case io.Closer:
		if t.target != os.Stdout {
			// nolint: errcheck
			defer t.target.(io.Closer).Close()
		}
	}
	err = tmpl.Execute(t.target, context)
	return err
}

// newGomplate -
func newGomplate(d *data.Data, leftDelim, rightDelim string) *gomplate {
	return &gomplate{
		leftDelim:  leftDelim,
		rightDelim: rightDelim,
		funcMap:    initFuncs(d),
	}
}

// RunTemplates - run all gomplate templates specified by the given configuration
func RunTemplates(o *Config) error {
	Metrics = newMetrics()
	defer runCleanupHooks()
	d, err := data.NewData(o.DataSources, o.DataSourceHeaders)
	if err != nil {
		return err
	}
	addCleanupHook(d.Cleanup)

	g := newGomplate(d, o.LDelim, o.RDelim)

	return g.runTemplates(o)
}

func (g *gomplate) runTemplates(o *Config) error {
	start := time.Now()
	tmpl, err := gatherTemplates(o)
	Metrics.GatherDuration = time.Since(start)
	if err != nil {
		Metrics.Errors++
		return err
	}
	Metrics.TemplatesGathered = len(tmpl)
	start = time.Now()
	defer func() { Metrics.TotalRenderDuration = time.Since(start) }()
	for _, t := range tmpl {
		tstart := time.Now()
		err := g.runTemplate(t)
		Metrics.RenderDuration[t.name] = time.Since(tstart)
		if err != nil {
			Metrics.Errors++
			return err
		}
		Metrics.TemplatesProcessed++
	}
	return nil
}
