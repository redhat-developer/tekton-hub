package polling

import (
	"context"

	"github.com/google/go-github/github"
)

// GetDirContents returns the contents of a directory
func GetDirContents(ctx context.Context, client *github.Client, owner string, repo string, path string, options *github.RepositoryContentGetOptions) ([]*github.RepositoryContent, error) {
	_, dirs, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

// GetFileContent returns the contents of a file
func GetFileContent(ctx context.Context, client *github.Client, owner string, repo string, path string, options *github.RepositoryContentGetOptions) (*github.RepositoryContent, error) {
	files, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, err
	}
	return files, nil
}
