-- Nettoyage des tables
DELETE FROM comments;
DELETE FROM posts;
DELETE FROM sessions;
DELETE FROM users;

-- Réinitialisation des IDs
DELETE FROM sqlite_sequence WHERE name IN ('users', 'posts', 'comments');

---------------------------------------------------------
-- 1) Utilisateurs
---------------------------------------------------------
INSERT INTO users (username, email, password) VALUES
('yves', 'yves@example.com', 'hashed_password_yves'),
('ewen', 'ewen@example.com', 'hashed_password_ewen'),
('benjamin', 'benjamin@example.com', 'hashed_password_benjamin');

---------------------------------------------------------
-- 2) Sessions (facultatif)
---------------------------------------------------------
INSERT INTO sessions (id, user_id, expires_at) VALUES
('session_yves', 1, DATETIME('now', '+1 day')),
('session_ewen', 2, DATETIME('now', '+1 day'));

---------------------------------------------------------
-- 3) Posts
---------------------------------------------------------
INSERT INTO posts (user_id, title, content, category) VALUES
(1, 'Bienvenue sur notre forum', 'Salut tout le monde, bienvenue sur le projet !', 'general'),
(2, 'Besoin d’aide en Go', 'Je comprends pas trop les goroutines, quelqu’un peut m’aider ?', 'programming'),
(3, 'Organisation du projet', 'On devrait se faire un point pour répartir les tâches.', 'team'),
(1, 'Tutoriel SQLite', 'Voici comment utiliser SQLite avec Go.', 'database'),
(2, 'Partage de ressources', 'Voici quelques liens utiles pour progresser en Go.', 'resources');

---------------------------------------------------------
-- 4) Commentaires
---------------------------------------------------------
INSERT INTO comments (post_id, user_id, content) VALUES
(1, 2, 'Merci Yves, c’est propre !'),
(1, 3, 'On avance bien les gars.'),
(2, 1, 'Oui je peux t’expliquer les goroutines.'),
(2, 3, 'Bonne question, ça m’intéresse aussi.'),
(3, 1, 'Bonne idée Benjamin.'),
(4, 2, 'Super tuto, merci Yves.'),
(5, 1, 'Très utile, merci Ewen.'),
(5, 3, 'Je vais regarder ça.');
