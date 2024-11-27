package model

import "strconv"

type Feature struct {
	Index   string `json:"-"`
	Keyword string `json:"key"`
	Url     string `json:"url"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
}

func FeatureFromList(rows []string) Feature {
	return Feature{
		Index:   rows[0],
		Keyword: rows[1],
		Url:     rows[2],
		Title:   rows[3],
		Desc:    rows[4],
	}
}

func ListFromFeature(data Feature) []string {
	return []string{
		data.Index,
		data.Keyword,
		data.Url,
		data.Title,
		data.Desc,
	}
}

func FeaturesFromList(rows [][]string) []Feature {
	result := make([]Feature, 0)
	for _, row := range rows {
		result = append(result, Feature{
			Keyword: row[1],
			Url:     row[2],
			Title:   row[3],
			Desc:    row[4],
		})
	}
	return result
}

func ListFromFeatures(data []Feature) [][]string {
	result := make([][]string, 0)
	for idx, feature := range data {
		result = append(result, []string{
			strconv.Itoa(idx + 1),
			feature.Keyword,
			feature.Url,
			feature.Title,
			feature.Desc,
		})
	}
	return result
}
