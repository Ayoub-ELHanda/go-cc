package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ayoubmcw/cc-go.git/pkg"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}

	user := os.Getenv("GITHUB_USER")
	token := os.Getenv("GITHUB_TOKEN")

	repos, err := pkg.FetchRepositories(user, token)
	if err != nil {
		log.Fatal(err)
	}

	// Chemin du répertoire de clonage
	cloneDir := "clone-repo"

	for _, repo := range repos {
		err := cloneRepository(repo.CloneURL, cloneDir)
		if err != nil {
			fmt.Printf("Erreur lors du clonage du dépôt %s : %v\n", repo.Name, err)
			continue // Passe au dépôt suivant en cas d'erreur de clonage
		}

		fmt.Printf("Le dépôt %s a été cloné avec succès.\n", repo.Name)

		// Effectuer un Git Pull sur la dernière branche modifiée
		err = gitPullLatestBranch(filepath.Join(cloneDir, repo.Name))
		if err != nil {
			fmt.Printf("Erreur lors du Git Pull pour le dépôt %s : %v\n", repo.Name, err)
		} else {
			fmt.Printf("Git Pull effectué avec succès pour le dépôt %s.\n", repo.Name)
		}
	}

	// Effectuer un Git Fetch pour récupérer toutes les références de branches
	for _, repo := range repos {
		err := gitFetchAll(filepath.Join(cloneDir, repo.Name))
		if err != nil {
			fmt.Printf("Erreur lors du Git Fetch pour le dépôt %s : %v\n", repo.Name, err)
		} else {
			fmt.Printf("Git Fetch effectué avec succès pour le dépôt %s.\n", repo.Name)
		}
	}

	// Créer une archive ZIP à la fin du traitement
	err = createZipArchive(cloneDir, "repositories.zip")
	if err != nil {
		fmt.Printf("Erreur lors de la création de l'archive ZIP : %v\n", err)
	} else {
		fmt.Println("L'archive ZIP a été créée avec succès.")
	}
}

func cloneRepository(cloneURL, cloneDir string) error {
	// Assurez-vous que le répertoire de clonage existe
	err := os.MkdirAll(cloneDir, 0755)
	if err != nil {
		return err
	}

	// Divisez le CloneURL pour obtenir le nom du dépôt
	parts := strings.Split(cloneURL, "/")
	repoName := parts[len(parts)-1]

	// Chemin complet du répertoire du dépôt cloné
	repoPath := fmt.Sprintf("%s/%s", cloneDir, repoName)

	// Utilisez la commande "git clone" pour cloner le dépôt
	cmd := exec.Command("git", "clone", cloneURL, repoPath)

	// Exécutez la commande
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func gitPullLatestBranch(repoPath string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoPath
	return cmd.Run()
}

func gitFetchAll(repoPath string) error {
	cmd := exec.Command("git", "fetch", "--all")
	cmd.Dir = repoPath
	return cmd.Run()
}

func createZipArchive(sourceDir, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorer les répertoires
		if info.IsDir() {
			return nil
		}

		// Créer un nouveau fichier dans l'archive
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		fileInArchive, err := archive.Create(relPath)
		if err != nil {
			return err
		}

		// Copier le contenu du fichier source dans l'archive
		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		_, err = io.Copy(fileInArchive, sourceFile)
		return err
	})
}
