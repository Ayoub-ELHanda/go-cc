package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync" // Importez le package sync pour utiliser WaitGroups

	"github.com/ayoubmcw/cc-go.git/pkg"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}

	username := os.Getenv("GITHUB_USER")
	token := os.Getenv("GITHUB_TOKEN")

	repos, err := pkg.FetchRepositoriesWithToken(username, token)
	if err != nil {
		log.Fatal(err)
	}
	err = writeReposToCSV(repos)
	if err != nil {
		log.Fatal(err)
	}

	// Chemin du répertoire de clonage
	cloneDir := "clone-repo"

	var wg sync.WaitGroup

	// Créez un canal pour collecter les erreurs de goroutines
	errChan := make(chan error, len(repos))

	// Utilisez une goroutine pour chaque dépôt pour le clonage
	for _, repo := range repos {
		wg.Add(1) // Incrémentation du compteur WaitGroup

		go func(repo pkg.Repository) {
			defer wg.Done() // Décrémentation du compteur WaitGroup lorsque la goroutine est terminée

			err := cloneRepository(repo.CloneURL, cloneDir)
			if err != nil {
				errChan <- err // Envoyer l'erreur au canal
				return
			}

			fmt.Printf("Le dépôt %s a été cloné avec succès.\n", repo.Name)

			// Effectuer un Git Pull sur la dernière branche modifiée
			err = gitPullLatestBranch(filepath.Join(cloneDir, repo.Name))
			if err != nil {
				errChan <- err // Envoyer l'erreur au canal
				return
			}

			fmt.Printf("Git Pull effectué avec succès pour le dépôt %s.\n", repo.Name)
		}(repo)
	}

	// Attendre que toutes les goroutines de clonage se terminent
	wg.Wait()

	// Fermez le canal d'erreur lorsque toutes les goroutines sont terminées
	close(errChan)

	// Parcourez le canal d'erreur pour vérifier s'il y a des erreurs
	for err := range errChan {
		fmt.Printf("Erreur : %v\n", err)
	}

	// Créer un gestionnaire de route HTTP pour le téléchargement des archives ZIP
	http.HandleFunc("/download-repo", func(w http.ResponseWriter, r *http.Request) {
		// Récupérer le nom du dépôt à télécharger depuis le paramètre "repo" dans l'URL
		repoName := r.URL.Query().Get("repo")

		// Compresser le dépôt et le renvoyer en tant que fichier ZIP
		err := createZipArchive(filepath.Join(cloneDir, repoName), repoName+".zip", w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erreur lors de la création de l'archive ZIP : %v", err), http.StatusInternalServerError)
		}
	})

	// Lancer le serveur HTTP
	port := ":8080" // Vous pouvez spécifier un port différent si nécessaire
	fmt.Printf("Serveur en cours d'exécution sur le port %s...\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
func writeReposToCSV(repos []pkg.Repository) error {
	file, err := os.Create("repositories.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Écrire l'en-tête du fichier CSV
	err = writer.Write([]string{"Name", "Clone URL", "Description", "Last Updated"})
	if err != nil {
		return err
	}

	// Écrire les informations de chaque dépôt dans le fichier CSV
	for _, repo := range repos {
		err = writer.Write([]string{repo.Name, repo.CloneURL, repo.Description, repo.UpdatedAt})
		if err != nil {
			return err
		}
	}

	fmt.Println("Les informations des dépôts ont été écrites dans repositories.csv")
	return nil
}

// ... (les fonctions cloneRepository, gitPullLatestBranch, gitFetchAll et createZipArchive restent inchangées)
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
func getRepositories(username, token string) ([]pkg.Repository, error) {
	// Utilisez le token API GitHub pour cloner également les dépôts privés (si le token est fourni)
	var repos []pkg.Repository
	var err error

	if token != "" {
		repos, err = pkg.FetchRepositoriesWithToken(username, token)
	} else {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return repos, nil
}

func createZipArchive(sourceDir, zipFileName string, w io.Writer) error {
	archive := zip.NewWriter(w)
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
