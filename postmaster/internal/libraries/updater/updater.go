package updater

import (
	"context"
	"log"
	"time"

	semver "github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v68/github"

	"github.com/0xdeafcafe/pillar-box/server/internal/utilities/ptr"
)

const Version string = "0.0.0-local"

type NewVersionAvailableFunc func(name, version, url string)

type Updater struct {
	githubClient *github.Client

	registeredNewVersionAvailableHandler NewVersionAvailableFunc
}

func New() *Updater {
	return &Updater{
		githubClient: github.NewClient(nil),
	}
}

func (u *Updater) RegisterNewVersionAvailableHandler(handler NewVersionAvailableFunc) {
	u.registeredNewVersionAvailableHandler = handler
}

func (u *Updater) StartBackgroundChecker() {
	semverVersion, err := semver.NewVersion(Version)
	if err != nil {
		log.Printf("failed to parse version: %v", err)
		return
	}

	go func() {
		for {
			release, _, err := u.githubClient.Repositories.GetLatestRelease(context.Background(), "0xdeafcafe", "pillar-box")
			if err != nil {
				if errGithub, ok := err.(*github.ErrorResponse); ok {
					if errGithub.Response.StatusCode == 404 {
						log.Println("latest release or repository not found, sleeping for 1 hour")
						time.Sleep(time.Hour)
						continue
					}
				}

				log.Printf("failed to get latest release: %v, sleeping for 1 hour", err)
				time.Sleep(time.Hour)
				continue
			}
			if release == nil {
				log.Println("no release found, sleeping for 1 hour")
				time.Sleep(time.Hour)
				continue
			}
			if release.TagName == nil {
				log.Println("no tag name found, sleeping for 1 hour")
				time.Sleep(time.Hour)
				continue
			}

			latestVersion, err := semver.NewVersion(*release.TagName)
			if err != nil {
				log.Printf("failed to parse latest version: %v, sleeping for 1 hour", err)
				time.Sleep(time.Hour)
				continue
			}

			if semverVersion.LessThan(latestVersion) {
				log.Printf("new version available: %s", latestVersion.String())

				if u.registeredNewVersionAvailableHandler != nil {
					u.registeredNewVersionAvailableHandler(
						ptr.ValueOrDefault(release.Name, "Release"),
						latestVersion.String(),
						ptr.ValueOrZero(release.HTMLURL),
					)
				}
			}

			log.Println("sleeping for 24 hours")
			time.Sleep(24 * time.Hour)
		}
	}()
}
