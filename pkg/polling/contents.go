package polling

import (
	"context"

	"github.com/google/go-github/github"
)

// Get contents of a directory
func GetDirContents(client *github.Client, ctx context.Context, owner string, repo string, path string, options *github.RepositoryContentGetOptions) ([]*github.RepositoryContent, error) {
	_, dirs, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

// Get contents of a file
func GetFileContent(client *github.Client, ctx context.Context, owner string, repo string, path string, options *github.RepositoryContentGetOptions) (*github.RepositoryContent, error) {
	files, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, err
	}
	return files, nil
}
