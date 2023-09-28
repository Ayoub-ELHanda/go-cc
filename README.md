# go-cc
Ce projet est une application Go qui vous permet de cloner les dépôts GitHub d'un utilisateur ou d'une organisation, d'effectuer un Git Pull sur la dernière branche modifiée, de récupérer toutes les références de branches avec Git Fetch, et de créer une archive ZIP des dépôts clonés.

## Utilisation

1. Configurez vos informations GitHub en créant un fichier `.env` (vous pouver copier les info sur ".env.dst") à la racine du projet avec les variables d'environnement suivantes :
 dans votre .env file modifier ca : 
  GITHUB_USER=VotreNomDUtilisateurGitHub
  GITHUB_TOKEN=VotreTokenGitHub (Optionnel pour les dépôts privés)

*vous pouver utiliser docker avec la commond 'docker run -d go-cc:latest'

3. Lancez l'application en exécutant `go run main.go`.
4. Accédez à l'interface web en ouvrant un navigateur et en allant à l'adresse `http://localhost:8080` (ou au port que vous avez spécifié).

5. Une fois le clonage terminé, vous pouvez télécharger les archives ZIP en allant à l'adresse `http://localhost:8080/download-repo` des dépôts clonés depuis l'interface web.

## Fonctionnalités
- génerer un csv pour avec une liste de l’ensemble des informations des dépo  récupérées  
- Clonage des dépôts publics de GitHub.
- Clonage des dépôts privés (si un token API GitHub est fourni).
- Git Pull sur la dernière branche modifiée.
- Git Fetch pour récupérer toutes les références de branches.
- Téléchargement des archives ZIP des dépôts clonés.

## Auteur

[Auteur du Projet]

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE.md](LICENSE.md) pour plus de détails.
