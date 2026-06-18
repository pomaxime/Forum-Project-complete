# ForumHub — Clone Reddit en Go

## 📖 Présentation

ForumHub est une plateforme de discussion inspirée de Reddit, développée en **Go (Golang)** avec une architecture modulaire. Les utilisateurs peuvent créer un compte, publier des sujets, commenter, réagir aux publications et filtrer le contenu par catégories.

Le projet met l'accent sur :

- 🔐 La sécurité des utilisateurs
- 🗄️ La persistance des données avec SQLite
- 🧩 Une architecture claire et maintenable
- 🎨 Une interface moderne inspirée des plateformes communautaires actuelles

---

## 🚀 Fonctionnalités

### Authentification

- Inscription utilisateur
- Connexion sécurisée
- Déconnexion
- Sessions persistantes
- Cookies HttpOnly

### Publications

- Création de posts
- Consultation des posts
- Affichage détaillé d'un sujet
- Filtrage par catégories

### Commentaires

- Ajout de commentaires
- Affichage des commentaires liés à un post

### Système de votes

- Like sur les publications
- Dislike sur les publications
- Like sur les commentaires
- Dislike sur les commentaires

### Profil utilisateur

- Consultation des publications personnelles
- Historique d'activité

---

## 🏗️ Architecture du projet

```text
/forum
│
├── config/                 # Configuration générale
├── database/               # Initialisation SQLite et schéma
├── handlers/               # Contrôleurs HTTP
├── middleware/             # Authentification et sécurité
├── models/                 # Structures métier
├── repository/             # Accès aux données
├── routes/                 # Déclaration des routes
├── services/               # Logique métier
├── sessions/               # Gestion des sessions
├── static/                 # CSS, images, ressources
├── templates/              # Templates HTML
├── utils/                  # Outils et validations
│
├── main.go                 # Point d'entrée
├── go.mod
└── go.sum