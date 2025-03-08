package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	gastromaticprocessor "trinkgeldApp/gastromaticProcessor"
	gastronoviprocessor "trinkgeldApp/gastronoviProcessor"
	"trinkgeldApp/models"

	"github.com/labstack/echo/v4"
	"github.com/pocketbase/pocketbase"
)

const (
	typeGastromaticReport = "gastromatic"
	typeGastronoviReport  = "gastronovi"
)

// App struct will hold PocketBase instance and other services
type AppContext struct {
	DB *pocketbase.PocketBase
}

func (a *AppContext) SubmitGastromatic(c echo.Context) error {

	fileType, err := checkFileType(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if fileType != typeGastromaticReport {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Invalid file type %s", fileType),
		})
	}

	filePath, err := upload(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	shifts, err := gastromaticprocessor.ProcessGastromaticFile(filePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	a.UploadShiftsForPeriod(shifts)

	return c.JSON(http.StatusOK, map[string]bool{
		"ok": true,
	})

}

func (a *AppContext) SubmitGastronovi(c echo.Context) error {

	fileType, err := checkFileType(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if fileType != typeGastronoviReport {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Invalid file type %s", fileType),
		})
	}

	location, err := getLocationFromRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	filePath, err := upload(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	tips, err := gastronoviprocessor.ProcessGastronoviFile(filePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	a.UploadTipsForPeriod(tips, location)

	return c.JSON(http.StatusOK, map[string]bool{
		"ok": true,
	})

}

// Returns the file name and type
func upload(c echo.Context) (string, error) {
	// Read form fields
	name := c.FormValue("name")

	// Source
	file, err := c.FormFile(name)
	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique file name
	timestamp := time.Now().Unix()
	ext := filepath.Ext(file.Filename)
	base := file.Filename[:len(file.Filename)-len(ext)]
	uniqueFileName := fmt.Sprintf("%s_%d%s", base, timestamp, ext)

	// Destination
	dst, err := os.Create("./cache/" + uniqueFileName)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return "./cache/" + uniqueFileName, nil
}

func checkFileType(c echo.Context) (string, error) {
	fileType := c.FormValue("type")

	switch fileType {
	case typeGastromaticReport:
		return typeGastromaticReport, nil
	case typeGastronoviReport:
		return typeGastronoviReport, nil
	default:
		return "", fmt.Errorf("Invalid file type %s", fileType)
	}
}

func getLocationFromRequest(c echo.Context) (string, error) {
	location := c.FormValue("location")

	switch location {
	case models.AltstadtLocationID:
		return models.AltstadtLocationID, nil
	case models.CampusLocationID:
		return models.CampusLocationID, nil
	case models.NordendLocationID:
		return models.NordendLocationID, nil
	case models.HdlLocationID:
		return models.HdlLocationID, nil
	default:
		return "", fmt.Errorf("Invalid location %s", location)
	}
}
