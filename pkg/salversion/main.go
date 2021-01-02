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

func GetHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling GET request")
	ctx := context.Background()
	version, err := salVersion(ctx)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(version))
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling POST request")
	ctx := context.Background()
	err := saveData(ctx, r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	version, err := salVersion(ctx)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(version))
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
		return err
	}
	return nil
}

func salVersion(ctx context.Context) (string, error) {
	var currentVersion SalVersion
	doc, found, err := firestore.GetDocument(ctx, "Settings", "CurrentVersion")
	if !found {
		// Don't have a current one, get the newest
		log.Info("No current version found")
		currentVersion, err := getSalVersion(ctx)
		if err != nil {
			return "", err
		}
		return currentVersion.CurrentVersion, nil
	}

	if err != nil {
		// Something went wrong, return an error
		return "", err
	}

	err = doc.DataTo(&currentVersion)
	if err != nil {
		return "", err
	}

	now := time.Now()
	oneDayAgo := now.Add(-24 * time.Hour)
	if currentVersion.LastChecked.Before(oneDayAgo) {
		currentVersion, err := getSalVersion(ctx)
		if err != nil {
			return "", err
		}
		return currentVersion.CurrentVersion, nil
	}

	return currentVersion.CurrentVersion, nil
}

func getSalVersion(ctx context.Context) (SalVersion, error) {
	log.Info("Getting Sal Version")
	url := "https://api.github.com/repos/salopensource/sal/releases"
	var salVersion SalVersion
	var githubReleases []GithubRelease

	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return salVersion, err
	}

	res, err := client.Do(req)
	if err != nil {
		return salVersion, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return salVersion, err
	}

	err = json.Unmarshal(body, &githubReleases)
	if err != nil {
		return salVersion, err
	}

	for _, item := range githubReleases {
		if item.Draft == false && item.PreRelease == false {
			salVersion.CurrentVersion = item.TagName
			salVersion.LastChecked = time.Now()

			m := structs.Map(salVersion)
			err := firestore.SetDocument(ctx, "Settings", "CurrentVersion", m)
			if err != nil {
				return salVersion, err
			}
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
