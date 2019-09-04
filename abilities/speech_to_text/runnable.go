package speech_to_text

import (
	"encoding/json"
	"net/http"

	"github.com/asticode/go-astibob"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// Message names
const (
	eventEnrichOptionsUpdatedMessage = "speech_to_text.enrich.options.updated"
)

type runnable struct {
	*astibob.BaseOperatable
	*astibob.BaseRunnable
	o RunnableOptions
}

type RunnableOptions struct {
	StoreSamples bool `toml:"store_samples"`
}

func NewRunnable(name string, o RunnableOptions) astibob.Runnable {
	// Create runnable
	r := &runnable{
		BaseOperatable: newBaseOperatable(),
		o:              o,
	}

	// Add routes
	r.BaseOperatable.AddRoute("/options/enrich", http.MethodPatch, r.updateEnrichOptions)
	r.BaseOperatable.AddRoute("/references/enrich", http.MethodGet, r.enrichReferences)

	// Set base runnable
	r.BaseRunnable = astibob.NewBaseRunnable(astibob.BaseRunnableOptions{
		Metadata: astibob.Metadata{
			Description: "Executes speech to text analysis when detecting silences in audio samples",
			Name:        name,
		},
	})
	return r
}

type EnrichOptions struct {
	StoreSamples bool `json:"store_samples"`
}

func (r *runnable) enrichReferences(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// Set content type
	rw.Header().Set("Content-Type", "application/json")

	// Write
	astibob.WriteHTTPData(rw, EnrichOptions{
		StoreSamples: r.o.StoreSamples,
	})
}

func (r *runnable) updateEnrichOptions(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// Set content type
	rw.Header().Set("Content-Type", "application/json")

	// Parse body
	var b EnrichOptions
	if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
		astibob.WriteHTTPError(rw, http.StatusBadRequest, errors.Wrap(err, "speech_to_text: parsing enrich options payload failed"))
		return
	}

	// Update options
	r.o.StoreSamples = b.StoreSamples

	// Create message
	m, err := r.newEnrichOptionsUpdatedMessage()
	if err != nil {
		astibob.WriteHTTPError(rw, http.StatusInternalServerError, errors.Wrap(err, "speech_to_text: creating enrich options updated message failed"))
		return
	}

	// Dispatch
	r.Dispatch(m)
}

func (r *runnable) newEnrichOptionsUpdatedMessage() (m *astibob.Message, err error) {
	// Create message
	m = astibob.NewMessage()

	// Set name
	m.Name = eventEnrichOptionsUpdatedMessage

	// Set to
	m.To = &astibob.Identifier{Type: astibob.UIIdentifierType}

	// Marshal
	if m.Payload, err = json.Marshal(EnrichOptions{StoreSamples: r.o.StoreSamples}); err != nil {
		err = errors.Wrap(err, "speech_to_text: marshaling payload failed")
		return
	}
	return
}
