package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type GoogleDrive struct {
	drive.Service
	CredentialDirPath string
}

// New initialize services.
func (gd *GoogleDrive) New() *drive.Service {
	var credentialsDirPath = "credentials"

	if gd.CredentialDirPath != "" {
		credentialsDirPath = gd.CredentialDirPath
	}
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", credentialsDirPath, "credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveScope, drive.DriveFileScope, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	var gdriveClient = &GoogleDrive{
		CredentialDirPath: credentialsDirPath,
	}
	client := gdriveClient.GetClient(config)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return srv
}

// Retrieve a token, saves the token, then returns the generated client.
func (gd *GoogleDrive) GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	var credentialsDirPath = "credentials"

	if gd.CredentialDirPath != "" {
		credentialsDirPath = gd.CredentialDirPath
	}
	tokFile := fmt.Sprintf("%s/%s", credentialsDirPath, "token.json")
	tok, err := gd.TokenFromFile(tokFile)
	if err != nil {
		tok = gd.GetTokenFromWeb(config)
		gd.SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func (gd *GoogleDrive) GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func (gd *GoogleDrive) TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (gd *GoogleDrive) SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// createDir create directory under particular
// Parent ID inside google drive.
func CreateDir(srv *drive.Service, parentId string, folderName string) (*drive.File, error) {
	d := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}
	dir, err := srv.Files.Create(d).Do()
	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}
	return dir, nil
}

// createFile create file(upload) into google drive.
func CreateFile(srv *drive.Service, name string, fileToUpload *os.File, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType:                     "application/tar+gzip",
		Name:                         name,
		Parents:                      []string{parentId},
		CopyRequiresWriterPermission: true,
		WritersCanShare:              true,
	}
	file, err := srv.Files.Create(f).Media(fileToUpload).Do()
	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}
	// Get downloadable link to open in browser
	file, err = srv.Files.Get(file.Id).Fields("id,description,webViewLink,webContentLink,properties,parents").Do()
	if err != nil {
		log.Println("Could not get file: " + err.Error())
		return nil, err
	}
	return file, nil
}
