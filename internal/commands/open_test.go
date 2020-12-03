package commands_test

import (
	"context"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/commands"
	filestorage "hellper/internal/file_storage"
	"hellper/internal/log"
	"hellper/internal/model"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type openCommandFixture struct {
	testName               string
	expectError            bool
	errorMessage           string
	ctx                    context.Context
	mockLogger             log.Logger
	mockClient             bot.Client
	mockIncidentRepository model.IncidentRepository
	mockServiceRepository  model.ServiceRepository
	mockFilestorage        filestorage.Driver
	triggerID              string
	mockDialogSubmission   bot.DialogSubmission
}

func (f *openCommandFixture) setup(t *testing.T) {
	var (
		loggerMock             = log.NewLoggerMock()
		clientMock             = bot.NewClientMock()
		incidentRepositoryMock = model.NewIncidentRepositoryMock()
		serviceRepositoryMock  = model.NewServiceRepositoryMock()
		filestorageMock        = filestorage.NewFileStorageMock()
	)

	f.ctx = context.Background()
	f.mockLogger = loggerMock
	f.mockClient = clientMock
	f.mockIncidentRepository = incidentRepositoryMock
	f.mockServiceRepository = serviceRepositoryMock
	f.mockFilestorage = filestorageMock

	serviceInstances := []*model.ServiceInstance{
		&model.ServiceInstance{ID: 1, Name: "Service1/llama"},
		&model.ServiceInstance{ID: 2, Name: "Service1/alpaca"},
		&model.ServiceInstance{ID: 3, Name: "Service1/horse"},
	}

	loggerMock.On("Info", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	clientMock.On("OpenDialog", f.triggerID, mock.AnythingOfType("slack.Dialog")).Return(nil)
	clientMock.On("AddPin", mock.AnythingOfType("string"), mock.AnythingOfType("slack.ItemRef")).Return(nil)
	clientMock.On("PostMessage", mock.AnythingOfType("string"), mock.AnythingOfType("[]slack.MsgOption")).Return("", "", nil)
	clientMock.On("JoinConversationContext", f.ctx, mock.AnythingOfType("string")).Return(new(slack.Channel), "", []string{}, nil)
	clientMock.On("InviteUsersToConversationContext", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]string")).Return(new(slack.Channel), nil)
	clientMock.On("SetTopicOfConversation", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(new(slack.Channel), nil)
	clientMock.On("GetUserInfoContext", f.ctx, mock.AnythingOfType("string")).Return(new(slack.User), nil)
	clientMock.On("CreateConversationContext", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(new(slack.Channel), nil)
	incidentRepositoryMock.On("InsertIncident", mock.AnythingOfType("*model.Incident")).Return(int64(1), nil)
	incidentRepositoryMock.On("AddPostMortemUrl", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	serviceRepositoryMock.On("ListServiceInstances", f.ctx).Return(serviceInstances, nil)
	filestorageMock.On("CreatePostMortemDocument", f.ctx, mock.AnythingOfType("string")).Return(string(""), nil)
}

func TestOpenIncidentDialog(t *testing.T) {
	var f openCommandFixture
	t.Run("Dialog created properly", func(t *testing.T) {
		f.setup(t)
		err := commands.OpenStartIncidentDialog(f.mockClient, f.mockServiceRepository, f.triggerID)

		if err != nil {
			t.Fatal("an error occurred, but was not expected", "error", err)
		}
	})
}

func TestStartIncidentByDialog(t *testing.T) {
	table := []openCommandFixture{
		{
			testName:    "When the form has all the data to open a new incident",
			expectError: false,
			mockDialogSubmission: bot.DialogSubmission{
				User: bot.User{
					ID: "UYGFQB9C0",
				},
				Submission: bot.Submission{
					IncidentTitle:       "Inc XYZ",
					ChannelName:         "inc-xyz",
					WarRoomURL:          "asd",
					SeverityLevel:       "2",
					Product:             "A",
					IncidentCommander:   "UYGFQB9C0",
					IncidentDescription: "Incident Resolved!",
				},
			},
		},
		{
			testName:    "When the form has no data",
			expectError: false,
			mockDialogSubmission: bot.DialogSubmission{
				User: bot.User{ID: ""},
				Submission: bot.Submission{
					IncidentTitle:       "",
					ChannelName:         "",
					WarRoomURL:          "",
					SeverityLevel:       "0",
					Product:             "",
					IncidentCommander:   "",
					IncidentDescription: "",
				},
			},
		},
		{
			testName:     "When SeverityLevel is not a number",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"High\": invalid syntax",
			mockDialogSubmission: bot.DialogSubmission{
				User: bot.User{ID: ""},
				Submission: bot.Submission{
					IncidentTitle:       "",
					ChannelName:         "",
					WarRoomURL:          "",
					SeverityLevel:       "High",
					Product:             "",
					IncidentCommander:   "",
					IncidentDescription: "",
				},
			},
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.StartIncidentByDialog(f.ctx, f.mockClient, f.mockLogger, f.mockIncidentRepository, f.mockFilestorage, f.mockDialogSubmission)
			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}

				assert.EqualError(t, err, f.errorMessage)
			} else {
				if err != nil {
					t.Fatalf("an error occurred, but was not expected:\n%s", err)
				}
			}
		})
	}
}
