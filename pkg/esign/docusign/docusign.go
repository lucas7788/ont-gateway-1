package docusign

import (
	"context"
	"fmt"

	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/v2.1/envelopes"
	"github.com/jfcote87/esign/v2.1/model"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
)

// DocuSign model
type DocuSign struct {
	cfg       *esign.JWTConfig
	envelopes *envelopes.Service
}

// New is ctor for DocuSign
func New() *DocuSign {
	conf := config.Load()

	cfg := &esign.JWTConfig{
		IntegratorKey: conf.EsignConfig.DocuConfig.IntegratorKey,
		KeyPairID:     conf.EsignConfig.DocuConfig.KeyPairID,
		PrivateKey:    conf.EsignConfig.DocuConfig.PrivateKey,
		AccountID:     conf.EsignConfig.DocuConfig.APIAccountID,
		IsDemo:        !conf.Prod,
	}

	cred, err := cfg.Credential(conf.EsignConfig.DocuConfig.APIUserID, nil, nil)
	if err != nil {
		panic(fmt.Sprintf("cfg.Credential:%v", err))
	}

	envelopes := envelopes.New(cred)
	return &DocuSign{envelopes: envelopes, cfg: cfg}
}

// SendEnvelope for send envelope
func (docu *DocuSign) SendEnvelope(ctx context.Context, envelopeDefinition *model.EnvelopeDefinition, uploads ...*esign.UploadFile) (*model.EnvelopeSummary, error) {
	return docu.envelopes.Create(envelopeDefinition, uploads...).Do(ctx)
}

// DownloadDocument for download documents by envelope id
func (docu *DocuSign) DownloadDocument(ctx context.Context, envID string) (*esign.Download, error) {
	return docu.envelopes.DocumentsGet("combined", envID).Certificate().Watermark().Do(ctx)
}

// UserConsentURL returns a consent url
func (docu *DocuSign) UserConsentURL() string {
	return docu.cfg.UserConsentURL("http://abc.com")
}
