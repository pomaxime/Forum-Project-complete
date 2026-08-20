INSERT OR IGNORE INTO users (username, email, password) VALUES
    ('alice',   'alice@example.com',   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
    ('bob',     'bob@example.com',     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
    ('charlie', 'charlie@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy');

INSERT OR IGNORE INTO posts (user_id, title, content, category) VALUES
    (1, 'Bienvenue sur Reddick !',
        'Ceci est le premier post du forum. N''hésitez pas à créer un compte et à participer !',
        'general'),
    (2, 'Go vs Rust : quel langage choisir en 2026 ?',
        'Je travaille sur un nouveau projet backend et j''hésite entre Go et Rust. Go semble plus simple à prendre en main mais Rust offre de meilleures garanties mémoire. Vos avis ?',
        'tech'),
    (1, 'Les meilleurs jeux indé du moment',
        'Voici ma liste des jeux indépendants à ne pas manquer cette année : Hollow Knight Silksong (enfin !), Hades II, et quelques autres pépites. Qu''est-ce que vous jouez en ce moment ?',
        'jeux');

INSERT OR IGNORE INTO comments (post_id, user_id, content) VALUES
    (1, 2, 'Super initiative ! Hâte de voir la communauté grandir.'),
    (1, 3, 'Merci pour ce forum, l''interface est vraiment propre.'),
    (2, 3, 'Go est mon choix de cœur pour un démarrage rapide. La concurrence avec les goroutines est vraiment intuitive.'),
    (2, 1, 'Rust pour la performance et la sécurité, Go pour la simplicité. Tout dépend de tes contraintes !'),
    (3, 2, 'Hades II est incroyable, le gameplay est encore plus fluide que le premier.');