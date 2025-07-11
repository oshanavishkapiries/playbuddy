package models

import "time"

// Movie represents a movie from TMDB API
type Movie struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       string  `json:"poster_path"`
	BackdropPath     string  `json:"backdrop_path"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
	Adult            bool    `json:"adult"`
	Video            bool    `json:"video"`
	OriginalLanguage string  `json:"original_language"`
	GenreIDs         []int   `json:"genre_ids"`
}

// TVShow represents a TV show from TMDB API
type TVShow struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	OriginalName     string  `json:"original_name"`
	Overview         string  `json:"overview"`
	FirstAirDate     string  `json:"first_air_date"`
	PosterPath       string  `json:"poster_path"`
	BackdropPath     string  `json:"backdrop_path"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
	OriginalLanguage string  `json:"original_language"`
	GenreIDs         []int   `json:"genre_ids"`
	OriginCountry    []string `json:"origin_country"`
}

// Season represents a TV season
type Season struct {
	ID           int    `json:"id"`
	SeasonNumber int    `json:"season_number"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	AirDate      string `json:"air_date"`
	EpisodeCount int    `json:"episode_count"`
}

// Episode represents a TV episode
type Episode struct {
	ID            int     `json:"id"`
	EpisodeNumber int     `json:"episode_number"`
	Name          string  `json:"name"`
	Overview      string  `json:"overview"`
	AirDate       string  `json:"air_date"`
	StillPath     string  `json:"still_path"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
}

// SearchResponse represents the TMDB search API response for movies
type SearchResponse struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

// TVSearchResponse represents the TMDB search API response for TV shows
type TVSearchResponse struct {
	Page         int      `json:"page"`
	Results      []TVShow `json:"results"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

// TVDetails represents detailed TV show information
type TVDetails struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Overview     string   `json:"overview"`
	FirstAirDate string   `json:"first_air_date"`
	LastAirDate  string   `json:"last_air_date"`
	Seasons      []Season `json:"seasons"`
	NumberOfSeasons int   `json:"number_of_seasons"`
	NumberOfEpisodes int  `json:"number_of_episodes"`
}

// SeasonDetails represents detailed season information
type SeasonDetails struct {
	ID           int       `json:"id"`
	SeasonNumber int       `json:"season_number"`
	Name         string    `json:"name"`
	Overview     string    `json:"overview"`
	AirDate      string    `json:"air_date"`
	Episodes     []Episode `json:"episodes"`
}

// GetReleaseYear returns the year from release date
func (m *Movie) GetReleaseYear() string {
	if m.ReleaseDate == "" {
		return "Unknown"
	}
	
	releaseTime, err := time.Parse("2006-01-02", m.ReleaseDate)
	if err != nil {
		return "Unknown"
	}
	
	return releaseTime.Format("2006")
}

// GetFirstAirYear returns the year from first air date
func (tv *TVShow) GetFirstAirYear() string {
	if tv.FirstAirDate == "" {
		return "Unknown"
	}
	
	airTime, err := time.Parse("2006-01-02", tv.FirstAirDate)
	if err != nil {
		return "Unknown"
	}
	
	return airTime.Format("2006")
}

// GetPosterURL returns the full poster URL
func (m *Movie) GetPosterURL() string {
	if m.PosterPath == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w500" + m.PosterPath
}

// GetBackdropURL returns the full backdrop URL
func (m *Movie) GetBackdropURL() string {
	if m.BackdropPath == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w1280" + m.BackdropPath
}

// GetPosterURL returns the full poster URL for TV shows
func (tv *TVShow) GetPosterURL() string {
	if tv.PosterPath == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w500" + tv.PosterPath
}

// GetBackdropURL returns the full backdrop URL for TV shows
func (tv *TVShow) GetBackdropURL() string {
	if tv.BackdropPath == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w1280" + tv.BackdropPath
} 