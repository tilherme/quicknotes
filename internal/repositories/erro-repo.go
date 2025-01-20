package repositories

type RepoErro struct {
	error
}

func newRepoErro(err error) error {
	return &RepoErro{error: err}
}
