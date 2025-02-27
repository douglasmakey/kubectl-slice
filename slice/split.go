package slice

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

const DefaultTemplateName = "{{.kind | lower}}-{{.metadata.name}}.yaml"

// Split is a Kubernetes Split instance. Each instance has its own template
// used to generate the resource names when saving to disk. Because of this,
// avoid reusing the same instance of Split
type Split struct {
	opts     Options
	log      *log.Logger
	template *template.Template
	data     *bytes.Buffer

	filesFound []yamlFile
	fileCount  int
}

// New creates a new Split instance with the options set
func New(opts Options) (*Split, error) {
	s := &Split{
		log: log.New(ioutil.Discard, "[debug] ", log.Lshortfile),
	}

	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}

	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}

	if opts.DebugMode {
		s.log.SetOutput(opts.Stderr)
	}

	s.opts = opts

	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

// Options holds the Split options used when splitting Kubernetes resources
type Options struct {
	Stdout io.Writer
	Stderr io.Writer

	InputFile       string // the name of the input file to be read
	OutputDirectory string // the path to the directory where the files will be stored
	OutputToStdout  bool   // if true, the output will be written to stdout instead of a file
	GoTemplate      string // the go template code to render the file names
	DryRun          bool   // if true, no files are created
	DebugMode       bool   // enables debug mode
	Quiet           bool   // disables all writing to stdout/stderr

	IncludedKinds    []string
	ExcludedKinds    []string
	IncludedNames    []string
	ExcludedNames    []string
	Included         []string
	Excluded         []string
	StrictKubernetes bool // if true, any YAMLs that don't contain at least an "apiVersion", "kind" and "metadata.name" will be excluded

	SortByKind bool // if true, it will sort the resources by kind
}
