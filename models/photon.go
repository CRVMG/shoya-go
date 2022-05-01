package models

import (
	"errors"
	"fmt"
	"time"
)

type PhotonValidateJoinJWTResponse struct {
	Time               string               `json:"time"`
	Valid              bool                 `json:"valid"`
	User               PhotonPropUser       `json:"user"`
	IP                 string               `json:"ip"`
	AvatarDict         PhotonPropAvatarDict `json:"avatarDict"`
	FallbackAvatarDict PhotonPropAvatarDict `json:"favatarDict"`
	WorldCapacity      int                  `json:"worldCapacity,omitempty"`
	WorldAuthor        string               `json:"worldAuthor,omitempty"`
	InstanceCreator    string               `json:"instanceCreator,omitempty"`
}

func (p *PhotonValidateJoinJWTResponse) FillFromUser(u *User) error {
	avatarImageUrl := u.CurrentAvatar.GetImageUrl()
	avatarImageThumbnailUrl := u.CurrentAvatar.GetThumbnailImageUrl()
	profilePicOverride := u.ProfilePicOverride

	if profilePicOverride != "" {
		avatarImageUrl = ""
		avatarImageThumbnailUrl = ""
	}
	p.User = PhotonPropUser{
		ID:                             u.ID,
		DisplayName:                    u.DisplayName,
		DeveloperType:                  u.DeveloperType,
		CurrentAvatarImageUrl:          avatarImageUrl,
		CurrentAvatarThumbnailImageUrl: avatarImageThumbnailUrl,
		UserIcon:                       u.UserIcon,
		LastPlatform:                   u.LastPlatform,
		Status:                         string(u.Status),
		StatusDescription:              u.StatusDescription,
		Bio:                            u.Bio,
		Tags:                           u.Tags,
		AllowAvatarCopying:             u.AllowAvatarCopying,
	}
	currAvAuthor, err := u.CurrentAvatar.GetAuthor()
	if err != nil {
		fmt.Println("error: avatar author was nil")
		return errors.New("avatar author was nil")
	}
	fbAvAuthor, err := u.FallbackAvatar.GetAuthor()
	if err != nil {
		fmt.Println("error: avatar author was nil")
		return errors.New("avatar author was nil")
	}
	p.AvatarDict = PhotonPropAvatarDict{
		ID:                u.CurrentAvatar.ID,
		AssetUrl:          u.CurrentAvatar.GetAssetUrl(),
		AuthorId:          u.CurrentAvatar.AuthorID,
		AuthorName:        currAvAuthor.DisplayName,
		UpdatedAt:         time.Unix(u.CurrentAvatar.CreatedAt, 0).UTC().Format(time.RFC3339Nano),
		Description:       u.CurrentAvatar.Description,
		ImageUrl:          avatarImageUrl,
		ThumbnailImageUrl: avatarImageThumbnailUrl,
		Name:              u.CurrentAvatar.Name,
		ReleaseStatus:     string(u.CurrentAvatar.ReleaseStatus),
		Version:           u.CurrentAvatar.Version,
		Tags:              u.CurrentAvatar.Tags,
		UnityPackages:     u.CurrentAvatar.GetUnityPackages(),
	}
	p.FallbackAvatarDict = PhotonPropAvatarDict{
		ID:                u.FallbackAvatar.ID,
		AssetUrl:          u.FallbackAvatar.GetAssetUrl(),
		AuthorId:          u.FallbackAvatar.AuthorID,
		AuthorName:        fbAvAuthor.DisplayName,
		UpdatedAt:         time.Unix(u.FallbackAvatar.CreatedAt, 0).UTC().Format(time.RFC3339Nano),
		Description:       u.FallbackAvatar.Description,
		ImageUrl:          u.FallbackAvatar.GetImageUrl(),
		ThumbnailImageUrl: u.FallbackAvatar.GetThumbnailImageUrl(),
		Name:              u.FallbackAvatar.Name,
		ReleaseStatus:     string(u.FallbackAvatar.ReleaseStatus),
		Version:           u.FallbackAvatar.Version,
		Tags:              u.FallbackAvatar.Tags,
		UnityPackages:     u.FallbackAvatar.GetUnityPackages(),
	}

	return nil
}

type PhotonPropUser struct {
	ID                             string            `json:"id"`
	DisplayName                    string            `json:"displayName"`
	DeveloperType                  string            `json:"developerType"`
	CurrentAvatarImageUrl          string            `json:"currentAvatarImageUrl"`
	CurrentAvatarThumbnailImageUrl string            `json:"currentAvatarThumbnailImageUrl"`
	UserIcon                       string            `json:"userIcon"`
	LastPlatform                   string            `json:"last_platform"`
	Status                         string            `json:"status"`
	StatusDescription              string            `json:"statusDescription"`
	Bio                            string            `json:"bio"`
	Tags                           []string          `json:"tags"`
	UnityPackages                  []APIUnityPackage `json:"unityPackages"`
	AllowAvatarCopying             bool              `json:"allowAvatarCopying"`
}

type PhotonPropAvatarDict struct {
	ID                string            `json:"id"`
	AssetUrl          string            `json:"assetUrl"`
	AuthorId          string            `json:"authorId"`
	AuthorName        string            `json:"authorName"`
	UpdatedAt         string            `json:"updated_at"`
	Description       string            `json:"description"`
	Featured          bool              `json:"featured"`
	ImageUrl          string            `json:"imageUrl"`
	ThumbnailImageUrl string            `json:"thumbnailImageUrl"`
	Name              string            `json:"name"`
	ReleaseStatus     string            `json:"releaseStatus"`
	Version           int               `json:"version"`
	Tags              []string          `json:"tags"`
	UnityPackages     []APIUnityPackage `json:"unityPackages"`
}

type PhotonConfig struct {
	MaxAccountsPerIPAddress int         `json:"maxAccsPerIp"`
	RateLimitList           map[int]int `json:"ratelimitList"`
	RatelimiterActive       bool        `json:"ratelimiterActive"`
}
