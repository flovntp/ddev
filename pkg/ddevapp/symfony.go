package ddevapp

import (
	"fmt"
	"github.com/drud/ddev/pkg/fileutil"
	"github.com/drud/ddev/pkg/nodeps"
	"os"
	"path/filepath"
)

func symfonyPostStartAction(app *DdevApp) error {
	fmt.Println("Beginning symfonyPostStartAction...")
	if !app.DisableSettingsManagement {
		if _, err := app.CreateSettingsFile(); err != nil {
			return fmt.Errorf("failed to write settings file %s: %v", app.SiteDdevSettingsFile, err)
		}
	}
	// We won't touch env if disable_settings_management: true
	if app.DisableSettingsManagement {
		return nil
	}
	_, envText, err := ReadProjectEnvFile(app)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to read .env file: %v", err)
	}
	port := "3306"
	dbConnection := "mysql"
	if app.Database.Type == nodeps.Postgres {
		fmt.Println("detected db is postgres")
		dbConnection = "pgsql"
		port = "5432"
	} else {
		fmt.Println("detected db is mysql")
	}
	envMap := map[string]string{
		"DATABASE_URL": dbConnection + "://db:db@db:" + port + "/db",
	}
	err = WriteProjectEnvFile(app, envMap, envText)
	if err != nil {
		return err
	} else {
		fmt.Println("Successfully wrote symfony info to .env file")
	}

	return nil
}

// getPHPUploadDir will return a custom upload dir if defined
func getSymfonyUploadDir(app *DdevApp) string {
	return app.UploadDir
}

// isSymfonyApp returns true if the app is of type symfondy
func isSymfonyApp(app *DdevApp) bool {
	return fileutil.FileExists(filepath.Join(app.AppRoot, "symfony.lock"))
}

func symfonyConfigOverrideAction(app *DdevApp) error {
	app.Database.Type = nodeps.Postgres
	app.Database.Version = nodeps.Postgres14
	return nil
}
