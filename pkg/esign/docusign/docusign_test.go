package docusign

import (
	"os"
	"testing"

	"context"
	"io/ioutil"

	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/v2.1/model"
	"github.com/stretchr/testify/assert"
)

func TestDocusign(t *testing.T) {
	env := &model.EnvelopeDefinition{
		EmailSubject: "[Go eSignagure SDK] - Please sign this doc",
		EmailBlurb:   "Please sign this test document",
		Status:       "sent",
		Documents: []model.Document{
			{
				Name:       "invite letter.pdf",
				DocumentID: "1",
			},
			{
				Name:       "contract.pdf",
				DocumentID: "2",
			},
		},
		Recipients: &model.Recipients{
			Signers: []model.Signer{
				{
					Email:             "652732310@qq.com",
					EmailNotification: nil,
					Name:              "J F Cote",
					RecipientID:       "1",
					Tabs: &model.Tabs{
						SignHereTabs: []model.SignHere{
							{
								TabBase: model.TabBase{
									DocumentID:  "1",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "signature",
									XPosition:  "192",
									YPosition:  "160",
								},
							},
						},
						DateSignedTabs: []model.DateSigned{
							{
								TabBase: model.TabBase{
									DocumentID:  "1",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "dateSigned",
									XPosition:  "334",
									YPosition:  "179",
								},
							},
						},
						TextTabs: []model.Text{
							{
								TabBase: model.TabBase{
									DocumentID:  "2",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "txtNote",
									XPosition:  "70",
									YPosition:  "564",
								},
								TabStyle: model.TabStyle{
									Name: "This is the tab tooltip",
								},
								Width:  "300",
								Height: "150",
							},
						},
					},
				},
			},
		},
	}

	// open files for upload
	f1, err := os.Open("letter.pdf")
	assert.Nil(t, err)
	f2, err := os.Open("contract.pdf")
	assert.Nil(t, err)

	url := docu.UserConsentURL()
	assert.True(t, url != "")

	ioutil.WriteFile("/tmp/consent.txt", []byte(url), 0644)

	_, err = docu.SendEnvelope(context.TODO(), env, &esign.UploadFile{
		ContentType: "application/pdf",
		FileName:    "invitation letter.pdf",
		ID:          "1",
		Reader:      f1,
	}, &esign.UploadFile{
		ContentType: "application/pdf",
		FileName:    "contract.pdf",
		ID:          "2",
		Reader:      f2,
	})
	assert.Nil(t, err)
}

var (
	docu *DocuSign
)

func TestMain(m *testing.M) {
	docu = New()
	m.Run()
}
