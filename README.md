# 🚀 ForumHub

> Une plateforme communautaire moderne inspirée de Reddit, développée en Go.

---

## 📖 Présentation

ForumHub est une application web communautaire permettant aux utilisateurs de partager des publications, échanger à travers des commentaires et interagir grâce à un système de votes.

Le projet s'inspire du fonctionnement de Reddit tout en proposant une identité visuelle propre et une architecture conçue pour être maintenable, sécurisée et évolutive.

L'objectif principal était de développer une plateforme complète en appliquant les bonnes pratiques du développement web moderne :

- 🔐 Sécurité des utilisateurs
- 🏗️ Architecture modulaire
- 🗄️ Gestion de base de données
- ⚡ Performance et simplicité grâce à Go
- 🎨 Interface utilisateur moderne

---

# ✨ Fonctionnalités

## 👤 Gestion des utilisateurs

- Création de compte
- Connexion sécurisée
- Déconnexion
- Gestion des sessions
- Cookies sécurisés

## 📝 Publications

- Création de posts
- Consultation des publications
- Affichage détaillé d'un sujet
- Filtrage par catégories
- Consultation des publications d'un utilisateur

## 💬 Commentaires

- Ajout de commentaires
- Consultation des commentaires associés à un post
- Interaction avec les commentaires

## 👍 Système de votes

Les utilisateurs peuvent :

- Like une publication
- Dislike une publication
- Like un commentaire
- Dislike un commentaire

Le score est mis à jour automatiquement afin de mettre en avant les contenus les plus appréciés.

## 👨‍💻 Profil utilisateur

Chaque utilisateur dispose d'un espace personnel permettant de consulter :

- Ses publications
- Son activité récente
- Ses interactions avec le forum

---

# 🛠️ Technologies utilisées

| Technologie | Utilisation |
|------------|-------------|
| Go (Golang) | Développement du backend |
| SQLite | Base de données |
| HTML | Structure des pages |
| CSS | Mise en forme |
| JavaScript | Interactivité |
| bcrypt | Chiffrement des mots de passe |
| UUID | Gestion des sessions |

---

# 🏗️ Architecture du projet

Le projet est organisé selon une architecture modulaire afin de séparer clairement les responsabilités de chaque composant.

```text
ForumHub/
│
├── config/
├── database/
├── handlers/
├── middleware/
├── models/
├── repository/
├── routes/
├── services/
├── sessions/
├── static/
├── templates/
│
├── main.go
├── go.mod
└── go.sum
```

---

## 📂 Description des dossiers

### ⚙️ config/

Contient les paramètres de configuration de l'application.

---

### 🗄️ database/

Gère :

- La connexion à SQLite
- L'initialisation de la base
- La création des tables

---

### 🌐 handlers/

Traite les requêtes HTTP reçues depuis l'interface utilisateur.

Exemples :

- Connexion
- Inscription
- Création de publication
- Commentaires
- Réactions

---

### 🛡️ middleware/

Ajoute des couches de sécurité entre l'utilisateur et l'application.

Exemples :

- Vérification des sessions
- Contrôle d'accès
- Protection des routes privées

---

### 📦 models/

Définition des structures de données utilisées dans l'application.

Exemples :

- User
- Post
- Comment
- Session

---

### 🗃️ repository/

Couche chargée des échanges avec la base de données.

Toutes les requêtes SQL sont centralisées ici.

---

### ⚙️ services/

Contient la logique métier du projet.

Cette couche permet de séparer les règles métier du code HTTP.

---

### 🛣️ routes/

Déclaration des routes de l'application.

Exemples :

- /
- /login
- /register
- /profile

---

### 🎨 static/

Contient les ressources statiques :

- CSS
- JavaScript
- Images
- Logos

---

### 📄 templates/

Contient les pages HTML affichées aux utilisateurs.

---

# 🚀 Installation

## 📋 Prérequis

Vérifier que Go est installé :

```bash
go version
```

Version recommandée :

```text
Go 1.24 ou supérieur
```

---

## 📥 Cloner le projet

```bash
git clone <url-du-repository>
```

```bash
cd ForumHub
```

---

## 📦 Installer les dépendances

```bash
go mod tidy
```

---

## ▶️ Lancer l'application

```bash
go run main.go
```

ou

```bash
go run .
```

