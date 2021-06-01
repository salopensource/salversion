package salversion

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/groob/plist"
	"github.com/salopensource/salversion/pkg/firestore"
)

type SalCheckin struct {
	ID          string    `plist:"id" structs:"id" firestore:"id,omitempty"`
	Machines    int64     `plist:"machines" structs:"machines" firestore:"machines,omitempty"`
	Plugins     []string  `plist:"plugins" structs:"plugins" firestore:"plugins,omitempty"`
	InstallType string    `plist:"install_type" structs:"install_type" firestore:"install_type,omitempty"`
	Database    string    `plist:"database" structs:"database" firestore:"database,omitempty"`
	Version     string    `plist:"version" strcuts:"version" firestore:"version,omitempty"`
	Date        time.Time `plist:"date" strcuts:"date" firestore:"date,omitempty"`
	IP          string    `plist:"ip_address" strcuts:"ip_address" firestore:"ip_address,omitempty"`
}

type SalVersion struct {
	CurrentVersion string    `json:"current_version" structs:"current_version" firestore:"current_version,omitempty"`
	LastChecked    time.Time `json:"last_checked" structs:"last_checked" firestore:"last_checked,omitempty"`
}

type GithubRelease struct {
	PreRelease bool   `json:"prerelease"`
	Draft      bool   `json:"draft"`
	TagName    string `json:"tag_name"`
}

func GetHandler(w http.ResponseWriter, r *http.Request, version string) {
	log.Info("Handling GET request")
	_, _ = w.Write([]byte(version))
}

func PostHandler(w http.ResponseWriter, r *http.Request, version string) {
	log.Info("Handling POST request")
	ctx := context.Background()
	err := saveData(ctx, r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(version))
}

func saveData(ctx context.Context, r *http.Request) error {
	log.Info("Saving data")
	var salCheckin SalCheckin
	err := r.ParseForm()
	if err != nil {
		return errors.Wrap(err, "parse post data")
	}
	x := r.Form.Get("data")
	err = plist.Unmarshal([]byte(x), &salCheckin)
	if err != nil {
		return errors.Wrap(err, "decode plist body from post")
	}

	salCheckin.Date = time.Now()
	salCheckin.IP = ReadUserIP(r)
	salCheckin.ID = uuid.New().String()

	m := structs.Map(salCheckin)
	err = firestore.SetDocument(ctx, "SalCheckins", salCheckin.ID, m)
	if err != nil {
		return errors.Wrap(err, "Write to Firestore")
	}
	return nil
}

// func salVersion(ctx context.Context) (string, error) {

// 	currentVersion, err := GetSalVersion(ctx)
// 	if err != nil {
// 		return "", errors.Wrap(err, "getSalVersion")
// 	}

// 	return currentVersion.CurrentVersion, nil
// }

func GetSalVersion(ctx context.Context) (SalVersion, error) {
	log.Info("Getting Sal Version")
	url := "https://api.github.com/repos/salopensource/sal/releases"
	var salVersion SalVersion
	var githubReleases []GithubRelease

	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return salVersion, errors.Wrap(err, "New Request to github")
	}

	res, err := client.Do(req)
	if err != nil {
		return salVersion, errors.Wrap(err, "Do request to github")
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return salVersion, errors.Wrap(err, "Read body from github releases")
	}

	err = json.Unmarshal(body, &githubReleases)
	if err != nil {
		log.Error(string(body))
		return salVersion, errors.Wrap(err, "Unmarshal json from github releases")
	}

	for _, item := range githubReleases {
		if !item.Draft && !item.PreRelease {
			salVersion.CurrentVersion = item.TagName
			salVersion.LastChecked = time.Now()

			// m := structs.Map(salVersion)
			// err := firestore.SetDocument(ctx, "Settings", "CurrentVersion", m)
			// if err != nil {
			// 	return salVersion, err
			// }
			return salVersion, nil
		}
	}

	return salVersion, nil
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
