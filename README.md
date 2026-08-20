# Forum Reddit Clone — Module Authentification & Sécurité

## Structure du projet

```
/forum
├── main.go                  ← Point d'entrée, routeur HTTP
├── go.mod / go.sum
├── main_test.go             ← Tests (register, login, logout, middleware, sessions)
│
├── database/
│   ├── db.go                ← Connexion SQLite
│   └── schema.go            ← Création des tables (users, sessions, posts)
│
├── handlers/
│   ├── auth_handler.go      ← Struct AuthHandler partagée
│   ├── register.go          ← Inscription
│   ├── login.go             ← Connexion + création session/cookie
│   ├── logout.go            ← Déconnexion + suppression session/cookie
│   └── post_handler.go      ← Index, création et affichage de posts
│
├── middleware/
│   └── auth.go              ← Vérification session + expiration + injection userID
│
├── utils/
│   ├── password.go          ← Hash bcrypt + vérification
│   └── validation.go        ← Validation formulaires (register, login, posts)
│
├── templates/               ← HTML (Go template)
│   ├── register.html
│   ├── login.html
│   ├── index.html
│   ├── create_post.html
│   └── post.html
│
└── static/
    └── style.css
```

## Installation

```bash
# 1. Cloner / copier le projet
cd forum

# 2. Installer les dépendances
go mod tidy

# 3. Lancer le serveur
go run main.go
# → http://localhost:8080
```

## Lancer les tests

```bash
go test -v ./...
```

## Ce que couvre ce projet

| Fonctionnalité       | Fichier                        |
|----------------------|--------------------------------|
| SQLite               | database/db.go + schema.go     |
| Tables users/sessions| database/schema.go             |
| Hash bcrypt          | utils/password.go              |
| Validation           | utils/validation.go            |
| Register             | handlers/register.go           |
| Login + Session      | handlers/login.go              |
| Logout               | handlers/logout.go             |
| Cookie HttpOnly      | handlers/login.go              |
| Middleware auth      | middleware/auth.go             |
| Routes protégées     | main.go                        |
| Protection SQL inj.  | Requêtes préparées (? partout) |
| Sessions expirées    | middleware/auth.go             |

## Points importants pour l'oral

**Pourquoi bcrypt ?**
bcrypt est une fonction de hachage lente par conception (coût configurable). Même si la base
de données est volée, le brute-force de tous les mots de passe prendrait des années.

**Différence cookie / session ?**
La session est stockée côté serveur (SQLite). Le cookie contient uniquement l'ID de session
(UUID opaque). Si le cookie est volé, on peut révoquer la session côté serveur.

**Pourquoi middleware ?**
Le middleware centralise la vérification auth en un seul endroit. On ne peut pas oublier de
protéger une route, et si la logique change (ex: ajouter 2FA), on le fait une seule fois.

**Comment éviter SQL Injection ?**
Toutes les requêtes utilisent des paramètres préparés (?) — jamais de concaténation de
strings avec des données utilisateur.

**Pourquoi HttpOnly sur le cookie ?**
Empêche JavaScript d'accéder au cookie, ce qui bloque les attaques XSS qui tenteraient
de voler la session.
