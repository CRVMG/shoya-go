package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/george/shoya-go/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type World struct {
	BaseModel
	AuthorID      string `json:"authorId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Capacity      int    `json:"capacity"`
	ImageID       string
	Image         File
	ReleaseStatus ReleaseStatus       `json:"releaseStatus" gorm:"default:'private'"`
	Tags          pq.StringArray      `json:"tags" gorm:"type:text[] NOT NULL;default: '{}'::text[]"`
	Version       int                 `json:"version" gorm:"type:bigint NOT NULL;default:0"`
	UnityPackages []WorldUnityPackage `json:"unityPackages" gorm:"foreignKey:BelongsToAssetID"`
}

func (w *World) BeforeCreate(*gorm.DB) (err error) {
	w.ID = "wrld_" + uuid.New().String()
	return
}

func GetWorldById(id string) (*World, error) {
	var w *World
	tx := config.DB.Preload(clause.Associations).
		Preload("Image").
		Preload("Image.Versions").
		Preload("Image.Versions.FileDescriptor").
		Preload("Image.Versions.DeltaDescriptor").
		Preload("Image.Versions.SignatureDescriptor").
		Preload("UnityPackages.File").
		Preload("UnityPackages.File.Versions").
		Preload("UnityPackages.File.Versions.FileDescriptor").
		Preload("UnityPackages.File.Versions.DeltaDescriptor").
		Preload("UnityPackages.File.Versions.SignatureDescriptor").
		Where("id = ?", id).First(&w)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrWorldNotFound
		}
		return nil, tx.Error
	}

	return w, nil
}

// GetAuthor returns a pointer to the world author's User.
func (w *World) GetAuthor() (*User, error) {
	var u User

	tx := config.DB.Where("id = ?", w.AuthorID).Find(&u)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &u, nil
}

// GetImageUrl returns the Url present in the Image field.
func (w *World) GetImageUrl() string {
	return w.Image.GetLatestVersion().GetFileUrl()
}

// GetThumbnailImageUrl returns the Url present in the Image field.
func (w *World) GetThumbnailImageUrl() string {
	return w.Image.GetLatestVersion().GetFileUrl()
}

// GetLatestAssetUrl iterates through a World's UnityPackages and returns the Url present in the File
// of the UnityPackage with the highest version number.
func (w *World) GetLatestAssetUrl() string {
	var assetUrl string
	maxVersion := 0
	for _, pkg := range w.UnityPackages {
		if pkg.Version >= maxVersion {
			assetUrl = pkg.File.GetLatestVersion().GetFileUrl()
		}
	}

	return assetUrl
}

// GetUnityPackages returns a list of APIUnityPackage.
func (w *World) GetUnityPackages(withAssetUrls bool) []APIUnityPackage {
	var pkgs []APIUnityPackage
	for _, pkg := range w.UnityPackages {
		pkgs = append(pkgs, *pkg.GetAPIUnityPackage(withAssetUrls))
	}

	return pkgs
}

func (w *World) GetAPIWorld() (*APIWorld, error) {
	a, err := w.GetAuthor()
	if err != nil {
		return nil, err
	}

	return &APIWorld{
		ID:                  w.ID,
		AuthorID:            a.ID,
		AuthorName:          a.DisplayName,
		Capacity:            w.Capacity,
		CreatedAt:           time.Unix(w.CreatedAt, 0).UTC().Format(time.RFC3339Nano),
		Description:         w.Description,
		Favorites:           0, // TODO: Implement favorites.
		Heat:                0, // Intentionally hardcoded to zero; Will not implement.
		ImageUrl:            w.GetImageUrl(),
		Instances:           [][]string{},
		PublicationDate:     time.Unix(w.CreatedAt, 0).Format(time.RFC3339),
		LabsPublicationDate: "", // Intentionally hardcoded to an empty string; Will not implement.
		Name:                w.Name,
		Occupants:           0,        // TODO: Implement instances + overall occupancy.
		Organization:        "vrchat", // It's *always* vrchat.
		PreviewYoutubeId:    "",       // TODO: This is almost never used, and is only available on the web. Low priority.
		PrivateOccupants:    0,        // TODO: Implement instances + overall occupancy.
		PublicOccupants:     0,        // TODO: Implement instances + overall occupancy.
		ReleaseStatus:       w.ReleaseStatus,
		Tags:                w.Tags,
		ThumbnailImageUrl:   w.GetThumbnailImageUrl(),
		Version:             w.Version,
		UnityPackages:       w.GetUnityPackages(false),
		Visits:              0, // TODO: Implement metrics.
		UpdatedAt:           time.Unix(w.UpdatedAt, 0).UTC().Format(time.RFC3339Nano),
	}, nil
}
func (w *World) GetAPIWorldWithPackages() (*APIWorldWithPackages, error) {
	a, err := w.GetAPIWorld()
	if err != nil {
		return nil, err
	}
	return &APIWorldWithPackages{
		APIWorld:      *a,
		AssetUrl:      w.GetLatestAssetUrl(),
		UnityPackages: w.GetUnityPackages(true),
	}, nil
}

type APIWorld struct {
	ID                  string            `json:"id"`
	AuthorID            string            `json:"authorId"`
	AuthorName          string            `json:"authorName"`
	Capacity            int               `json:"capacity"`
	CreatedAt           string            `json:"created_at"`
	Description         string            `json:"description"`
	Favorites           int               `json:"favorites"`
	Featured            bool              `json:"featured"`
	Namespace           string            `json:"namespace"`
	Popularity          int               `json:"popularity"`
	Heat                int               `json:"heat"`
	ImageUrl            string            `json:"imageUrl"`
	Instances           [][]string        `json:"instances"`
	PublicationDate     string            `json:"publicationDate"`
	LabsPublicationDate string            `json:"labsPublicationDate"`
	Name                string            `json:"name"`
	Occupants           int               `json:"occupants"`
	Organization        string            `json:"organization"`
	PluginUrlObject     interface{}       `json:"pluginUrlObject"`
	PreviewYoutubeId    string            `json:"previewYoutubeId"`
	PrivateOccupants    int               `json:"privateOccupants"`
	PublicOccupants     int               `json:"publicOccupants"`
	ReleaseStatus       ReleaseStatus     `json:"releaseStatus"`
	Tags                []string          `json:"tags"`
	ThumbnailImageUrl   string            `json:"thumbnailImageUrl"`
	Version             int               `json:"version"`
	Visits              int               `json:"visits"`
	UnityPackages       []APIUnityPackage `json:"unityPackages"`
	UpdatedAt           string            `json:"updated_at"`
}

type APIWorldWithPackages struct {
	APIWorld
	AssetUrl              string            `json:"assetUrl"`
	PluginUrl             string            `json:"pluginUrl"`      // Required for older builds of the game (e.g.: 333)
	AssetUrlObject        interface{}       `json:"assetUrlObject"` // Always an empty object.
	UnityPackages         []APIUnityPackage `json:"unityPackages"`
	UnityPackageUrlObject interface{}       `json:"unityPackageUrlObject"` // Always an empty object.
}
