package domain

type TDocument struct {
	Url            string
	PubDate        uint64
	FetchTime      uint64
	Text           string
	FirstFetchTime uint64 // заполняется процессором
}