Une fois démarrée, l'application est accessible depuis :

```text
http://localhost:8080
```

---

# 🔨 Compiler le projet

### Windows

```bash
go build -o forumhub.exe
```

Exécution :

```bash
forumhub.exe
```

---

### Linux / MacOS

```bash
go build -o forumhub
```

Exécution :

```bash
./forumhub
```

---

# 🧪 Tests

Exécuter tous les tests :

```bash
go test ./...
```

Mode détaillé :

```bash
go test -v ./...
```

Mesurer la couverture :

```bash
go test -cover ./...
```

---

# 🗄️ Base de données

Le projet utilise SQLite pour stocker l'ensemble des informations nécessaires au fonctionnement du forum.

### 👤 Utilisateurs

- Nom d'utilisateur
- Adresse email
- Mot de passe chiffré

### 📝 Publications

- Titre
- Contenu
- Catégorie
- Auteur
- Date de création

### 💬 Commentaires

- Contenu
- Auteur
- Publication associée

### 🔑 Sessions

- Identifiant de session
- Utilisateur associé
- Date d'expiration

### 👍 Réactions

- Likes
- Dislikes

---

# 🔒 Sécurité

La sécurité a été intégrée dès la conception du projet.

## 🔐 Chiffrement des mots de passe

Les mots de passe ne sont jamais stockés en clair.

Ils sont automatiquement protégés grâce à **bcrypt**, un algorithme de hachage reconnu et largement utilisé.

---

## 🛡️ Protection contre les injections SQL

Toutes les requêtes utilisent des paramètres préparés afin d'empêcher l'exécution de requêtes malveillantes.

Exemple :

```sql
SELECT * FROM users WHERE email = ?
```

---

## 🍪 Gestion sécurisée des sessions

Après connexion :

1. Une session unique est créée.
2. Un cookie est envoyé au navigateur.
3. Chaque requête vérifie la validité de cette session.

---

## 🚫 Cookies HttpOnly

Les cookies de session sont inaccessibles depuis JavaScript.

Cette mesure réduit fortement les risques liés aux attaques XSS.

---

## 🧱 Headers de sécurité

Plusieurs protections HTTP sont appliquées :

- Content-Security-Policy
- X-Frame-Options
- Referrer-Policy
- X-Content-Type-Options

---

# ⚙️ Fonctionnement général

## Création d'un compte

1. L'utilisateur remplit le formulaire.
2. Le mot de passe est chiffré.
3. Les données sont enregistrées.
4. Le compte est créé.

---

## Connexion

1. Vérification des identifiants.
2. Création d'une session.
3. Attribution d'un cookie sécurisé.

---

## Création d'une publication

1. Vérification de l'authentification.
2. Validation des données.
3. Enregistrement en base de données.
4. Affichage sur le forum.

---

## Publication d'un commentaire

1. Sélection d'un post.
2. Saisie du commentaire.
3. Enregistrement en base.
4. Affichage immédiat sous la publication.

---

# 🎯 Choix techniques

## Pourquoi Go ?

- Très performant
- Faible consommation mémoire
- Compilation rapide
- Gestion native de la concurrence
- Déploiement simplifié

---

## Pourquoi SQLite ?

- Léger
- Rapide à mettre en place
- Aucune installation serveur nécessaire
- Adapté aux projets de petite et moyenne taille

---

## Pourquoi une architecture modulaire ?

Cette organisation permet :

- Une maintenance plus simple
- Une meilleure lisibilité du code
- Une évolution facilitée du projet
- Une séparation claire des responsabilités

---

# 🚧 Améliorations futures

Parmi les évolutions envisageables :

- 🔔 Notifications en temps réel
- 🖼️ Upload d'images
- 🌙 Mode sombre
- 🔍 Recherche avancée
- 📬 Messagerie privée
- 🚩 Signalement de contenu
- 👮 Outils de modération
- 🌍 API REST publique
- 🐳 Dockerisation du projet

---

# 👥 Équipe

Projet réalisé dans le cadre d'un apprentissage du développement web et de la cybersécurité.

L'objectif était de mettre en pratique les concepts fondamentaux de :

- Développement backend
- Gestion des bases de données
- Sécurisation d'applications web
- Architecture logicielle

---

# 📄 Licence

Projet réalisé à des fins pédagogiques et d'apprentissage.