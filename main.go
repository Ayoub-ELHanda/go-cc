package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

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

	err = writeReposToCSV(repos)
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
