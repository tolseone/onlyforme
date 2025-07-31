package domain

//type TDocument struct {
//	Url            string
//	PubDate        uint64
//	FetchTime      uint64
//	Text           string
//	FirstFetchTime uint64 // заполняется процессором
//}

type Document struct {
	URL            string // Аббревиатуры в Go пишутся в верхнем регистре
	PubDate        uint64 // Более понятное название
	FetchTime      uint64
	Text           string // "Text" слишком общее
	FirstFetchTime uint64
}
