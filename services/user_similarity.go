package services

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pclubiitk/puppylove_tags/models"
	"gorm.io/gorm"
)

type UserSimilarityService struct {
	db *gorm.DB
}

func NewUserSimilarityService(db *gorm.DB) *UserSimilarityService {
	return &UserSimilarityService{db: db}
}

func (s *UserSimilarityService) UpdateUser(userID string, tags []int) error {
	if err := s.db.Where("user_id = ?", userID).Delete(&models.UserTag{}).Error; err != nil {
		return err
	}
	for _, tag := range tags {
		ut := models.UserTag{
			UserID: userID,
			Tag:    tag,
		}
		if err := s.db.Create(&ut).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *UserSimilarityService) GetUserTags(userID string) ([]int, error) {
	var uts []models.UserTag
	if err := s.db.Where("user_id = ?", userID).Find(&uts).Error; err != nil {
		return nil, err
	}
	var tags []int
	for _, ut := range uts {
		tags = append(tags, ut.Tag)
	}
	return tags, nil
}

type UserScore struct {
	UserID string
	Score  float64
}

func (s *UserSimilarityService) QuerySimilar(userID string, offset, limit int) ([]string, error) {
	queryTags, err := s.GetUserTags(userID)
	if err != nil {
		return nil, err
	}
	if len(queryTags) == 0 {
		return []string{}, nil
	}

	tagStrings := make([]string, len(queryTags))
	for i, tag := range queryTags {
		tagStrings[i] = strconv.Itoa(tag)
	}
	inClause := strings.Join(tagStrings, ",")

	type Candidate struct {
		UserID      string
		CommonCount int
	}

	var candidates []Candidate
	rawQuery := fmt.Sprintf(
		`SELECT user_id, COUNT(*) AS common_count
		 FROM user_tags
		 WHERE tag IN (%s) AND user_id <> ?
		 GROUP BY user_id`, inClause)
	if err := s.db.Raw(rawQuery, userID).Scan(&candidates).Error; err != nil {
		return nil, err
	}

	candidateMap := make(map[string]int)
	for _, c := range candidates {
		candidateMap[c.UserID] = c.CommonCount
	}

	var ranked []UserScore
	for candidate, commonCount := range candidateMap {
		candidateTags, err := s.GetUserTags(candidate)
		if err != nil {
			return nil, err
		}
		unionCount := len(queryTags) + len(candidateTags) - commonCount
		if unionCount == 0 {
			continue
		}
		score := float64(commonCount) / float64(unionCount)
		ranked = append(ranked, UserScore{UserID: candidate, Score: score})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	if offset >= len(ranked) {
		return []string{}, nil
	}
	end := offset + limit
	if end > len(ranked) {
		end = len(ranked)
	}
	var result []string
	for _, r := range ranked[offset:end] {
		result = append(result, r.UserID)
	}
	return result, nil
}
