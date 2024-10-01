package repo

type Repo interface {
	GetData(data RepoData) (*RepoData, error)
	SetData(data RepoData) error

	LoginUser(user, pass string) error
	RegisterUser(user, pass string) error

	DownloadStorage(user string) ([]RepoData, error)

	Ping() error
}
