package updater

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	semver "github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v68/github"

	"github.com/0xdeafcafe/pillar-box/server/internal/utilities/ptr"
)

const (
	Version = "0.0.0-local"

	githubOwner = "0xdeafcafe"
	githubRepo  = "pillar-box"
)

type GetPrereleasePreferenceFunc func() bool
type NewVersionAvailableFunc func(name, version, url string)

type Updater struct {
	githubClient *github.Client

	registeredGetPrereleasePreferenceHandler GetPrereleasePreferenceFunc
	registeredNewVersionAvailableHandler     NewVersionAvailableFunc
}

func New() *Updater {
	return &Updater{
		githubClient: github.NewClient(nil),
	}
}

func (u *Updater) RegisterGetPrereleasePreferenceHandler(handler GetPrereleasePreferenceFunc) {
	u.registeredGetPrereleasePreferenceHandler = handler
}

func (u *Updater) RegisterNewVersionAvailableHandler(handler NewVersionAvailableFunc) {
	u.registeredNewVersionAvailableHandler = handler
}

func (u *Updater) CheckForUpdates() error {
	semverVersion, err := semver.NewVersion(Version)
	if err != nil {
		return fmt.Errorf("failed to parse version: %w", err)
	}

	prerelease := false
	if u.registeredGetPrereleasePreferenceHandler != nil {
		prerelease = u.registeredGetPrereleasePreferenceHandler()
	}

	release, err := u.getGitHubRelease(prerelease)
	if err != nil {
		if errGithub, ok := err.(*github.ErrorResponse); ok && errGithub.Response.StatusCode == 404 {
			return errors.New("no release found")
		}

		return err
	}
	if release == nil {
		return errors.New("no release available")
	}
	if release.TagName == nil {
		return errors.New("release tag name is nil")
	}

	latestVersion, err := semver.NewVersion(*release.TagName)
	if err != nil {
		return err
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

	return nil
}

func (u *Updater) StartBackgroundChecker() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from panic: %v", r)
			}
		}()

		for {
			if err := u.CheckForUpdates(); err != nil {
				log.Printf("failed or unable to check for updates, sleeping for an hour: %v", err)
				time.Sleep(time.Hour)

				continue
			}

			log.Println("no update available, sleeping for 24 hours")
			time.Sleep(24 * time.Hour)
		}
	}()
}

func (u *Updater) getGitHubRelease(prerelease bool) (*github.RepositoryRelease, error) {
	if prerelease {
		releases, _, err := u.githubClient.Repositories.ListReleases(context.Background(), githubOwner, githubRepo, nil)
		if err != nil {
			return nil, err
		}
		if len(releases) == 0 {
			return nil, nil
		}

		return releases[0], nil
	}

	// Fetch latest release if not pre-release
	release, _, err := u.githubClient.Repositories.GetLatestRelease(context.Background(), githubOwner, githubRepo)
	if err != nil {
		return nil, err
	}

	return release, nil
}
